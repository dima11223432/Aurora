import axios from "axios";

const api = axios.create({
  baseURL: "http://localhost:5000/api",
  timeout: 10000,
  headers: {
    "Content-Type": "application/json",
  },
});
export const analyticsAPI = {
    getPortfolio: async () => {
        try {
            const response = await api.get("/analytics/portfolio");
            return response.data;
        } catch (error) {
        console.error("Ошибка загрузки портфеля:", error);
        throw error;
        }
 },
};