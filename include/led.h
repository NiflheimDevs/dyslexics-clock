#include <FastLED.h>

#define FPSerial Serial1
#define LED_PIN 27
#define NUM_LEDS 121
extern uint8_t BRIGHTNESS;

extern CRGB leds[NUM_LEDS];

extern CRGB gClockColor;

extern bool is_birthday;

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

void showTime(uint8_t m, uint8_t h);
void showHours(uint8_t h);
void lightWord(const uint8_t word[2]);
void showHappyBirthdayAnimation(void *pvParameters);
