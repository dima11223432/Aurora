import axios from "axios";
import { useEffect, useState } from "react";
import { routes } from "./config/api";
import { useNavigate } from "react-router-dom";

export default function AdminPanel() {
  const [parsingChannels, setParsingChannels] = useState([]);
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [addedChannel, setAddedChannel] = useState("");
  const [addedChannelCategory, setAddedChannelCategory] = useState("");
  const [deletedChannel, setDeletedChannel] = useState("");
  const navigate = useNavigate();

  useEffect(() => {
    const checkIsAdmin = async () => {
      try {
        const token = localStorage.getItem("token");
        if (!token) {
          navigate("/404");
          return;
        }

        const resp = await axios.get(routes.isAdmin, {
          headers: { Authorization: `Bearer ${token}` },
        });

        if (resp.data.isAdmin === true || resp.data.is_admin === true) {
          setIsLoggedIn(true);
        } else {
          navigate("/404");
        }
      } catch (e) {
        console.error("Ошибка проверки прав:", e);
        navigate("/404");
      }
    };

    checkIsAdmin();
  }, [navigate]);

  useEffect(() => {
    const fetchParsingChannels = async () => {
      if (!isLoggedIn) return;

      try {
        const token = localStorage.getItem("token");
        const resp = await axios.get(
          routes.getAllDefaultParsingChannelsWithCategories,
          { headers: { Authorization: `Bearer ${token}` } },
        );

        const formatted = Object.entries(resp.data.channels || {}).map(
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

  const AddNewParsingChannel = async () => {
    try {
      const token = localStorage.getItem("token");
      if (!token) {
        navigate("/404");
        return;
      }

      await axios.post(
        routes.addNewDefaultParsingChannel,
        { channel_username: addedChannel, category: addedChannelCategory },
        { headers: { Authorization: `Bearer ${token}` } },
      );

      setAddedChannel("");
      setAddedChannelCategory("");
    } catch (e) {
      console.error(e);
    }
  };

  const DeleteParsingChannel = async () => {
    try {
      const token = localStorage.getItem("token");
      if (!token) {
        navigate("/404");
        return;
      }

      await axios.post(
        routes.deleteDefaultParsingChannel,
        { channel_username: deletedChannel },
        { headers: { Authorization: `Bearer ${token}` } },
      );

      setDeletedChannel("");
    } catch (e) {
      console.error(e);
    }
  };

  if (!isLoggedIn) {
    return null;
  }

  return (
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
            <div className="flex flex-col sm:flex-row gap-3">
              <div className="flex justify-between flex-1 gap-2">
                <input
                  id="channel-input"
                  className="flex-1 bg-white/5 border border-[#0fd2f5]/20 rounded-2xl px-5 py-3 text-white focus:outline-none focus:border-[#0fd2f5] focus:ring-1 focus:ring-[#0fd2f5] transition-all placeholder:text-gray-600 w-full"
                  placeholder="durov"
                  value={addedChannel}
                  onChange={(e) => setAddedChannel(e.target.value)}
                />
                <select
                  className="flex-1 bg-white/5 border border-[#0fd2f5]/20 rounded-2xl px-5 py-3 text-white focus:outline-none focus:border-red-400/50 appearance-none cursor-pointer w-full"
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
                onClick={AddNewParsingChannel}
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
                onClick={DeleteParsingChannel}
                className="border border-red-500/50 text-red-400 font-semibold py-3 px-6 rounded-2xl hover:bg-red-500/10 active:scale-95 transition-all"
              >
                Удалить
              </button>
            </div>
          </section>
        </div>
      </div>
    </div>
  );
}
