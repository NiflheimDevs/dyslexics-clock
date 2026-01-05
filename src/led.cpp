#include "led.h"
#include "Arduino.h"
#include "FastLED.h"

CRGB leds[NUM_LEDS];
CRGB gClockColor = CRGB::Green;
uint8_t BRIGHTNESS = 185;

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

void lightWord(const uint8_t word[2]) {
  for (uint8_t i = word[0]; i <= word[1]; i++) {
    leds[i] = gClockColor;
  }
}

void showTime(uint8_t m, uint8_t h) {
  FastLED.clear();
  FastLED.show();
  FastLED.setBrightness(BRIGHTNESS);

  lightWord(IT);
  lightWord(IS);

  uint8_t block = (m + 2) / 5;

  switch (block) {
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
    showHours(h);
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

void showHours(uint8_t h) {
  uint8_t h12 = h % 12;
  if (h12 == 0)
    h12 = 12;

  switch (h12) {
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

void showHappyBirthdayAnimation(void *pvParameters) {
  CRGB colors[] = {CRGB::Red,   CRGB::Green, CRGB::Blue,
                   CRGB::Yellow, CRGB::Purple, CRGB::Orange};
  int num_colors = sizeof(colors) / sizeof(colors[0]);
  int color_index = 0;

  for (;;) { // Infinite loop
    // Light up HAPPY
    for (uint8_t i = HAPPY[0]; i <= HAPPY[1]; i++) {
      leds[i] = colors[color_index];
    }

    // Light up BIRTHDAY
    for (uint8_t i = BIRTHDAY[0]; i <= BIRTHDAY[1]; i++) {
      leds[i] = colors[color_index];
    }

    FastLED.show();
    vTaskDelay(pdMS_TO_TICKS(500)); // On time

    FastLED.clear();
    FastLED.show();
    vTaskDelay(pdMS_TO_TICKS(500)); // Off time

    color_index = (color_index + 1) % num_colors; // Change color
  }
}