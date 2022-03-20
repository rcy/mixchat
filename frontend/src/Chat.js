import { gql, useMutation, useQuery, useSubscription } from '@apollo/client';
import { useState, useEffect, useRef } from 'react'

const STATION_MESSAGES = gql`
  query StationMessages($stationId: Int!) {
    allMessages(condition: { stationId: $stationId }, last: 100) {
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

const STATION_MESSAGES_SUBSCRIPTION = gql`
  subscription TopicStationMessages($topic: String!) {
    listen(topic: $topic) {
      relatedNodeId
      relatedNode {
        ... on Message {
          id
          body
          nick
        }
      }
    }
  }
`;

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

function ChatInput({ onSubmit }) {
  const inputEl = useRef(null)
  const [input, setInput] = useState('')

  useEffect(() => inputEl.current.focus(), [])

  function change(ev) {
    setInput(ev.target.value)
  }

  function submit(ev) {
    ev.preventDefault()
    onSubmit(input)
    setInput('')
  }

  return (
    <form onSubmit={submit}>
      <input
        type="text"
        onChange={change}
        value={input}
        ref={inputEl}
      />
    </form>
  )
}

export default function Chat({ stationId }) {
  const messagesEl = useRef(null)

  const { data, loading, subscribeToMore } = useQuery(STATION_MESSAGES, { variables: { stationId } })
  const [ postMessage ] = useMutation(POST_STATION_MESSAGE, { variables: { stationId } })

  useEffect(() => {
    messagesEl?.current?.scrollTo(100000,100000)
  }, [data])
  
  useEffect(() => {
    subscribeToMore({
      document: STATION_MESSAGES_SUBSCRIPTION,
      variables: { topic: `station:${stationId}:messages` },
      updateQuery: (prev, { subscriptionData }) => {
        if (!subscriptionData.data.listen.relatedNode) return prev;
        const newNode = subscriptionData.data.listen.relatedNode;
        const prevEdges = prev.allMessages.edges
        return Object.assign({}, prev, {
          allMessages: {
            edges: [
              ...prevEdges,
              {
                ___typename: 'MessagesEdge',
                node: newNode,
              }
            ]
          }
        });
      }
    })
  }, [subscribeToMore])
  
  if (loading) {
    return null
  }

  const messages = data.allMessages.edges

  async function submit(input) {
    await postMessage({ variables: { nick: 'bob', body: input } })
  }

  return (
    <article style={{ height: '100%' }}>
      <header>
      </header>

      <main style={{overflowY: 'scroll' }} ref={messagesEl} className="messages">
        {messages.map(({node}) => (
          <div key={node.id}><b>{node.nick}</b>: {node.body}</div>
        ))}
      </main>

      <footer>
        <ChatInput onSubmit={submit} />
      </footer>
    </article>
  )
}
