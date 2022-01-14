import logo from './banana.png';
import './App.css';
import { useQuery, gql } from '@apollo/client';
import { useEffect, useState } from 'react';

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
        <audio controls>
	  <source src="https://stream.djfullmoon.com/emb.ogg" type="audio/ogg" />
 	  <source src="https://stream.djfullmoon.com/emb.mp3" type="audio/mp3" />
        </audio>
        <CurrentTrack station="emb" />
      </header>
    </div>
  );
}

export default App;
