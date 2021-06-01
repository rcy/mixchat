const irc = require('irc-upd')
const liquidsoap = require('./liquidsoap')
const youtubeDownload = require('./youtube.js')
const { pushRequest } = require('./source.js')

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

  return client
}

const handlers = {
  add: async ({ client, args, to, from }) => {
    const url = args

    client.say(to, `${from}: ripping ${url}...`)

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
  skip: async ({ client, args, to, from }) => {
    const result = await liquidsoap("dynlist.skip")

    client.say(to, `${from}: ${result}`)
  },
  help: async ({ client, args, to, from }) => {
    const commands = Object.keys(handlers).map(k => `!${k}`).join(' ')
    client.say(to, `${from}: ${commands}`)
  },
}

const fetch = require('node-fetch')
const convert = require('xml-js');
async function fetchXspf() {
  const raw = await fetch('http://radio.nonzerosoftware.com:8000/emb.ogg.xspf')
  const text = await raw.text()
  return JSON.parse(convert.xml2json(text, { arrayNotation: true }))
}
