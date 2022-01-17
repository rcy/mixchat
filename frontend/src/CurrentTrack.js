import { useEffect, useState } from 'react';
import { useQuery, gql } from '@apollo/client';

export default function CurrentTrack({ stationId }) {
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
