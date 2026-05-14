import { useState, useEffect } from "react";
import axios from "axios";

export default function PriorityChannels() {
  const [parsingChannels, setParsingChannels] = useState([]);
  const [isLoggedIn, setIsLoggedIn] = useState(false);

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
          }
        );
        setParsingChannels(resp.data.channels || []);
      } catch (e) {
        console.log(e);
      }
    };
    fetchParsingChannels();
  }, [isLoggedIn]);

  const setPriorityChannels = async () => {
    const selected = Array.from(
      document.querySelectorAll('input[type="checkbox"]:checked')
    ).map((checkbox) => checkbox.value);
    console.log(selected);
  };

  return (
    <ul>
      {parsingChannels.length === 0 ? (
        <li>Нет доступных каналов</li>
      ) : (
        parsingChannels.map((channel, idx) => (
          <li key={idx}>
            <input
              type="checkbox"
              onClick={setPriorityChannels}
            />
          </li>
        ))
      )}
    </ul>
  );
}