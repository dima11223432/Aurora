import React, { useEffect, useState, useRef } from "react";

const Herozone = () => {
    const [isLoading, setIsLoading] = useState(false);
      const [error, setError] = useState(null);
      const [telegramUser, setTelegramUser] = useState(null);
      const widgetContainerRef = useRef(null);
      const API_URL = "https://27dc-213-176-17-134.ngrok-free.app/v1/login";
    
      useEffect(() => {
        window.onTelegramAuth = async (user) => {
          console.log("Telegram auth success:", user);
          setIsLoading(true);
    
          try {
            fetch(API_URL, {
              method: "POST",
              headers: {
                "Content-Type": "application/json",
              },
              body: JSON.stringify({
                telegram_id: user.id,
                username: user.username || "",
                first_name: user.first_name || "",
                last_name: user.last_name || "",
                is_admin: false,
                app_id: 1,
              }),
            })
              .then((res) => res.json())
              .then((data) => {
                alert("JWT Token:" + data.token);
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

    const items = [
        {
            title: 'ИИ аналитика',
            description: 'Продвинутые модели машинного обучения анализируют рыночные тенденции',
        },
        {
            title: 'Анализ в реальном времени',
            description: 'Текущие рыночные данные и мгновенное обновление прогнозов',
        },
        {
            title: 'Высокая точность',
            description: '*% точности прогнозов на основе исторических данных и текущих рыночных условий',
        },
        {
            title: 'Безопасный и приватный',
            description: 'Ваши данные зашифрованы и защищены',
        }
    ];

    return (
        <div className="min-h-screen bg-gradient-to-br from-[#0A0F1F] via-[#0F1A2F] to-[#02B7DB] flex items-center justify-center p-4 sm:p-6">
            <div className="max-w-3xl w-full bg-[rgba(20,25,50,0.7)] backdrop-blur-md rounded-[3rem] p-6 sm:p-8 md:p-12 shadow-2xl border border-white/5">
                <img
                    src="/Aurora-logo.png"
                    alt="Aurora logo"
                    className="w-20 h-20 mb-2 mx-auto"
                />
                
                <p className="text-[#95bec7] text-base text-center sm:text-lg mb-8 sm:mb-12 pl-4 to-transparent">
                    Добро пожаловать в <span className="text-[#0fd2f5] font-bold">Aurora</span><br /> ИИ-поддерживаемые прогнозы TON и T-Investments
                </p>

                <div className="space-y-4 sm:space-y-5 mb-8 sm:mb-12">
                    {items.map((item) => (
                        <div key={item.id} className="flex gap-3 sm:gap-4 items-start group">
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
                <div className="w-full max-w-md mx-auto mb-16 animate-fadeInUp [animation-delay:400ms]">
                <div className="relative p-8 rounded-3xl border border-primary/30 bg-gradient-to-br from-[#0a1e3c]/70 to-[#0d1f2d]/50 backdrop-blur-xl shadow-[0_20px_50px_-15px_rgba(0,200,255,0.25)]">
                    <div className="absolute -inset-0.5 bg-gradient-to-r from-primary/20 to-transparent rounded-3xl blur-xl opacity-50" />
                    <div className="relative">
                    <p className="text-white/60 text-center mb-6 text-sm">
                        {!telegramUser
                        ? "Войдите через Telegram"
                        : "Управляйте своим профилем"}
                    </p>

                    {!telegramUser ? (
                        <>
                        {isLoading ? (
                            <div className="flex justify-center py-4">
                            <span className="inline-block w-6 h-6 border-2 border-white/30 border-t-white rounded-full animate-spin"></span>
                            </div>
                        ) : (
                            <div
                            id="telegram-login-container"
                            ref={widgetContainerRef}
                            className="flex justify-center items-center min-h-[60px]"
                            />
                        )}
                        </>
                    ) : (
                        <div className="flex flex-col items-center">
                        <div className="flex items-center gap-4 mb-4">
                            {telegramUser.photo_url ? (
                            <img
                                src={telegramUser.photo_url}
                                alt="Profile"
                                className="w-16 h-16 rounded-full ring-2 ring-primary/50"
                            />
                            ) : (
                            <div className="w-16 h-16 rounded-full bg-primary/20 ring-2 ring-primary/50 flex items-center justify-center text-2xl font-bold text-primary">
                                {telegramUser.first_name?.charAt(0)}
                            </div>
                            )}
                            <div className="text-left">
                            <h4 className="text-white font-semibold text-lg">
                                {telegramUser.first_name} {telegramUser.last_name}
                            </h4>
                            {telegramUser.username && (
                                <p className="text-white/50">@{telegramUser.username}</p>
                            )}
                            </div>
                        </div>
                        <button
                            onClick={handleLogout}
                            className="px-6 py-2 bg-red-500/20 border border-red-500/50 rounded-full text-red-400 text-sm hover:bg-red-500/30 transition-colors"
                        >
                            Выйти
                        </button>
                        </div>
                    )}

                    {error && (
                        <div className="mt-4 p-3 bg-red-500/10 border border-red-500/50 rounded-xl text-red-400 text-sm text-center">
                        <span className="mr-2">⚠️</span>
                        {error}
                        </div>
                    )}
                    </div>
                </div>
            </div>
        </div>
            
    </div>
    );
};

export default Herozone;
