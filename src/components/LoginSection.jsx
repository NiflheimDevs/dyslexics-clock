import { useState } from "react";
import { useMutation } from "@tanstack/react-query";
import { login } from "../api/services";
// eslint-disable-next-line no-unused-vars
import { motion } from "framer-motion";
import { MdLogin, MdAlarm, MdVisibility, MdVisibilityOff } from "react-icons/md";

const Login = ({ onLoginSuccess }) => {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);

  const mutation = useMutation({
    mutationFn: login,
    onSuccess: (response) => {
      const token = response.data.token;
      localStorage.setItem("token", token);
      localStorage.setItem(
        "user",
        JSON.stringify(response.data.user || { username })
      );
      onLoginSuccess();
    },
    onError: () => {
      alert("نام کاربری یا رمز نادرست است");
    },
  });

  const handleSubmit = (e) => {
    e.preventDefault();
    mutation.mutate({ username, password });
  };

  return (
    <div className="min-h-screen bg-linear-to-br from-gray-900 via-blue-900/20 to-gray-900 flex items-center justify-center p-6">
      <motion.div
        initial={{ opacity: 0, y: 30, scale: 0.9 }}
        animate={{ opacity: 1, y: 0, scale: 1 }}
        transition={{ duration: 0.6, ease: "easeOut" }}
        className="relative w-full max-w-md"
      >
        <div className="absolute inset-0 bg-linear-to-r from-blue-600/20 to-blue-600/20 blur-3xl -z-10" />
        <div className="backdrop-blur-2xl bg-white/5 border border-white/20 rounded-3xl shadow-2xl p-10">
          <motion.div
            animate={{ rotate: [0, 10, -10, 0] }}
            transition={{ repeat: Infinity, duration: 6 }}
            className="flex justify-center mb-8"
          >
            <div className="p-6 bg-linear-to-br from-blue-500 to-blue-600 rounded-3xl shadow-2xl">
              <MdAlarm size={48} className="text-white" />
            </div>
          </motion.div>
          <h1 className="text-4xl font-bold text-center text-white mb-10 bg-clip-text bg-linear-to-r from-blue-400 to-blue-400">
            Dyslexics Clock
          </h1>
          <form onSubmit={handleSubmit} className="space-y-6">
            <motion.div>
              <input
                type="text"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                placeholder="نام کاربری"
                className="w-full px-6 py-5 text-lg bg-white/10 backdrop-blur-md border border-white/30 rounded-2xl text-white placeholder-white/50 focus:outline-none focus:ring-4 focus:ring-blue-500/50 transition-all"
                required
              />
            </motion.div>

            <motion.div className="relative">
              <input
                type={showPassword ? "text" : "password"}
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="رمز عبور"
                className="w-full px-6 py-5 text-lg bg-white/10 backdrop-blur-md border border-white/30 rounded-2xl text-white placeholder-white/50 focus:outline-none focus:ring-4 focus:ring-blue-500/50 transition-all"
                required
              />
              <button
                type="button"
                onClick={() => setShowPassword(!showPassword)}
                className="absolute left-5 top-1/2 -translate-y-1/2 text-blue-600/95 hover:text-blue-600 transition-colors"
              >
                {showPassword ? (
                  <MdVisibilityOff size={26} />
                ) : (
                  <MdVisibility size={26} />
                )}
              </button>
            </motion.div>

            <motion.button
              whileHover={{ scale: 1.02 }}
              whileTap={{ scale: 0.95 }}
              type="submit"
              disabled={mutation.isPending}
              className="w-full cursor-pointer py-5 bg-linear-to-r from-blue-600 to-blue-600 text-white font-bold text-xl rounded-2xl shadow-2xl flex items-center justify-center gap-3 disabled:opacity-70"
            >
              <MdLogin size={28} />
              {mutation.isPending ? "در حال ورود..." : "ورود"}
            </motion.button>
          </form>
        </div>
      </motion.div>
    </div>
  );
};

export default Login;