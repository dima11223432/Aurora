import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";

// https://vite.dev/config/
export default defineConfig({
  plugins: [react(), tailwindcss()],
  server: {
    allowedHosts: [
      "207d-213-176-17-134.ngrok-free.app",
      // Или разрешить все хосты (менее безопасно)
      // 'all'
    ],
    // Также можно указать хост
    host: true, // Прослушивать все сетевые интерфейсы
  },
});
