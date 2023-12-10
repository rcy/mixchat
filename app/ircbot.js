const PubSub = require('pubsub-js');
const irc = require('irc-upd')
//const liquidsoap = require('./liquidsoap')
//const youtubeDownload = require('./youtube.js')
//const { pushRequest } = require('./source.js')
const { countListeners } = require('./icecast.js')
const { formatDuration } = require('./util.js')

module.exports = function ircBot(host, nick, pgClient, options) {
  console.log('ircBot connecting', { host, nick, options })
  const client = new irc.Client(host, nick, options)

  client.addListener('error', function(message) {
    console.log('irc error: ', message);
  });

  client.addListener('join', function(channel, nick, message) {
    if (client.nick !== nick) {
      setTimeout(() => client.say(channel, `Hello ${nick}!`), 1000 * (10 + 50 * Math.random()))
    }
  });

  client.addListener('message', async function (from, to, message) {
    console.log(from + ' => ' + to + ': ' + message);

    const match = message.trim().match(/^!(.*)/);

    if (match) {
      const tokens = match[1].trim().split(/\s+/)
      const data = { via: 'irc', from, to, tokens, host, nick }
      const { rows: [{ station_id }] } = await pgClient.query("select station_id from irc_channels where channel = $1 and server = $2", [to, host]);

      await pgClient.query('insert into events (station_id, name, data) values ($1::integer, $2::text, $3::jsonb) returning *', [station_id, 'IRC_COMMAND', data])
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

      console.log({ row })

      if (row.result_name !== 'IRC_RESPONSE') {
        // could be WEB_RESPONSE, ignore for now.
        // Later broadcast this to irc as a web user adding track
        return
      }

      //const message = `${row.event_data.from}: ${row.event_id}.${row.result_id} ${JSON.stringify(row.result_data)}`
      let message = `${row.event_data.from}: ${row.result_data.message}`
      if (row.result_data.error) {
        message += ` id=${row.event_id} !error ${row.result_id}`
      }
      client.say(row.event_data.to, message)
    }
  })
  pgClient.query("LISTEN result");

  PubSub.subscribe('NOW', async function(msg, data) {
    try {
      console.log('RECV', msg, data.station, data.filename)
      const nowPlayingData = data
      const count = await countListeners(data.station)

      if (count > 0) {
        const { rows: [{ channel, id: stationId }] } = await pgClient.query(
          "select stations.id, channel from irc_channels join stations on stations.id = station_id where slug = $1",
          [data.station]
        );

        const { rows } = await pgClient.query(
          "select id from tracks where station_id = $1 and filename = $2",
          [stationId, data.filename]
        );
        if (rows.length) {
          const trackId = rows[0].id

          announceNowPlaying({ client, to: channel, count, nowPlayingData: {...nowPlayingData, trackId} })
        } else {
          client.say(channel, `playing system file: ${data.filename} (not in db)`);
        }
      }
    } catch(e) {
      console.error(e)
    }
  })

  return client
}

function ifString(x) {
  return typeof x === 'string'
}

async function announceNowPlaying({ client, to, count, nowPlayingData }) {
  const { artist, album, title, duration, station, trackId } = nowPlayingData

  const str = [
    `${count} listener${count === 1 ? "" : "s"}:`,
    [
      [artist, title].filter(ifString).join(', '),
      `[${trackId}]`,
      formatDuration(duration),
    ].filter(ifString).join(' - ') || "???"
  ].join(' ')

  client.say(to, str)
}
