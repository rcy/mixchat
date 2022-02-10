import { useEffect, useState } from 'react';
import { useQuery, gql } from '@apollo/client';
import Metadata from './Metadata.js';

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
        &gt;&gt;&nbsp;
      <b>
        <Metadata metadata={edges[0].node.trackByTrackId?.metadata} />
      </b>
      {edges.slice(1).map(({ node }) => (
        <div key={node.id}>
          <span className="track-list-item">
            <MetadataLink metadata={node.trackByTrackId?.metadata} />
            {' '}
            [{node.trackByTrackId.id}]
            {' '}
            <Metadata metadata={node.trackByTrackId?.metadata} />
            {' '}
          </span>
        </div>
      ))}
      </div>
    )
  } else {
    return 'no edges'
  }
}


function MetadataLink({ metadata }) {
  const link = metadata?.native?.vorbis?.find(x => x.id === "PURL")?.value
  if (!link) {
    return null
  }

  const url = new URL(link)

  // www.youtube.com -> youtube, etc
  const host = url.host.replace(/^(.+\.)?(.+)\..+$/, '$2')

  const shortcode = ({
    bandcamp: 'bc',
    soundcloud: 'sc',
    tiktok: 'tt',
    twitter: 'tw',
    youtube: 'yt',
    vimeo: 'vm'
  })[host] || '??';

  return (
    <a
      href={link}
      target="_blank"
    >{shortcode}</a>
  )
}
