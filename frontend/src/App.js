import logo from './banana.png';
import { Outlet, Link, Routes, Route } from "react-router-dom";
import StationPage from './StationPage.js';
import StationList from './StationList.js';
import CreateStation from './CreateStation.js';

function App() {
  return (
    <div>
      <header>
        <h1>
          DJFULLMOON
        </h1>
      </header>

      <div>
        <StationList />
      </div>
      <div>
        <CreateStation />
      </div>
      <div>
        <Routes>
          <Route path=":slug" element={<StationPage />} />
        </Routes>
      </div>
    </div>
  );
}

export default App;
