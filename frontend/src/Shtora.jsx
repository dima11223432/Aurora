import React, { useState, useEffect } from "react";
import axios from "axios";

const Shtora = () => {
  const [parsingChannels, setParsingChannels] = useState([]);
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [isOpen, setIsOpen] = useState(false);
  const [selectedChannels, setSelectedChannels] = useState([]);

  const getUserPriorityChannels = async () => {
    const TOKEN = localStorage.getItem("token");

    try {
      const response = await fetch(
        "http://localhost:8081/v1/get_user_priority_channels",
        {
          method: "GET",
          headers: {
            Authorization: `Bearer ${TOKEN}`,
            "Content-Type": "application/json",
          },
        }
      );

      const data = await response.json();
      console.log(data);
      const userPriorityChannels = data.channels || [];
      setSelectedChannels(userPriorityChannels);
      return data;
    } catch (error) {
      console.error("Error: ", error);
      return null;
    }
  };

  const setPriorityChannelsRequest = async (channels) => {
    const TOKEN = localStorage.getItem("token");
    if (!TOKEN) return;

    try {
      const response = await fetch(
        "http://localhost:8081/v1/set_priority_channels",
        {
          method: "POST",
          headers: {
            Authorization: `Bearer ${TOKEN}`,
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            priority_channels: channels,
          }),
        }
      );

      const data = await response.json();
      console.log("Set priority channels response:", data);
      return data;
    } catch (error) {
      console.error("Error setting priority channels:", error);
      throw error;
    }
  };
  const handleCheckboxChange = async (channel) => {
    let newSelectedChannels;
    if (selectedChannels.includes(channel)) {
      newSelectedChannels = selectedChannels.filter((ch) => ch !== channel);
    } else {
      newSelectedChannels = [...selectedChannels, channel];
    }
    setSelectedChannels(newSelectedChannels);
    await setPriorityChannelsRequest(newSelectedChannels);
  };

  useEffect(() => {
    const login = async () => {
      try {
        const res = await axios.post("http://localhost:8081/v1/login", {
          telegram_id: 123456789,
          username: "john_doe",
          first_name: "John",
          last_name: "Doe",
          is_admin: false,
          app_id: 1,
        });
        localStorage.setItem("token", res.data.token);
        setIsLoggedIn(true);
      } catch (e) {
        console.error("Login error:", e);
      }
    };
    login();
  }, []);

  useEffect(() => {
    const fetchParsingChannels = async () => {
      try {
        if (!isLoggedIn) return;
        const token = localStorage.getItem("token");
        const resp = await axios.post(
          "http://localhost:8081/v1/get_all_parsing_channels",
          {},
          {
            headers: {
              Authorization: `Bearer ${token}`,
            },
          },
        );
        const channels = resp.data.channels || [];
        setParsingChannels(channels);
        await getUserPriorityChannels();
      } catch (e) {
        console.log(e);
      }
    };
    fetchParsingChannels();
  }, [isLoggedIn]);

  return (
    <div className="fixed top-4 left-4 z-50">
      <button
        onClick={() => setIsOpen((prev) => !prev)}
        className="bg-gradient-to-r from-blue-700 via-blue-600 to-blue-800 hover:from-blue-800 hover:to-blue-900 text-white font-bold py-2 px-5 rounded-xl shadow-2xl focus:outline-none transition-all duration-300 hover:scale-105 hover:shadow-blue-500/40"
        style={{ backdropFilter: "blur(2px)" }}
      >
        {isOpen ? "Скрыть" : "Скрыть"}
      </button>
      <div
        className={`mt-2 w-72 max-w-xs ${isOpen ? "shtora-animate-in" : "shtora-animate-out pointer-events-none opacity-0"} bg-gray-900/90 rounded-2xl shadow-2xl border border-blue-700 p-5 absolute left-0 top-14 transition-all duration-500 backdrop-blur-xl`}
        style={{ boxShadow: "0 8px 32px 0 rgba(31, 38, 135, 0.37)" }}
      >
        <h3 className="text-xl font-bold text-blue-400 mb-3 drop-shadow">
          Каналы для парсинга
        </h3>
        <ul className="max-h-60 overflow-y-auto space-y-1">
          {parsingChannels.length === 0 ? (
            <li className="text-gray-400 text-sm italic animate-pulse">
              Нет доступных каналов
            </li>
          ) : (
            parsingChannels.map((channel, idx) => {
              const channelName = typeof channel === "string"
                ? channel
                : channel?.name || JSON.stringify(channel);
              const isChecked = selectedChannels.includes(channelName);
              return (
                <li
                  key={idx}
                  className="py-2 px-3 rounded-lg hover:bg-blue-700/70 hover:scale-[1.03] hover:shadow-lg text-white cursor-pointer transition-all duration-200 flex items-center gap-2 group"
                  style={{ backdropFilter: "blur(1px)" }}
                >
                  <input
                    type="checkbox"
                    id={`channel-${idx}`}
                    checked={isChecked}
                    onChange={() => handleCheckboxChange(channelName)}
                    className="w-4 h-4 text-blue-600 bg-gray-100 border-gray-300 rounded focus:ring-blue-500 focus:ring-2"
                  />
                  <label 
                    htmlFor={`channel-${idx}`}
                    className="flex-1 group-hover:text-blue-200 transition-colors duration-200 cursor-pointer"
                  >
                    {channelName}
                  </label>
                  <span className="ml-auto opacity-0 group-hover:opacity-100 text-xs text-blue-300 transition-opacity duration-200">
                    →
                  </span>
                </li>
              );
            })
          )}
        </ul>
        {/* {selectedChannels.length > 0 && (
          <div className="mt-3 pt-2 border-t border-blue-700/50 text-xs text-blue-300">
    
          </div>
        )} */}
      </div>
      <style>{`
        .shtora-animate-in {
          opacity: 1;
          pointer-events: auto;
          transform: translateY(0) scale(1);
          filter: drop-shadow(0 8px 32px rgba(31,38,135,0.37));
          transition: opacity 0.5s cubic-bezier(.4,2,.6,1), transform 0.5s cubic-bezier(.4,2,.6,1);
        }
        .shtora-animate-out {
          opacity: 0;
          pointer-events: none;
          transform: translateY(-30px) scale(0.95);
          filter: blur(2px);
          transition: opacity 0.4s, transform 0.4s, filter 0.4s;
        }
      `}</style>
    </div>
  );
};

export default Shtora;