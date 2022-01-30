import { useState } from 'react'

export default function Chat({ stationId }) {
  const [messages, setMessages] = useState([])
  const [input, setInput] = useState('')

  function addMessage({ nick, text, id }) {
    setMessages([...messages, { nick, text, id }])
  }

  function change(ev) {
    setInput(ev.target.value)
  }

  function submit(ev) {
    ev.preventDefault()
    addMessage({ text: input, nick: 'rcy', id: Math.random() })
    setInput('')
  }

  return (
    <article style={{ height: '100%' }}>
      <header>
      </header>

      <main style={{overflowY: 'scroll' }}>
        {messages.map(m => (
          <div key={m.id}><b>{m.nick}</b>: {m.text}</div>
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
