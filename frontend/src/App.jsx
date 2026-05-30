import React from "react";
import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import "./index.css";
import Analytics from "./Analytics";
import ErrorPage from "./ErrorPage";
import Herozone from "./Herozone";
import RecommendationFeed from "./RecommendationFeed";
import AdminPanel from "./AdminPanel";
import Settings from "./Settings";

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Herozone />} />
        <Route path="/analytics" element={<Analytics />} />
        <Route path="/feed" element={<RecommendationFeed />} />
        <Route path="/settings" element={<Settings />}></Route>
        <Route path="/admin" element={<AdminPanel />} />
        <Route path="/404" element={<ErrorPage />}></Route>
        <Route path="*" element={<Navigate to="/404" replace />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
