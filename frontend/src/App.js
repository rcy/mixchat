import logo from './banana.png';
import './App.css';
import { Outlet, Link, Routes, Route } from "react-router-dom";
import StationPage from './StationPage.js';
import StationList from './StationList.js';

function App() {
  return (
    <div>
      <header>
        <p style={{ float: 'right' }}>
          <a href="https://twitter.com/rcyeske" target="_blank">tw</a>
          {'/'}
          <a href="https://github.com/rcy/djfullmoon" target="_blank">gh</a>
        </p>
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
