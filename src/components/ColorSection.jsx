import { MdPalette } from "react-icons/md";
import { HexColorPicker } from "react-colorful";
import { useDeviceColor, useUpdateColor } from "../hooks/useAlarms";
// eslint-disable-next-line no-unused-vars
import { motion } from "framer-motion";

const ColorSection = () => {
  const { data: color, isLoading } = useDeviceColor();
  const mutation = useUpdateColor();
  const currentColor = color;

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
          animate={{ rotate: [0, 12, -12, 0] }}
          transition={{ repeat: Infinity, duration: 6 }}
          className="flex justify-center mb-4"
        >
          <div className="p-5 bg-linear-to-br from-blue-500 to-blue-600 rounded-3xl shadow-2xl">
            <MdPalette size={38} className="text-white" />
          </div>
        </motion.div>

        <h2 className="sm:text-xl text-lg font-bold text-center text-white mb-8 bg-clip-text bg-linear-to-r from-blue-400 to-blue-400">
          رنگ دستگاه
        </h2>

        <div className="mb-8 flex justify-center">
          <motion.div
            animate={{
              boxShadow: [
                "0 0 30px rgba(168, 85, 247, 0.35)",
                "0 0 50px rgba(59, 130, 246, 0.55)",
                "0 0 30px rgba(168, 85, 247, 0.35)",
              ],
            }}
            transition={{ duration: 4, repeat: Infinity }}
            className="p-4 bg-white/10 backdrop-blur-xl rounded-3xl border border-white/20"
          >
            <HexColorPicker
              color={currentColor}
              onChange={(c) => mutation.mutate(c)}
              className="w-56! h-56!"   
            />
          </motion.div>
        </div>

        {mutation.isPending && (
          <motion.p
            animate={{ opacity: [0.6, 1, 0.6] }}
            transition={{ duration: 1.2, repeat: Infinity }}
            className="text-center text-cyan-400 font-medium my-3 mb-6"
          >
            در حال اعمال رنگ...
          </motion.p>
        )}
        
        <motion.div
          className="p-5 bg-white/10 backdrop-blur-md border border-white/30 rounded-2xl text-center"
        >
          <p className="text-white/70 sm:text-sm text-xs mb-1">رنگ فعلی</p>
          <p className="sm:text-2xl text-xl font-mono tracking-wider text-white">
            {isLoading ? "..." : currentColor.toUpperCase()}
          </p>
        </motion.div>
      </div>
    </motion.div>
  );
};

export default ColorSection;