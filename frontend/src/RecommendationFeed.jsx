import { useEffect, useState } from "react";
import axios from "axios";
import RecommendationCard from "./RecommendationCard";
import Footer from "./Footer";
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
        const res = await axios.post("http://localhost:8081/v1/login", {
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
        "http://localhost:8081/v1/get_recommendated_posts",
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
    <div className="min-h-screen flex flex-col gap-6 bg-gradient-to-br from-[#0A0F1F] via-[#0F1A2F] to-[#02B7DB] flex items-center justify-center p-4 sm:p-6">
      {recommendatedPosts.map((post, index) => (
        <RecommendationCard key={index} {...post} />
      ))}

      <Footer />
    </div>
  );
}
