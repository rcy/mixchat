import { useEffect, useRef, useState } from 'react';
import { useParams } from "react-router-dom";
import { useQuery, gql } from '@apollo/client';
import AudioControl from './AudioControl.js';
import RecentTracks from './RecentTracks.js';
import AddTrack from './AddTrack.js';

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

  if (loading) {
    return "spinner"
  }
  
  const station = data.stationBySlug

  const channel = station?.ircChannelByStationId?.channel

  return (
    <div>
      <h1>{station.slug}</h1>
      {channel &&
       <p>
         <a
           href={`https://web.libera.chat/${channel}`}
           target="_blank"
           rel="noopener noreferrer"
         >join the {channel} chat</a>
       </p>}
      <hr/>

      <AudioControl stationSlug={station.slug} />
      <h3>Add Track</h3><hr/>
      <AddTrack stationId={station.id} />

      <h3>Last {count} Tracks</h3><hr/>
      <RecentTracks stationId={station.id} count={count} />
    </div>
  )
}

export default StationPage;
