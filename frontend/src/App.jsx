import React from "react";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import "./index.css";
import Analytics from "./Analytics";
import Herozone from "./Herozone";
import RecommendationFeed from "./RecommendationFeed";

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Herozone />} />
        <Route path="/analytics" element={<Analytics />} />
        <Route path="/feed" element={<RecommendationFeed />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
