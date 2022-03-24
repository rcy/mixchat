import { useEffect, useState } from 'react';
import { useQuery, gql } from '@apollo/client';
import Metadata from './Metadata.js';
import MetadataLink from './MetadataLink.js';

export default function RecentTracks({ stationId, count = 10 }) {
  const { loading, error, data } = useQuery(gql`
    query RecentlyPlayed($stationId: Int!, $count: Int!) {
      allTrackEvents(condition: { stationId: $stationId, action: "played"}, orderBy: CREATED_AT_DESC, first: $count) {
        edges {
          node {
            createdAt
            stationId
            id
            action
            trackByTrackId {
              id
              filename
              metadata
            }
          }
        }
      }
    }
  `, {
    pollInterval: 10000,
    variables: { stationId, count }
  });

  if (loading) {
    return 'spinner'
  }

  const { edges } = data.allTrackEvents

  if (edges.length) {
    return (
      <div>
        {edges.slice(1).map(({ node }) => (
          <div key={node.id}>
            <span className="track-list-item">
              <span className="timestamp">{new Intl.DateTimeFormat("en", { timeStyle: 'short' }).format(new Date(node.createdAt))}</span>
              {' '}
              <MetadataLink metadata={node.trackByTrackId?.metadata} />
              {' '}
              <Metadata metadata={node.trackByTrackId?.metadata} />
            </span>
          </div>
        ))}
      </div>
    )
  } else {
    return 'no edges'
  }
}
