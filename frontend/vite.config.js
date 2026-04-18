import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

import tailwindcss from "@tailwindcss/vite";
export default defineConfig({
  plugins: [react(), tailwindcss()],

  server: {
    allowedHosts: ["06d4-213-176-17-134.ngrok-free.app"],
    host: true,
  },
});
