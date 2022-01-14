import logo from './banana.png';
//import './App.css';
import { useQuery, gql } from '@apollo/client';
import { useEffect, useState } from 'react';
import { Outlet, Link, Routes, Route } from "react-router-dom";
import StationPage from './StationPage.js';

function CurrentTrack({ station }) {
  const [track, setTrack] = useState()

  const { loading, error, data } = useQuery(gql`
    query RecentlyPlayed {
      allTrackEvents(condition: { stationId: 1, action: "played"}, orderBy: CREATED_AT_DESC, first: 5) {
        edges {
          node {
            createdAt
            stationId
            id
            action
            trackByTrackId {
              id
              filename
            }
          }
        }
      }
    }
  `, {
    pollInterval: 10000,
  });

  useEffect(() => {
    if (!loading) {
      setTrack(data.allTrackEvents.edges[0].node.trackByTrackId)
    }
  }, [data]);

  return <p>{track?.filename.replace(/^\/media\//,'').replace(/\.ogg$/,'')}</p>
}

function App() {
  return (
    <div className="App">
      <Link to="/emb">emb</Link> | <Link to="/tlm">tlm</Link>
      <Outlet />
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
        <h1>
          DJFULLMOON
        </h1>
        <p>
          <a
            className="App-link"
            href="https://web.libera.chat/#djfullmoon"
            target="_blank"
            rel="noopener noreferrer"
          >join the chat</a>
        </p>
        <Routes>
          <Route path=":slug" element={<StationPage />} />
        </Routes>

        <CurrentTrack station="emb" />
      </header>
    </div>
  );
}

export default App;
