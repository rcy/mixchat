const PubSub = require('pubsub-js');
const irc = require('irc-upd')
//const liquidsoap = require('./liquidsoap')
//const youtubeDownload = require('./youtube.js')
//const { pushRequest } = require('./source.js')
const { countListeners } = require('./icecast.js')
const { formatDuration } = require('./util.js')

let nowPlayingData = {}

module.exports = function ircBot(host, nick, pgClient, options) {
  const client = new irc.Client(host, nick, options)

  client.addListener('error', function(message) {
    console.log('irc error: ', message);
  });

  client.addListener('message', async function (from, to, message) {
    console.log(from + ' => ' + to + ': ' + message);

    const match = message.trim().match(/^!(.*)/);

    if (match) {
      const tokens = match[1].trim().split(/\s+/)
      const data = { via: 'irc', from, to, tokens, host, nick }
      await pgClient.query('insert into events (name, data) values ($1::text, $2::jsonb) returning *', ['IRC_COMMAND', data])
    }
  });

  pgClient.on('notification', async function(msg) {
    console.log('NOTIFICATION', msg)
    if (msg.channel === 'result') {
      const res = await pgClient.query(`
select
   events.name as event_name,
   events.id as event_id,
   events.data as event_data,
   results.name as result_name,
   results.id as result_id,
   results.data as result_data
 from results 
 join events on results.event_id = events.id
 where results.id = $1
      `, [msg.payload])
      const row = res.rows[0];
      //const message = `${row.event_data.from}: ${row.event_id}.${row.result_id} ${JSON.stringify(row.result_data)}`
      let message = `${row.event_data.from}: ${row.result_data.message}`
      if (row.result_data.error) {
        message += ` !? ${row.result_id}`
      }
      client.say(row.event_data.to, message)
    }
  })
  pgClient.query("LISTEN result");

  PubSub.subscribe('NOW', async function(msg, data) {
    console.log('RECV', msg, data.station, data.filename)
    nowPlayingData = data

    const count = await countListeners(data.station)

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
  const { artist, album, title, duration, station } = nowPlayingData
  const count = await countListeners(station)

  const str = [
    `${station} ${count} listener${count === 1 ? "" : "s"}:`,
    [
      [artist, title].filter(ifString).join(', '),
      formatDuration(duration),
    ].filter(ifString).join(' - ') || "???"
  ].join(' ')

  client.say(to, str)
}
