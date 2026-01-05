import { useState, useEffect } from "react";
import { Clock, BellOff, Sun } from "lucide-react";
import { motion } from "framer-motion";
import { 
  useDeviceBrightness,
  useUpdateBrightness,
  useSnoozeAlarm,
  useStopAlarm,
  useBezanBekob
} from "../hooks/useAlarms";

const AlarmControls = () => {
  const [brightness, setBrightness] = useState(128);
  
  // Query to get current brightness
  const { data: currentBrightness } = useDeviceBrightness();
  
  // Mutations
  const updateBrightnessMutation = useUpdateBrightness();
  const snoozeAlarmMutation = useSnoozeAlarm();
  const stopAlarmMutation = useStopAlarm();
  const bezanBekoبMutation = useBezanBekob();

  // Update local brightness when data is fetched
  useEffect(() => {
    if (currentBrightness !== undefined) {
      setBrightness(currentBrightness);
    }
  }, [currentBrightness]);

  const handleStop = () => {
    stopAlarmMutation.mutate();
  };

  const handleSnooze = () => {
    snoozeAlarmMutation.mutate();
  };

  const handleBezan = async () => {
    try {
      console.log("بزن بکوب!");
      bezanBekoبMutation.mutate();
    } catch (error) {
      console.error("Error:", error);
    }
  };

  const handleBrightnessChange = (e) => {
    const value = parseInt(e.target.value);
    setBrightness(value);
    console.log("Brightness set to:", value);
  };

  const handleSubmit = () => {
    updateBrightnessMutation.mutate(brightness);
  };

  return (
    <motion.div
      initial={{ opacity: 0, y: 30, scale: 0.9 }}
      animate={{ opacity: 1, y: 0, scale: 1 }}
      transition={{ duration: 0.6, ease: "easeOut" }}
      className="relative w-full max-w-md"
    >
      <div className="absolute inset-0 bg-linear-to-r from-blue-600/20 to-blue-600/20 blur-3xl -z-10" />

      <div className="backdrop-blur-2xl bg-white/5 border border-white/20 rounded-3xl shadow-2xl p-8">
        <motion.div
          animate={{ rotate: [0, 10, -10, 0] }}
          transition={{ repeat: Infinity, duration: 6 }}
          className="flex justify-center mb-6"
        >
          <div className="p-5 bg-gradient-to-br from-blue-500 to-blue-600 rounded-3xl shadow-2xl">
            <Clock size={38} className="text-white" />
          </div>
        </motion.div>

        <h2 className="text-2xl font-bold text-center text-white mb-8">
          کنترل آلارم
        </h2>

        <div className="space-y-4">
          <div className="p-5 bg-white/5 backdrop-blur-xl border border-white/20 rounded-2xl">
            <div className="flex items-center justify-between mb-3">
              <div className="flex items-center gap-2">
                <Sun className="w-5 h-5 text-yellow-400" />
                <span className="text-white/90 font-medium">روشنایی</span>
              </div>
              <span className="text-white/70 text-sm font-mono bg-white/10 px-3 py-1 rounded-lg">
                {brightness}
              </span>
            </div>
            <input
              type="range"
              min="0"
              max="255"
              value={brightness}
              onChange={handleBrightnessChange}
              className="w-full h-2 bg-white/10 rounded-lg appearance-none cursor-pointer slider"
              dir="rtl"
              style={{
                background: `linear-gradient(to left, rgb(59, 130, 246) 0%, rgb(59, 130, 246) ${(brightness / 255) * 100}%, rgba(255,255,255,0.1) ${(brightness / 255) * 100}%, rgba(255,255,255,0.1) 100%)`
              }}
            />
            <div className="flex w-full mt-4 justify-center items-center">
              <motion.button
                whileHover={{ scale: 1.02 }}
                whileTap={{ scale: 0.95 }}
                type="button"
                onClick={handleSubmit}
                disabled={updateBrightnessMutation.isPending}
                className="px-6 sm:text-lg text-md cursor-pointer py-4 bg-gradient-to-r from-blue-600 to-blue-700 text-white font-semibold rounded-2xl shadow-2xl flex items-center justify-center gap-2 disabled:opacity-70 hover:from-blue-700 hover:to-blue-800 transition-all"
              >
                {updateBrightnessMutation.isPending ? "..." : "تغییر روشنایی"}
              </motion.button>
            </div>
          </div>

          {/* Bezan Bekob Button */}
          <motion.button
            whileHover={{ scale: 1.02 }}
            whileTap={{ scale: 0.95 }}
            onClick={handleBezan}
            disabled={bezanBekoبMutation.isPending}
            className="w-full group cursor-pointer flex items-center justify-center gap-3 px-6 py-5 bg-white/5 backdrop-blur-xl border border-white/20 rounded-2xl text-white/80 hover:text-white transition-all duration-300 shadow-xl hover:shadow-2xl disabled:opacity-70"
          >
            <span className="font-medium text-lg">
              {bezanBekoبMutation.isPending ? "..." : "🎉 بزن بکوب 💣"}
            </span>
          </motion.button>

          {/* Control Buttons Grid */}
          <div className="grid grid-cols-2 gap-4 mt-6">
            <motion.button
              whileHover={{ scale: 1.05 }}
              whileTap={{ scale: 0.95 }}
              onClick={handleSnooze}
              disabled={snoozeAlarmMutation.isPending}
              className="flex flex-col items-center justify-center gap-2 p-5 bg-blue-500/20 border border-blue-400/30 rounded-2xl text-blue-400 hover:bg-blue-500/30 transition-all disabled:opacity-70"
            >
              <Clock className="w-6 h-6" />
              <span className="text-sm font-medium">
                {snoozeAlarmMutation.isPending ? "..." : "اسنوز 5 دقیقه"}
              </span>
            </motion.button>

            <motion.button
              whileHover={{ scale: 1.05 }}
              whileTap={{ scale: 0.95 }}
              onClick={handleStop}
              disabled={stopAlarmMutation.isPending}
              className="flex flex-col items-center justify-center gap-2 p-5 bg-red-500/20 border border-red-400/30 rounded-2xl text-red-400 hover:bg-red-500/30 transition-all disabled:opacity-70"
            >
              <BellOff className="w-6 h-6" />
              <span className="text-sm font-medium">
                {stopAlarmMutation.isPending ? "..." : "قطع آلارم"}
              </span>
            </motion.button>
          </div>
        </div>
      </div>

      <style jsx>{`
        .slider::-webkit-slider-thumb {
          appearance: none;
          width: 20px;
          height: 20px;
          border-radius: 50%;
          background: white;
          cursor: pointer;
          box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
          transition: all 0.2s ease;
        }

        .slider::-webkit-slider-thumb:hover {
          transform: scale(1.2);
          box-shadow: 0 4px 12px rgba(59, 130, 246, 0.5);
        }

        .slider::-moz-range-thumb {
          width: 20px;
          height: 20px;
          border-radius: 50%;
          background: white;
          cursor: pointer;
          border: none;
          box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
          transition: all 0.2s ease;
        }

        .slider::-moz-range-thumb:hover {
          transform: scale(1.2);
          box-shadow: 0 4px 12px rgba(59, 130, 246, 0.5);
        }
      `}</style>
    </motion.div>
  );
};

export default AlarmControls;