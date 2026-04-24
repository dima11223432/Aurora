import axios from "axios";
import { useEffect, useState, useCallback } from "react";

export default function AdminPanel() {
  const [parsingChannels, setParsingChannels] = useState([]);
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [addedChannel, setAddedChannel] = useState("");
  const [deletedChannel, setDeletedChannel] = useState("");

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
    if (!isLoggedIn) return;

    const fetchParsingChannels = async () => {
      try {
        const token = localStorage.getItem("token");
        const resp = await axios.post(
          "http://localhost:8081/v1/get_all_parsing_channels",
          {},
          { headers: { Authorization: `Bearer ${token}` } },
        );

        setParsingChannels(resp.data.channels || []);
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
        "http://localhost:8081/v1/add_new_parsing_channel",
        { channel_username: addedChannel },
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
        "http://localhost:8081/v1/delete_parsing_channel",
        { channel_username: deletedChannel },
        { headers: { Authorization: `Bearer ${token}` } },
      );
    } catch (e) {
      console.log(e);
    }
  };

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
                parsingChannels.map((channel, index) => (
                  <div
                    key={index}
                    className="bg-white/5 border border-[#0fd2f5]/20 rounded-xl px-4 py-3 text-white flex items-center gap-2 hover:border-[#0fd2f5]/40 transition-colors"
                  >
                    <span className="text-[#0fd2f5]">@</span>
                    {channel}
                  </div>
                ))
              ) : (
                <p className="text-[#95bec7] italic opacity-60">
                  Загрузка каналов...
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
              <input
                id="channel-input"
                className="flex-1 bg-white/5 border border-[#0fd2f5]/20 rounded-2xl px-5 py-3 text-white focus:outline-none focus:border-[#0fd2f5] focus:ring-1 focus:ring-[#0fd2f5] transition-all placeholder:text-gray-600"
                placeholder="durov"
                value={addedChannel}
                onChange={(e) => setAddedChannel(e.target.value)}
              />
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
                {parsingChannels.map((channel, index) => (
                  <option key={index} value={channel} className="bg-[#0F1A2F]">
                    {channel}
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
  );
}
