#include "main.h"
#include "Arduino.h"
#include "ArduinoJson.h"
#include "ArduinoJson/Deserialization/DeserializationError.hpp"
#include "ArduinoJson/Document/JsonDocument.hpp"
#include "HardwareSerial.h"
#include "RTClib.h"
#include "WiFi.h"
#include "WiFiGeneric.h"
#include "WiFiManager.h"
#include "WiFiType.h"
#include "fl/str.h"
#include "freertos/projdefs.h"
#include "led.h"
#include <cstdint>
#include <ctime>

#define DEBUGMODE true

Alarm *prev = NULL;
WiFiManager wifiManager;
volatile bool start_portal = false;

struct MqttMsg {
  String topic;
  String payload;
};

QueueHandle_t mqttQueue;

bool first_time_online = false;
bool wifiConnected = false;
const uint8_t sub_topics_count = 10;
String sub_topics[sub_topics_count] = {"devices/" + DEVICEID + "/alarms/create",
                                       "devices/" + DEVICEID + "/alarms/update",
                                       "devices/" + DEVICEID + "/alarms/delete",
                                       "devices/" + DEVICEID + "/alarms",
                                       "devices/" + DEVICEID + "/color",
                                       "devices/" + DEVICEID + "/volume",
                                       "devices/" + DEVICEID + "/ring",
                                       "devices/" + DEVICEID + "/silent",
                                       "devices/" + DEVICEID + "/snooze",
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
  Serial.println("mamad");
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
  Serial.println("[mbCallback] Entry");
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
  Serial.println("[mbCallback] Topic: " + stringTopic);
  Serial.println("[mbCallback] Message: " + messageString);

  if (stringTopic == sub_topics[0]) {
    Serial.println("[mbCallback] Action: add_alarm");
    add_alarm(messageString);
  } else if (stringTopic == sub_topics[1]) {
    Serial.println("[mbCallback] Action: update_alarm");
    update_alarm(messageString);
  } else if (stringTopic == sub_topics[2]) {
    Serial.println("[mbCallback] Action: delete_alarm");
    delete_alarm(messageString);
  } else if (stringTopic == sub_topics[3]) {
    Serial.println("[mbCallback] Action: add_alarms_batch");
    add_alarms_batch(messageString);
  } else if (stringTopic == sub_topics[4]) {
    Serial.println("[mbCallback] Action: set_color");
    set_color(messageString);
  } else if (stringTopic == sub_topics[5]) {
    Serial.println("[mbCallback] Action: set_volume");
    set_volume(messageString);
  } else if (stringTopic == sub_topics[6]) {
    Serial.println("[mbCallback] Action: ring");
    ring(messageString);
  } else if (stringTopic == sub_topics[7]) {
    Serial.println("[mbCallback] Action: silent");
    player.stop();
    Serial.println("Silenced");
  } else if (stringTopic == sub_topics[8]) {
    Serial.println("[mbCallback] Action: snooze");
    Snooze();
  } else if (stringTopic == sub_topics[9]) {
    Serial.println("[mbCallback] Action: sync_time");
    sync_time(messageString);
  } else {
    Serial.println("[mbCallback] No matching topic found.");
  }
  Serial.println("[mbCallback] Exit");
}
Alarm *parseAlarmFromJson(const char *jsonString) {
  Serial.println("[parseAlarmFromJson] Entry");
  JsonDocument doc;

  DeserializationError error = deserializeJson(doc, jsonString);
  if (error) {
    Serial.print("[parseAlarmFromJson] deserializeJson() failed: ");
    Serial.println(error.c_str());
    return nullptr;
  }

  if (!doc["id"].is<uint32_t>()) {
    Serial.println("[parseAlarmFromJson] 'id' is not a uint32_t.");
    return nullptr;
  }

  if (!doc["time"].is<const char *>()) {
    Serial.println("[parseAlarmFromJson] 'time' is not a const char*.");
    return nullptr;
  }

  Alarm *alarm = new Alarm();
  if (!alarm) {
    Serial.println("[parseAlarmFromJson] Failed to allocate memory for Alarm.");
    return nullptr;
  }

  alarm->id = doc["id"].as<uint32_t>();
  Serial.print("[parseAlarmFromJson] Alarm ID: ");
  Serial.println(alarm->id);

  // Parse ISO8601 time (most common format from Go)
  const char *timeStr = doc["time"];
  Serial.print("[parseAlarmFromJson] Time string: ");
  Serial.println(timeStr);
  // We expect something like: "2025-04-10T14:30:00"  or with Z / offset
  alarm->timestamp = DateTime(timeStr); // RTClib's DateTime can parse ISO8601

  if (!alarm->timestamp.isValid()) {
    Serial.println("[parseAlarmFromJson] Failed to parse timestamp.");
    delete alarm;
    return nullptr;
  }
  Serial.println("[parseAlarmFromJson] Timestamp parsed successfully.");

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
  Serial.print("[parseAlarmFromJson] repeating_days_count: ");
  Serial.println(alarm->repeating_days_count);

  Serial.println("[parseAlarmFromJson] Exit (Success)");
  return alarm;
}
void add_alarm(String messageString) {
  Serial.println("[add_alarm] Entry");
  Alarm *alarm = parseAlarmFromJson(messageString.c_str());
  if (alarm == nullptr) {
    Serial.println("[add_alarm] Failed to parse alarm from JSON");
    return;
  }
  alarmHeap.insert(alarm);
  Serial.println("[add_alarm] Exit");
}

void update_alarm(String messageString) {
  Serial.println("[update_alarm] Entry");
  Alarm *updated_alarm = parseAlarmFromJson(messageString.c_str());
  if (updated_alarm == nullptr) {
    Serial.println("[update_alarm] Failed to parse alarm from JSON");
    return;
  }
  alarmHeap.update_alarm(updated_alarm->id, updated_alarm);
  Serial.println("[update_alarm] Exit");
}
void delete_alarm(String messageString) {
  Serial.println("[delete_alarm] Entry");
  JsonDocument doc;

  DeserializationError error = deserializeJson(doc, messageString.c_str());
  if (error) {
    Serial.println("[delete_alarm] Failed to parse delete alarm from JSON in "
                   "deserializeJson");
    return;
  }

  if (!doc["id"].is<uint32_t>()) {
    Serial.println("[delete_alarm] Failed to parse alarm from JSON. no id");
    return;
  }

  if (!doc["deleted_at"].is<const char *>()) {
    Serial.println(
        "[delete_alarm] Failed to parse alarm from JSON. no deleted_at");
    return;
  }

  uint32_t id = doc["id"].as<uint32_t>();
  Serial.print("[delete_alarm] Deleting alarm with id: ");
  Serial.println(id);
  // Parse ISO8601 time (most common format from Go)
  const char *timeStr = doc["deleted_at"];
  Serial.print("[delete_alarm] deleted_at string: ");
  Serial.println(timeStr);
  // We expect something like: "2025-04-10T14:30:00"  or with Z / offset
  DateTime timestamp = DateTime(timeStr); // RTClib's DateTime can parse ISO8601

  if (!timestamp.isValid()) {
    Serial.println("[delete_alarm] Failed to parse alarm from JSON. invalid "
                   "deleted_at time");
    return;
  }

  alarmHeap.remove_alarm(id);
  Serial.println("[delete_alarm] Exit");
}

void add_alarms_batch(String messageString) {
  Serial.println("[add_alarms_batch] Entry");
  JsonDocument doc;

  DeserializationError error = deserializeJson(doc, messageString.c_str());
  if (error) {
    Serial.println("[add_alarms_batch] Failed to parse batch from JSON in "
                   "deserializeJson");
    return;
  }
  JsonArray alarms = doc.as<JsonArray>();
  Serial.print("[add_alarms_batch] Number of alarms to add: ");
  Serial.println(alarms.size());
  for (JsonVariantConst v : alarms) {
    add_alarm(v.as<String>());
  }
  Serial.println("[add_alarms_batch] Exit");
}
void set_color(String messageString) {
  Serial.println("[set_color] Entry");
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
    Serial.println("[set_color] Invalid color format. Expected #RRGGBB");
  }

  DateTime now = rtc.now();
  Serial.println("[set_color] Showing time");
  showTime(now.minute(), now.hour());

  Serial.println("[set_color] Exit");
}

void set_volume(String messageString) {
  Serial.println("[set_volume] Entry");
  volume = messageString.toInt();
  Serial.println("setting volume to " + messageString);
  Serial.println("[set_volume] Exit");
}

void ring(String messageString) {
  Serial.println("[ring] Entry");
  AlarmStart();
  Serial.println("[ring] Exit");
}

void sync_time(String messageString) {
  Serial.println("[sync_time] Entry");
  unsigned long unixTime = messageString.toInt();

  DateTime time(unixTime);
    time = time + TimeSpan(0, 3, 30, 0);

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
  Serial.println("[sync_time] Exit");
}

void setupDfPlayer() {
  FPSerial.begin(9600, SERIAL_8N1, 16, 17);
  // while (!player.begin(FPSerial)) {
  if (player.begin(FPSerial)) {
        Serial.println("DFPlayer Mini online!");
    delay(100);
  }
  Serial.println("DFPlayer Mini not online!!!!!!!!");
  // Serial.println("DFPlayer Mini online!");
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

  if (prev != NULL) {
    Alarm *alarm = copy_alarm_for_snooze(prev);
    alarmHeap.insert(alarm);
    if (DEBUGMODE)
      Serial.println("snoozed alarm added");
  }
}

void ARDUINO_ISR_ATTR Stop() {
  player.stop();
  Serial.println("Stop");
}

void AlarmStart() {
  if (DEBUGMODE) {
    Serial.println("Alarm Started!");
  }

  player.volume(volume);
  player.play(2);
}

void wifiProcessor(void *args) {
  Serial.println("running wifi processor thread");
  for (;;) {
    if (start_portal) {
      Serial.println(
          "[wifiProcessor] Starting config portal because of disconnect.");
      wifiManager.startConfigPortal("DyslexicClock-Setup", "12345678");
      start_portal = false;
    }
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
        connect_mqtt();
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

void connect_mqtt() {
  if (client.connect(DEVICEID.c_str())) {
    for (uint8_t i = 0; i < sub_topics_count; i++) {
      client.subscribe(sub_topics[i].c_str(), 1);
    }
    client.loop();
  }
}

void WiFiEvent(WiFiEvent_t event, WiFiEventInfo_t info) {
  Serial.print("WiFi Event:");
  switch (event) {
  case ARDUINO_EVENT_WIFI_STA_START:
    Serial.println("WiFi, started");
    break;
  case ARDUINO_EVENT_WIFI_STA_CONNECTED:
    Serial.println("WiFi, connected");
    start_portal = false;
    break;
  case ARDUINO_EVENT_WIFI_STA_DISCONNECTED:
    Serial.println("WiFi, disconnected!");
    wifiConnected = false;
    client.disconnect();
    Serial.println("disconnceted from mqtt");
    start_portal = true;
    break;
  case ARDUINO_EVENT_WIFI_STA_GOT_IP:
    Serial.print("WiFi, got IP");
    Serial.println(WiFi.localIP());
    wifiConnected = true;
    connect_mqtt();
    break;
  default:
    break;
  }
}

void print_alarm(Alarm *alarm) {
  Serial.printf("ID: %d\nis_repeatable: %d\nHour: %d\nMinute: %d\n", alarm->id,
                alarm->is_repeat, alarm->timestamp.hour(),
                alarm->timestamp.minute());
  Serial.print("repeating days: ");
  for (uint8_t i; i < alarm->repeating_days_count; i++) {
    Serial.printf("%d, ", alarm->repeating_days[i]);
  }
  Serial.println("___________________________");

  return;
}

Alarm *copy_alarm_for_snooze(Alarm *snoozed_alarm) {
  Serial.println("[copy_alarm_for_snooze] Entry");
  Alarm *snoozed_alarm_copy = new Alarm();
  snoozed_alarm_copy->id = snoozed_alarm->id;
  snoozed_alarm_copy->timestamp =
      snoozed_alarm->timestamp + TimeSpan(0, 0, 1, 0);
  snoozed_alarm_copy->is_repeat = false;
  snoozed_alarm_copy->repeating_days_count =
      snoozed_alarm->repeating_days_count;
  for (uint8_t i = 0; i < snoozed_alarm->repeating_days_count; i++) {
    snoozed_alarm_copy->repeating_days[i] = snoozed_alarm->repeating_days[i];
  }
  Serial.println("[copy_alarm_for_snooze] Exit");
  return snoozed_alarm_copy;
}

void loop() {
  Serial.println("[loop] Entry");
  DateTime now = rtc.now();

  Serial.println("[loop] Showing time");
  showTime(now.minute(), now.hour());

  if (DEBUGMODE) {
    Serial.println("[loop] Printing alarms (DEBUGMODE)");
    for (uint8_t i = 0; i < alarmHeap.size(); i++) {
      print_alarm(alarmHeap.alarms[i]);
    }
  }
  if (!alarmHeap.empty()) {
    Serial.println("[loop] Alarm heap is not empty");
    if (prev != NULL) {
      Serial.println("[loop] Re-inserting repeating alarm");
      alarmHeap.insert(prev);
      prev = NULL;
    }
    Alarm *next = alarmHeap.get_top();
    Serial.println("[loop] Got top alarm from heap");
    DateTime next_time = alarmHeap.get_next_occurrence(next);
    Serial.println("[loop] Got next occurrence");
    if (now.minute() == next_time.minute() && now.hour() == next_time.hour()) {
      Serial.print("ALARM TRIGGERED! ID: ");
      Serial.println(next->id);

      alarmHeap.pop_top();
      Serial.println("[loop] Popped alarm from heap");
      // process alarm
      if (next->is_repeat) {
        Serial.println("[loop] Alarm is repeating");
        prev = next;
      } else {
        Serial.println("[loop] Alarm is not repeating, deleting it");
        delete next;
      }
      AlarmStart();
    }
  } else {
    Serial.println("[loop] Alarm heap is empty");
  }
  Serial.println("[loop] Starting delay");
  vTaskDelay(pdMS_TO_TICKS(59000));
}
