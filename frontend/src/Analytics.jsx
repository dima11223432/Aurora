import { useEffect, useRef } from "react";

export default function Analytics() {
  const containerRef = useRef(null);

  useEffect(() => {
    const script = document.createElement("script");
    script.src = "https://s3.tradingview.com/tv.js";
    script.async = true;

    script.onload = () => {
      new window.TradingView.widget({
        width: "100%",
        height: "500",
        symbol: "USDT",
        locale: "ru",
        theme: "Dark",
        container_id: containerRef.current.id,
      });
    };

    document.body.appendChild(script);

    return () => {
      if (containerRef.current) containerRef.current.innerHTML = "";
    };
  }, []);
  
  return (
    <div
      className="min-h-screen text-white"
      style={{ backgroundColor: "#151925" }}
    >
      <header className="bg-black/20 backdrop-blur-sm">
        <div className="container mx-auto px-4 py-4">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2">
              <div className="flex items-center gap-2">
                <img
                  src="/Aurora-logo.png"
                  alt="Aurora Logo"
                  className="w-12 h-12 object-contain"
                />
                <span
                  className="text-2xl font-bold bg-clip-text text-transparent pt-6 block"
                  style={{
                    backgroundImage:
                      "linear-gradient(to right, #36DEF4, #208390)",
                    WebkitBackgroundClip: "text",
                  }}
                >
                  urora
                </span>
              </div>
            </div>

            <button
              className="px-6 py-2 rounded-lg font-semibold transition-all transform hover:scale-105 shadow-lg text-white"
              style={{
                background: "linear-gradient(to right, #208390, #36DEF4)",
                border: "none",
                cursor: "pointer",
              }}
              onMouseEnter={(e) =>
                (e.target.style.background =
                  "linear-gradient(to right, #6bedfeff, #54f1ffff)")
              }
              onMouseLeave={(e) =>
                (e.target.style.background =
                  "linear-gradient(to right, #208390, #36DEF4)")
              }
            >
              Connect Wallet
            </button>
          </div>
        </div>
      </header>

      <main className="container mx-auto px-4 py-6 pb-32">
        <div className="grid grid-cols-2 gap-4 mb-6">
          <div className="bg-grey backdrop-blur-sm rounded-xl p-4 border border-cyan-500 h-32 flex flex-col justify-between">
            <p className="text-yellow-400 text-lg">Капитал</p>
            <p className="text-4xl font-bold text-cyan-400 self-end">$12,999</p>
          </div>
          <div className="bg-grey backdrop-blur-sm rounded-xl p-4 border border-cyan-500 h-32 flex flex-col justify-between">
            <p className="text-yellow-400 text-lg">Точность</p>
            <p className="text-4xl font-bold text-cyan-400 self-end">86,2%</p>
          </div>
        </div>

        <div id="tradingview_widget"></div>
        <div
          id="tradingview_widget"
          ref={containerRef}
          style={{ width: "100%", height: "500px" }}
        ></div>

        <div className="bg-grey backdrop-blur-sm rounded-xl p-5 border border-cyan-500 mb-6">
          <h2 className="text-yellow-400 text-2xl font-bold mb-2">
            Рекомендация ИИ
          </h2>
          <div className="mb-3 flex items-center">
            <span className="bg-green-500/20 text-xl font-bold text-green-400 font-bold px-6 py-2 rounded-full">
              BUY
            </span>
            <span className="text-xl text-gray-300 ml-3 font-extrabold">
              Акция - YNDX
            </span>
          </div>
          <p className="text-gray-300 text-sm leading-relaxed">
            Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do
            eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim
            ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut
            aliquip ex ea commodo consequat. Duis aute irure dolor in
            reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla
            pariatur.
          </p>
        </div>
      </main>

  <div className="fixed bottom-0 left-0 right-0 z-50">
  <div className="fixed bottom-0 left-0 right-0 z-40 h-48 bg-gradient-to-t from-black/80 via-black/40 to-transparent">
  <div className="fixed bottom-0 left-0 right-0 z-50">
  <div className="container mx-auto px-4 py-2">
    <div 
      className="rounded-full p-3 border border-cyan-500"
      style={{ backgroundColor: "#151925" }}
    >
      <div className="grid grid-cols-3 gap-2 mb-2">
        
        {/* Аналитика */}
        <button className="flex flex-col items-center justify-center py-2 px-1 text-cyan-300 hover:text-cyan hover:bg-cyan-600/20 rounded-2xl transition-colors">
          <div className="w-8 h-8 mb-1 flex items-center justify-center">
            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
            </svg>
          </div>
          <span className="text-xs">Аналитика</span>
        </button>
        
        {/* Новости */}
        <button className="flex flex-col items-center justify-center py-2 px-1 text-cyan-300 hover:text-cyan hover:bg-cyan-600/20 rounded-2xl transition-colors">
          <div className="w-8 h-8 mb-1 flex items-center justify-center">
            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M19 20H5a2 2 0 01-2-2V6a2 2 0 012-2h10a2 2 0 012 2v1m2 13a2 2 0 01-2-2V7m2 13a2 2 0 002-2V9a2 2 0 00-2-2h-2m-4-3H9M7 16h6M7 8h6v4H7V8z" />
            </svg>
          </div>
          <span className="text-xs">Новости</span>
        </button>
        
        {/* Портфель */}
        <button className="flex flex-col items-center justify-center py-2 px-1 text-cyan-300 hover:text-cyan hover:bg-cyan-600/20 rounded-2xl transition-colors">
          <div className="w-8 h-8 mb-1 flex items-center justify-center">
            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M3 10h18M7 15h1m4 0h1m-7 4h12a3 3 0 003-3V8a3 3 0 00-3-3H6a3 3 0 00-3 3v8a3 3 0 003 3z" />
            </svg>
          </div>
          <span className="text-xs">Портфель</span>
        </button>
      </div>
      
    </div>
  </div>
  </div>
</div>
</div>
    </div>
  );
}