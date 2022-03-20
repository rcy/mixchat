import { gql, useMutation, useQuery, useSubscription } from '@apollo/client';
import { useState, useEffect, useRef } from 'react'
import useLocalStorage from 'react-localstorage-hook'
import { uniqueNamesGenerator, adjectives, animals } from 'unique-names-generator'
import Linkify from 'linkify-react';

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

function ChatInput({ onSubmit, target }) {
  const inputEl = useRef(null)
  const [input, setInput] = useState('')

  useEffect(() => inputEl.current.focus(), [])

  function change(ev) {
    setInput(ev.target.value)
  }

  function submit(ev) {
    ev.preventDefault()
    const submission = input.trim()
    if (submission.length) {
      onSubmit(submission)
      setInput('')
    }
  }

  return (
    <form onSubmit={submit}>
      <input
        type="text"
        onChange={change}
        value={input}
        ref={inputEl}
        placeholder={`Message ${target}`}
      />
    </form>
  )
}

export default function Chat({ stationId, stationSlug }) {
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
        <ColorTest />
      </header>

      <main style={{overflowY: 'scroll' }} ref={messagesEl} className="messages">
        {messages.map(({node}) => {
          const result = (
            <div key={node.id} className="message">
              <Time message={node} prevMessage={prevNode} /> {' '}
              <Nick message={node} prevMessage={prevNode}/> {' '}
              <Body message={node} />
            </div>
          )
          prevNode = node
          return result
        })}
      </main>

      <footer>
        {nick ? <ChatInput onSubmit={submit} target={stationSlug} /> : <SetNick onSubmit={setNick} />}
      </footer>
    </article>
  )
}

function Body({ message }) {
  return (
    <Linkify options={{target: "_blank"}}>
      {message.body}
    </Linkify>
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
  useEffect(() => { onSubmit(gennick()) }, [])
  return null
}

function SetNick3({ onSubmit }) {
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

// from https://sashamaps.net/docs/resources/20-colors/
const colors0 = ['#e6194B', '#3cb44b', '#ffe119', '#4363d8', '#f58231', '#911eb4', '#42d4f4', '#f032e6', '#bfef45', '#fabed4', '#469990', '#dcbeff', '#9A6324', '#fffac8', '#800000', '#aaffc3', '#808000', '#ffd8b1', '#000075', '#a9a9a9', '#ffffff', '#000000']
const colors = ['#e6194B', '#3cb44b', '#ffe119', '#4363d8', '#f58231', '#911eb4', '#42d4f4', '#f032e6', '#bfef45', '#fabed4', '#469990', '#9A6324', '#800000', '#000075', '#a9a9a9', '#000000']

function ColorTest() {
  return colors.map(c => <span key={c} style={{ color: c }}>{c} </span>)
}

function hashString(s, max) {
  let total = 0
  for (let i = 0; i < s.length; i++) {
    total += s.charCodeAt(i)
  }
  return max ? total % max : total
}

function Nick({ message, prevMessage }) {
  if (message.nick === prevMessage?.nick) {
    return <span>&nbsp;</span>;
  }
  const color = colors[hashString(message.nick, colors.length)]
  return (<span style={{ color }}><b>{message.nick}</b></span>)
}

function Time({ message, prevMessage }) {
  const md = new Date(message.createdAt)
  const pd = prevMessage && new Date(prevMessage?.createdAt)

  const mt = new Intl.DateTimeFormat("en", { timeStyle: 'short' }).format(md);
  const pt = pd && new Intl.DateTimeFormat("en", { timeStyle: 'short' }).format(pd);

  const newTime = mt !== pt
  const newNick = message.nick !== prevMessage?.nick

  const show = newNick //|| newTime

  return <span
           style={{
             opacity: '50%',
             //visibility: mt === pt ? 'hidden' : 'visible',
             visibility: show ? 'visible' : 'hidden',
             //fontSize: '30%',
           }}
         >{mt}</span>
}
