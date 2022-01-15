import logo from './banana.png';
import './App.css';
import { Outlet, Link, Routes, Route } from "react-router-dom";
import StationPage from './StationPage.js';

function App() {
  return (
    <div className="App">
      <header>
        <h1>
          DJFULLMOON
        </h1>
      </header>

      <Routes>
        <Route path=":slug" element={<StationPage />} />
      </Routes>

      <Link to="/emb">emb</Link> | <Link to="/tlm">tlm</Link>

    </div>
  );
}

export default App;
