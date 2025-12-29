#include "main.h"

String sub_topics[9] = {
    "devices/" + String(DEVICEID) + "/alarms/create",
    "devices/" + String(DEVICEID) + "/alarms/update",
    "devices/" + String(DEVICEID) + "/alarms/delete",
    "devices/" + String(DEVICEID) + "/alarms",
    "devices/" + String(DEVICEID) + "/color",
    "devices/" + String(DEVICEID) + "/volume",
    "devices/" + String(DEVICEID) + "/ring",
    "devices/" + String(DEVICEID) + "/silence",
    "devices/time"};

String pub_topics[4] = {
    "devices/" + String(DEVICEID) + "/status",
    "devices/alarms",
    "devices/" + String(DEVICEID) + "/ringing",
    "devices/" + String(DEVICEID) + "/log"};

String getTopic(ActionPublish action)
{
  switch (action)
  {
  case ACTION_STATUS:
    return pub_topics[0];
  case ACTION_GET_ALL_ALARMS:
    return pub_topics[1];
  case ACTION_RINGING:
    return pub_topics[2];
  case ACTION_LOG:
    return pub_topics[3];
  default:
    return "";
  }
}

// char topicBuffer[4][64];
// void setup_topics() {
//   snprintf(topicBuffer[0], 64, "devices/%d/status", DEVICEID);
//   snprintf(topicBuffer[1], 64, "devices/alarms", DEVICEID);
//   snprintf(topicBuffer[2], 64, "devices/%d/ringing", DEVICEID);
//   snprintf(topicBuffer[3], 64, "devices/%d/log", DEVICEID);
// }

const uint8_t snoozePin = 26;
volatile int32_t stopPin = 25;

CRGB leds[NUM_LEDS];
RTC_DS3231 rtc;
CRGB gClockColor = CRGB::Green;
uint8_t volume = 26;
DFRobotDFPlayerMini player;
const char *mqtt_server = "api.dyslexics-clock.niflheimdevs.ir";
WiFiClient espClient;
PubSubClient client(espClient);
long lastMsg = 0;
char msg[50];
int value = 0;

const uint8_t PROGMEM IT[] = {0, 1};
const uint8_t PROGMEM IS[] = {3, 4};
const uint8_t PROGMEM HAPPY[] = {6, 10};
const uint8_t PROGMEM QUARTER[] = {14, 20};
const uint8_t PROGMEM TWENTY[] = {22, 27};
const uint8_t PROGMEM FIVE[] = {29, 32};
const uint8_t PROGMEM TO[] = {33, 34};
const uint8_t PROGMEM TEN[] = {36, 38};
const uint8_t PROGMEM HALF[] = {40, 43};
const uint8_t PROGMEM PAST[] = {44, 47};
const uint8_t PROGMEM h_SEVEN[] = {49, 53};
const uint8_t PROGMEM h_THREE[] = {55, 59};
const uint8_t PROGMEM h_TWO[] = {60, 62};
const uint8_t PROGMEM h_ONE[] = {63, 65};
const uint8_t PROGMEM h_FOUR[] = {66, 69};
const uint8_t PROGMEM h_FIVE[] = {70, 73};
const uint8_t PROGMEM h_SIX[] = {74, 76};
const uint8_t PROGMEM h_TWELVE[] = {77, 82};
const uint8_t PROGMEM h_NINE[] = {84, 87};
const uint8_t PROGMEM h_EIGHT[] = {88, 92};
const uint8_t PROGMEM h_ELEVEN[] = {93, 98};
const uint8_t PROGMEM OCLOCK[] = {99, 104};
const uint8_t PROGMEM h_TEN[] = {107, 109};
const uint8_t PROGMEM BIRTHDAY[] = {112, 120};

AlarmHeap::AlarmHeap() {}

AlarmHeap::AlarmHeap(std::vector<Alarm> list)
{
  for (int i = 0; i < list.size(); i++)
  {
    this->alarms.push_back(list[i]);
  }
  for (int i = alarms.size() / 2; i >= 0; i--)
  {
    this->sift_down(i);
  }
}

Alarm AlarmHeap::get_top()
{
  this->check_exception();
  return alarms[0];
}

void AlarmHeap::pop_top()
{
  this->check_exception();
  alarms[0] = alarms[alarms.size() - 1];
  alarms.pop_back();
  this->sift_down(0);
}

void AlarmHeap::insert(Alarm alarm)
{
  alarms.push_back(alarm);
  this->sift_up(alarms.size() - 1);
}

bool AlarmHeap::update_alarm(uint32_t alarm_id, Alarm new_alarm)
{
  for (int i = 0; i < alarms.size(); i++)
  {
    if (alarms[i].id == alarm_id)
    {
      Alarm old_alarm = alarms[i];
      alarms[i] = new_alarm;

      if (should_come_before(new_alarm, old_alarm))
      {
        sift_up(i);
      }
      else
      {
        sift_down(i);
      }
      return true;
    }
  }
  return false;
}

bool AlarmHeap::remove_alarm(uint32_t alarm_id)
{
  for (int i = 0; i < alarms.size(); i++)
  {
    if (alarms[i].id == alarm_id)
    {
      alarms[i] = alarms[alarms.size() - 1];
      alarms.pop_back();
      if (i < alarms.size())
      {
        sift_down(i);
        sift_up(i);
      }
      return true;
    }
  }
  return false;
}

int AlarmHeap::size()
{
  return alarms.size();
}

bool AlarmHeap::empty()
{
  return alarms.empty();
}

DateTime AlarmHeap::get_next_occurrence(const Alarm &alarm)
{
  DateTime now = get_current_time();

  if (!alarm.is_repeat)
  {
    return alarm.timestamp;
  }

  uint8_t alarm_hour = alarm.timestamp.hour();
  uint8_t alarm_minute = alarm.timestamp.minute();
  uint8_t alarm_second = alarm.timestamp.second();

  bool has_today = false;
  for (uint8_t i = 0; i < alarm.repeating_days_count; i++)
  {
    if (alarm.repeating_days[i] == now.dayOfTheWeek())
    {
      has_today = true;
      break;
    }
  }

  if (has_today)
  {
    DateTime today_alarm = DateTime(now.year(), now.month(), now.day(), alarm_hour, alarm_minute, alarm_second);
    if (today_alarm >= now)
    {
      return today_alarm;
    }
  }

  for (int days_ahead = 1; days_ahead <= 7; days_ahead++)
  {
    uint8_t next_wday = (now.dayOfTheWeek() + days_ahead) % 7;
    for (uint8_t i = 0; i < alarm.repeating_days_count; i++)
    {
      if (alarm.repeating_days[i] == next_wday)
      {
        DateTime next = now + TimeSpan(days_ahead * 86400L);
        return DateTime(next.year(), next.month(), next.day(), alarm_hour, alarm_minute, alarm_second);
      }
    }
  }
  return now + TimeSpan(7 * 86400L);
}

void AlarmHeap::check_exception()
{
  if (alarms.size() == 0)
  {
    Serial.println("ERROR: Alarm heap is empty!");
    while (1)
      ;
  }
}

void AlarmHeap::swap(int i1, int i2)
{
  Alarm temp = alarms[i1];
  alarms[i1] = alarms[i2];
  alarms[i2] = temp;
}

bool AlarmHeap::should_come_before(const Alarm &a, const Alarm &b)
{
  DateTime next_a = get_next_occurrence(a);
  DateTime next_b = get_next_occurrence(b);
  return next_a < next_b;
}

DateTime AlarmHeap::get_current_time()
{
  return rtc.now();
}

void AlarmHeap::sift_down(int index)
{
  int left = index * 2 + 1;
  int right = index * 2 + 2;
  if (left >= alarms.size())
    return;

  int target_index;
  if (right >= alarms.size())
  {
    target_index = left;
  }
  else
  {
    target_index = should_come_before(alarms[left], alarms[right]) ? left : right;
  }

  if (should_come_before(alarms[target_index], alarms[index]))
  {
    this->swap(target_index, index);
    this->sift_down(target_index);
  }
}

void AlarmHeap::sift_up(int index)
{
  if (index == 0)
    return;
  int parent = ((index % 2 == 0) ? ((index / 2) - 1) : (index / 2));
  if (should_come_before(alarms[index], alarms[parent]))
  {
    this->swap(parent, index);

    this->sift_up(parent);
  }
}

Alarm create_alarm(uint32_t id, uint32_t device_id, uint32_t timestamp, bool is_repeat = false, Weekday *repeating_days = nullptr, uint8_t repeating_days_count = 0)
{
  Alarm alarm;
  alarm.id = id;
  alarm.device_id = device_id;
  alarm.timestamp = timestamp;
  alarm.is_repeat = is_repeat;
  alarm.repeating_days_count = repeating_days_count;

  if (is_repeat && repeating_days != nullptr)
  {
    for (uint8_t i = 0; i < repeating_days_count && i < 7; i++)
    {
      alarm.repeating_days[i] = repeating_days[i];
    }
  }

  return alarm;
}

void lightWord(const uint8_t word[2])
{
  for (uint8_t i = word[0]; i <= word[1]; i++)
  {
    leds[i] = gClockColor;
  }
}

void showTime(uint8_t m, uint8_t h)
{
  FastLED.clear();

  lightWord(IT);
  lightWord(IS);

  uint8_t block = (m + 2) / 5;

  switch (block)
  {
  case 0: // 00–02 → o'clock
    showHours(h);
    lightWord(OCLOCK);
    break;

  case 1: // 03–07 → FIVE PAST
    lightWord(FIVE);
    lightWord(PAST);
    showHours(h);
    break;

  case 2: // 08–12 → TEN PAST
    lightWord(TEN);
    lightWord(PAST);
    showHours(h);
    break;

  case 3: // 13–17 → QUARTER PAST
    lightWord(QUARTER);
    lightWord(PAST);
    showHours(h);
    break;

  case 4: // 18–22 → TWENTY PAST
    lightWord(TWENTY);
    lightWord(PAST);
    showHours(h);
    break;

  case 5: // 23–27 → TWENTY FIVE PAST
    lightWord(TWENTY);
    lightWord(FIVE);
    lightWord(PAST);
    break;

  case 6: // 28–32 → HALF PAST
    lightWord(HALF);
    lightWord(PAST);
    showHours(h);
    break;

  case 7: // 33–37 → TWENTY FIVE TO
    lightWord(TWENTY);
    lightWord(FIVE);
    lightWord(TO);
    showHours(h + 1);
    break;

  case 8: // 38–42 → TWENTY TO
    lightWord(TWENTY);
    lightWord(TO);
    showHours(h + 1);
    break;

  case 9: // 43–47 → QUARTER TO
    lightWord(QUARTER);
    lightWord(TO);
    showHours(h + 1);
    break;

  case 10: // 48–52 → TEN TO
    lightWord(TEN);
    lightWord(TO);
    showHours(h + 1);
    break;

  case 11: // 53–57 → FIVE TO
    lightWord(FIVE);
    lightWord(TO);
    showHours(h + 1);
    break;

  case 12: // 58–59 → o'clock
    showHours(h + 1);
    lightWord(OCLOCK);
    break;
  }

  FastLED.show();
}

void showHours(uint8_t h)
{
  uint8_t h12 = h % 12;
  if (h12 == 0)
    h12 = 12;

  switch (h12)
  {
  case 1:
    lightWord(h_ONE);
    break;
  case 2:
    lightWord(h_TWO);
    break;
  case 3:
    lightWord(h_THREE);
    break;
  case 4:
    lightWord(h_FOUR);
    break;
  case 5:
    lightWord(h_FIVE);
    break;
  case 6:
    lightWord(h_SIX);
    break;
  case 7:
    lightWord(h_SEVEN);
    break;
  case 8:
    lightWord(h_EIGHT);
    break;
  case 9:
    lightWord(h_NINE);
    break;
  case 10:
    lightWord(h_TEN);
    break;
  case 11:
    lightWord(h_ELEVEN);
    break;
  case 12:
    lightWord(h_TWELVE);
    break;
  }
}

AlarmHeap alarmHeap;

void setup()
{
  Serial.begin(115200);
  // setupWifi();
  setupRTC();
  setupLED();
  setupDfPlayer();
  setupTouch();
}

void setupRTC()
{
  if (!rtc.begin())
  {
    Serial.println("Couldn't find RTC");
    Serial.flush();
    while (1)
      delay(10);
  }

  if (rtc.lostPower())
  {
    Serial.println("RTC lost power, let's set the time!");
    rtc.adjust(DateTime(F(__DATE__), F(__TIME__)));
  }
}

void setupMqtt()
{
  client.setServer(mqtt_server, 11883);
  client.setCallback(mbCallback);
}

void mbCallback(char *topic, byte *message, unsigned int length)
{
  Serial.print("Message arrived on topic: ");
  Serial.print(topic);
  Serial.print(". Message: ");
  String messageTemp;

  for (int i = 0; i < length; i++)
  {
    Serial.print((char)message[i]);
    messageTemp += (char)message[i];
  }
  Serial.println();

  // Feel free to add more if statements to control more GPIOs with MQTT

  // If a message is received on the topic esp32/output, you check if the message is either "on" or "off".
  // Changes the output state according to the message

  String stringTopic = String(topic);

  if (stringTopic == sub_topics[0])
  {
    // Handle sub_topics[0]
  }
  else if (stringTopic == sub_topics[1])
  {
    // Handle sub_topics[1]
  }
  else if (stringTopic == sub_topics[2])
  {
    // Handle sub_topics[2]
  }
  else if (stringTopic == sub_topics[3])
  {
    // Handle sub_topics[3]
  }
  else if (stringTopic == sub_topics[4])
  {
    // Handle sub_topics[4]
  }
  else if (stringTopic == sub_topics[5])
  {
    // Handle sub_topics[5]
  }
  else if (stringTopic == sub_topics[6])
  {
    // Handle sub_topics[6]
  }
  else if (stringTopic == sub_topics[7])
  {
    // Handle sub_topics[7]
  }
  else if (stringTopic == sub_topics[8])
  {
    // Handle sub_topics[8]
  }
  else
  {
    Serial.println("No matching topic found.");
  }
}

void setupDfPlayer()
{
  FPSerial.begin(9600, SERIAL_8N1, 16, 17);
  if (player.begin(FPSerial))
  {
    Serial.println("DFPlayer Mini online!");
  }
  else
  {
    Serial.println("Unable to begin DFPlayer:");
  }
}

void setupLED()
{
  FastLED.addLeds<WS2812B, LED_PIN, GRB>(leds, NUM_LEDS);
  FastLED.setBrightness(BRIGHTNESS);
  FastLED.clear();
  FastLED.show();
}

void setupTouch()
{
  pinMode(snoozePin, INPUT_PULLUP);
  attachInterrupt(snoozePin, Snooze, RISING);
  pinMode(stopPin, INPUT_PULLUP);
  attachInterrupt(stopPin, Stop, RISING);
  Serial.println("Touch Connected to Interrupt");
}

void ARDUINO_ISR_ATTR Snooze()
{
  player.stop();
  Serial.println("Snooze");
}

void ARDUINO_ISR_ATTR Stop()
{
  player.stop();
  Serial.println("Stop");
}

void AlarmStart()
{
  if (player.begin(FPSerial))
  {
    Serial.println("DFPlayer Mini online!");
    player.volume(volume);
    player.play(1);
  }
}

void setupWifi()
{
  WiFi.onEvent(WiFiEvent);

  WiFiManager wifiManager;
  wifiManager.setDebugOutput(true);
  wifiManager.setConfigPortalTimeout(180);

  bool connected = wifiManager.autoConnect("DyslexicClock-Setup", "12345678");
  if (!connected)
  {
    Serial.println("Wifi, failed!");
  }
  WiFi.setAutoReconnect(true);
}

void WiFiEvent(WiFiEvent_t event, WiFiEventInfo_t info)
{
  Serial.print("WiFi Event:");
  switch (event)
  {
  case ARDUINO_EVENT_WIFI_STA_START:
    Serial.println("WiFi, started");
    break;
  case ARDUINO_EVENT_WIFI_STA_CONNECTED:
    Serial.println("WiFi, connected");
    break;
  case ARDUINO_EVENT_WIFI_STA_DISCONNECTED:
    Serial.println("WiFi, disconnected!");
    break;
  case ARDUINO_EVENT_WIFI_STA_GOT_IP:
    Serial.print("WiFi, got IP");
    Serial.println(WiFi.localIP());
    break;
  default:
    break;
  }
}

void loop()
{
  DateTime now = rtc.now();
  showTime(now.minute(), now.hour());
  if (!alarmHeap.empty())
  {
    Alarm next = alarmHeap.get_top();
    DateTime next_time = alarmHeap.get_next_occurrence(next);
    if (now >= next_time)
    {
      Serial.print("ALARM TRIGGERED! ID: ");
      Serial.println(next.id);

      alarmHeap.pop_top();
      if (next.is_repeat)
      {
        alarmHeap.insert(next);
      }
      AlarmStart();
    }
  }
  delay(60000);
}
