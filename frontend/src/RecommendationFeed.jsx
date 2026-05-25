import { useEffect, useState, useRef } from "react";
import axios from "axios";
import RecommendationCard from "./RecommendationCard";
import Footer from "./Footer";
import { routes } from "./config/api";
import { TonConnectButton } from "@tonconnect/ui-react";
import { useNavigate } from "react-router-dom";

export default function RecommendationFeed() {
  const [recommendatedPosts, setRecommendedPosts] = useState([]);
  const [cursor, setCursor] = useState({ score: 0, id: "" });
  const [isLoading, setIsLoading] = useState(false);
  const [hasMore, setHasMore] = useState(true);
  const navigate = useNavigate();

  const stateRef = useRef({ isLoading, hasMore, cursor });

  useEffect(() => {
    stateRef.current = { isLoading, hasMore, cursor };
  }, [isLoading, hasMore, cursor]);

  const fetchPosts = async (isRefresh = false) => {
    const currentIsLoading = isRefresh ? isLoading : stateRef.current.isLoading;
    const currentHasMore = isRefresh ? true : stateRef.current.hasMore;

    if (!isRefresh && (currentIsLoading || !currentHasMore)) return;

    setIsLoading(true);

    try {
      const token = localStorage.getItem("token");
      if (!token) {
        // navigate("/404");
        return;
      }

      const currentCursor = isRefresh
        ? { score: 0, id: "" }
        : stateRef.current.cursor;

      const response = await axios.post(
        routes.getRecommendatedPosts,
        { cursor: currentCursor },
        { headers: { Authorization: `Bearer ${token}` } },
      );

      const newPosts = response.data.posts || [];
      const nextCursor = response.data.nextCursor;

      setRecommendedPosts((prev) =>
        isRefresh ? newPosts : [...prev, ...newPosts],
      );

      setCursor(nextCursor);

      if (nextCursor === null || newPosts.length === 0) {
        setHasMore(false);
      } else if (isRefresh) {
        setHasMore(true);
      }
    } catch (e) {
      console.error("Ошибка при запросе постов:", e);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchPosts();
  }, []);

  useEffect(() => {
    const handleScroll = () => {
      const bottom =
        window.innerHeight + window.scrollY >=
        document.documentElement.offsetHeight - 150;
      if (bottom) {
        fetchPosts();
      }
    };
    window.addEventListener("scroll", handleScroll);
    return () => {
      window.removeEventListener("scroll", handleScroll);
    };
  }, [cursor, isLoading, hasMore, isLoggedIn]);

  const EmptyFeedState = () => (
    <div className="flex flex-col items-center justify-center min-h-[400px] text-center px-4">
      <svg 
        className="w-24 h-24 mb-6 text-gray-600"
        fill="none" 
        stroke="currentColor" 
        viewBox="0 0 24 24"
      >
        <path 
          strokeLinecap="round" 
          strokeLinejoin="round" 
          strokeWidth={1.5} 
          d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" 
        />
      </svg>
      <p className="text-gray-300 text-lg mb-2">
        Нет новостей для показа
      </p>
      <p className="text-gray-400 text-sm max-w-md">
        Для того, чтобы получать аналитику, откройте раздел "Настройки" 
        и отметьте галочкой нужные вам каналы
      </p>
    </div>
  );
    return (
    <div className="bg-gray-900 min-h-screen flex flex-col p-4 sm:p-6">
      <div className="flex justify-between items-center mb-6">
        <button
          className="text-cyan-300 border border-cyan-300/30 hover:bg-cyan-300/20 px-4 py-2 rounded-2xl transition-all font-medium text-sm active:scale-95 disabled:opacity-50"
          onClick={() => fetchPosts(true)}
          disabled={isLoading}
        >
          {isLoading && !recommendatedPosts.length
            ? "Загрузка..."
            : "Обновить ленту"}
        </button>

        <TonConnectButton />
      </div>

      <div className="flex-grow flex flex-col items-center gap-6">
        {!isLoading && recommendatedPosts.length === 0 ? (
          <EmptyFeedState />
        ) : (
          <>
            {recommendatedPosts.map((post, index) => (
              <RecommendationCard key={index} {...post} />
            ))}
            {isLoading && <p className="text-white">Загрузка...</p>}
          </>
        )}
      </div>
      <Footer />
    </div>
  );
}
