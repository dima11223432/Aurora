import axios from "axios";
import { useEffect, useState, useCallback } from "react";
import { routes } from "./config/api";

export default function AdminPanel() {
  const [parsingChannels, setParsingChannels] = useState([]);
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [addedChannel, setAddedChannel] = useState("");
  const [addedChannelCategory, setAddedChannelCategory] = useState("");
  const [deletedChannel, setDeletedChannel] = useState("");
  const [categories, setCategories] = useState([]);
  сonst [isShtoraOpen, setIsShtoraOpen] = useState(false);

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
    if (!isLoggedIn) return;

    const fetchParsingChannels = async () => {
      try {
        const token = localStorage.getItem("token");
        const resp = await axios.get(
          routes.getAllDefaultParsingChannelsWithCategories,
          { headers: { Authorization: `Bearer ${token}` } },
        );

        const formatted = Object.entries(resp.data.channels).map(
          ([cat, data]) => ({
            category: cat,
            usernames: data.usernames || [],
          }),
        );

        setParsingChannels(formatted);
      } catch (e) {
        console.error("Fetch error:", e);
      }
    };

    fetchParsingChannels();
  }, [isLoggedIn]);

  const AddNewParsingChannel = () => {
    try {
      const token = localStorage.getItem("token");
      axios.post(
        routes.addNewDefaultParsingChannel,
        { channel_username: addedChannel, category: addedChannelCategory },
        { headers: { Authorization: `Bearer ${token}` } },
      );
    } catch (e) {
      console.log(e);
    }
  };

  const DeleteParsingChannel = () => {
    try {
      const token = localStorage.getItem("token");
      axios.post(
        routes.deleteDefaultParsingChannel,
        { channel_username: deletedChannel },
        { headers: { Authorization: `Bearer ${token}` } },
      );
    } catch (e) {
      console.log(e);
    }
  };

  return (
    <>
    <div className="min-h-screen bg-gradient-to-br from-[#0A0F1F] via-[#0F1A2F] to-[#02B7DB] flex items-center justify-center p-4 sm:p-6 font-sans">
      <div className="max-w-2xl w-full bg-[rgba(20,25,50,0.7)] backdrop-blur-md rounded-[2.5rem] p-6 sm:p-8 md:p-10 shadow-2xl border border-white/5">
        <h1 className="text-2xl sm:text-3xl font-bold text-white mb-8 text-center">
          Управление <span className="text-[#0fd2f5]">Парсингом</span>
        </h1>

        <div className="space-y-8">
          <section>
            <h2 className="text-[#95bec7] text-sm uppercase tracking-widest mb-4 ml-1">
              Текущие каналы:
            </h2>
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
              {parsingChannels.length > 0 ? (
                parsingChannels.map((item) => (
                  <div
                    key={item.category}
                    className="bg-white/5 p-4 rounded-2xl border border-white/10"
                  >
                    <p className="text-[#0fd2f5] font-bold text-sm mb-2 uppercase">
                      {item.category}
                    </p>
                    <div className="flex flex-wrap gap-2">
                      {item.usernames.map((user) => (
                        <span
                          key={user}
                          className="bg-[#0fd2f5]/10 text-white text-xs px-3 py-1 rounded-full border border-[#0fd2f5]/20"
                        >
                          @{user}
                        </span>
                      ))}
                    </div>
                  </div>
                ))
              ) : (
                <p className="text-gray-500 text-sm italic">
                  Каналы не найдены...
                </p>
              )}
            </div>
          </section>
          <hr className="border-white/5" />

          <section className="flex flex-col gap-4">
            <label
              htmlFor="channel-input"
              className="text-[#95bec7] text-sm ml-1"
            >
              Добавить новый канал
            </label>
            <div className="grid grid-cls-2 sm:flex-row gap-3">
              <div className="flex justify-between">
                <input
                  id="channel-input"
                  className="flex-1 bg-white/5 border border-[#0fd2f5]/20 rounded-2xl px-5 py-3 text-white focus:outline-none focus:border-[#0fd2f5] focus:ring-1 focus:ring-[#0fd2f5] transition-all placeholder:text-gray-600"
                  placeholder="durov"
                  value={addedChannel}
                  onChange={(e) => setAddedChannel(e.target.value)}
                />
                <select
                  className="flex-1 bg-white/5 border border-[#0fd2f5]/20 rounded-2xl px-5 py-3 text-white focus:outline-none focus:border-red-400/50 appearance-none cursor-pointer"
                  onChange={(e) => setAddedChannelCategory(e.target.value)}
                  value={addedChannelCategory}
                >
                  <option value="" className="bg-[#0F1A2F]">
                    -- Выберите категорию --
                  </option>

                  {parsingChannels.length > 0 &&
                    parsingChannels.map((item) => (
                      <option
                        key={item.category}
                        value={item.category}
                        className="bg-[#0F1A2F]"
                      >
                        {item.category}
                      </option>
                    ))}
                </select>
              </div>
              <button
                onClick={() => AddNewParsingChannel()}
                className="bg-[#0fd2f5] text-[#0A0F1F] font-bold py-3 px-6 rounded-2xl hover:bg-white active:scale-95 transition-all shadow-lg shadow-[#0fd2f5]/20"
              >
                Добавить
              </button>
            </div>
          </section>

          <section className="flex flex-col gap-4">
            <label className="text-[#95bec7] text-sm ml-1">
              Удалить telegram канал
            </label>
            <div className="flex flex-col sm:flex-row gap-3">
              <select
                className="flex-1 bg-white/5 border border-[#0fd2f5]/20 rounded-2xl px-5 py-3 text-white focus:outline-none focus:border-red-400/50 appearance-none cursor-pointer"
                onChange={(e) => setDeletedChannel(e.target.value)}
                value={deletedChannel}
              >
                <option value="" className="bg-[#0F1A2F]">
                  -- Выберите канал --
                </option>
                {parsingChannels
                  .flatMap((item) => item.usernames)
                  .map((username, index) => (
                    <option
                      key={index}
                      value={username}
                      className="bg-[#0F1A2F]"
                    >
                      {username}
                    </option>
                  ))}
              </select>
              <button
                onClick={() => DeleteParsingChannel()}
                className="border border-red-500/50 text-red-400 font-semibold py-3 px-6 rounded-2xl hover:bg-red-500/10 active:scale-95 transition-all"
              >
                Удалить
              </button>
            </div>
          </section>
        </div>
      </div>
    </div>
    <button
        onClick={() => setIsShtoraOpen(!isShtoraOpen)}
        className="fixed top-4 right-4 z-50 text-white font-bold py-2 px-5 rounded-xl"
        style={{ background: "linear-gradient(to right, #208390, #36DEF4)" }}
      >
        {isShtoraOpen ? "Скрыть ТГК" : "Показать ТГК"}
      </button>

      {isShtoraOpen && (
        <div className="fixed top-20 right-4 z-40 w-80 bg-gray-900/90 rounded-2xl shadow-2xl border border-blue-700 p-5 backdrop-blur-xl">
          <h3 className="text-xl font-bold text-blue-400 mb-3">
            Каналы для парсинга
          </h3>
          <ul className="max-h-96 overflow-y-auto space-y-2">
            {parsingChannels.length === 0 ? (
              <li className="text-gray-400 text-sm italic">Нет доступных каналов</li>
            ) : (
              parsingChannels.flatMap((item) => item.usernames).map((username, idx) => (
                <li key={idx} className="py-2 px-3 rounded-lg hover:bg-blue-700/70 text-white flex items-center gap-2">
                  <input type="checkbox" className="w-4 h-4 rounded" />
                  <span>@{username}</span>
                </li>
              ))
            )}
          </ul>
        </div>
      )}
    </>
  );
}
