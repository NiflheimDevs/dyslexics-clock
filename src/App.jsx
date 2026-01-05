import { useState, useEffect } from "react";
import { MdAccessAlarm, MdPalette, MdVolumeUp, MdLogout, MdSettings } from "react-icons/md";
import AlarmSection from "./components/AlarmSection";
import ColorSection from "./components/ColorSection";
import VolumeSection from "./components/VolumeSection";
import BezanBekob from "./components/BezanBekob";
import Login from "./components/LoginSection";
// eslint-disable-next-line no-unused-vars
import { motion } from "framer-motion";

function App() {
  const [isLoggedIn, setIsLoggedIn] = useState(!!localStorage.getItem("token"));
  const [activeTab, setActiveTab] = useState("alarm");

  useEffect(() => {
    const checkToken = () => setIsLoggedIn(!!localStorage.getItem("token"));
    window.addEventListener("storage", checkToken);
    return () => window.removeEventListener("storage", checkToken);
  }, []);


  const handleLogout = () => {
    localStorage.clear();
    setIsLoggedIn(false);
  };

  if (!isLoggedIn) {
    return <Login onLoginSuccess={() => setIsLoggedIn(true)} />;
  }

  return (
    <>
      <div className="fixed inset-0 bg-linear-to-br from-gray-900 via-blue-900/30 to-gray-900 z-[-1]" />
      <div className="min-h-screen">
        <div className="sm:px-6 px-4">
          <header className="flex items-center justify-between mb-8 pt-4">
            <motion.button
              whileHover={{ scale: 1.02 }}
              whileTap={{ scale: 0.95 }}
              onClick={handleLogout}
              className="group cursor-pointer flex items-center justify-center sm:gap-3 gap-1.5 px-6 py-4 bg-white/5 backdrop-blur-xl border border-white/20 rounded-2xl text-white/80 hover:text-white transition-all duration-300 shadow-xl"
            >
              <MdLogout className="transition-transform sm:text-[20px] text-[18px] duration-500" />
              <span className="font-medium sm:text-md text-sm Vazirmatn">خروج</span>
            </motion.button>

            <motion.h1
              initial={{ opacity: 0, y: -30 }}
              animate={{ opacity: 1, y: 0 }}
              className="text-2xl sm:text-3xl font-extrabold bg-clip-text text-transparent bg-linear-to-r from-blue-400 via-cyan-400 to-blue-400"
            >
              Dyslexics Clock
            </motion.h1>
          </header>

          {/* تب‌ها */}
          <nav className="flex justify-center mb-10">
            <div className="bg-white/10 backdrop-blur-2xl border border-white/20 rounded-3xl p-3 shadow-2xl">
              <div className="flex gap-4 flex-wrap justify-center">
                {[
                  {
                    id: "BezanBekob",
                    label: "مدیریت",
                    icon: MdSettings,
                    gradient: "bg-blue-600",
                  },
                  {
                    id: "alarm",
                    label: "آلارم",
                    icon: MdAccessAlarm,
                    gradient: "bg-blue-600",
                  },
                  {
                    id: "color",
                    label: "رنگ",
                    icon: MdPalette,
                    gradient: "bg-blue-600",
                  },
                  {
                    id: "volume",
                    label: "صدا",
                    icon: MdVolumeUp,
                    gradient: "bg-blue-600",
                  },
                ].map((tab) => (
                  <motion.button
                    key={tab.id}
                    whileHover={{ scale: 1.08 }}
                    whileTap={{ scale: 0.95 }}
                    onClick={() => setActiveTab(tab.id)}
                    className={`relative cursor-pointer px-3 sm:px-6 py-3 rounded-2xl font-bold sm:text-lg text-md flex items-center sm:gap-3 gap-2 transition-all duration-300 ${
                      activeTab === tab.id
                        ? "text-white shadow-2xl"
                        : "text-white/60 hover:text-white/90"
                    }`}
                  >
                    {activeTab === tab.id && (
                      <motion.div
                        layoutId="activeTab"
                        className={`absolute inset-0 bg-linear-to-r ${tab.gradient} rounded-2xl -z-10 shadow-lg`}
                        initial={false}
                        transition={{
                          type: "spring",
                          stiffness: 400,
                          damping: 30,
                        }}
                      />
                    )}
                    <tab.icon className="sm:text-[28px] text-[22px]" />
                    {tab.label}
                  </motion.button>
                ))}
              </div>
            </div>
          </nav>

          <motion.main
            key={activeTab}
            initial={{ opacity: 0, y: 30, scale: 0.95 }}
            animate={{ opacity: 1, y: 0, scale: 1 }}
            exit={{ opacity: 0, y: -30 }}
            transition={{ duration: 0.5, ease: "easeOut" }}
            className="flex justify-center"
          >
            {activeTab === "BezanBekob" && <BezanBekob />}
            {activeTab === "alarm" && <AlarmSection />}
            {activeTab === "color" && <ColorSection />}
            {activeTab === "volume" && <VolumeSection />}
          </motion.main>

        </div>
      </div>
    </>
  );
}

export default App;
