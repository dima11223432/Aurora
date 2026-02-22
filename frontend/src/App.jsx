import { BrowserRouter, Routes, Route } from "react-router-dom";
import './styles/App.css'
import Landing from "./Authorisation";
function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Landing/>} />
      </Routes>
    </BrowserRouter>
  );
}

export default App
