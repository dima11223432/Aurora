import { useState } from 'react'
import reactLogo from './assets/react.svg'
import viteLogo from '/vite.svg'
import './App.css'

import { BrowserRouter, Routes, Route } from 'react-router-dom'
import HeroZone from './Herozone.jsx'
import Home from './Home.jsx'


function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<HeroZone />} />
         <Route path="/home" element={<Home />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App
