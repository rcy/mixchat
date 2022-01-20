import { useEffect, useState, useRef } from 'react';
import { gql, useMutation} from '@apollo/client';

const CREATE_STATION = gql`
  mutation CreateStation($name: String!, $slug: String!) {
    createStation(input: {station: {name: $name, slug: $slug}}) {
      clientMutationId
    }
  }
`

const STATION_SLUG_MAX_LENGTH=20;

export default function CreateStation() {
  const [createStation, { data, loading, error }] = useMutation(CREATE_STATION)
  const [state, setState] = useState('START')
  const [variables, setVariables] = useState({})

  async function submit(value) {
    const name = value.trim()
    const slug = name.toLowerCase().replace(/[^a-z0-9]/g, '').slice(0,STATION_SLUG_MAX_LENGTH)

    setVariables({ name, slug })
    setState('CONFIRM')
  }

  function cancel() {
    setVariables({})
    setState('START')
  }

  async function confirm() {
    const result = await createStation({ variables })
    console.log({ result })
    if (!error) {
      cancel()
    }
  }

  return (
    <div>
      {state === 'START' &&
       <Button
         onClick={() => setState('EDIT')}
       />}

      {state === 'EDIT' &&
       <Form
         defaultValue={variables.name || ''}
         onCancel={cancel}
         onSubmit={submit}
       />}

      {state === 'CONFIRM' &&
       <div>
         Create a new station named "{variables.name}" (with slug <code>{variables.slug}</code>)? {' '}
      <button onClick={confirm}>confirm</button> {' '}
      <a href="#" onClick={() => setState('EDIT')}>edit</a> {' '}
      <a href="#" onClick={cancel}>cancel</a>
       </div>
      }
    </div>
  )
}

function Form({ onCancel, onSubmit, defaultValue }) {
  const [input, setInput] = useState(defaultValue)
  const [slug, setSlug] = useState('')

  const inputEl = useRef(null);

  function handleCancel(ev) {
    ev.preventDefault()
    onCancel()
  }

  function submit(ev) {
    ev.preventDefault()
    onSubmit(input)
  }

  function keyup(ev) {
    if (ev.key === 'Escape') {
      onCancel()
    }
  }

  useEffect(() => inputEl.current.focus(), [])

  return (
    <div>
      <form onSubmit={submit}>
        <input
          ref={inputEl}
          type="text"
          placeholder="station name"
          value={input}
          onChange={(ev) => setInput(ev.currentTarget.value)}
          onKeyUp={keyup}
        />
        <button
          disabled={!input?.trim().length}
        >create</button>
        {' '}
        <a href="" onClick={handleCancel}>cancel</a>
      </form>
    </div>
  )
}

function Button({ onClick }) {
  return <button onClick={onClick}>create station</button>
}
