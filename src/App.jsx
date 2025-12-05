import { useState, useEffect } from "react";
import { MdAccessAlarm, MdPalette, MdLogout } from "react-icons/md";
import AlarmSection from "./components/AlarmSection";
import ColorSection from "./components/ColorSection";
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
    <div className="min-h-screen bg-linear-to-br from-gray-900 via-blue-900/30 to-gray-900">
      <div className="sm:px-6 px-4">
        <header className="flex items-center justify-between mb-12 pt-4">
          <motion.h1
            initial={{ opacity: 0, y: -30 }}
            animate={{ opacity: 1, y: 0 }}
            className="text-2xl sm:text-3xl font-extrabold bg-clip-text text-transparent bg-linear-to-r from-blue-400 via-blue-400 to-blue-400"
          >
            Dyslexics Clock
          </motion.h1>
          <motion.button
            whileHover={{ scale: 1.02 }}
            whileTap={{ scale: 0.95 }}
            onClick={handleLogout}
            className="group cursor-pointer flex items-center gap-3 px-6 py-4 bg-white/5 backdrop-blur-xl border border-white/20 rounded-2xl text-white/80 hover:text-white transition-all duration-300 shadow-xl"
          >
            <MdLogout
              size={20}
              className="transition-transform duration-500"
            />
            <span className="font-medium text-md hidden sm:block">خروج</span>
          </motion.button>
        </header>
        <nav className="flex justify-center mb-4">
          <div className="bg-white/10 backdrop-blur-2xl border border-white/20 rounded-3xl p-3 shadow-2xl">
            <div className="flex gap-4">
              {[
                {
                  id: "alarm",
                  label: "آلارم",
                  icon: MdAccessAlarm,
                  gradient: "from-blue-600 to-blue-600",
                },
                {
                  id: "color",
                  label: "رنگ",
                  icon: MdPalette,
                  gradient: "from-cyan-600 to-blue-600",
                },
              ].map((tab) => (
                <motion.button
                  key={tab.id}
                  whileHover={{ scale: 1.08 }}
                  whileTap={{ scale: 0.95 }}
                  onClick={() => setActiveTab(tab.id)}
                  className={`relative cursor-pointer px-6 py-3 rounded-2xl font-bold sm:text-lg text-md flex items-center gap-3 transition-all duration-300 ${
                    activeTab === tab.id
                      ? "text-white shadow-2xl"
                      : "text-white/60 hover:text-white/90"
                  }`}
                >
                  {activeTab === tab.id && (
                    <motion.div
                      layoutId="activeTab"
                      className={`absolute inset-0 bg-linear-to-r ${tab.gradient} rounded-2xl -z-10`}
                      initial={false}
                      transition={{
                        type: "spring",
                        stiffness: 400,
                        damping: 30,
                      }}
                    />
                  )}
                  <tab.icon size={28} />
                  {tab.label}
                </motion.button>
              ))}
            </div>
          </div>
        </nav>

        <motion.main
          key={activeTab}
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          exit={{ opacity: 0, y: -20 }}
          transition={{ duration: 0.4 }}
          className="flex justify-center py-4"
        >
          {activeTab === "alarm" ? <AlarmSection /> : <ColorSection />}
        </motion.main>
      </div>
    </div>
  );
}

export default App;
