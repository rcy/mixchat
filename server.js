const express = require('express')
const bodyParser = require('body-parser')
const net = require('net');
const youtubedl = require('youtube-dl-exec')
const irc = require('irc')
const fetch = require('node-fetch')
const convert = require('xml-js');

const app = express()
const port = 3010

app.get('/client', (req, res) => {
  console.log(client)
  res.sendStatus(200)
})

app.get('/', (req, res) => {
  res.send('Hello World!')
})

var jsonParser = bodyParser.json()

app.post('/youtube', jsonParser, async (req, res) => {
  try {
    const { filename, data } = requestYoutube(req.body.id)
    res.status(200).json({ filename, data })
  } catch(e) {
    console.error(e)
    res.sendStatus(500, { error: e.message })
  }
})

async function requestYoutube(id) {
  const filename = await download(id)
  const data = await liquidsoap(`request.push ${filename}`)
  return { filename, data }
}

app.listen(port, () => {
  console.log(`Example app listening at http://localhost:${port}`)
})

async function liquidsoap(command, host = '172.17.0.1', port = 1234) {
  const client = new net.Socket();

  return new Promise((resolve, reject) => {
    client.connect(port, host, function() {
      client.write(command + '\n');
    })

    client.on('data', function(data) {
      resolve(data.toString())
      client.write('quit\n')
      client.destroy();
    });

    client.on('close', function() {
      console.log('liquidsoap telnet connection closed');
    });

    client.on('error', function(e) {
      console.error(e)
      reject(e)
    });
  })
}

async function download(id) {
  const output = await youtubedl(`https://www.youtube.com/watch?v=${id}`, {
    quiet: true,
    extractAudio: true,
    audioFormat: 'vorbis',
    //dumpSingleJson: true,
    // noWarnings: true,
    noCallHome: true,
    // noCheckCertificate: true,
    // preferFreeFormats: true,
    youtubeSkipDashManifest: true,
    //output: '/media/%(id)s.%(ext)s',
    //    referer: 'https://example.com',
    addMetadata: true,
    restrictFilenames: true,
    exec: "mv {} /media && echo {}", // output
  })
  console.log(output)
  return `/media/${output}`
}

const freenodeClient = new irc.Client('irc.freenode.net', 'djfullmoon', {
  channels: ["#emb-radio"],
  debug: true,
  sasl: true,
  userName: 'djfullmoon',
  password: 'JJyf376fGgbPnfcz9',
})

addIrcClientListeners(freenodeClient)

const liberaClient = new irc.Client('irc.libera.chat', 'djfullmoon', {
  port: 6697,
  secure: true,
  channels: ["#emb-radio"],
  debug: true,
  //  sasl: true,
  userName: 'djfullmoon',
  password: '6MdYaHHKBpSzEwLAa',
})

addIrcClientListeners(liberaClient)

function addIrcClientListeners(client) {
  client.addListener('error', function(message) {
    console.log('irc error: ', message);
  });

  client.addListener('message', async function (from, to, message) {
    console.log(from + ' => ' + to + ': ' + message);

    const match = message.trim().match(/^!(\w+)\s*(.*)/);

    if (match) {
      const command = match[1].toLowerCase()
      const args = match[2]

      const handler = handlers[command] || handlers.help
      if (handler) {
        try {
          await handler({ client, args, from, to, message })
        } catch(e) {
          client.say(to, `${from}: ${message}: ${e.message}`)
        }
      }
    }
  });
}

const handlers = {
  add: async ({ client, args, to, from }) => {
    const id = args

    client.say(to, `${from}: looking for ${id}...`)

    const filename = await download(id)
    const data = await liquidsoap(`request.push ${filename}`)

    client.say(to, `${from}: ${filename} requested.`)
  },
  now: async ({ client, args, to, from }) => {
    const xspf = await fetchXspf()

    const result =
      xspf.elements[0].elements
          .find(e => e.name === 'trackList')
          .elements[0].elements
          .filter(e => e.name === 'creator' || e.name === 'title')
          .map(e => e.elements[0].text)
          .join(' | ');

    client.say(to, `Now playing: ${result}`)
  },
  who: async({ client, args, to }) => {
    const xspf = await fetchXspf()

    const result =
      xspf.elements[0].elements
          .find(e => e.name === 'trackList')
          .elements[0].elements
          .find(e => e.name === 'annotation')
          .elements[0].text
          .split('\n')
          .find(e => e.match('Current Listeners'))

    client.say(to, result)
  },
  echo: async ({ client, args, to, from }) => {
    client.say(to, `${from}: ${args}`)
  },
  help: async ({ client, args, to, from }) => {
    const commands = Object.keys(handlers).map(k => `!${k}`).join(' ')
    client.say(to, `${from}: ${commands}`)
  }
}

async function fetchXspf() {
  const raw = await fetch('http://radio.nonzerosoftware.com:8000/emb.ogg.xspf')
  const text = await raw.text()
  return JSON.parse(convert.xml2json(text, { arrayNotation: true }))
}
