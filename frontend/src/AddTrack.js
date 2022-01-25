import { useState } from 'react';
import { gql, useMutation, useQuery} from '@apollo/client';

const ADD_TRACK = gql`
  mutation PostWebEvent($stationId: Int!, $command: String!, $args: String!) {
    postWebEvent(input: { stationId: $stationId, command: $command, args: $args }) {
      result
      eventId
    }    
  }
`

export default function AddTrack({ stationId }) {
  const [addTrack, { data, loading, error }] = useMutation(ADD_TRACK, { 
    variables: { stationId, command: 'add' } 
  });

  const [url, setUrl] = useState('')
  const [eventIds, setEventIds] = useState([])

  async function submit(ev) {
    ev.preventDefault()
    console.log('submit', ev) 
    const result = await addTrack({ variables: { args: url } })
    console.log({ result })
    setEventIds([result.data.postWebEvent.eventId, ...eventIds])
    setUrl('')
  }

  function change(ev) {
    ev.preventDefault()
    setUrl(ev.currentTarget.value)
  }

  return (
    <div>
      <form onSubmit={submit}>
        <label>
          Paste url here to add music from youtube, bandcamp, twitter, tiktok, twitter, archive.org, etc:
          <br/>
          <input
            type="text"
            onChange={change}
            value={url}
            placeholder="https://"
            style={{ width: '300px' }}
          />
        </label>
        <button>add</button>
      </form>

      {eventIds.map(eventId => (
        <div key={eventId}>
          <Results eventId={eventId} />
        </div>
      ))}
    </div>
  )
}

const RESULTS = gql`
  query Results($eventId: Int!) {
    allResults(condition:{ eventId: $eventId }) {
      edges {
        node {
          id
          name
          data
          eventId
        }
      }
    }
  }
`

function Results({ eventId }) {
  const { data, error, loading } = useQuery(RESULTS, {
    variables: { eventId },
    pollInterval: 1000,
  })

  if (loading) return 'spinner';

  console.log({ resultsData: data })

  return (
    <ul>
      {data.allResults.edges.map(({ node }) => (
        <li key={node.id}>
          {node.data.message}
        </li>
      ))}
    </ul>
  )
}
