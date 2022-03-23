import { useEffect, useState } from 'react';
import { useQuery, gql } from '@apollo/client';
import Metadata from './Metadata.js';

export default function NowPlaying({ stationId }) {
  const { loading, error, data } = useQuery(gql`
    query NowPlaying($stationId: Int!) {
      allTrackEvents(condition: { stationId: $stationId, action: "played"}, orderBy: CREATED_AT_DESC, first: 1) {
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
    variables: { stationId }
  });

  if (loading) {
    return '---'
  }

  const { edges } = data.allTrackEvents

  if (edges.length) {
    return (
      <Metadata metadata={edges[0].node.trackByTrackId?.metadata} />
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
