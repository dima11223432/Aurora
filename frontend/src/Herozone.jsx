import React, { useEffect, useState, useRef } from "react";
import { useNavigate } from "react-router-dom";
import { API_BASE_URL_SECURE, routes } from "./config/api";

function Herozone() {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState(null);
  const [telegramUser, setTelegramUser] = useState(null);
  const widgetContainerRef = useRef(null);

  const navigate = useNavigate();
  const API_URL = routes.loginSecure;

  useEffect(() => {
    const tg = window.Telegram?.WebApp;
    if (!tg) {
      setTimeout(() => {
        setError("не удалось загруить данные из telegram");
      }, 0);
      return;
    }

    tg.ready();
    tg.expand();

    const user = tg.initDataUnsafe?.user;

    if (!user || !user.id) {
      console.error("Данные пользователя не найдены");
      setTimeout(() => {
        setError("Не удалось получить данные пользователя");

        setIsLoading(false);
      }, 0);
      return;
    }

    console.log("Telegram user data:", user);
    setTimeout(() => {
      setIsLoading(true);
    }, 0);

    fetch(API_URL, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        telegram_id: user.id,
        username: user.username || " ",
        first_name: user.first_name || " ",
        last_name: user.last_name || " ",
        is_admin: false,
        app_id: 1,
      }),
    })
      .then((res) => res.json())
      .then((data) => {
        console.log("JWT Token получен");
        localStorage.setItem("jwt", data.token);
      })
      .catch((err) => {
        console.error("Login API error:", err);
        setError("Ошибка авторизации на сервере");
      })
      .finally(() => {
        setIsLoading(false);
      });
  }, []);

  const items = [
    {
      title: "ИИ аналитика",
      description:
        "Продвинутые модели машинного обучения анализируют рыночные тенденции",
    },
    {
      title: "Анализ в реальном времени",
      description: "Текущие рыночные данные и мгновенное обновление прогнозов",
    },
    {
      title: "Высокая точность",
      description:
        "*% точности прогнозов на основе исторических данных и текущих рыночных условий",
    },
    {
      title: "Безопасный и приватный",
      description: "Ваши данные зашифрованы и защищены",
    },
  ];

  return (
    <div className="min-h-screen bg-gradient-to-br from-[#0A0F1F] via-[#0F1A2F] to-[#02B7DB] flex items-center justify-center p-4 sm:p-6">
      <div className="max-w-3xl w-full bg-[rgba(20,25,50,0.7)] backdrop-blur-md rounded-[3rem] p-6 sm:p-8 md:p-12 shadow-2xl border border-white/5">
        <img
          src="./assets/Aurora-logo.png"
          alt="Aurora logo"
          className="w-20 h-20 mb-2 mx-auto"
        />

        <p className="text-[#95bec7] text-base text-center sm:text-lg mb-8 sm:mb-12 pl-4 to-transparent">
          Добро пожаловать в{" "}
          <span className="text-[#0fd2f5] font-bold">Aurora</span>
          <br /> ИИ-поддерживаемые прогнозы TON и T-Investments
        </p>

        <div className="space-y-4 sm:space-y-5 mb-8 sm:mb-12">
          {items.map((item) => (
            <div
              key={item.id}
              className="flex gap-3 sm:gap-4 items-start group"
            >
              <div className="flex-1 bg-white/5 backdrop-blur-sm p-5 sm:p-6 rounded-xl border border-[#0fd2f5]/20 hover:border-[#0fd2f5]/50 transition-all duration-300 hover:shadow-lg hover:shadow-[#0fd2f5]/10">
                <div className="flex gap-3 sm:gap-4 items-center">
                  <div className="w-10 h-10 sm:w-10 sm:h-10 rounded bg-[#0fd2f5]/20 border border-[#0fd2f5]/50 flex-shrink-0"></div>
                  <div className="flex-1">
                    <h3 className="text-left text-white text-lg sm:text-xl md:text-2xl font-semibold mb-1">
                      {item.title}
                    </h3>
                    <p className="text-left text-[#95bec7] text-sm sm:text-base opacity-90 leading-relaxed">
                      {item.description}
                    </p>
                  </div>
                </div>
              </div>
            </div>
          ))}
        </div>
        <div>
          <button
            onClick={() => navigate("/feed")}
            className="w-full flex justify-center items-center gap-2 bg-[#0fd2f5] text-[#0A0F1F] font-bold text-lg py-4 px-8 rounded-full shadow-lg shadow-[#0fd2f5]/20 hover:bg-white hover:shadow-[#0fd2f5]/40 active:scale-95 transition-all duration-300 transform"
          >
            Начать
          </button>
        </div>
      </div>
    </div>
  );
}

export default Herozone;
