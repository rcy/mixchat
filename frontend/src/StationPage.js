import { useEffect, useRef, useState } from 'react';
import { useParams } from "react-router-dom";
import { useQuery, gql } from '@apollo/client';

function AudioControl({ stationSlug }) {
  const audioRef = useRef(null)

  useEffect(() => {
    // <source dataFormat="ogg" src={`https://stream.djfullmoon.com/${slug}.ogg`} type="audio/ogg" />
    // <source dataFormat="mp3" src={`https://stream.djfullmoon.com/${slug}.mp3`} type="audio/mp3" />
    audioRef.current.src = `https://stream.djfullmoon.com/${stationSlug}.ogg`
  }, [stationSlug])

  return (
    <audio controls ref={audioRef} autoPlay />
  )
}

function CurrentTrack({ stationId }) {
  const [track, setTrack] = useState()

  const { loading, error, data } = useQuery(gql`
    query RecentlyPlayed($stationId: Int!) {
      allTrackEvents(condition: { stationId: $stationId, action: "played"}, orderBy: CREATED_AT_DESC, first: 5) {
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
    variables: { stationId }
  });

  useEffect(() => {
    if (!loading && data) {
      setTrack(data.allTrackEvents.edges[0].node.trackByTrackId)
    }
  }, [stationId, data, loading]);

  return <p>{track?.filename.replace(/^\/media\//,'').replace(/\.ogg$/,'')}</p>
}


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
      <AudioControl stationSlug={station.slug} />
      <CurrentTrack stationId={station.id} />
      <p>
        <a
          className="App-link"
          href={`https://web.libera.chat/${channel}`}
          target="_blank"
          rel="noopener noreferrer"
        >join the {channel} chat</a>
      </p>
    </div>
  )
}

export default StationPage;
