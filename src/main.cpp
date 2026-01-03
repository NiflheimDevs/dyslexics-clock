#include "main.h"
#include "Arduino.h"
#include "ArduinoJson.h"
#include "ArduinoJson/Deserialization/DeserializationError.hpp"
#include "ArduinoJson/Document/JsonDocument.hpp"
#include "HardwareSerial.h"
#include "WiFi.h"
#include "WiFiManager.h"
#include "freertos/projdefs.h"
#include "led.h"
#include <ctime>

WiFiManager wifiManager;

struct MqttMsg {
  String topic;
  String payload;
};

QueueHandle_t mqttQueue;

bool first_time_online = false;
bool wifiConnected = false;
const uint8_t sub_topics_count = 9;
String sub_topics[sub_topics_count] = {"devices/" + DEVICEID + "/alarms/create",
                                       "devices/" + DEVICEID + "/alarms/update",
                                       "devices/" + DEVICEID + "/alarms/delete",
                                       "devices/" + DEVICEID + "/alarms",
                                       "devices/" + DEVICEID + "/color",
                                       "devices/" + DEVICEID + "/volume",
                                       "devices/" + DEVICEID + "/ring",
                                       "devices/" + DEVICEID + "/silence",
                                       "devices/time"};

String pub_topics[6] = {"devices/" + DEVICEID + "/status",
                        "devices/alarms",
                        "devices/volume",
                        "devices/color",
                        "devices/" + DEVICEID + "/ringing",
                        "devices/" + DEVICEID + "/log"};

String getTopic(ActionPublish action) {
  switch (action) {
  case ACTION_STATUS:
    return pub_topics[0];
  case ACTION_GET_ALL_ALARMS:
    return pub_topics[1];
  case ACTION_RINGING:
    return pub_topics[4];
  case ACTION_LOG:
    return pub_topics[5];
  case ACTION_GET_COLOR:
    return pub_topics[3];
  case ACTION_GET_VOLUME:
    return pub_topics[2];
  default:
    return "";
  }
}

const uint8_t snoozePin = 26;
volatile int32_t stopPin = 25;

RTC_DS3231 rtc;
uint8_t volume = 26;
DFRobotDFPlayerMini player;
const char *mqtt_server = "api.dyslexics-clock.niflheimdevs.ir";
WiFiClient espClient;
PubSubClient client(espClient);
long lastMsg = 0;
char msg[50];
int value = 0;

Alarm *create_alarm(uint32_t id, uint32_t device_id, uint32_t timestamp,
                    bool is_repeat = false, Weekday *repeating_days = nullptr,
                    uint8_t repeating_days_count = 0) {
  Alarm *alarm = new Alarm();
  alarm->id = id;
  alarm->timestamp = timestamp;
  alarm->is_repeat = is_repeat;

  if (is_repeat && repeating_days != nullptr) {
    for (uint8_t i = 0; i < repeating_days_count && i < 7; i++) {
      alarm->repeating_days[i] = repeating_days[i];
    }
  }

  return alarm;
}

AlarmHeap alarmHeap;

void setup() {
  Serial.begin(115200);
  setupWifi();
  setupRTC();
  setupLED();
  setupDfPlayer();
  setupTouch();
  setupMQTT();
}

void setupRTC() {
  if (!rtc.begin()) {
    Serial.println("Couldn't find RTC");
    Serial.flush();
    while (1)
      delay(10);
  }

  if (rtc.lostPower()) {
    Serial.println("RTC lost power, let's set the time!");
    rtc.adjust(DateTime(F(__DATE__), F(__TIME__)));
  }
}

void setupMQTT() {

  mqttQueue = xQueueCreate(10,             // queue length (number of messages)
                           sizeof(MqttMsg) // size of each message
  );
  client.setServer(mqtt_server, 11883);
  client.setCallback(mbCallback);
  xTaskCreate(mqttTask, "mqtttask", 4096, NULL, 3, NULL);
}

void mbCallback(char *topic, byte *message, unsigned int length) {
  Serial.print("Message arrived on topic: ");
  Serial.print(topic);
  Serial.print(". Message: ");
  String messageString;

  for (int i = 0; i < length; i++) {
    Serial.print((char)message[i]);
    messageString += (char)message[i];
  }
  Serial.println();

  String stringTopic = String(topic);

  if (stringTopic == sub_topics[0]) {
    add_alarm(messageString);
  } else if (stringTopic == sub_topics[1]) {
    update_alarm(messageString);
  } else if (stringTopic == sub_topics[2]) {
    delete_alarm(messageString);
  } else if (stringTopic == sub_topics[3]) {
    add_alarms_batch(messageString);
  } else if (stringTopic == sub_topics[4]) {
    set_color(messageString);
  } else if (stringTopic == sub_topics[5]) {
    set_volume(messageString);
  } else if (stringTopic == sub_topics[6]) {
    ring(messageString);
  } else if (stringTopic == sub_topics[7]) {
    player.stop();
  } else if (stringTopic == sub_topics[8]) {
    sync_time(messageString);
  } else {
    Serial.println("No matching topic found.");
  }
}
Alarm *parseAlarmFromJson(const char *jsonString) {
  JsonDocument doc;

  DeserializationError error = deserializeJson(doc, jsonString);
  if (error) {
    return nullptr;
  }

  if (!doc["id"].is<uint32_t>()) {
    return nullptr;
  }

  if (!doc["time"].is<const char *>()) {
    return nullptr;
  }

  Alarm *alarm = new Alarm();
  if (!alarm)
    return nullptr;

  alarm->id = doc["id"].as<uint32_t>();

  // Parse ISO8601 time (most common format from Go)
  const char *timeStr = doc["time"];
  // We expect something like: "2025-04-10T14:30:00"  or with Z / offset
  alarm->timestamp = DateTime(timeStr); // RTClib's DateTime can parse ISO8601

  if (!alarm->timestamp.isValid()) {
    delete alarm;
    return nullptr;
  }

  // is_repeat – defaults to false if missing
  alarm->is_repeat = doc["is_repeat"] | false;

  // repeating_days – array of 0..6 (Sunday=0)
  // We clear everything first
  for (int i = 0; i < 7; i++) {
    alarm->repeating_days[i] = Weekday::SUNDAY; // or false equivalent
  }

  alarm->repeating_days_count = 0;
  if (doc["days"].is<JsonArray>()) {
    JsonArray days = doc["days"];
    for (JsonVariantConst v : days) {
      int day = v.as<int>();
      if (day >= 0 && day <= 6) {
        alarm->repeating_days[day] = static_cast<Weekday>(day);
        alarm->repeating_days_count++;
      }
    }
  }

  return alarm;
}
void add_alarm(String messageString) {
  Alarm *alarm = parseAlarmFromJson(messageString.c_str());
  if (alarm == nullptr) {
    Serial.println("Failed to parse alarm from JSON");
    return;
  }
  alarmHeap.insert(alarm);
}

void update_alarm(String messageString) {
  Alarm *updated_alarm = parseAlarmFromJson(messageString.c_str());
  if (updated_alarm == nullptr) {
    Serial.println("Failed to parse alarm from JSON");
    return;
  }
  alarmHeap.update_alarm(updated_alarm->id, updated_alarm);
}
void delete_alarm(String messageString) {
  JsonDocument doc;

  DeserializationError error = deserializeJson(doc, messageString.c_str());
  if (error) {
    Serial.println("Failed to parse delete alarm from JSON in deserializeJson");
    return;
  }

  if (!doc["id"].is<uint32_t>()) {
    Serial.println("Failed to parse alarm from JSON. no id");
    return;
  }

  if (!doc["deleted_at"].is<const char *>()) {
    Serial.println("Failed to parse alarm from JSON. no time");
    return;
  }

  uint32_t id = doc["id"].as<uint32_t>();
  // Parse ISO8601 time (most common format from Go)
  const char *timeStr = doc["time"];
  // We expect something like: "2025-04-10T14:30:00"  or with Z / offset
  DateTime timestamp = DateTime(timeStr); // RTClib's DateTime can parse ISO8601

  if (!timestamp.isValid()) {
    Serial.println("Failed to parse alarm from JSON. invalid time");
    return;
  }

  alarmHeap.remove_alarm(id);
}

void add_alarms_batch(String messageString) {
  JsonDocument doc;

  DeserializationError error = deserializeJson(doc, messageString.c_str());
  if (error) {
    Serial.println("Failed to parse delete alarm from JSON in deserializeJson");
    return;
  }
  JsonArray alarms = doc["alarms"].as<JsonArray>();
  for (JsonVariantConst v : alarms) {
    add_alarm(v.as<String>());
  }
}
void set_color(String messageString) {
  if (messageString.length() == 7 && messageString[0] == '#') {
    unsigned long hexValue =
        strtoul(messageString.substring(1).c_str(), NULL, 16);

    uint8_t r = (hexValue >> 16) & 0xFF;
    uint8_t g = (hexValue >> 8) & 0xFF;
    uint8_t b = hexValue & 0xFF;

    gClockColor = CRGB(r, g, b);

    Serial.print("Color set to: ");
    Serial.print(messageString);
    Serial.print(" (R:");
    Serial.print(r);
    Serial.print(" G:");
    Serial.print(g);
    Serial.print(" B:");
    Serial.print(b);
    Serial.println(")");
  } else {
    Serial.println("Invalid color format. Expected #RRGGBB");
  }
}

void set_volume(String messageString) {
  volume = messageString.toInt();
  Serial.println("setting volume to " + messageString);
}

void ring(String messageString) { AlarmStart(); }

void sync_time(String messageString) {
  unsigned long unixTime = messageString.toInt();

  DateTime time(unixTime);

  rtc.adjust(time);

  Serial.print("Time set to: ");
  Serial.print(time.year());
  Serial.print("-");
  Serial.print(time.month());
  Serial.print("-");
  Serial.print(time.day());
  Serial.print(" ");
  Serial.print(time.hour());
  Serial.print(":");
  Serial.print(time.minute());
  Serial.print(":");
  Serial.println(time.second());
}

void setupDfPlayer() {
  FPSerial.begin(9600, SERIAL_8N1, 16, 17);
  if (player.begin(FPSerial)) {
    Serial.println("DFPlayer Mini online!");
  } else {
    Serial.println("Unable to begin DFPlayer:");
  }
}

void setupLED() {
  FastLED.addLeds<WS2812B, LED_PIN, GRB>(leds, NUM_LEDS);
  FastLED.setBrightness(BRIGHTNESS);
  FastLED.clear();
  FastLED.show();
}

void setupTouch() {
  pinMode(snoozePin, INPUT_PULLUP);
  attachInterrupt(snoozePin, Snooze, RISING);
  pinMode(stopPin, INPUT_PULLUP);
  attachInterrupt(stopPin, Stop, RISING);
  Serial.println("Touch Connected to Interrupt");
}

void ARDUINO_ISR_ATTR Snooze() {
  player.stop();
  Serial.println("Snooze");
}

void ARDUINO_ISR_ATTR Stop() {
  player.stop();
  Serial.println("Stop");
}

void AlarmStart() {
  if (player.begin(FPSerial)) {
    Serial.println("DFPlayer Mini online!");
    player.volume(volume);
    player.play(1);
  }
}

void wifiProcessor(void *args) {
  Serial.println("running wifi processor thread");
  for (;;) {
    wifiManager.process();
    delay(500);
  }
}

void mqttTask(void *pv) {
  MqttMsg msg;
  for (;;) {
    if (wifiConnected) {
      if (client.connected()) { // already connected
        if (!first_time_online) {
          first_time_online = true;
          client.publish(getTopic(ACTION_GET_COLOR).c_str(), DEVICEID.c_str());
          client.publish(getTopic(ACTION_GET_VOLUME).c_str(), DEVICEID.c_str());
          client.publish(getTopic(ACTION_GET_ALL_ALARMS).c_str(),
                         DEVICEID.c_str());
        }
        while (xQueueReceive(mqttQueue, &msg, 0)) {
          client.publish(msg.topic.c_str(), msg.payload.c_str());
        }
        client.loop();
      } else { // connect and sub
        if (client.connect(DEVICEID.c_str())) {
          for (uint8_t i = 0; i < sub_topics_count; i++) {
            client.subscribe(sub_topics[i].c_str(), 1);
          }
          client.loop();
        }
      }
      vTaskDelay(pdMS_TO_TICKS(10));
    }
  }
}

void setupWifi() {
  WiFi.onEvent(WiFiEvent);

  wifiManager.setDebugOutput(true);
  wifiManager.setConfigPortalTimeout(0); // no timeout
  wifiManager.setConfigPortalBlocking(false);

  bool connected = wifiManager.autoConnect("DyslexicClock-Setup", "12345678");
  if (!connected) {
    Serial.println("Wifi, failed!");
  }
  WiFi.setAutoReconnect(true);

  xTaskCreate(wifiProcessor, "wifiProcessor", 4096, NULL, 2, NULL);
}

void ConnectMqtt() {}

void WiFiEvent(WiFiEvent_t event, WiFiEventInfo_t info) {
  Serial.print("WiFi Event:");
  switch (event) {
  case ARDUINO_EVENT_WIFI_STA_START:
    Serial.println("WiFi, started");
    break;
  case ARDUINO_EVENT_WIFI_STA_CONNECTED:
    Serial.println("WiFi, connected");
    break;
  case ARDUINO_EVENT_WIFI_STA_DISCONNECTED:
    Serial.println("WiFi, disconnected!");
    wifiConnected = false;
    client.disconnect();
    Serial.println("disconnceted from mqtt");
    break;
  case ARDUINO_EVENT_WIFI_STA_GOT_IP:
    Serial.print("WiFi, got IP");
    Serial.println(WiFi.localIP());
    wifiConnected = true;

    break;
  default:
    break;
  }
}

void loop() {
  DateTime now = rtc.now();

  showTime(now.minute(), now.hour());
  if (!alarmHeap.empty()) {
    Alarm *next = alarmHeap.get_top();
    DateTime next_time = alarmHeap.get_next_occurrence(next);
    if (now >= next_time) {
      Serial.print("ALARM TRIGGERED! ID: ");
      Serial.println(next->id);

      alarmHeap.pop_top();
      // process alarm
      if (next->is_repeat) {
        alarmHeap.insert(next);
      } else {
      }

      AlarmStart();
    }
  }
  vTaskDelay(pdMS_TO_TICKS(60000));
}
