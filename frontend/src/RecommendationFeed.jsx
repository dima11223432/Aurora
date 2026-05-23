import { useEffect, useState } from "react";
import axios from "axios";
import RecommendationCard from "./RecommendationCard";
import Footer from "./Footer";
import Shtora from "./Shtora";
import { routes } from "./config/api";
import { TonConnectButton } from "@tonconnect/ui-react";
export default function RecommendationFeed() {
  const [recommendatedPosts, setRecommendedPosts] = useState([]);
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [cursor, setCursor] = useState({
    score: 0,
    id: "",
  });
  const [isLoading, setIsLoading] = useState(false);
  const [hasMore, setHasMore] = useState(true);
  useEffect(() => {
    const login = async () => {
      try {
        const res = await axios.post(routes.login, {
          telegram_id: 123456789,
          username: "john_doe",
          first_name: "John",
          last_name: "Doe",
          is_admin: false,
          app_id: 1,
        });

        localStorage.setItem("token", res.data.token);
        setIsLoggedIn(true);
      } catch (e) {
        console.error("Login error:", e);
      }
    };

    login();
  }, []);
  const fetchPosts = async () => {
    if (isLoading || !hasMore) return;

    setIsLoading(true);

    try {
      const token = localStorage.getItem("token");

      const response = await axios.post(
        routes.getRecommendatedPosts,
        { cursor },
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        },
      );

      setRecommendedPosts((prev) => [...prev, ...response.data.posts]);
      setCursor(response.data.nextCursor);
      if (response.data.nextCursor === null) {
        setHasMore(false);
      }
    } catch (e) {
      console.error(e);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    if (isLoggedIn) {
      fetchPosts();
    }
  }, [isLoggedIn]);

  useEffect(() => {
    const handleScroll = () => {
      const bottom =
        window.innerHeight + window.scrollY >=
        document.documentElement.offsetHeight - 100;
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
      <div className="flex justify-end mb-4">
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
      <Shtora />
      <Footer />
    </div>
  );
}
