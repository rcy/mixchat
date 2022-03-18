import logo from './banana.png';
import './App.css';
import { Link, Routes, Route } from "react-router-dom";
import StationPage from './StationPage.js';
import StationList from './StationList.js';

function App() {
  return (
    <article style={{ height: '100%' }}>
      <header>
        <h1>
          <Link to="/">mixchat</Link>
        </h1>
      </header>

      <main style={{ overflowY: 'hidden' }}>
        <Routes>
          <Route path="/" element={<StationList />} />
          <Route path=":slug/*" element={<StationPage />} />
        </Routes>
      </main>

      <footer>
      </footer>
    </article>
  );
}

export default App;
