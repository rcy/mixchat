import { useEffect, useRef } from 'react';

export default function AudioControl({ stationSlug }) {
  const audioRef = useRef(null)
  const oggSourceRef = useRef(null)
  const mp3SourceRef = useRef(null)

  useEffect(() => {
    const oggStreamUrl = `${process.env.REACT_APP_ICECAST_URL}/${stationSlug}.ogg`
    const mp3StreamUrl = `${process.env.REACT_APP_ICECAST_URL}/${stationSlug}.mp3`
    oggSourceRef.current.src = oggStreamUrl
    mp3SourceRef.current.src = mp3StreamUrl
  }, [stationSlug])

  return (
    <audio controls ref={audioRef} autoPlay>
      <source ref={oggSourceRef} type="audio/mp3" />
      <source ref={mp3SourceRef} type="audio/ogg" />
    </audio>
  )
}
