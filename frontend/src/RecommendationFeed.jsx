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
          is_admin: true,
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
  return (
    // {/* <div className="min-h-screen flex flex-col gap-6 bg-gradient-to-br from-[#0A0F1F] via-[#0F1A2F] to-[#02B7DB] flex items-center justify-center p-4 sm:p-6"> */}

    <div className="bg-gray-900 min-h-screen flex flex-col p-4 sm:p-6">
      <div className="flex justify-end mb-4">
        <TonConnectButton />
      </div>
      <div className="flex-grow flex flex-col items-center gap-6">
        {recommendatedPosts.map((post, index) => (
          <RecommendationCard key={index} {...post} />
        ))}

        {isLoading && <p className="text-white">Загрузка...</p>}
      </div>
      <Shtora />
      <Footer />
    </div>
  );
}
