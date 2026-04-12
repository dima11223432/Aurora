<<<<<<< HEAD
import React from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import './index.css';
import Analytics from './Analytics';
import Herozone from './Herozone';
import News from './News';
=======
import React from "react";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import "./index.css";
import Analytics from "./Analytics";
import Herozone from "./Herozone";
>>>>>>> dd4ed17a23ab31af3c1a81b040fbadfef0b9d7e3

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Herozone />} />
<<<<<<< HEAD
        <Route path="/analytics" element={<Analytics/>}/>
        <Route path="/news" element={<News/>}/>
=======
        <Route path="/analytics" element={<Analytics />} />
>>>>>>> dd4ed17a23ab31af3c1a81b040fbadfef0b9d7e3
      </Routes>
    </BrowserRouter>
  );
}

export default App;
