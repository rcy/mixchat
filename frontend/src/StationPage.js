import { useEffect, useRef, useState } from 'react';
import { useParams, useSearchParams } from "react-router-dom";
import { useQuery, useMutation, gql } from '@apollo/client';
import AudioControl from './AudioControl.js';
import RecentPlayedTracks from './RecentPlayedTracks.js';
import RecentAddedTracks from './RecentAddedTracks.js';
import NowPlaying from './NowPlaying.js';
import Settings from './Settings.js';
import AddTrack from './AddTrack.js';
import Chat from './Chat.js';
import { Outlet, Routes, Route } from "react-router-dom";
import { isMobile } from 'react-device-detect';
import Link from './Link'

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

function StationPage() {
  const [ postMessage ] = useMutation(POST_STATION_MESSAGE)
  const params = useParams()
  const [search] = useSearchParams()

  const [count, setCount] = useState(100)

  const { loading, error, data } = useQuery(gql`
    query StationBySlug($slug: String!) {
      stationBySlug(slug: $slug) {
        id
        slug
        ircChannelByStationId {
          id
          channel
        }
      }
    }`, { variables: { slug: params.slug } });

  const [tab,setTab] = useState('chat')

  if (error) {
    return error.message
  }

  if (loading) {
    return "spinner"
  }
  
  const showAudio = search.get('noaudio') !== '1'

  const station = data.stationBySlug

  const channel = station?.ircChannelByStationId?.channel

  // return (
  //   <article style={{ height: '100%' }}>
  //     <div>
  //       <h1>{station.slug}</h1>
  //       <AudioControl stationSlug={station.slug} />
  // 
  //       <h3>Add Track</h3><hr/>
  //       <AddTrack stationId={station.id} />
  // 
  //       <h3>Chat</h3><hr/>
  //     </div>
  // 
  //     <main style={{ overflowY: 'hidden' }}>
  //       <Chat stationId={station.id} />
  //     </main>
  // 
  //     <footer>
  //     </footer>
  //   </article>
  // )

  function clickTab(ev) {
    ev.preventDefault();
    setTab(ev.target.href.split('#')[1])
  }
  
  return (
    <article style={{ height: '100%' }}>
      <div>
        <div>
          {showAudio && <AudioControl stationSlug={station.slug} />}
        </div>

        {/* <h3 className="banner">TASTESLIKEME NEXT SHOW: ((( Interesting Music ))) Saturday 8PM</h3> */}
        <div className="now-playing">Now Playing: <NowPlaying stationId={station.id} /></div>
        <div className="menubar">
          <a href="#chat" onClick={clickTab}>chat</a>
          <a href="#library" onClick={clickTab}>library</a>
          <a href="#settings" onClick={clickTab}>settings</a>
        </div>
      </div>
      <main style={{ overflowY: 'hidden' }}>
        <Tab active={tab} id="chat">
          <Chat stationId={station.id} stationSlug={station.slug} postMessage={postMessage} />
        </Tab>
        <Tab active={tab} id="library">
          <article style={{height: '100%' }}>
            <div></div>
            <div style={{ overflowY: 'scroll' }}>
              <h1>Recently Played</h1>
              <RecentPlayedTracks stationId={station.id} count={25} />

              <h1>Recently Added</h1>
              <RecentAddedTracks stationId={station.id} count={25} />

              <h1>Add track to library</h1>
              <AddTrack stationId={station.id} />
            </div>
            <div></div>
          </article>
        </Tab>
        <Tab active={tab} id="add">
          <article style={{height: '100%' }}>
            <div></div>
            <div style={{ overflowY: 'scroll' }}>
              <AddTrack stationId={station.id} />
            </div>
            <div></div>
          </article>
        </Tab>
        <Tab active={tab} id="settings">
          <article style={{height: '100%' }}>
            <div></div>
            <div style={{ overflowY: 'scroll' }}>
              <Settings stationId={station.id} postMessage={postMessage} />
            </div>
            <div></div>
          </article>
        </Tab>
      </main>
      <footer>
        {isMobile || <div style={{height: '50px'}}/>}
      </footer>
    </article>
  )
}

function Tab({ children, active, id }) {
  return <div style={{display: active === id ? 'inline' : 'none'}}>{children}</div>
}

export default StationPage;
