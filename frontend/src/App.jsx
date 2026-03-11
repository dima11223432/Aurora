import React from 'react';
import './index.css';

function App() {
  return (
    <div 
      className="min-h-screen text-white"
      style={{ backgroundColor: '#151925' }}
    >
      <header className=" bg-black/20 backdrop-blur-sm">
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
                    backgroundImage: 'linear-gradient(to right, #36DEF4, #208390)',
                    WebkitBackgroundClip: 'text'
                  }}
                >
                urora
                </span>
              </div>
            </div>

            <button 
              className="px-6 py-2 rounded-lg font-semibold transition-all transform hover:scale-105 shadow-lg text-white"
              style={{ 
                background: 'linear-gradient(to right, #208390, #36DEF4)',
                border: 'none',
                cursor: 'pointer'
              }}
              onMouseEnter={(e) => e.target.style.background = 'linear-gradient(to right, #6bedfeff, #54f1ffff)'}
              onMouseLeave={(e) => e.target.style.background = 'linear-gradient(to right, #208390, #36DEF4)'}
            >
            Connect Wallet
          </button>
          </div>
        </div>
      </header>

      <main className="container mx-auto px-4 py-6">
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

        <div className="bg- backdrop-blur-sm rounded-xl p-4 border border-cyan-500 mb-6">
          <h2 className="text-yellow-400 text-lg font-semibold mb-3">График Цен TON</h2>
          
          <div className="flex justify-between text-xs text-gray-400 mb-2">
            <span>2.6</span>
            <span>1.95</span>
            <span>1.3</span>
            <span>0.65</span>
            <span>0</span>
          </div>

          <div className="h-32 relative mb-4">
            <svg className="w-full h-full" viewBox="0 0 300 100" preserveAspectRatio="none">
              <path
                d="M0,80 Q30,70 60,50 T120,30 T180,45 T240,20 T300,40"
                stroke="#36DEF4"
                strokeWidth="2"
                fill="none"
              />
              <path
                d="M0,80 Q30,70 60,50 T120,30 T180,45 T240,20 T300,40"
                stroke="url(#gradient)"
                strokeWidth="2"
                fill="none"
                opacity="0.3"
              />
            </svg>
          </div>

          <div className="flex justify-between text-xs text-gray-400">
            <span>10:00</span>
            <span>11:00</span>
            <span>12:00</span>
            <span>13:00</span>
            <span>14:00</span>
            <span>15:00</span>
            <span>17:00</span>
          </div>
        </div>

        <div className="bg-grey backdrop-blur-sm rounded-xl p-5 border border-cyan-500 mb-6">
          <h2 className="text-yellow-400 text-2xl font-bold mb-2">Рекомендация ИИ</h2>
          <div className="mb-3 flex items-center">
            <span className="bg-green-500/20 text-xl font-bold text-green-400 font-bold px-6 py-2 rounded-full">BUY</span>
            <span className="text-xl text-gray-300 ml-3 font-extrabold">Акция - YNDX</span>
          </div>
          <p className="text-gray-300 text-sm leading-relaxed">
            Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur.
          </p>
        </div>

        <div className="grid grid-cols-3 gap-2 bg-grey backdrop-blur-sm rounded-xl p-2 border border-cyan-500">
          <button className="py-3 text-center text-cyan-300 hover:text-cyan hover:bg-cyan-600/20 rounded-lg transition-colors">
            Аналитика
          </button>
          <button className="py-3 text-center text-cyan-300 hover:text-cyan hover:bg-cyan-600/20 rounded-lg transition-colors">
            Новости
          </button>
          <button className="py-3 text-center text-cyan-300 hover:text-cyan hover:bg-cyan-600/20 rounded-lg transition-colors">
            Портфель
          </button>
        </div>
      </main>
    </div>
  );
}

export default App;