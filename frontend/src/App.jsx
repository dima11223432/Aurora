import { BrowserRouter, Routes, Route } from "react-router-dom";
import './App.css'

import HeroZone from './Herozone.jsx'

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<HeroZone />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App
