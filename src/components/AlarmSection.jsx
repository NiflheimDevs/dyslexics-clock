import { useState } from "react";
import { MdAccessAlarm, MdDelete, MdAdd } from "react-icons/md";
import { useAlarms, useCreateAlarm, useDeleteAlarm } from "../hooks/useAlarms";
// eslint-disable-next-line no-unused-vars
import { AnimatePresence, motion } from "framer-motion";

const daysOfWeek = [
  { id: 6, label: "ش" },
  { id: 0, label: "ی" },
  { id: 1, label: "د" },
  { id: 2, label: "س" },
  { id: 3, label: "چ" },
  { id: 4, label: "پ" },
  { id: 5, label: "ج" },
];

const AlarmSection = () => {
  const [time, setTime] = useState("");
  const [repeatType, setRepeatType] = useState("once");
  const [selectedDays, setSelectedDays] = useState([]);

  const { data, isLoading } = useAlarms();
  const alarms = Array.isArray(data) ? data : data?.data || data?.alarms || [];
  const createAlarm = useCreateAlarm();
  const deleteAlarm = useDeleteAlarm();
  console.log(alarms);

  const addAlarm = (e) => {
    e.preventDefault();
    if (!time) return;

    createAlarm.mutate(
      {
        time: `2025-12-05T${time}:00Z`,
        is_repeat: repeatType !== "once",
        days:
          repeatType === "weekly"
            ? selectedDays
            : repeatType === "daily"
            ? [0, 1, 2, 3, 4, 5, 6]
            : [],
      },
      {
        onSuccess: () => {
          setTime("");
          setRepeatType("once");
          setSelectedDays([]);
        },
      }
    );
  };

  const toggleDay = (id) => {
    setSelectedDays((prev) =>
      prev.includes(id) ? prev.filter((d) => d !== id) : [...prev, id]
    );
  };

  return (
    <motion.div
      initial={{ opacity: 0, y: 30, scale: 0.9 }}
      animate={{ opacity: 1, y: 0, scale: 1 }}
      transition={{ duration: 0.6, ease: "easeOut" }}
      className="relative w-full max-w-md"
    >
      <div className="absolute inset-0 bg-linear-to-r from-blue-600/20 to-blue-600/20 blur-3xl -z-10" />

      <div className="backdrop-blur-2xl bg-white/5 border border-white/20 rounded-3xl shadow-2xl sm:p-8 p-4">
        <motion.div
          animate={{ rotate: [0, 10, -10, 0] }}
          transition={{ repeat: Infinity, duration: 6 }}
          className="flex justify-center mb-4"
        >
          <div className="p-5 bg-linear-to-br from-blue-500 to-blue-600 rounded-3xl shadow-2xl">
            <MdAccessAlarm size={38} className="text-white" />
          </div>
        </motion.div>

        <h2 className="sm:text-xl text-lg font-bold text-center text-white mb-8 bg-clip-text bg-linear-to-r from-blue-400 to-blue-400">
          آلارم‌ها
        </h2>

        <form onSubmit={addAlarm} className="space-y-5 mb-8">
          <motion.div>
            <input
              type="time"
              value={time}
              onChange={(e) => setTime(e.target.value)}
              required
              className="w-full px-6 py-5 sm:text-lg text-md bg-white/10 backdrop-blur-md border border-white/30 rounded-2xl text-white placeholder-white/50 focus:outline-none focus:ring-4 focus:ring-blue-500/50 transition-all"
            />
          </motion.div>

          <div className="flex gap-3">
            <select
              value={repeatType}
              onChange={(e) => {
                setRepeatType(e.target.value);
                if (e.target.value !== "weekly") setSelectedDays([]);
              }}
              className="flex-1 px-5 py-5 text-base bg-white/10 backdrop-blur-md border border-white/30 rounded-2xl text-white focus:outline-none focus:ring-4 focus:ring-blue-500/50"
            >
              <option value="once" className="bg-gray-900 sm:text-lg text-md">
                یکبار
              </option>
              <option value="daily" className="bg-gray-900 sm:text-lg text-md">
                روزانه
              </option>
              <option value="weekly" className="bg-gray-900 sm:text-lg text-md">
                هفتگی
              </option>
            </select>

            <motion.button
              whileHover={{ scale: 1.02 }}
              whileTap={{ scale: 0.95 }}
              type="submit"
              disabled={!time || createAlarm.isPending}
              className="px-7 sm:text-md text-sm cursor-pointer py-5 bg-linear-to-r from-blue-600 to-blue-600 text-white font-bold rounded-2xl shadow-2xl flex items-center justify-center gap-2 disabled:opacity-70"
            >
              <MdAdd size={24} />
              {createAlarm.isPending ? "..." : "افزودن"}
            </motion.button>
          </div>

          <AnimatePresence>
            {repeatType === "weekly" && (
              <motion.div
                initial={{ opacity: 0, y: -10 }}
                animate={{ opacity: 1, y: 0 }}
                exit={{ opacity: 0, y: -10 }}
                className="flex flex-wrap gap-2 justify-center pt-2"
              >
                {daysOfWeek.map((day) => (
                  <motion.button
                    key={day.id}
                    whileTap={{ scale: 0.9 }}
                    type="button"
                    onClick={() => toggleDay(day.id)}
                    className={`w-10 h-10 rounded-full font-bold text-sm transition-all ${
                      selectedDays.includes(day.id)
                        ? "bg-linear-to-r from-blue-500 to-blue-500 text-white shadow-lg"
                        : "bg-white/10 text-white/70 border border-white/20"
                    }`}
                  >
                    {day.label}
                  </motion.button>
                ))}
              </motion.div>
            )}
          </AnimatePresence>
        </form>

        <div className="space-y-4 max-h-65 overflow-y-auto scrollbar-thin scrollbar-thumb-white/20 scrollbar-track-transparent pl-2">
          <AnimatePresence>
            {isLoading ? (
              <p className="text-center text-white/50 py-6">
                در حال بارگذاری...
              </p>
            ) : alarms.length === 0 ? (
              <p className="text-center text-white/40 italic py-4">
                هنوز آلارمی تنظیم نکردی
              </p>
            ) : (
              alarms.map((alarm, i) => (
                <motion.div
                  key={alarm.id}
                  initial={{ opacity: 0, x: -30 }}
                  animate={{ opacity: 1, x: 0 }}
                  exit={{ opacity: 0, x: 30 }}
                  transition={{ delay: i * 0.05 }}
                  className="group relative bg-white/10 backdrop-blur-md border border-white/20 rounded-2xl p-5 flex items-center justify-between overflow-hidden"
                >
                  <div className="flex items-center gap-4">
                    <div className="p-3 bg-linear-to-br from-blue-500 to-blue-600 rounded-2xl">
                      <MdAccessAlarm size={24} className="text-white" />
                    </div>
                    <div>
                      <p className="sm:text-xl text-lg font-bold text-white">
                        {new Date(alarm.time).toISOString()?.split('T')[1]?.slice(0,5)}
                      </p>
                      <p className="sm:text-sm text-xs text-white/60">
                        {alarm.is_repeat
                          ? alarm.days?.length === 7
                            ? "هر روز"
                            : alarm.days
                                ?.sort((a, b) => a - b)
                                .map(
                                  (dayId) =>
                                    daysOfWeek.find((d) => d.id === dayId)
                                      ?.label
                                )
                                .join(" ") || "—"
                          : "یکبار"}
                      </p>
                    </div>
                  </div>
                  <motion.button
                    whileHover={{ scale: 1.02 }}
                    whileTap={{ scale: 0.9 }}
                    onClick={() => deleteAlarm.mutate(alarm.id)}
                    className="p-3 bg-red-600/20 rounded-full text-red-400 cursor-pointer transition-all"
                  >
                    <MdDelete size={20} />
                  </motion.button>
                </motion.div>
              ))
            )}
          </AnimatePresence>
        </div>
      </div>
    </motion.div>
  );
};

export default AlarmSection;
