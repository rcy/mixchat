const express = require('express')
const bodyParser = require('body-parser')
const net = require('net');
const youtubedl = require('youtube-dl-exec')
const irc = require('irc')

const app = express()
const port = 3010

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

async function liquidsoap(command, host = '127.0.0.1', port = 1234) {
  const client = new net.Socket();

  return new Promise((resolve, reject) => {
    client.connect(port, host, function() {
      client.write(command + '\n');

      client.on('data', function(data) {
        resolve(data)
        client.write('quit\n')
        client.destroy();
      });

      client.on('close', function() {
        console.log('liquidsoap telnet connection closed');
      });
    })
  })
}

async function download(id) {
  const output = await youtubedl(`https://www.youtube.com/watch?v=${id}`, {
    extractAudio: true,
    audioFormat: 'vorbis',
    //dumpSingleJson: true,
    // noWarnings: true,
    noCallHome: true,
    // noCheckCertificate: true,
    // preferFreeFormats: true,
    youtubeSkipDashManifest: true,
    output: '/tmp/%(id)s.%(ext)s',
    //    referer: 'https://example.com',
  })
  console.log(output)
  return `/tmp/${id}.ogg`
}

const client = new irc.Client('irc.freenode.net', 'mcm00n', { channels: ["#emb-radio"] })

client.addListener('message', async function (from, to, message) {
  console.log(from + ' => ' + to + ': ' + message);

  const match = message.match(/!request (.+)/)
  if (match) {
    const id = match[1]

    try {
      client.say(to, `${id}: downloading...`)
      const filename = await download(id)
      client.say(to, `${id}: downloading...done`)

      client.say(to, `${id}: requesting ${filename}...`)
      const data = await liquidsoap(`request.push ${filename}`)
      client.say(to, `${id}: requesting ${filename}...done`)
    } catch(e) {
      return client.say(to, `${id}: ${e.message}`)
    }
  } else {
    console.log('unhandled message', message)
  }
});
