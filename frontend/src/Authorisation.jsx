import React, { useEffect, useState, useRef } from "react";
import "./styles/Landing.css";
import Logo from "./assets/Aurora.png";

export const Landing = () => {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState(null);
  const [telegramUser, setTelegramUser] = useState(null);
  const widgetContainerRef = useRef(null);

  useEffect(() => {
    window.onTelegramAuth = async (user) => {
      console.log("Telegram auth success:", user);
      setIsLoading(true);

      try {
        fetch("https://24c2-2605-e440-9-00-3a.ngrok-free.app/v1/login", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            telegram_id: 12345,
            username: "",
            first_name: "",
            last_name: "",
            is_admin: false,
            app_id: 1,
          }),
        })
          .then((res) => res.json())
          .then((data) => {
            console.log("JWT Token:", data.token);
            localStorage.setItem("jwt", data.token);
          })
          .catch(console.error);
      } catch (err) {
        console.error("Login API error:", err);
        setError("Ошибка авторизации на сервере");
      } finally {
        setIsLoading(false);
      }
    };

    const script = document.createElement("script");
    script.src = "https://telegram.org/js/telegram-widget.js?22";
    script.setAttribute("data-telegram-login", "AuroraFinances_bot");
    script.setAttribute("data-size", "large");
    script.setAttribute("data-onauth", "onTelegramAuth(user)");
    script.setAttribute("data-request-access", "write");
    script.async = true;

    script.onload = () => {
      console.log("Telegram widget loaded successfully");
      setIsLoading(false);
    };

    script.onerror = () => {
      console.error("Failed to load Telegram widget");
      setError("Не удалось загрузить виджет Telegram");
      setIsLoading(false);
    };

    if (widgetContainerRef.current) {
      widgetContainerRef.current.innerHTML = "";
      widgetContainerRef.current.appendChild(script);
    }

    const savedUser = localStorage.getItem("telegramUser");
    if (savedUser) {
      setTelegramUser(JSON.parse(savedUser));
    }

    return () => {
      if (widgetContainerRef.current) {
        widgetContainerRef.current.innerHTML = "";
      }
      delete window.onTelegramAuth;
    };
  }, []);

  const handleLogout = () => {
    localStorage.removeItem("telegramUser");
    setTelegramUser(null);

    if (widgetContainerRef.current) {
      widgetContainerRef.current.innerHTML = "";

      const newScript = document.createElement("script");
      newScript.src = "https://telegram.org/js/telegram-widget.js?22";
      newScript.setAttribute("data-telegram-login", "AuroraFinances_bot");
      newScript.setAttribute("data-size", "large");
      newScript.setAttribute("data-onauth", "onTelegramAuth(user)");
      newScript.setAttribute("data-request-access", "write");
      newScript.async = true;

      widgetContainerRef.current.appendChild(newScript);
    }
  };

  return (
    <div className="landing-container">
      <div className="gradient-bg"></div>
      <main className="landing-main">
        <header className="landing-header">
          <div className="aurora-logo">
            <img src={Logo} alt="Aurora Logo" className="logo-image" />
          </div>
          <h1 className="aurora-name">Aurora</h1>
        </header>

        <section className="landing-hero">
          <h2 className="hero-title">
            Добро пожаловать в <span className="highlight">Aurora</span>
          </h2>
          <p className="hero-description">
            AI-сделанные TON & T-Инвестиционные прогнозы
          </p>
        </section>

        <section className="features-grid">
          <div className="feature-card">
            <div className="feature-icon">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor">
                <path
                  d="M12 2L15 10H23L17 15L19 23L12 18L5 23L7 15L1 10H9L12 2Z"
                  strokeWidth="1.5"
                />
              </svg>
            </div>
            <h3 className="feature-title">AI-сделанные Прогнозы</h3>
            <p className="feature-description">
              Продвинутые модели машинного обучения анализируют тренды рынка
            </p>
          </div>

          <div className="feature-card">
            <div className="feature-icon">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor">
                <path
                  d="M3 12L5 10L9 14L15 3L17 5"
                  strokeWidth="1.5"
                  strokeLinecap="round"
                />
                <path
                  d="M21 6L21 18M6 21H18"
                  strokeWidth="1.5"
                  strokeLinecap="round"
                />
              </svg>
            </div>
            <h3 className="feature-title">Анализ в реальном времени</h3>
            <p className="feature-description">
              Данные рынка в реальном времени и мгновенные обновления прогнозов
            </p>
          </div>

          <div className="feature-card">
            <div className="feature-icon">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor">
                <circle cx="12" cy="12" r="9" strokeWidth="1.5" />
                <circle cx="12" cy="12" r="5" strokeWidth="1.5" />
                <circle cx="12" cy="12" r="1.5" fill="currentColor" />
              </svg>
            </div>
            <h3 className="feature-title">Высокая точность</h3>
            <p className="feature-description">
              Уровень точности прогнозирования до 94%
            </p>
          </div>

          <div className="feature-card">
            <div className="feature-icon">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor">
                <path
                  d="M12 2L2 6V12C2 18 12 22 12 22C12 22 22 18 22 12V6L12 2Z"
                  strokeWidth="1.5"
                />
              </svg>
            </div>
            <h3 className="feature-title">Безопасность & Конфиденциальность</h3>
            <p className="feature-description">
              Ваши данные зашифрованы и защищены
            </p>
          </div>
        </section>

        <div className="telegram-auth-section">
          {!telegramUser ? (
            <>
              {isLoading ? (
                <div className="loading-widget">
                  <span className="loading-spinner"></span>
                </div>
              ) : (
                <div
                  id="telegram-login-container"
                  ref={widgetContainerRef}
                  className="telegram-widget-wrapper"
                />
              )}
            </>
          ) : (
            <div className="user-profile">
              <div className="profile-header">
                {telegramUser.photo_url ? (
                  <img
                    src={telegramUser.photo_url}
                    alt="Profile"
                    className="profile-photo"
                  />
                ) : (
                  <div className="profile-avatar">
                    {telegramUser.first_name.charAt(0)}
                  </div>
                )}
                <div className="profile-info">
                  <h3>
                    {telegramUser.first_name} {telegramUser.last_name}
                  </h3>
                  {telegramUser.username && (
                    <p className="profile-username">@{telegramUser.username}</p>
                  )}
                </div>
              </div>
              <button className="logout-btn" onClick={handleLogout}>
                Выйти
              </button>
            </div>
          )}
        </div>

        {error && (
          <div className="error-message">
            <span>⚠️</span>
            <p>{error}</p>
          </div>
        )}

        <footer className="landing-footer">
          <p className="footer-main">Начните прогнозировать сейчас</p>
          <p className="footer-secondary">
            На базе продвинутого ИИ • Не является финансовым советом
          </p>

          {/* Виджет Telegram слева снизу */}
        </footer>
      </main>
    </div>
  );
};
