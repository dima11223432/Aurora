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
    <div className="p-4">
      <h1 className="text-xl font-bold mb-4">Telegram каналы для парса:</h1>

      <div>
        {parsingChannels.length > 0 ? (
          parsingChannels.map((channel, index) => (
            <div key={index}>{channel}</div>
          ))
        ) : (
          <p className="text-gray-500 italic">Загрузка каналов...</p>
        )}
      </div>

      <div className="flex flex-col gap-2 max-w-xs">
        <label htmlFor="channel-input">Добавить новый канал:</label>
        <input
          id="channel-input"
          className="border border-black p-1"
          placeholder="durov"
          value={addedChannel}
          onChange={(e) => setAddedChannel(e.target.value)}
        />
        <button
          onClick={() => {
            AddNewParsingChannel();
          }}
        >
          Добавить новый канал
        </button>
      </div>
      <div>
        <p>Удалить telegram канал:</p>
        <select
          onChange={(e) => {
            setDeletedChannel(e.target.value);
          }}
        >
          <option value="">--Выберите канал--</option>
          {parsingChannels.map((channel, index) => (
            <option key={index} value={channel}>
              {channel}
            </option>
          ))}
        </select>
        <button
          onClick={() => {
            DeleteParsingChannel();
          }}
        >
          Удалить
        </button>
      </div>
    </div>
  );
}
