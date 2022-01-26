import logo from './banana.png';
import './App.css';
import { Outlet, Link, Routes, Route } from "react-router-dom";
import StationPage from './StationPage.js';
import StationList from './StationList.js';

function App() {
  return (
    <div>
      <header>
        <h1>
          <Link to="/">DJFULLMOON</Link>
        </h1>
      </header>

      <Routes>
        <Route path="/" element={<StationList />} />
        <Route path=":slug" element={<StationPage />} />
      </Routes>
    </div>
  );
}

export default App;
