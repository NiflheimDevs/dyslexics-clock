import { useEffect, useState } from "react";
import { motion } from "framer-motion";
import { Calendar } from "lucide-react";
import { useBirthDate, useUpdateBirthDate } from "../hooks/useAlarms";

const BirthDateCard = () => {
  const { data: birthDate, isLoading } = useBirthDate();
  const updateBirthDateMutation = useUpdateBirthDate();

  const [date, setDate] = useState("");
  const [newdate, setNewDate] = useState(date);

  useEffect(() => {
    if (birthDate) {
      setDate(birthDate.split("T")[0]);
      setNewDate(birthDate.split("T")[0]);
    }
  }, [birthDate]);

  const handleUpdate = () => {
    updateBirthDateMutation.mutate({
      data: { birthDate: newdate },
    });
  };

  return (
    <motion.div
      initial={{ opacity: 0, y: 30, scale: 0.95 }}
      animate={{ opacity: 1, y: 0, scale: 1 }}
      transition={{ duration: 0.6, ease: "easeOut" }}
      className="w-full max-w-md"
    >
      <div className="backdrop-blur-2xl bg-white/5 border border-white/20 rounded-3xl shadow-2xl p-8">
        {/* Header */}
        <div className="flex items-center justify-center gap-3 mb-6">
          <div className="p-4 bg-linear-to-br from-blue-500 to-blue-600 rounded-2xl shadow-xl">
            <Calendar className="text-white" />
          </div>
          <h2 className="text-2xl font-bold text-white">
            تاریخ تولد 🎂
          </h2>
        </div>

        {/* Content */}
        {isLoading ? (
          <p className="text-center text-white/60">در حال دریافت...</p>
        ) : (
          <div className="space-y-4">
            {/* Current Date */}
            <div className="text-center text-white/80 text-sm">
              تاریخ تولد
            </div>
            <div className="text-center font-mono text-lg text-white bg-white/10 py-3 rounded-xl">
              {date || "—"}
            </div>

            {/* Date Picker */}
            <div className="mt-6">
              <label className="block text-center text-white/80 mb-2 text-sm">
                تغییر تاریخ تولد
              </label>
              <input
                type="date"
                value={newdate}
                onChange={(e) => setNewDate(e.target.value)}
                className="w-full px-4 py-3 rounded-xl bg-white/5 border border-white/20 text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>

            {/* Update Button */}
            <motion.button
              whileHover={{ scale: 1.03 }}
              whileTap={{ scale: 0.95 }}
              onClick={handleUpdate}
              disabled={updateBirthDateMutation.isPending}
              className="w-full mt-6 py-4 bg-linear-to-r cursor-pointer from-blue-600 to-blue-600 text-white font-semibold rounded-2xl shadow-xl disabled:opacity-70 transition-all"
            >
              {updateBirthDateMutation.isPending
                ? "در حال بروزرسانی..."
                : "ذخیره تغییرات"}
            </motion.button>
          </div>
        )}
      </div>
    </motion.div>
  );
};

export default BirthDateCard;
