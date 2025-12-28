#include <FastLED.h>
#include "RTClib.h"
#include <DFRobotDFPlayerMini.h>
#include <WiFi.h>
#include <WiFiManager.h>
#include <vector>
#include <Arduino.h>
#include <PubSubClient.h>

#define FPSerial Serial1
#define LED_PIN 27
#define NUM_LEDS 121
#define BRIGHTNESS 185

const uint8_t DEVICEID = 1;

extern String sub_topics[9];
extern String pub_topics[4];

enum ActionPublish
{
    ACTION_STATUS,
    ACTION_GET_ALL_ALARMS,
    ACTION_RINGING,
    ACTION_LOG
};

// Simple Weekday enum for Arduino
enum Weekday
{
    SUNDAY = 0,
    MONDAY = 1,
    TUESDAY = 2,
    WEDNESDAY = 3,
    THURSDAY = 4,
    FRIDAY = 5,
    SATURDAY = 6
};

struct Alarm
{
    uint32_t id;
    uint32_t device_id;
    DateTime timestamp;
    bool is_repeat;
    Weekday repeating_days[7];    // Array of repeating days
    uint8_t repeating_days_count; // How many days are actually set
};

class AlarmHeap
{
public:
    AlarmHeap();
    AlarmHeap(std::vector<Alarm> list);
    Alarm get_top();
    void pop_top();
    void insert(Alarm alarm);
    bool update_alarm(uint32_t alarm_id, Alarm new_alarm);
    bool remove_alarm(uint32_t alarm_id);
    int size();
    bool empty();
    DateTime get_next_occurrence(const Alarm &alarm);

private:
    std::vector<Alarm> alarms;
    void check_exception();
    void swap(int i1, int i2);
    bool should_come_before(const Alarm &a, const Alarm &b);
    DateTime get_current_time();
    void sift_down(int index);
    void sift_up(int index);
};

extern CRGB leds[NUM_LEDS];
extern RTC_DS3231 rtc;
extern CRGB gClockColor;
extern uint8_t volume;
extern DFRobotDFPlayerMini player;
extern const char *mqtt_server;
extern WiFiClient espClient;
extern PubSubClient client;
extern long lastMsg;
extern char msg[50];
extern int value;

extern const uint8_t PROGMEM IT[];
extern const uint8_t PROGMEM IS[];
extern const uint8_t PROGMEM HAPPY[];
extern const uint8_t PROGMEM QUARTER[];
extern const uint8_t PROGMEM TWENTY[];
extern const uint8_t PROGMEM FIVE[];
extern const uint8_t PROGMEM TO[];
extern const uint8_t PROGMEM TEN[];
extern const uint8_t PROGMEM HALF[];
extern const uint8_t PROGMEM PAST[];
extern const uint8_t PROGMEM h_SEVEN[];
extern const uint8_t PROGMEM h_THREE[];
extern const uint8_t PROGMEM h_TWO[];
extern const uint8_t PROGMEM h_ONE[];
extern const uint8_t PROGMEM h_FOUR[];
extern const uint8_t PROGMEM h_FIVE[];
extern const uint8_t PROGMEM h_SIX[];
extern const uint8_t PROGMEM h_TWELVE[];
extern const uint8_t PROGMEM h_NINE[];
extern const uint8_t PROGMEM h_EIGHT[];
extern const uint8_t PROGMEM h_ELEVEN[];
extern const uint8_t PROGMEM OCLOCK[];
extern const uint8_t PROGMEM h_TEN[];
extern const uint8_t PROGMEM BIRTHDAY[];

extern const uint8_t snoozePin;
extern volatile int32_t stopPin;

extern AlarmHeap alarmHeap;

String getTopic(ActionPublish action);
void setup();
void setupRTC();
void setupMqtt();
void mbCallback(char *topic, byte *message, unsigned int length);
void setupDfPlayer();
void setupLED();
void setupTouch();
void ARDUINO_ISR_ATTR Snooze();
void ARDUINO_ISR_ATTR Stop();
void AlarmStart();
void setupWifi();
void WiFiEvent(WiFiEvent_t event, WiFiEventInfo_t info);
void loop();
void showTime(uint8_t m, uint8_t h);
void showHours(uint8_t h);
void lightWord(const uint8_t word[2]);
Alarm create_alarm(uint32_t id, uint32_t device_id, uint32_t timestamp, bool is_repeat, Weekday *repeating_days, uint8_t repeating_days_count);
