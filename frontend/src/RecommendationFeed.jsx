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
  }, []);

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
