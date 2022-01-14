import { useEffect } from 'react';
import { useParams } from "react-router-dom";

function StationPage() {
  const params = useParams()

  useEffect(() => {
    console.log('station changed', params.slug)
  }, [params])
  
  return (
    <audio controls>
      <source src={`https://stream.djfullmoon.com/${params.slug}.ogg`} type="audio/ogg" />
      <source src={`https://stream.djfullmoon.com/${params.slug}.mp3`} type="audio/mp3" />
    </audio>
  )
}

export default StationPage;
