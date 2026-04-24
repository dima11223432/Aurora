import React from 'react';

const NewsSection = () => {
  const articles = [
    {
      id: 1,
      source: "РИА Новости",
      title: "Lorem ipsum dolor",
      text: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur.",
      badges: ["YNDX", "YNDX"],
    },
    {
      id: 2,
      source: "РИА Новости",
      title: "Lorem ipsum dolor",
      text: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur.",
      badges: ["YNDX", "YNDX"],
    },
    {
      id: 3,
      source: "РИА Новости",
      title: "Lorem ipsum dolor",
      text: "ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris",
      badges: ["YNDX", "YNDX"],
    }
  ];

  const ArticleCard = ({ article, isLast }) => (
    <article className={`pb-8 ${!isLast ? 'mb-8 border-b border-gray-700' : ''}`}>
      <div 
        className="w-full mb-3 rounded-full flex items-center gap-2 px-3 py-1.5"
        style={{
          background: "linear-gradient(to right, #36DEF4, #208390)",
        }}
      >
        <div 
          className="w-5 h-5 rounded-full"
          style={{ backgroundColor: "#182E3A" }}
        ></div>
        <span className="text-xs text-white font-medium tracking-wide">
          {article.source}
        </span>
      </div>
      
      <h3 className="text-2xl font-bold text-white mb-4 tracking-tight">
        {article.title}
      </h3>
      
      <p className="text-gray-300 leading-relaxed mb-5 text-base">
        {article.text}
      </p>
      
      <div className="flex flex-wrap gap-4">
        <div 
          className="inline-flex items-center gap-1.5 px-3 py-1 rounded-md"
          style={{
            border: "1px solid",
            borderColor: "#3FFF5B",
            backgroundColor: "transparent"
          }}
        >
          <span style={{ color: "#3FFF5B" }} className="font-semibold text-sm">
            {article.badges[0]}
          </span>
          <svg 
            className="w-3.5 h-3.5" 
            fill="none" 
            stroke="currentColor" 
            viewBox="0 0 24 24"
            style={{ color: "#3FFF5B" }}
          >
            <path 
              strokeLinecap="round" 
              strokeLinejoin="round" 
              strokeWidth={2} 
              d="M5 15l7-7 7 7" 
            />
          </svg>
        </div>

        <div 
          className="inline-flex items-center gap-1.5 px-3 py-1 rounded-md"
          style={{
            border: "1px solid",
            borderColor: "#FF0000",
            backgroundColor: "transparent"
          }}
        >
          <span style={{ color: "#FF0000" }} className="font-semibold text-sm">
            {article.badges[1]}
          </span>
          <svg 
            className="w-3.5 h-3.5" 
            fill="none" 
            stroke="currentColor" 
            viewBox="0 0 24 24"
            style={{ color: "#FF0000" }}
          >
            <path 
              strokeLinecap="round" 
              strokeLinejoin="round" 
              strokeWidth={2} 
              d="M19 9l-7 7-7-7" 
            />
          </svg>
        </div>
      </div>
    </article>
  );

  return (
    <div className="max-w-3xl mx-auto px-5 py-10">
      {articles.map((article, index) => (
        <ArticleCard 
          key={article.id}
          article={article}
          isLast={index === articles.length - 1}
        />
      ))}
    </div>
  );
};

const Header = () => (
  <header className="bg-black/20 backdrop-blur-sm sticky top-0 z-10">
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
                backgroundImage: "linear-gradient(to right, #36DEF4, #208390)",
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
);

const App = () => {
  return (
    <div 
      className="min-h-screen font-sans antialiased"
      style={{ backgroundColor: "#151925" }}
    >
      <Header />
      <main>
        <NewsSection />
      </main>
    </div>
  );
};

export default App;