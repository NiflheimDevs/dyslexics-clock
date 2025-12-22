#include <FastLED.h>
#include "RTClib.h"
#include <DFRobotDFPlayerMini.h>
#include <WiFi.h>       
#include <WiFiManager.h>

#define FPSerial Serial1
#define LED_PIN     16       
#define NUM_LEDS    121      
#define BRIGHTNESS  100

const uint8_t snoozePin = 26;
volatile int32_t stopPin = 25;
//HAHA

CRGB leds[NUM_LEDS];
RTC_DS3231 rtc;
CRGB gClockColor = 0xFFFFFF;
DFRobotDFPlayerMini player;

const int8_t PROGMEM IT[]       = { 0, 1 };
const int8_t PROGMEM IS[]       = { 3, 4 };
const int8_t PROGMEM HAPPY[]    = { 6, 10 };
const int8_t PROGMEM QUARTER[]  = { 14, 20 };
const int8_t PROGMEM TWENTY[]   = { 22, 27 };
const int8_t PROGMEM FIVE[]     = { 29, 32 };
const int8_t PROGMEM TO[]       = { 33, 34 };
const int8_t PROGMEM TEN[]      = { 36, 38 };
const int8_t PROGMEM HALF[]     = { 40, 43 };
const int8_t PROGMEM PAST[]     = { 44, 47 };
const int8_t PROGMEM h_SEVEN[]  = { 49, 53 };
const int8_t PROGMEM h_THREE[]  = { 55, 59 };
const int8_t PROGMEM h_TWO[]    = { 60, 62 };
const int8_t PROGMEM h_ONE[]    = { 63, 65 };
const int8_t PROGMEM h_FOUR[]   = { 66, 69 };
const int8_t PROGMEM h_FIVE[]   = { 70, 73 };
const int8_t PROGMEM h_SIX[]    = { 74, 76 };
const int8_t PROGMEM h_TWELVE[] = { 77, 82 };
const int8_t PROGMEM h_NINE[]   = { 84, 87 };
const int8_t PROGMEM h_EIGHT[]  = { 88, 92 };
const int8_t PROGMEM h_ELEVEN[] = { 93, 98 };
const int8_t PROGMEM OCLOCK[]   = { 99, 104 };
const int8_t PROGMEM h_TEN[]    = { 107, 109 };
const int8_t PROGMEM BIRTHDAY[] = { 112, 120 };

void lightWord(const uint8_t word[2]) {
  for (uint8_t i = word[0]; i <= word[1]; i++) {
    leds[i] = gClockColor;
  }
}

void showTime(uint8_t m, uint8_t h) {  
  FastLED.clear();  

  lightWord(IT);
  lightWord(IS);  

  uint8_t block = (m + 2) / 5;

  switch (block) {
    case 0:   // 00–02 → o'clock
      showHours(h);
      lightWord(OCLOCK);
      break;

    case 1:   // 03–07 → FIVE PAST
      lightWord(FIVE);
      lightWord(PAST);
      showHours(h);
      break;

    case 2:   // 08–12 → TEN PAST
      lightWord(TEN);
      lightWord(PAST);
      showHours(h);
      break;

    case 3:   // 13–17 → QUARTER PAST
      lightWord(QUARTER);
      lightWord(PAST);
      showHours(h);
      break;

    case 4:   // 18–22 → TWENTY PAST
      lightWord(TWENTY);
      lightWord(PAST);
      showHours(h);
      break;

    case 5:   // 23–27 → TWENTY FIVE PAST
      lightWord(TWENTY);
      lightWord(FIVE);
      lightWord(PAST);
      break;

    case 6:   // 28–32 → HALF PAST
      lightWord(HALF);
      lightWord(PAST);
      showHours(h);
      break;

    case 7:   // 33–37 → TWENTY FIVE TO
      lightWord(TWENTY);
      lightWord(FIVE);
      lightWord(TO);
      showHours(h + 1);
      break;

    case 8:   // 38–42 → TWENTY TO
      lightWord(TWENTY);
      lightWord(TO);
      showHours(h + 1);
      break;

    case 9:   // 43–47 → QUARTER TO
      lightWord(QUARTER);
      lightWord(TO);
      showHours(h + 1);
      break;

    case 10:  // 48–52 → TEN TO
      lightWord(TEN);
      lightWord(TO);
      showHours(h + 1);
      break;

    case 11:  // 53–57 → FIVE TO
      lightWord(FIVE);
      lightWord(TO);
      showHours(h + 1);
      break;

    case 12:  // 58–59 → o'clock
      showHours(h + 1);
      lightWord(OCLOCK);
      break;
  }     
  
  FastLED.show();
}

void showHours(uint8_t h) {      
  uint8_t h12 = h % 12;
  if (h12 == 0) h12 = 12; 

  switch (h12) {
    case 1:  lightWord(h_ONE);     break;
    case 2:  lightWord(h_TWO);     break;
    case 3:  lightWord(h_THREE);   break;
    case 4:  lightWord(h_FOUR);    break;
    case 5:  lightWord(h_FIVE);    break;
    case 6:  lightWord(h_SIX);     break;
    case 7:  lightWord(h_SEVEN);   break;
    case 8:  lightWord(h_EIGHT);   break;
    case 9:  lightWord(h_NINE);    break;
    case 10: lightWord(h_TEN);     break;
    case 11: lightWord(h_ELEVEN);  break;
    case 12: lightWord(h_TWELVE);  break;
  }
}

void setup() {
  Serial.begin(115200);
  setupWifi();
  setupRTC();
  setupLED();
  setupDfPlayer();
  setupTouch();
}

void setupRTC() {
  if (! rtc.begin()) {
    Serial.println("Couldn't find RTC");
    Serial.flush();
    while (1) delay(10);
  }

  if (rtc.lostPower()) {
    Serial.println("RTC lost power, let's set the time!");
    rtc.adjust(DateTime(F(__DATE__), F(__TIME__)));
  }
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
  gClockColor = 0xFFFFFF;
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

void setupWifi() {
  WiFi.onEvent(WiFiEvent);
  
  WiFiManager wifiManager;
  wifiManager.setDebugOutput(true);
  wifiManager.setConfigPortalTimeout(180);
  
  if (!wifiManager.autoConnect("ESP32-Hotspot", "12345678")) {
    Serial.println("Wifi, failed!");
    ESP.restart();
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
      break;
    case ARDUINO_EVENT_WIFI_STA_DISCONNECTED:
      Serial.println("WiFi, disconnected!");
      WiFi.reconnect(); 
      break;
    case ARDUINO_EVENT_WIFI_STA_GOT_IP:
      Serial.print("WiFi, got IP");
      Serial.println(WiFi.localIP());
      break;
    default:
      break;
  }
}

void loop() { 
  DateTime now = rtc.now();  
  showTime(now.minute(), now.hour());
  delay(60000);
}