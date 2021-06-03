const PubSub = require('pubsub-js');
const irc = require('irc-upd')
const liquidsoap = require('./liquidsoap')
const youtubeDownload = require('./youtube.js')
const { pushRequest } = require('./source.js')
const { countListeners, fetchXspf } = require('./icecast.js')

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

  PubSub.subscribe('NOW', function(msg, data) {
    console.log('RECV', msg, data)
    const { artist, album, title } = data
    client.say(options.channels[0], `Now playing: ${artist} | ${album} | ${title}`)
  })

  return client
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
}
