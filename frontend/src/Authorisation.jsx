import React, { useEffect, useState } from "react";
import "./styles/Landing.css";
import Logo from "./assets/Aurora.png";

export const Landing = () => {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState(null);
  const [telegramUser, setTelegramUser] = useState(null);

  useEffect(() => {
    window.onTelegramAuth = (user) => {
      console.log("Telegram auth success:", user);
      setTelegramUser(user);
      setIsLoading(false);

      localStorage.setItem("telegramUser", JSON.stringify(user));
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
    };

    script.onerror = () => {
      console.error("Failed to load Telegram widget");
      setError("Не удалось загрузить виджет Telegram");
    };

    document.body.appendChild(script);

    const savedUser = localStorage.getItem("telegramUser");
    if (savedUser) {
      setTelegramUser(JSON.parse(savedUser));
    }

    return () => {
      if (document.body.contains(script)) {
        document.body.removeChild(script);
      }
      delete window.onTelegramAuth;
    };
  }, []);

  const handleLogout = () => {
    localStorage.removeItem("telegramUser");
    setTelegramUser(null);
  };

  const openTelegramBot = () => {
    const botUsername = "AuroraFinances_bot";
    window.open(`https://t.me/${botUsername}`, "_blank");
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

        <div className="cta-section">
          {!telegramUser ? (
            <>
              <div
                id="telegram-login-container"
                style={{
                  display: "flex",
                  justifyContent: "center",
                  margin: "20px 0",
                  minHeight: "58px",
                }}
              />

              <button
                className="get-started-btn"
                onClick={openTelegramBot}
                disabled={isLoading}
                style={{ marginTop: "10px" }}
              >
                <svg
                  className="telegram-icon"
                  viewBox="0 0 24 24"
                  fill="currentColor"
                >
                  <path d="M12 0C5.373 0 0 5.373 0 12s5.373 12 12 12 12-5.373 12-12S18.627 0 12 0zm5.894 8.221l-1.97 9.28c-.145.658-.537.818-1.084.508l-3-2.21-1.446 1.394c-.14.146-.357.292-.611.292-.005 0-.01 0-.016 0l.213-3.053 5.56-5.023c.242-.213-.054-.328-.375-.115l-6.869 4.332-2.961-.924c-.643-.204-.657-.643.136-.953l11.566-4.458c.529-.196 1.083.128.897.983z" />
                </svg>
                Открыть бота
              </button>
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

          {error && (
            <div className="error-message">
              <span>⚠️</span>
              <p>{error}</p>
            </div>
          )}
        </div>

        <footer className="landing-footer">
          <p className="footer-main">Начните прогнозировать сейчас</p>
          <p className="footer-secondary">
            На базе продвинутого ИИ • Не является финансовым советом
          </p>
        </footer>
      </main>
    </div>
  );
};
