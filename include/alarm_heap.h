#include "alarm.h"
#include <vector>

class AlarmHeap
{
public:
    AlarmHeap();
    AlarmHeap(std::vector<Alarm*> list);
    Alarm* get_top();
    void pop_top();
    void insert(Alarm* alarm);
    bool update_alarm(uint32_t alarm_id, Alarm* new_alarm);
    bool remove_alarm(uint32_t alarm_id);
    int size();
    bool empty();
    DateTime get_next_occurrence(const Alarm *alarm);
    std::vector<Alarm*> alarms;
private:
    void check_exception();
    void swap(int i1, int i2);
    bool should_come_before(const Alarm *a, const Alarm *b);
    DateTime get_current_time();
    void sift_down(int index);
    void sift_up(int index);
};
