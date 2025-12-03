import { useState, useEffect } from "react";
import { MdAccessAlarm, MdPalette } from "react-icons/md";
import mqtt from "mqtt";
import AlarmSection from "./components/AlarmSection";
import ColorSection from "./components/ColorSection";

function App() {
  const [activeTab, setActiveTab] = useState("alarm");
  const [client, setClient] = useState(null);

  useEffect(() => {
    const mqttClient = mqtt.connect("wss://test.mosquitto.org:8081", {
      clientId: "react-alarm-app-" + Math.random().toString(16).substr(2, 8),
    });

    mqttClient.on("connect", () => {
      console.log("MQTT متصل شد!");
      setClient(mqttClient);
    });

    mqttClient.on("error", (err) => console.error("خطای MQTT:", err));

    return () => {
      if (mqttClient) mqttClient.end();
    };
  }, []);

  const publishMessage = (topic, message) => {
    if (client && client.connected) {
      client.publish(topic, JSON.stringify(message), { qos: 1 }, (err) => {
        if (err) console.error("خطای ارسال:", err);
        else console.log(`پیام ارسال شد به ${topic}:`, message);
      });
    } else {
      console.warn("MQTT وصل نیست!");
    }
  };

  return (
    <div className="min-h-screen p-4 sm:p-6 lg:p-8 custom-bg transition-colors duration-300 font-sans antialiased">
      {/* Header */}
      <header className="flex flex-col sm:flex-row justify-between items-center gap-4 mb-8 text-center">
        <h1 className="text-2xl sm:text-3xl font-bold text-gray-100 tracking-tight">
          Word Clock
        </h1>

        <div className="flex items-center gap-4">
          {/* وضعیت MQTT */}
          <span
            className={`text-sm font-medium px-3 py-1 rounded-full ${
              client && client.connected
                ? "bg-green-800 text-green-200"
                : "bg-red-800 text-red-200"
            }`}
          >
            {client && client.connected ? "Connected" : "Not Connected"}
          </span>
        </div>
      </header>

      {/* Navigation Tabs */}
      <nav className="flex flex-col sm:flex-row justify-center mb-8 gap-2 sm:gap-0 items-center">
        <div className="rounded-xl shadow-lg border border-gray-600 flex w-fit">
          <button
            onClick={() => setActiveTab("alarm")}
            className={`px-6 py-3 cursor-pointer rounded-lg font-semibold transition-all duration-300 transform hover:scale-105 ${
              activeTab === "alarm"
                ? "bg-gray-700 text-gray-300 hover:bg-gray-600"
                : "bg-alarm text-white shadow-xl"
            } animate-fade-in`}
          >
            <MdAccessAlarm className="inline mr-2" size={20} /> آلارم
          </button>
          <button
            onClick={() => setActiveTab("color")}
            className={`px-6 py-3 cursor-pointer rounded-lg font-semibold transition-all duration-300 transform hover:scale-105 ${
              activeTab === "color"
                ? "bg-gray-700 text-gray-300 hover:bg-gray-600"
                : "bg-color text-white shadow-xl"
            } animate-fade-in`}
          >
            <MdPalette className="inline mr-2" size={20} /> رنگ
          </button>
        </div>
      </nav>

      {/* Main Content */}
      <main className="max-w-md mx-auto animate-fade-in">
        {activeTab === "alarm" ? (
          <AlarmSection publishMessage={publishMessage} />
        ) : (
          <ColorSection publishMessage={publishMessage} />
        )}
      </main>
    </div>
  );
}

export default App;
