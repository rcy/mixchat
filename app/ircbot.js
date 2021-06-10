const PubSub = require('pubsub-js');
const irc = require('irc-upd')
const liquidsoap = require('./liquidsoap')
const youtubeDownload = require('./youtube.js')
const { pushRequest } = require('./source.js')
const { countListeners, fetchXspf } = require('./icecast.js')
const { formatDuration } = require('./util.js')

let nowPlayingData = {}

module.exports = function ircBot(host, nick, options) {
  const client = new irc.Client(host, nick, options)

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

  PubSub.subscribe('NOW', async function(msg, data) {
    console.log('RECV', msg, data)
    nowPlayingData = data

    const count = await countListeners()

    if (count > 0) {
      announceNowPlaying(client, options.channels[0])
    }
  })

  return client
}

function ifString(x) {
  return typeof x === 'string'
}

async function announceNowPlaying(client, to) {
  const { artist, album, title, duration } = nowPlayingData
  const count = await countListeners()

  const str = [
    `${count} listening to:`,
    [
      [artist, title].filter(ifString).join(', '),
      formatDuration(duration),
    ].filter(ifString).join(' - ') || "???"
  ].join(' ')

  client.say(to, str)
}

const handlers = {
  add: async ({ client, args, to, from }) => {
    const url = args

    //client.say(to, `${from}: ripping...`)

    const filename = await youtubeDownload(url)
    pushRequest(filename)

    client.say(to, `${from}: queued ${filename}`)
  },
  now: async ({ client, args, to, from }) => {
    await announceNowPlaying(client, to)
  },
  who: async({ client, args, to }) => {
    const numListeners = await countListeners()
    client.say(to, numListeners)
  },
  echo: async ({ client, args, to, from }) => {
    client.say(to, `${from}: ${args}`)
  },
  skip: async ({ client, args, to, from }) => {
    const result = await liquidsoap("dynlist.skip")

    client.say(to, `${from}: ${result}`)
  },
  help: async ({ client, args, to, from }) => {
    const commands = Object.keys(handlers).map(k => `!${k}`).join(' ')
    client.say(to, `${from}: ${commands}`)
  },
  src: async ({ client, args, to, from }) => {
    client.say(to, `${from}: https://github.com/rcy/djfullmoon`)
  },
}
