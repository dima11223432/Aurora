import React from 'react';

const Herozone = () => {
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

                <div className="text-center">
                    <button className="bg-[#0fd2f5] text-black text-lg sm:text-xl md:text-2xl font-semibold px-6 sm:px-8 md:px-5 py-3 sm:py-4 rounded-full shadow-[0_10px_25px_-5px_#014b5c] hover:shadow-[0_15px_30px_-5px_#014b5c] hover:scale-105 active:scale-95 transition-all duration-200">
                        Начать работу →
                    </button>
                </div>
            </div>
            
        </div>
    );
};

export default Herozone;