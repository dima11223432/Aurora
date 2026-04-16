import { useEffect, useState } from "react";
import axios from "axios";
import RecommendationCard from "./RecommendationCard";
export default function RecommendationFeed() {
  const [recommendatedPosts, setRecommendedPosts] = useState([]);
  const [cursor, setCursor] = useState({
    score: 0,
    id: "",
  });
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
      } catch (e) {
        console.error("Login error:", e);
      }
    };

    login();
  }, []);
  useEffect(() => {
    const fetchPosts = async () => {
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
        console.log(response.data);
      } catch (e) {
        console.error(e);
      }
    };
    fetchPosts();
  }, []);
  return (
    <div className="min-h-screen flex flex-col gap-6 bg-gradient-to-br from-[#0A0F1F] via-[#0F1A2F] to-[#02B7DB] flex items-center justify-center p-4 sm:p-6">
      {recommendatedPosts.map((post, index) => (
        <RecommendationCard key={index} {...post} />
      ))}
    </div>
  );
}
