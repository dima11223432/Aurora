import axios from "axios";
import { useEffect, useState } from "react";
import { routes } from "./config/api";
import Footer from "./Footer";
import UserCustomChannelCard from "./UserCustomChannelCard";
import BaseChannelCard from "./BaseChannelCard";
import { useNavigate } from "react-router-dom";

export default function Settings() {
  const navigate = useNavigate();
  const [addedChannel, setAddedChannel] = useState("");
  const [parsingChannels, setParsingChannels] = useState({});
  const [selectedChannels, setSelectedChannels] = useState([]);
  const [userCustomParsingChannels, setUserCustomParsingChannels] = useState(
    [],
  );

  const [isLoading, setIsLoading] = useState(true);

  const getAllUserCustomParsingChannels = async () => {
    const TOKEN = localStorage.getItem("token");
    if (!TOKEN) return;

    try {
      const response = await axios.get(routes.getAllUserCustomParsingChannels, {
        headers: { Authorization: `Bearer ${TOKEN}` },
      });
      const data = response.data.channels || [];
      setUserCustomParsingChannels(data);
      return data;
    } catch (e) {
      console.error("Ошибка получения кастомных каналов:", e);
    }
  };

  const getUserPriorityChannels = async () => {
    const TOKEN = localStorage.getItem("token");
    if (!TOKEN) return;

    try {
      const response = await fetch(routes.getUserPriorityChannels, {
        method: "GET",
        headers: {
          Authorization: `Bearer ${TOKEN}`,
          "Content-Type": "application/json",
        },
      });
      const data = await response.json();
      const userPriorityChannels = data.channels || [];
      setSelectedChannels(userPriorityChannels);
      return data;
    } catch (error) {
      console.error("Error priority channels: ", error);
      return null;
    }
  };

  const deleteUserCustomParsingChannelRequest = async (channelUsername) => {
    const TOKEN = localStorage.getItem("token");
    if (!TOKEN) return;

    try {
      const resp = await axios.post(
        routes.deleteUserCustomParsingChannel,
        { channel_username: channelUsername },
        { headers: { Authorization: `Bearer ${TOKEN}` } },
      );

      const data = resp.data;
      console.log("Успешно удален канал", data);

      await getAllUserCustomParsingChannels();
      return data;
    } catch (e) {
      console.log("Ошибка при удалении кастомного канала:", e);
    }
  };

  const AddNewParsingChannel = async () => {
    const TOKEN = localStorage.getItem("token");
    if (!TOKEN || !addedChannel.trim()) return;

    try {
      await axios.post(
        routes.addNewUserCustomParsingChannel,
        { channel_username: addedChannel.trim() },
        { headers: { Authorization: `Bearer ${TOKEN}` } },
      );
      setAddedChannel("");
      await getAllUserCustomParsingChannels();
    } catch (e) {
      console.error("Ошибка при добавлении канала:", e);
    }
  };

  const deletePriorityChannelsRequest = async (channels) => {
    const TOKEN = localStorage.getItem("token");
    if (!TOKEN) return;
    try {
      await axios.post(
        routes.deletePriorityChannels,
        { channels },
        { headers: { Authorization: `Bearer ${TOKEN}` } },
      );
    } catch (e) {
      console.error(e);
    }
  };

  const setPriorityChannelsRequest = async (channels) => {
    const TOKEN = localStorage.getItem("token");
    if (!TOKEN) return;
    try {
      await fetch(routes.setPriorityChannels, {
        method: "POST",
        headers: {
          Authorization: `Bearer ${TOKEN}`,
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ priority_channels: channels }),
      });
    } catch (error) {
      console.error(error);
    }
  };

  const handleCheckboxChange = async (channel) => {
    let newSelectedChannels;
    if (selectedChannels.includes(channel)) {
      newSelectedChannels = selectedChannels.filter((ch) => ch !== channel);
      setSelectedChannels(newSelectedChannels);
      await deletePriorityChannelsRequest([channel]);
    } else {
      newSelectedChannels = [...selectedChannels, channel];
      setSelectedChannels(newSelectedChannels);
      await setPriorityChannelsRequest(newSelectedChannels);
    }
  };

  useEffect(() => {
    const token = localStorage.getItem("token");

    if (!token) {
      navigate("/404");
      return;
    }

    const fetchAllData = async () => {
      try {
        setIsLoading(true);
        const resp = await axios.get(
          routes.getAllDefaultParsingChannelsWithCategories,
          { headers: { Authorization: `Bearer ${token}` } },
        );
        setParsingChannels(resp.data.channels || {});

        await Promise.all([
          getUserPriorityChannels(),
          getAllUserCustomParsingChannels(),
        ]);
      } catch (e) {
        console.error("Ошибка инициализации данных:", e);
      } finally {
        setIsLoading(false);
      }
    };

    fetchAllData();
  }, []);

  const hasParsingChannels = Object.keys(parsingChannels).length > 0;

  if (isLoading) {
    return (
      <div className="min-h-screen bg-[#0A0F1F] flex items-center justify-center">
        <p className="text-cyan-400 text-lg animate-pulse">
          Загрузка настроек...
        </p>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-[#0A0F1F] via-[#0F1A2F] to-[#02B7DB] flex items-center justify-center p-4 sm:p-6 pb-24 sm:pb-32 font-sans">
      <div className="max-w-2xl w-full bg-[rgba(20,25,50,0.7)] backdrop-blur-md rounded-[2.5rem] p-6 sm:p-8 md:p-10 shadow-2xl border border-white/5">
        <div className="space-y-8">
          <div>
            <h3 className="text-xl font-bold text-cyan-400 mb-3 drop-shadow">
              Рекомендуем:
            </h3>
            <div
              className="max-h-60 overflow-y-auto space-y-4 pr-2
                            [&::-webkit-scrollbar]:w-1.5
                            [&::-webkit-scrollbar-thumb]:bg-[#0fd2f5]/30 
                            [&::-webkit-scrollbar-thumb]:rounded-full
                            [&::-webkit-scrollbar-track]:bg-transparent"
            >
              {!hasParsingChannels ? (
                <p className="text-gray-400 text-sm italic animate-pulse">
                  Нет доступных каналов
                </p>
              ) : (
                Object.entries(parsingChannels).map(
                  ([category, categoryData]) => (
                    <div key={category} className="mb-3">
                      <h4 className="text-xs font-bold text-cyan-500/70 uppercase px-1 mb-2 tracking-wider">
                        {category}
                      </h4>
                      <ul className="space-y-2">
                        {categoryData.usernames?.map((channelName, idx) => {
                          const isChecked =
                            selectedChannels.includes(channelName);
                          return (
                            <BaseChannelCard
                              key={channelName}
                              idx={idx}
                              isChecked={isChecked}
                              channelName={channelName}
                              handleCheckboxChange={handleCheckboxChange}
                            />
                          );
                        })}
                      </ul>
                    </div>
                  ),
                )
              )}
            </div>
          </div>

          <hr className="border-white/5" />

          <section>
            <h3 className="text-xl font-bold text-cyan-400 mb-3 drop-shadow">
              Ваши каналы:
            </h3>

            {userCustomParsingChannels &&
            userCustomParsingChannels.length > 0 ? (
              <ul
                className="grid grid-cols-1 gap-3 max-h-52 overflow-y-auto pr-2 
                           scrollbar-thin scrollbar-thumb-[#0fd2f5]/20 scrollbar-track-transparent 
                           [&::-webkit-scrollbar]:w-1.5
                           [&::-webkit-scrollbar-thumb]:bg-[#0fd2f5]/30 
                           [&::-webkit-scrollbar-thumb]:rounded-full
                           [&::-webkit-scrollbar-track]:bg-transparent"
              >
                {userCustomParsingChannels.map((channel, idx) => {
                  const channelName =
                    typeof channel === "string"
                      ? channel
                      : channel?.channel_username ||
                        channel?.name ||
                        JSON.stringify(channel);

                  const isChecked = selectedChannels.includes(channelName);

                  return (
                    <UserCustomChannelCard
                      key={channelName || idx}
                      idx={idx}
                      isChecked={isChecked}
                      channelName={channelName}
                      handleCheckboxChange={handleCheckboxChange}
                      deleteUserCustomParsingChannel={
                        deleteUserCustomParsingChannelRequest
                      }
                    />
                  );
                })}
              </ul>
            ) : (
              <div className="bg-white/5 border border-dashed border-white/10 rounded-2xl px-5 py-4 text-center">
                <p className="text-gray-500 text-sm italic">
                  Каналы не найдены...
                </p>
              </div>
            )}
          </section>

          <hr className="border-white/5" />

          <section className="flex flex-col gap-4">
            <label
              htmlFor="channel-input"
              className="text-[#95bec7] text-sm ml-1"
            >
              Добавить новый канал
            </label>
            <div className="flex flex-col sm:flex-row gap-3">
              <input
                id="channel-input"
                className="flex-1 bg-white/5 border border-[#0fd2f5]/20 rounded-2xl px-5 py-3 text-white focus:outline-none focus:border-[#0fd2f5] focus:ring-1 focus:ring-[#0fd2f5] transition-all placeholder:text-gray-600"
                placeholder="durov"
                value={addedChannel}
                onChange={(e) => setAddedChannel(e.target.value)}
              />
              <button
                onClick={AddNewParsingChannel}
                className="bg-[#0fd2f5] text-[#0A0F1F] font-bold py-3 px-6 rounded-2xl hover:bg-white active:scale-95 transition-all shadow-lg shadow-[#0fd2f5]/20"
              >
                Добавить
              </button>
            </div>
          </section>
        </div>
      </div>
      <Footer />
    </div>
  );
}
