import { MdVolumeUp } from "react-icons/md";
// eslint-disable-next-line no-unused-vars
import { motion } from "framer-motion";
import { useDeviceVolume, useUpdateVolume } from "../hooks/useAlarms";
import { useState } from "react";

const VolumeSection = () => {
  const { data: volume = 15, isLoading } = useDeviceVolume();
  const mutation = useUpdateVolume();

  // برای اسلایدر نرم‌تر و کنترل‌شده
  const [localVolume, setLocalVolume] = useState(volume);

  // وقتی کاربر اسلایدر رو رها کرد، به سرور بفرست
  const handleChangeComplete = (value) => {
    mutation.mutate(value);
  };

  return (
    <motion.div
      initial={{ opacity: 0, y: 30, scale: 0.9 }}
      animate={{ opacity: 1, y: 0, scale: 1 }}
      transition={{ duration: 0.6, ease: "easeOut" }}
      className="relative w-full max-w-md"
    >
      <div className="absolute inset-0 bg-linear-to-r from-blue-600/20 to-blue-600/20 blur-3xl -z-10" />

      <div className="backdrop-blur-2xl ltr bg-white/10 border border-white/20 rounded-3xl shadow-2xl sm:p-8 p-4">
        {/* آیکون متحرک صدا */}
        <motion.div
          animate={{ rotate: [0, 10, -10, 0] }}
          transition={{ repeat: Infinity, duration: 5, ease: "easeInOut" }}
          className="flex justify-center mb-4"
        >
          <div className="p-5 bg-linear-to-br from-blue-500 to-blue-600 rounded-3xl shadow-2xl">
            <MdVolumeUp size={40} className="text-white" />
          </div>
        </motion.div>

        <h2 className="sm:text-xl text-lg font-bold text-center text-white mb-8 bg-clip-text bg-linear-to-r from-blue-400 to-blue-400">
          صدای دستگاه
        </h2>

        <div className="relative mb-10 px-8">
          <motion.div
            animate={{
              background: [
                "linear-gradient(90deg, #a855f7 0%, #ec4899 100%)",
                "linear-gradient(90deg, #ec4899 0%, #a855f7 100%)",
                "linear-gradient(90deg, #a855f7 0%, #ec4899 100%)",
              ],
            }}
            transition={{
              duration: 6,
              repeat: Infinity,
              repeatType: "reverse",
            }}
            className="absolute inset-x-8 -top-4 h-2 rounded-full opacity-30 blur-xl"
          />

          <input
            type="range"
            min="0"
            max="30"
            value={isLoading ? 15 : localVolume}
            onChange={(e) => setLocalVolume(Number(e.target.value))}
            onMouseUp={(e) => handleChangeComplete(Number(e.target.value))}
            onTouchEnd={(e) => handleChangeComplete(Number(e.target.value))}
            className="w-full h-10 bg-white/10 rounded-full appearance-none cursor-pointer slider-thumb-glow"
            style={{
              background: `linear-gradient(to right, 
                oklch(54.6% 0.245 262.881) 0%, 
                oklch(54.6% 0.245 262.881) ${(localVolume / 30) * 100}%, 
                rgba(255,255,255,0.12) ${(localVolume / 30) * 100}%, 
                rgba(255,255,255,0.12) 100%)`,
            }}
          />
        </div>

        {mutation.isPending && (
          <motion.p
            animate={{ opacity: [0.5, 1, 0.5] }}
            transition={{ duration: 1.5, repeat: Infinity }}
            className="text-center text-blue-400 font-medium my-4"
          >
            در حال تغییر صدا...
          </motion.p>
        )}

        <motion.div className="p-5 bg-white/10 backdrop-blur-md border border-white/30 rounded-2xl text-center">
          <p className="text-white/70 sm:text-sm text-xs mb-1">
            میزان صدای فعلی
          </p>
          <p className="sm:text-2xl text-xl font-mono tracking-widest text-white">
            {isLoading ? "..." : localVolume} / 30
          </p>
        </motion.div>
      </div>
    </motion.div>
  );
};

export default VolumeSection;
