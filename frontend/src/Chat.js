import { gql, useMutation, useQuery} from '@apollo/client';
import { useState } from 'react'

const STATION_MESSAGES = gql`
  query StationMessages($stationId: Int!) {
    allMessages(condition: { stationId: $stationId }) {
      edges {
        node {
          id
	  body
          nick
        }
      }
    }
  }
`

const POST_STATION_MESSAGE = gql`
  mutation PostStationMessage($stationId: Int!, $body: String!, $nick: String!) {
    createMessage(
      input: { message: {stationId: $stationId, body: $body, nick: $nick}}
    ) {
      message {
        id
      }
    }
  }
`

export default function Chat({ stationId }) {
  const [input, setInput] = useState('')

  const { data, loading, error } = useQuery(STATION_MESSAGES, { variables: { stationId } })
  const [ postMessage ] = useMutation(POST_STATION_MESSAGE, { variables: { stationId } })

  if (loading) {
    return null
  }

  const messages = data.allMessages.edges
  console.log(messages)

  function change(ev) {
    setInput(ev.target.value)
  }

  async function submit(ev) {
    ev.preventDefault()
    console.log('submit')
    const result = await postMessage({ variables: { nick: 'bob', body: input } })
    console.log('post result', result)
    setInput('')
  }

  return (
    <article style={{ height: '100%' }}>
      <header>
      </header>

      <main style={{overflowY: 'scroll' }}>
        {messages.map(({node}) => (
          <div key={node.id}><b>{node.nick}</b>: {node.body}</div>
        ))}
      </main>

      <footer>
        <form onSubmit={submit}>
          <input
            type="text"
            onChange={change}
            value={input}
          />
        </form>
      </footer>
    </article>
  )
}
