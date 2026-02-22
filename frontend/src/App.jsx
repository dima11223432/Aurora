import { BrowserRouter, Routes, Route } from "react-router-dom";
import './styles/App.css'
import { Landing } from "./Authorisation.jsx";

import { BrowserRouter, Routes, Route } from 'react-router-dom'
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
