import { gql, useMutation, useQuery, useSubscription } from '@apollo/client';
import { useState, useEffect, useRef } from 'react'
import useLocalStorage from 'react-localstorage-hook'
import { uniqueNamesGenerator, adjectives, colors, animals, names, languages } from 'unique-names-generator'

const STATION_MESSAGES = gql`
  query StationMessages($stationId: Int!) {
    allMessages(condition: { stationId: $stationId }, last: 100) {
      edges {
        node {
          id
	  body
          nick
          createdAt
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
          createdAt
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
        placeholder='Type your message here...'
      />
    </form>
  )
}

export default function Chat({ stationId }) {
  const messagesEl = useRef(null)

  const { data, loading, subscribeToMore } = useQuery(STATION_MESSAGES, { variables: { stationId } })
  const [ postMessage ] = useMutation(POST_STATION_MESSAGE, { variables: { stationId } })
  const [nick, setNick] = useLocalStorage('nick', null)

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
    await postMessage({ variables: { nick, body: input } })
  }

  let prevNode = null

  return (
    <article style={{ height: '100%' }}>
      <header>
      </header>

      <main style={{overflowY: 'scroll' }} ref={messagesEl} className="messages">
        {messages.map(({node}) => {
          const result = (
            <div key={node.id}>
              <Time message={node} prevMessage={prevNode} /> <b>{node.nick}</b>: {node.body}
            </div>
          )
          prevNode = node
          return result
        })}
      </main>

      <footer>
        {nick ? <ChatInput onSubmit={submit} /> : <SetNick onSubmit={setNick} />}
      </footer>
    </article>
  )
}

const shortAdjectives = adjectives.filter(a => a.length <= 5)
const shortAnimals = animals.filter(a => a.length <= 5)

console.log(shortAnimals.length * shortAdjectives.length)

function gennick() {
  const nick = uniqueNamesGenerator({
    dictionaries: [shortAdjectives, shortAnimals],
    length: 2,
    separator: '',
    style: 'lowercase'
  })
  return nick
}

function SetNick({ onSubmit }) {
  const [nicks, setNicks] = useState([])

  function shuffle(ev) {
    ev?.preventDefault()
    setNicks([gennick(), gennick(), gennick()])
  }

  useEffect(shuffle, [])

  return (
    <div>
      <b>To start chatting, choose your nickname: </b>
      {
        nicks.map(n => (
          <span key={n}>
            <a href="" onClick={(ev) => { ev.preventDefault(); onSubmit(n) }}>{n}</a>
            {' '}
          </span>
        ))
      }
      <a onClick={shuffle} href="">[more]</a>
    </div>
  )
}

function Time({ message, prevMessage }) {
  const md = new Date(message.createdAt)
  const pd = prevMessage && new Date(prevMessage?.createdAt)

  const mt = new Intl.DateTimeFormat("en", { timeStyle: 'short' }).format(md);
  const pt = pd && new Intl.DateTimeFormat("en", { timeStyle: 'short' }).format(pd);

  return <span
           style={{
             opacity: '50%',
             visibility: mt === pt ? 'hidden' : 'visible',
             fontSize: '30%',
           }}
         >{mt}</span>
}
