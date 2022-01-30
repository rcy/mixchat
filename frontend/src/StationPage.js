import { useEffect, useRef, useState } from 'react';
import { useParams } from "react-router-dom";
import { useQuery, gql } from '@apollo/client';
import AudioControl from './AudioControl.js';
import RecentTracks from './RecentTracks.js';
import AddTrack from './AddTrack.js';
import Chat from './Chat.js';
import { Outlet, Link, Routes, Route } from "react-router-dom";

function StationPage() {
  const params = useParams()
  const [count, setCount] = useState(100)

  const { loading, error, data } = useQuery(gql`
    query StationBySlug($slug: String!) {
      stationBySlug(slug: $slug) {
        id
        slug
        ircChannelByStationId {
          id
          channel
        }
      }
    }`, { variables: { slug: params.slug } });

  if (error) {
    return error.message
  }

  if (loading) {
    return "spinner"
  }
  
  const station = data.stationBySlug

  const channel = station?.ircChannelByStationId?.channel

  return (
    <article style={{ height: '100%' }}>
      <div>
        <h1>{station.slug}</h1>
        <AudioControl stationSlug={station.slug} />

        <h3>Add Track</h3><hr/>
        <AddTrack stationId={station.id} />

        <h3>Chat</h3><hr/>
      </div>

      <main style={{ overflowY: 'hidden' }}>
        <Chat stationId={station.id} />
      </main>

      <footer>
      </footer>
    </article>
  )
}

export default StationPage;
