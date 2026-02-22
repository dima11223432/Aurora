import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";

// https://vite.dev/config/
export default defineConfig({
  plugins: [react(), tailwindcss()],
  server: {
    allowedHosts: [
      "0920-2605-e440-9-00-3a.ngrok-free.app",
      // Или разрешить все хосты (менее безопасно)
      // 'all'
    ],
    // Также можно указать хост
    host: true, // Прослушивать все сетевые интерфейсы
  },
});
