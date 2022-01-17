import { useEffect, useRef } from 'react';

export default function AudioControl({ stationSlug }) {
  const audioRef = useRef(null)

  useEffect(() => {
    // <source dataFormat="ogg" src={`https://stream.djfullmoon.com/${slug}.ogg`} type="audio/ogg" />
    // <source dataFormat="mp3" src={`https://stream.djfullmoon.com/${slug}.mp3`} type="audio/mp3" />
    audioRef.current.src = `https://stream.djfullmoon.com/${stationSlug}.ogg`
  }, [stationSlug])

  return (
    <audio controls ref={audioRef} autoPlay />
  )
}

