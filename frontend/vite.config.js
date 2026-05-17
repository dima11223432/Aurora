import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

import tailwindcss from "@tailwindcss/vite";
export default defineConfig({
  plugins: [react(), tailwindcss()],

  server: {
    allowedHosts: ["0745-2a0d-d904-2-d5-00-2.ngrok-free.app"],
    host: true,
  },
});
