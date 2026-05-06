import React from "react";
import ReactDOM from "react-dom/client";
import App from "./App";
import "./index.css";
import Footer from "./Footer";
import { TonConnectUIProvider } from "@tonconnect/ui-react";

ReactDOM.createRoot(document.getElementById("root")).render(
  <div>
    <TonConnectUIProvider manifestUrl="${routes.}/tonconnect-manifest.json">
      <App />
    </TonConnectUIProvider>
  </div>,
);
