#include <alarm_heap.h>

extern RTC_DS3231 rtc;

AlarmHeap::AlarmHeap() {}

AlarmHeap::AlarmHeap(std::vector<Alarm *> list) {
  for (int i = 0; i < list.size(); i++) {
    this->alarms.push_back(list[i]);
  }
  for (int i = alarms.size() / 2; i >= 0; i--) {
    this->sift_down(i);
  }
}

Alarm *AlarmHeap::get_top() {
  this->check_exception();
  return alarms[0];
}

void AlarmHeap::pop_top() {
  this->check_exception();
  alarms[0] = alarms[alarms.size() - 1];
  alarms.pop_back();
  this->sift_down(0);
}

void AlarmHeap::insert(Alarm *alarm) {
  for (int i = 0; i < alarms.size(); i++) {
    if (alarms[i]->id == alarm->id) {
      alarms[i] = alarm;
      return;
    }
  }
  alarms.push_back(alarm);
  this->sift_up(alarms.size() - 1);
}

bool AlarmHeap::update_alarm(uint32_t alarm_id, Alarm *new_alarm) {
  for (int i = 0; i < alarms.size(); i++) {
    if (alarms[i]->id == alarm_id) {
      Alarm *old_alarm = alarms[i];
      alarms[i] = new_alarm;

      if (should_come_before(new_alarm, old_alarm)) {
        sift_up(i);
      } else {
        sift_down(i);
      }
      return true;
    }
  }
  return false;
}

bool AlarmHeap::remove_alarm(uint32_t alarm_id) {
  for (int i = 0; i < alarms.size(); i++) {
    if (alarms[i]->id == alarm_id) {
      alarms[i] = alarms[alarms.size() - 1];
      alarms.pop_back();
      if (i < alarms.size()) {
        sift_down(i);
        sift_up(i);
      }
      return true;
    }
  }
  return false;
}

int AlarmHeap::size() { return alarms.size(); }

bool AlarmHeap::empty() { return alarms.empty(); }

DateTime AlarmHeap::get_next_occurrence(const Alarm *alarm) {
  DateTime now = get_current_time();

  if (!alarm->is_repeat) {
    return alarm->timestamp;
  }

  uint8_t alarm_hour = alarm->timestamp.hour();
  uint8_t alarm_minute = alarm->timestamp.minute();
  uint8_t alarm_second = alarm->timestamp.second();

  bool has_today = false;
  for (uint8_t i = 0; i < alarm->repeating_days_count; i++) {
    if (alarm->repeating_days[i] == now.dayOfTheWeek()) {
      has_today = true;
      break;
    }
  }

  if (has_today) {
    DateTime today_alarm = DateTime(now.year(), now.month(), now.day(),
                                    alarm_hour, alarm_minute, alarm_second);
    if (today_alarm >= now) {
      return today_alarm;
    }
  }

  for (int days_ahead = 1; days_ahead <= 7; days_ahead++) {
    uint8_t next_wday = (now.dayOfTheWeek() + days_ahead) % 7;
    for (uint8_t i = 0; i < alarm->repeating_days_count; i++) {
      if (alarm->repeating_days[i] == next_wday) {
        DateTime next = now + TimeSpan(days_ahead * 86400L);
        return DateTime(next.year(), next.month(), next.day(), alarm_hour,
                        alarm_minute, alarm_second);
      }
    }
  }
  return now + TimeSpan(7 * 86400L);
}

void AlarmHeap::check_exception() {
  if (alarms.size() == 0) {
    Serial.println("ERROR: Alarm heap is empty!");
    while (1)
      ;
  }
}

void AlarmHeap::swap(int i1, int i2) {
  Alarm *temp = alarms[i1];
  alarms[i1] = alarms[i2];
  alarms[i2] = temp;
}

bool AlarmHeap::should_come_before(const Alarm *a, const Alarm *b) {
  DateTime next_a = get_next_occurrence(a);
  DateTime next_b = get_next_occurrence(b);
  return next_a < next_b;
}

DateTime AlarmHeap::get_current_time() { return rtc.now(); }

void AlarmHeap::sift_down(int index) {
  int left = index * 2 + 1;
  int right = index * 2 + 2;
  if (left >= alarms.size())
    return;

  int target_index;
  if (right >= alarms.size()) {
    target_index = left;
  } else {
    target_index =
        should_come_before(alarms[left], alarms[right]) ? left : right;
  }

  if (should_come_before(alarms[target_index], alarms[index])) {
    this->swap(target_index, index);
    this->sift_down(target_index);
  }
}

void AlarmHeap::sift_up(int index) {
  if (index == 0)
    return;
  int parent = ((index % 2 == 0) ? ((index / 2) - 1) : (index / 2));
  if (should_come_before(alarms[index], alarms[parent])) {
    this->swap(parent, index);

    this->sift_up(parent);
  }
}
