import { useEffect, useRef, useState } from 'react';
import { useParams } from "react-router-dom";
import { useQuery, gql } from '@apollo/client';
import AudioControl from './AudioControl.js';
import CurrentTrack from './CurrentTrack.js';

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
      <h2>{station.slug}</h2>
      <p>
        <a
          href={`https://web.libera.chat/${channel}`}
          target="_blank"
          rel="noopener noreferrer"
        >join the {channel} chat</a>
      </p>
      <AudioControl stationSlug={station.slug} />
      <CurrentTrack stationId={station.id} />
    </div>
  )
}

export default StationPage;
