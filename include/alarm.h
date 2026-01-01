#include "weekday.h"
#include "RTClib.h"
#include <cstdint>
struct Alarm
{
    uint32_t id;
    DateTime timestamp;
    bool is_repeat;
    Weekday repeating_days[7];    // Array of repeating days
    uint8_t repeating_days_count; // How many days are actually set
};

