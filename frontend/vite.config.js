import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

export default defineConfig({
  plugins: [react()],

  server: {
    allowedHosts: [
      "06d4-213-176-17-134.ngrok-free.app",
      // Или разрешить все хосты (менее безопасно)
      // 'all'
    ],
    // Также можно указать хост
    host: true, // Прослушивать все сетевые интерфейсы
  },
});
