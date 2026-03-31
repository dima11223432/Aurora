import React from "react";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import "./index.css";
import Analytics from "./Analytics";
import Herozone from "./Herozone";

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Herozone />} />
        <Route path="/analytics" element={<Analytics />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
