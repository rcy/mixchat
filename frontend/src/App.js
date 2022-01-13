import logo from './banana.png';
import './App.css';

function App() {
  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
        <p>
          dj Fu LL Moon
        </p>
        <a
          className="App-link"
          href="https://reactjs.org"
          target="_blank"
          rel="noopener noreferrer"
        >
        </a>
        <audio controls>
	  <source src="https://stream.djfullmoon.com/emb.ogg" type="audio/ogg" />
 	  <source src="https://stream.djfullmoon.com/emb.mp3" type="audio/mp3" />
        </audio>
        <br/>
        <a
          className="App-link"
          href="https://web.libera.chat/#djfullmoon"
          target="_blank"
          rel="noopener noreferrer"
        >join the chat!</a>
      </header>
    </div>
  );
}

export default App;
