import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    allowedHosts: [
      "2f15-2a01-ecc0-c80-21-00-2.ngrok-free.app",
      // Или разрешить все хосты (менее безопасно)
      // 'all'
    ],
    // Также можно указать хост
    host: true, // Прослушивать все сетевые интерфейсы
  },
});
