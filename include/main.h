#include "RTClib.h"
#include <Arduino.h>
#include <DFRobotDFPlayerMini.h>
#include <PubSubClient.h>
#include <WiFi.h>
#include <WiFiManager.h>
#include <alarm_heap.h>

const String DEVICEID = "1";

extern String sub_topics[10];
extern String pub_topics[6];

enum ActionPublish {
  ACTION_STATUS,
  ACTION_GET_ALL_ALARMS,
  ACTION_GET_COLOR,
  ACTION_GET_VOLUME,
  ACTION_RINGING,
  ACTION_LOG
};

extern RTC_DS3231 rtc;
extern uint8_t volume;
extern DFRobotDFPlayerMini player;
extern const char *mqtt_server;
extern WiFiClient espClient;
extern PubSubClient client;
extern long lastMsg;
extern char msg[50];
extern int value;

extern const uint8_t snoozePin;
extern volatile int32_t stopPin;

extern AlarmHeap alarmHeap;

void connect_mqtt();
String getTopic(ActionPublish action);
void setup();
void setupRTC();
void setupMQTT();
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
Alarm *create_alarm(uint32_t id, uint32_t device_id, uint32_t timestamp,
                    bool is_repeat, Weekday *repeating_days,
                    uint8_t repeating_days_count);

Alarm *copy_alarm_for_snooze(Alarm *snoozed_alarm);
void add_alarm(String messageString);
void update_alarm(String messageString);
void delete_alarm(String messageString);
void add_alarms_batch(String messageString);
void set_color(String messageString);
void set_volume(String messageString);
void ring(String messageString);
void sync_time(String messageString);

void mqttTask(void *pv);
