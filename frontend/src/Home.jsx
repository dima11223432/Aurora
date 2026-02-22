import React from "react";

const Home = () => {
  return (
    <div className="relative w-full min-h-screen overflow-x-hidden bg-gradient-to-br from-[#0a1428] via-[#0d1f2d] via-25% via-[#132238] via-50% to-[#0a1428] text-white/87">
      {/* Gradient Background */}
      <div
        className="fixed top-0 left-0 w-full h-full pointer-events-none z-0"
        style={{
          background: `
               radial-gradient(circle at 20% 50%, rgba(0, 200, 255, 0.08) 0%, transparent 50%),
               radial-gradient(circle at 80% 80%, rgba(100, 200, 255, 0.05) 0%, transparent 50%)
             `,
        }}
      />

      <main className="relative z-10 flex flex-col items-center justify-center min-h-screen px-6 py-8 mx-auto max-w-7xl sm:px-8 md:px-6">
        {/* Header */}
        <header className="flex flex-col items-center gap-4 mb-8">
          <div className="w-20 h-20 text-[#00c8ff] drop-shadow-[0_0_20px_rgba(0,200,255,0.5)] sm:w-16 sm:h-16">
            <svg
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              className="w-full h-full"
            >
              <path
                d="M12 2L15 10H23L17 15L19 23L12 18L5 23L7 15L1 10H9L12 2Z"
                strokeWidth="1.5"
              />
            </svg>
          </div>
          <h1 className="text-4xl font-bold bg-gradient-to-r from-[#00c8ff] to-[#00a8d8] bg-clip-text text-transparent tracking-wider sm:text-3xl">
            Aurora
          </h1>
        </header>

        {/* Hero Section */}
        <section className="mb-12 text-center">
          <h2 className="text-5xl font-bold mb-4 leading-tight tracking-tight sm:text-4xl md:text-3xl">
            Добро пожаловать в{" "}
            <span className="bg-gradient-to-r from-[#00c8ff] to-[#00a8d8] bg-clip-text text-transparent">
              Aurora
            </span>
          </h2>
          <p className="text-2xl text-white/70 font-light tracking-wide sm:text-xl md:text-lg">
            AI-сделанные TON & T-Инвестиционные прогнозы
          </p>
        </section>
        {/* Features Grid - Центральное расположение квадратом */}
        <div className="grid grid-cols-2 gap-4 w-full max-w-2xl mx-auto mb-12 lg:gap-4 md:grid-cols-2">
          {/* Feature Card 1 */}
          <div className="p-8 border border-[#00c8ff]/20 rounded-2xl bg-gradient-to-br from-[#0a1e3c]/80 to-[#0d1f2d]/60 backdrop-blur-md transition-all duration-300 hover:border-[#00c8ff]/50 hover:from-[#0a2846]/90 hover:to-[#0d2837]/80 hover:-translate-y-1 hover:shadow-[0_10px_30px_rgba(0,200,255,0.15)] cursor-pointer sm:p-6 aspect-square flex flex-col">
            <div className="w-12 h-12 text-[#00c8ff] mb-4 flex-shrink-0 sm:w-10 sm:h-10">
              <svg
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                className="w-full h-full"
              >
                <path
                  d="M12 2L15 10H23L17 15L19 23L12 18L5 23L7 15L1 10H9L12 2Z"
                  strokeWidth="1.5"
                />
              </svg>
            </div>
            <h3 className="text-xl font-semibold text-white mb-3 sm:text-lg">
              AI-сделанные Прогнозы
            </h3>
            <p className="text-white/65 leading-relaxed sm:text-sm">
              Продвинутые модели машинного обучения анализируют тренды рынка
            </p>
          </div>
{/* Feature Card 2 */}
          <div className="p-8 border border-[#00c8ff]/20 rounded-2xl bg-gradient-to-br from-[#0a1e3c]/80 to-[#0d1f2d]/60 backdrop-blur-md transition-all duration-300 hover:border-[#00c8ff]/50 hover:from-[#0a2846]/90 hover:to-[#0d2837]/80 hover:-translate-y-1 hover:shadow-[0_10px_30px_rgba(0,200,255,0.15)] cursor-pointer sm:p-6 aspect-square flex flex-col">
            <div className="w-12 h-12 text-[#00c8ff] mb-4 flex-shrink-0 sm:w-10 sm:h-10">
              <svg
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                className="w-full h-full"
              >
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
            <h3 className="text-xl font-semibold text-white mb-3 sm:text-lg">
              Анализ в реальном времени
            </h3>
            <p className="text-white/65 leading-relaxed sm:text-sm">
              Данные рынка в реальном времени и мгновенные обновления прогнозов
            </p>
          </div>

          {/* Feature Card 3 */}
          <div className="p-8 border border-[#00c8ff]/20 rounded-2xl bg-gradient-to-br from-[#0a1e3c]/80 to-[#0d1f2d]/60 backdrop-blur-md transition-all duration-300 hover:border-[#00c8ff]/50 hover:from-[#0a2846]/90 hover:to-[#0d2837]/80 hover:-translate-y-1 hover:shadow-[0_10px_30px_rgba(0,200,255,0.15)] cursor-pointer sm:p-6 aspect-square flex flex-col">
            <div className="w-12 h-12 text-[#00c8ff] mb-4 flex-shrink-0 sm:w-10 sm:h-10">
              <svg
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                className="w-full h-full"
              >
                <circle cx="12" cy="12" r="9" strokeWidth="1.5" />
                <circle cx="12" cy="12" r="5" strokeWidth="1.5" />
                <circle cx="12" cy="12" r="1.5" fill="currentColor" />
              </svg>
            </div>
            <h3 className="text-xl font-semibold text-white mb-3 sm:text-lg">
              Высокая точность
            </h3>
            <p className="text-white/65 leading-relaxed sm:text-sm">
              Уровень точности прогнозирования до 94%
            </p>
          </div>


{/* Feature Card 4 */}
          <div className="p-8 border border-[#00c8ff]/20 rounded-2xl bg-gradient-to-br from-[#0a1e3c]/80 to-[#0d1f2d]/60 backdrop-blur-md transition-all duration-300 hover:border-[#00c8ff]/50 hover:from-[#0a2846]/90 hover:to-[#0d2837]/80 hover:-translate-y-1 hover:shadow-[0_10px_30px_rgba(0,200,255,0.15)] cursor-pointer sm:p-6 aspect-square flex flex-col">
            <div className="w-12 h-12 text-[#00c8ff] mb-4 flex-shrink-0 sm:w-10 sm:h-10">
              <svg
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                className="w-full h-full"
              >
                <path
                  d="M12 2L2 6V12C2 18 12 22 12 22C12 22 22 18 22 12V6L12 2Z"
                  strokeWidth="1.5"
                />
              </svg>
            </div>
            <h3 className="text-xl font-semibold text-white mb-3 sm:text-lg">
              Безопасность & Конфиденциальность
            </h3>
            <p className="text-white/65 leading-relaxed sm:text-sm">
              Ваши данные зашифрованы и защищены
            </p>
          </div>
        </div>

        {/* Footer */}
        <footer className="text-center">
          <p className="text-white/75 mb-3 tracking-wide sm:text-sm">
            Начните прогнозировать сейчас
          </p>
          <p className="text-white/50 text-sm sm:text-xs">
            На базе продвинутого ИИ • Не является финансовым советом
          </p>
        </footer>
      </main>
    </div>
  );
}
export default Home;