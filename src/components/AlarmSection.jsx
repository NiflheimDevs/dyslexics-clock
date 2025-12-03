import { useState, useEffect } from 'react';
import { MdAccessAlarm, MdDelete } from 'react-icons/md';

const daysOfWeek = [
  { id: 0, label: 'یکشنبه' },
  { id: 1, label: 'دوشنبه' },
  { id: 2, label: 'سه‌شنبه' },
  { id: 3, label: 'چهارشنبه' },
  { id: 4, label: 'پنجشنبه' },
  { id: 5, label: 'جمعه' },
  { id: 6, label: 'شنبه' },
];

const AlarmSection = ({ publishMessage }) => {
  const [alarms, setAlarms] = useState(JSON.parse(localStorage.getItem('alarms')) || []);
  const [time, setTime] = useState('');
  const [repeatType, setRepeatType] = useState('once');
  const [selectedDays, setSelectedDays] = useState([]);

  const notify = (label) => {
    if ('Notification' in window && Notification.permission === 'granted') {
      new Notification('آلارم!', { body: label });
    }
  };

  useEffect(() => {
    localStorage.setItem('alarms', JSON.stringify(alarms));
    const interval = setInterval(() => {
      const now = new Date();
      const currentMinutes = now.getHours() * 60 + now.getMinutes();
      const currentDay = now.getDay();

      alarms.forEach(alarm => {
        const [hour, minute] = alarm.time.split(':').map(Number);
        const alarmMinutes = hour * 60 + minute;

        if (currentMinutes === alarmMinutes) {
          let shouldRing = true;

          if (alarm.repeatType === 'once' && alarm.rung) {
            shouldRing = false;
          } else if (alarm.repeatType === 'weekly') {
            shouldRing = selectedDays.includes(currentDay); // Note: selectedDays is global state, but should be alarm.days
            // Correction: use alarm.days
            shouldRing = alarm.days.includes(currentDay);
          }
          // For 'daily', always ring if time matches

          if (shouldRing) {
            notify(alarm.label);
            if (alarm.repeatType === 'once') {
              setAlarms(prev => prev.map(a => a.id === alarm.id ? { ...a, rung: true } : a));
            }
          }
        }
      });
    }, 60000); // Every minute

    return () => clearInterval(interval);
  }, [alarms]); // Dependency on alarms to re-run if alarms change

  const getRepeatLabel = (type, days = []) => {
    switch (type) {
      case 'daily': return ' (روزانه)';
      case 'weekly': {
        if (days.length === 0) return ' (هفتگی - بدون روز)';
        const dayLabels = days.map(d => daysOfWeek.find(day => day.id === d)?.label.slice(0, 2) || '').join('، ');
        return ` (هفتگی: ${dayLabels})`;
      }
      default: return ' (یکبار)';
    }
  };

  const addAlarm = (e) => {
    e.preventDefault();
    if (time) {
      const alarmDate = new Date(`1970-01-01T${time}:00`);
      const baseLabel = `آلارم ${alarmDate.toLocaleTimeString('fa-IR')}`;
      const repeatLabel = getRepeatLabel(repeatType, repeatType === 'weekly' ? selectedDays : []);
      const newAlarm = {
        id: Date.now(),
        time,
        label: baseLabel + repeatLabel,
        rung: repeatType === 'once' ? false : undefined, // Only for once
        repeatType,
        days: repeatType === 'weekly' ? selectedDays : [],
      };
      setAlarms([...alarms, newAlarm]);
      setTime('');
      setRepeatType('once');
      setSelectedDays([]);
      publishMessage('alarm/set', {
        time: time,
        label: newAlarm.label,
        repeatType,
        days: newAlarm.days,
      });
    }
  };

  const deleteAlarm = (id) => {
    setAlarms(alarms.filter(a => a.id !== id));
  };

  const updateSelectedDays = (dayId) => {
    setSelectedDays(prev => 
      prev.includes(dayId) ? prev.filter(d => d !== dayId) : [...prev, dayId]
    );
  };

  return (
    <div className="bg-gray-800 p-6 rounded-xl shadow-lg border border-gray-600 transition-all duration-300 hover:shadow-xl animate-fade-in">
      <h2 className="text-xl font-bold mb-6 text-center text-gray-100 flex items-center justify-center gap-2">
        <MdAccessAlarm size={24} className="text-alarm" /> تنظیم آلارم
      </h2>
      
      <form className="space-y-4 mb-6" onSubmit={addAlarm}>
        <input
          type="time"
          value={time}
          onChange={(e) => setTime(e.target.value)}
          className="w-full sm:block hidden p-3 border border-gray-600 rounded-lg bg-gray-700 text-gray-100 text-right placeholder-gray-400 focus:ring-2 focus:ring-alarm focus:border-transparent focus:outline-none transition-all duration-200"
          placeholder="زمان را انتخاب کنید"
        />
        <div className="sm:hidden relative">
          <input
            type="time"
            value={time}
            onChange={(e) => setTime(e.target.value)}
            className="w-full p-3 pr-24 border border-gray-600 rounded-lg bg-gray-700 text-gray-100 text-right placeholder-gray-400 focus:ring-2 focus:ring-alarm focus:border-transparent focus:outline-none transition-all duration-200"
          />
          {!time && (
            <span className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 pointer-events-none text-sm">
              زمان را انتخاب کنید
            </span>
          )}
        </div>
        
        <div className="flex flex-col sm:flex-row gap-3">
          <select
            value={repeatType}
            onChange={(e) => {
              setRepeatType(e.target.value);
              if (e.target.value !== 'weekly') setSelectedDays([]);
            }}
            className="flex-1 p-3 border border-gray-600 rounded-lg bg-gray-700 text-gray-100 text-right focus:ring-2 focus:ring-alarm focus:border-transparent focus:outline-none transition-all duration-200"
          >
            <option value="once">یکبار</option>
            <option value="daily">روزانه</option>
            <option value="weekly">هفتگی (روزهای خاص)</option>
          </select>
          
          <button
            type="submit"
            className="px-6 py-3 bg-alarm cursor-pointer bg-gray-700 font-semibold rounded-lg text-gray-300 hover:bg-gray-600 focus:ring-2 focus:ring-alarm focus:outline-none transition-all duration-200 shadow-md hover:shadow-lg disabled:opacity-50 disabled:cursor-not-allowed"
            disabled={!time}
          >
            اضافه کن و ارسال
          </button>
        </div>

        {repeatType === 'weekly' && (
          <div className="flex flex-wrap gap-2 justify-center">
            {daysOfWeek.map(day => (
              <label key={day.id} className="flex items-center space-x-1 space-x-reverse cursor-pointer">
                <input
                  type="checkbox"
                  checked={selectedDays.includes(day.id)}
                  onChange={() => updateSelectedDays(day.id)}
                  className="hidden"
                />
                <span className={`px-3 py-1 rounded-full text-sm font-medium transition-all duration-200 ${
                  selectedDays.includes(day.id)
                    ? 'bg-gray-700 text-gray-400 hover:bg-gray-600'
                    : 'bg-alarm bg-gray-800 p-6 rounded-xl shadow-lg border border-gray-600'
                }`}>
                  {day.label}
                </span>
              </label>
            ))}
          </div>
        )}
      </form>
      
      <ul className="space-y-3">
        {alarms.map(alarm => (
          <li
            key={alarm.id}
            className="flex justify-between items-center p-4 bg-gray-700 rounded-lg hover:bg-gray-600 transition-all duration-200 shadow-sm hover:shadow-md"
          >
            <div className="flex items-center gap-3">
              <MdAccessAlarm 
                size={20} 
                className={`text-${(alarm.repeatType === 'once' && alarm.rung) ? 'gray-500' : 'alarm'}`} 
              />
              <span className={`font-medium ${(alarm.repeatType === 'once' && alarm.rung) ? 'text-gray-500 line-through' : 'text-gray-100'}`}>
                {alarm.label} {alarm.repeatType === 'once' && alarm.rung ? '(زنگ خورده)' : ''}
              </span>
            </div>
            <button
              onClick={() => deleteAlarm(alarm.id)}
              className="p-2 text-red-400 cursor-pointer hover:text-red-300 hover:bg-red-900 rounded-full transition-all duration-200 focus:ring-2 focus:ring-red-500 focus:outline-none"
              aria-label="حذف آلارم"
            >
              <MdDelete size={18} />
            </button>
          </li>
        ))}
      </ul>
      {!alarms.length && (
        <p className="text-center text-gray-400 mt-6 italic animate-fade-in">
          هیچ آلارمی تنظیم نشده. یکی اضافه کنید!
        </p>
      )}
    </div>
  );
};

export default AlarmSection;