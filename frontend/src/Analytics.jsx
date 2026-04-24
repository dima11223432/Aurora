import { useEffect, useRef } from "react";
import RecommendationFeed from "./RecommendationFeed";

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
            <p className="text-yellow-400 text-lg">Точность</p>
            <p className="text-4xl font-bold text-cyan-400 self-end">86,2%</p>
          </div>
        </div>

        <div id="tradingview_widget"></div>
        <div id="tradingview_widget" ref={containerRef}></div>
        <RecommendationFeed />
      </main>

      <div className="fixed bottom-0 left-0 right-0 z-50">
        <div className="fixed bottom-0 left-0 right-0 z-40 h-48 bg-gradient-to-t from-black/80 via-black/40 to-transparent"></div>
      </div>
    </div>
  );
}
