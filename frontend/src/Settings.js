import { useState } from 'react'
import useLocalStorage from 'react-localstorage-hook'

export default function Settings({ stationId, postMessage }) {
  return (
    <div>
      <h1>Settings</h1>
      <em>sometimes things need to be changed</em>
      
      <section>
        <NickSettings stationId={stationId} postMessage={postMessage} />
      </section>
    </div>
  )
}

function NickSettings({ stationId, postMessage }) {
  const [nick, setNick] = useLocalStorage('nick', null)
  const [input, setInput] = useState('')

  function change(ev) {
    setInput(ev.target.value.trim())
  }

  async function submit(ev) {
    ev.preventDefault()
    if (input.length > 0) {
      await postMessage({ variables: { stationId, nick, body: `** is now known as ${input} **` } })
      setNick(input)
    }
  }

  return (
    <form onSubmit={submit} style={{ display: 'flex' }}>
      <label>
        Nickname:
        <input type="text" placeholder="new nick" defaultValue={nick} onChange={change}/>
      </label>
      {input && input !== nick && <button>Submit</button>}
    </form>
  )
}
