import React, { useState, useEffect } from "react";
import axios from "axios";
import { routes } from "./config/api";
import Channel from "./Channel";
import BaseChannelCard from "./Channel";

const Shtora = () => {
  const [parsingChannels, setParsingChannels] = useState([]);
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [isOpen, setIsOpen] = useState(false);
  const [selectedChannels, setSelectedChannels] = useState([]);
  const [customChannel, setCustomChannel] = useState("");
  const [userCustomParsingChannels, setUserCustomParsingChannels] = useState(
    [],
  );
  const [isAddingChannel, setIsAddingChannel] = useState(false);

  const getAllUserCustomParsingChannels = async () => {
    const TOKEN = localStorage.getItem("token");
    if (!TOKEN) return;

    try {
      const responce = await axios.get(routes.getAllUserCustomParsingChannels, {
        headers: {
          Authorization: `Bearer ${TOKEN}`,
        },
      });
      const data = await responce.data.channels;
      setUserCustomParsingChannels(data);
      return data;
    } catch (e) {
      console.error(e);
    }
  };

  const addNewUserCustomParsingChannelRequest = async (channelUsername) => {
    const TOKEN = localStorage.getItem("token");
    if (!TOKEN || !channelUsername.trim()) return;

    try {
      const response = await axios.post(
        routes.addNewUserCustomParsingChannel,
        { channel_username: channelUsername.trim() },
        {
          headers: {
            Authorization: `Bearer ${TOKEN}`,
          },
        },
      );

      const data = await response.json();
      console.log("Успешно добавлен канал", data);
      return data;
    } catch (error) {
      if (error.response) {
        const { code, message } = error.response.data;

        if (code === "AlreadyExists" || error.response.status === 409) {
          console.error("Такой канал уже есть!");
        }
      } else {
        console.error("Network error", error.message);
      }
      throw error;
    }
  };

  const handleAddCustomChannel = async () => {
    if (!customChannel.trim()) return;
    setIsAddingChannel(true);
    try {
      await addNewUserCustomParsingChannelRequest(customChannel);
      setCustomChannel("");
      await getUserPriorityChannels();
    } catch (e) {
      console.error(e);
    } finally {
      setIsAddingChannel(false);
    }
  };
  const getUserPriorityChannels = async () => {
    const TOKEN = localStorage.getItem("token");

    try {
      const response = await fetch(routes.getUserPriorityChannels, {
        method: "GET",
        headers: {
          Authorization: `Bearer ${TOKEN}`,
          "Content-Type": "application/json",
        },
      });

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

  const deletePriorityChannelsRequest = async (channels) => {
    const TOKEN = localStorage.getItem("token");
    if (!TOKEN) return;

    try {
      const responce = await axios.post(
        routes.deletePriorityChannels,
        { channels: channels },
        {
          headers: {
            Authorization: `Bearer ${TOKEN}`,
          },
        },
      );
    } catch (e) {
      console.error(e);
    }
  };

  const setPriorityChannelsRequest = async (channels) => {
    const TOKEN = localStorage.getItem("token");
    if (!TOKEN) return;

    try {
      const response = await fetch(routes.setPriorityChannels, {
        method: "POST",
        headers: {
          Authorization: `Bearer ${TOKEN}`,
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          priority_channels: channels,
        }),
      });

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
      await deletePriorityChannelsRequest([channel]);
    } else {
      newSelectedChannels = [...selectedChannels, channel];
      await setPriorityChannelsRequest(newSelectedChannels);
    }
    setSelectedChannels(newSelectedChannels);
  };

  useEffect(() => {
    const login = async () => {
      try {
        const res = await axios.post(routes.login, {
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
        const resp = await axios.get(routes.getAllDefaultParsingChannels, {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });
        const channels = resp.data.channels || [];
        setParsingChannels(channels);
        await getUserPriorityChannels();
      } catch (e) {
        console.log(e);
      }
    };
    fetchParsingChannels();
    getAllUserCustomParsingChannels();
  }, [isLoggedIn]);

  return (
    <div className="fixed top-4 left-4 z-50">
      <button
        onClick={() => setIsOpen((prev) => !prev)}
        className="bg-gradient-to-r from-cyan-300 via-cyan-400 to-cyan-500 hover:from-cyan-300 hover:to-cyan-600 text-white font-bold py-2 px-5 rounded-xl shadow-2xl focus:outline-none transition-all duration-300 hover:scale-105 hover:shadow-cyan-500/40"
        style={{ backdropFilter: "blur(2px)" }}
      >
        {isOpen ? "Скрыть" : "Каналы"}
      </button>
      <div
        className={`mt-2 w-72 max-w-xs ${isOpen ? "shtora-animate-in" : "shtora-animate-out pointer-events-none opacity-0"} bg-gray-900/90 rounded-2xl shadow-2xl border border-cyan-400 p-5 absolute left-0 top-14 transition-all duration-500 backdrop-blur-xl`}
        style={{ boxShadow: "0 8px 32px 0 rgba(31, 38, 135, 0.37)" }}
      >
        <h3 className="text-xl font-bold text-cyan-400 mb-3 drop-shadow">
          Каналы для парсинга
        </h3>
        <ul className="max-h-60 overflow-y-auto space-y-1">
          {parsingChannels.length === 0 ? (
            <li className="text-gray-400 text-sm italic animate-pulse">
              Нет доступных каналов
            </li>
          ) : (
            parsingChannels.map((channel, idx) => {
              const channelName =
                typeof channel === "string"
                  ? channel
                  : channel?.name || JSON.stringify(channel);
              const isChecked = selectedChannels.includes(channelName);
              return (
                <BaseChannelCard
                  idx={idx}
                  isChecked={isChecked}
                  channelName={channelName}
                  handleCheckboxChange={handleCheckboxChange}
                />
              );
            })
          )}
        </ul>
        <div className="mt-3 pt-3 border-t border-cyan-700/50">
          <p className="text-xl font-bold text-cyan-400">Ваши личные каналы:</p>
        </div>
        <ul className="max-h-60 overflow-y-auto space-y-1">
          {userCustomParsingChannels.length > 0 &&
            userCustomParsingChannels.map((channel, idx) => {
              const channelName =
                typeof channel === "string"
                  ? channel
                  : channel?.name || JSON.stringify(channel);
              const isChecked = selectedChannels.includes(channelName);
              return (
                <BaseChannelCard
                  idx={idx}
                  isChecked={isChecked}
                  channelName={channelName}
                  handleCheckboxChange={handleCheckboxChange}
                />
              );
            })}
        </ul>
        <div className="mt-3 pt-3 border-t border-cyan-700/50">
          <div className="flex gap-2">
            <input
              type="text"
              value={customChannel}
              onChange={(e) => setCustomChannel(e.target.value)}
              onKeyDown={(e) => e.key === "Enter" && handleAddCustomChannel()}
              placeholder="Новый канал (username)"
              className="flex-1 bg-gray-800/80 text-white text-sm px-3 py-2 rounded-lg border border-cyan-700/50 focus:border-cyan-400 focus:outline-none placeholder-gray-500"
            />
            <button
              onClick={handleAddCustomChannel}
              disabled={isAddingChannel || !customChannel.trim()}
              className="bg-cyan-600 hover:bg-cyan-500 disabled:bg-gray-600 disabled:cursor-not-allowed text-white text-sm px-3 py-2 rounded-lg transition-colors duration-200"
            >
              {isAddingChannel ? "..." : "➕"}
            </button>
          </div>
        </div>
        {/* {selectedChannels.length > 0 && (
          <div className="mt-3 pt-2 border-t border-cyan-700/50 text-xs text-cyan-300">
    
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
