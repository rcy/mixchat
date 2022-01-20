import { useEffect, useRef, useState } from 'react';
import { useParams } from "react-router-dom";
import { useQuery, gql } from '@apollo/client';
import AudioControl from './AudioControl.js';
import CurrentTrack from './CurrentTrack.js';
import AddTrack from './AddTrack.js';

function StationPage() {
  const params = useParams()

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
      <hr/>

      <AudioControl stationSlug={station.slug} />

      <h3>Now Playing</h3><hr/>
      <CurrentTrack stationId={station.id} />

      <h3>Add Track</h3><hr/>
      <AddTrack stationId={station.id} />
    </div>
  )
}

export default StationPage;
