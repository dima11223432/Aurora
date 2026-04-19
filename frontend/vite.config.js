import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

import tailwindcss from "@tailwindcss/vite";
export default defineConfig({
  plugins: [react(), tailwindcss()],

  server: {
    allowedHosts: ["b566-2001-41d0-ab02-00-4-0-19.ngrok-free.app"],
    host: true,
  },
});
