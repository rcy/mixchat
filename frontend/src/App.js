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
          <a href="https://github.com/rcy/mixchat" target="_blank">github</a>
        </p>
        <h1>
          <Link to="/">MIXCHAT</Link>
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
