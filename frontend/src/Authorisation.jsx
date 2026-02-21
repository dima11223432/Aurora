import React, { useEffect, useState } from 'react';
import './styles/Landing.css';
import Logo from './assets/Aurora.png';

export const Landing = () => {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState(null);

  useEffect(() => {
    const script = document.createElement('script');
    script.src = 'https://telegram.org/js/telegram-web-app.js';
    script.async = true;
    document.body.appendChild(script);

    return () => {
      if (document.body.contains(script)) {
        document.body.removeChild(script);
      }
    };
  }, []);

  const handleTelegramLogin = async () => {
    setIsLoading(true);
    setError(null);

    try {
      if (window.Telegram?.WebApp) {
        const tg = window.Telegram.WebApp;
        tg.ready();

        const userData = tg.initDataUnsafe?.user;
        
        if (userData) {
          const response = await fetch('/api/auth/telegram/register', {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify({
              telegramId: userData.id,
              firstName: userData.first_name,
              lastName: userData.last_name || '',
              username: userData.username || '',
              photoUrl: userData.photo_url || '',
              language: userData.language_code || 'en',
            }),
          });

          if (response.ok) {
            const data = await response.json();
            console.log('Registration successful:', data);
          } else {
            setError('Ошибка при регистрации. Попробуйте еще раз.');
          }
        } else {
          setError('Не удалось получить данные Telegram.');
        }
      } else {
        const botUsername = 'ты мне так и не скинул юзернейм бота :(';
        window.open(`https://t.me/${botUsername}`, '_blank');
      }
    } catch (err) {
      setError('Ошибка подключения. Проверьте интернет соединение.');
      console.error('Login error:', err);
    } finally {
      setIsLoading(false);
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
                <path d="M12 2L15 10H23L17 15L19 23L12 18L5 23L7 15L1 10H9L12 2Z" strokeWidth="1.5"/>
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
                <path d="M3 12L5 10L9 14L15 3L17 5" strokeWidth="1.5" strokeLinecap="round"/>
                <path d="M21 6L21 18M6 21H18" strokeWidth="1.5" strokeLinecap="round"/>
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
                <circle cx="12" cy="12" r="9" strokeWidth="1.5"/>
                <circle cx="12" cy="12" r="5" strokeWidth="1.5"/>
                <circle cx="12" cy="12" r="1.5" fill="currentColor"/>
              </svg>
            </div>
            <h3 className="feature-title">Высокая точность</h3>
            <p className="feature-description">
              Уровень точности прогнозирования бла бла бла%
            </p>
          </div>

          <div className="feature-card">
            <div className="feature-icon">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor">
                <path d="M12 2L2 6V12C2 18 12 22 12 22C12 22 22 18 22 12V6L12 2Z" strokeWidth="1.5"/>
              </svg>
            </div>
            <h3 className="feature-title">Безопасность & Конфиденциальность</h3>
            <p className="feature-description">
              Ваши данные зашифрованы и защищены
            </p>
          </div>
        </section>

        <div className="cta-section">
          <button 
            className="get-started-btn"
            onClick={handleTelegramLogin}
            disabled={isLoading}
          >
            {isLoading ? (
              <>
                <span className="loading-spinner"></span>
                Загрузка...
              </>
            ) : (
              <>
                <svg className="telegram-icon" viewBox="0 0 24 24" fill="currentColor">
                  <path d="M12 0C5.373 0 0 5.373 0 12s5.373 12 12 12 12-5.373 12-12S18.627 0 12 0zm5.894 8.221l-1.97 9.28c-.145.658-.537.818-1.084.508l-3-2.21-1.446 1.394c-.14.146-.357.292-.611.292-.005 0-.01 0-.016 0l.213-3.053 5.56-5.023c.242-.213-.054-.328-.375-.115l-6.869 4.332-2.961-.924c-.643-.204-.657-.643.136-.953l11.566-4.458c.529-.196 1.083.128.897.983z"/>
                </svg>
                Начать работу с Telegram
              </>
            )}
          </button>
          {error && <p className="error-message">{error}</p>}
        </div>

        <footer className="landing-footer">
          <p className="footer-main">
            Начните прогнозировать сейчас
          </p>
          <p className="footer-secondary">
            На базе продвинутого ИИ • Не является финансовым советом
          </p>
        </footer>
      </main>
    </div>
  );
};
