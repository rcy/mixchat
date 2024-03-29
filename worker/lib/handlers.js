const { execSync } = require('child_process')
const liquidsoap = require('../liquidsoap.js')
const youtubeDownload = require('../youtube.js')

const handlers = {
  'help': async function(args, { helpers, insertResult }) {
    const commands = Object.keys(handlers).map(k => `!${k}`).join(' ')
    await insertResult({ message: commands })
  },
  'error': async function(args, { helpers, insertResult }) {
    const { rows } = await helpers.query("select * from results where id = $1", [args[0]])
    if (rows[0]) {
      insertResult({ message: JSON.stringify(rows[0].data) })
    }
  },
  echo: async function(args, { insertResult }) {
    insertResult({ message: `${args.join(' ')}` })
  },
  skip: async function(args, { event, helpers, insertResult }) {
    try {
      const { rows: [{ slug }] } = await helpers.query("select slug from stations where id = $1", [event.station_id]);

      const result = await liquidsoap(`dynlist_${slug}.skip`)
      if (result[0] === 'Skipped!' && result[1] === 'END') {
        await helpers.query("insert into track_events (station_id, action, track_id) values ($1, 'skipped', (select track_id from plays where action = 'played' order by created_at desc limit 1))", [event.station_id])
      } else {
        await insertResult({ error: result, message: 'Skip failed' })
      }
    } catch(e) {
      await insertResult({ error: e, message: e.message })
    }
  },
  yeet: async function(args, { event, helpers, insertResult }) {
    const trackId = args[0]
    const { rows } = await helpers.query("select * from tracks where station_id = $1 and id = $2", [event.station_id, trackId])
    await helpers.query("insert into track_events (station_id, action, track_id) values ($1, 'yeeted', $2)", [event.station_id, trackId])
    await insertResult({ message: `yeeted ${trackId} ${rows[0].filename} into oblivion`})
  },
  now: async function(args, { event, helpers, insertResult }) {
    const { rows } = await helpers.query(`
select tracks.id, filename, plays.created_at as started_at
from plays
join tracks on track_id = tracks.id
where station_id = $1
  and plays.action = 'played'
order by plays.created_at DESC limit 1
    `, [event.station_id]);
    const track = rows[0]
    if (track) {
      const message = `${track.id} ${track.filename}`
      await insertResult({ message })
    } else {
      await insertResult({ message: 'something outside the collection is playing' })
    }
  },
  add: async function(args, { event, helpers, insertResult }) {
    await insertResult({ message: 'adding track...' })

    const url = args.join(' ')
    console.log(`ripping ${url}`)

    let filename;

    // download
    try {
      filename = await youtubeDownload(url)
      //await insertResult({ status: 'added', filename })
    } catch(e) {
      console.error(e)

      await insertResult({ status: 'error', error: e, message: e.message || "something bad happened, will retry"})

      throw e
    }

    // add track to db
    try {
      const { rows } = await helpers.query("insert into tracks (station_id, filename, event_id, bucket) values ($1::integer, $2::text, $3::integer, current_bucket($1::integer)) returning id", [event.station_id, filename, event.id]);
      const track_id = rows[0].id
      await insertResult({ filename, track_id, message: `queued track ${track_id} ${filename}` })
    } catch(e) {
      console.error(e)
      await insertResult({ status: 'error', message: e.message, error: e, code: e.code, detail: e.detail })
    }
  },
  sleep: async function(args, { insertResult }) {
    const ms = +args;

    await new Promise(resolve => setTimeout(resolve, ms))

    insertResult({ message: `slept for ${ms} milliseconds` })
  },
  pom: async function(_args, { insertResult }) {
    try {
      const message = execSync('/usr/games/pom').toString()
      insertResult({ message })
    } catch(e) {
      insertResult({ status: 'error', message: e.message, error: e, code: e.code, detail: e.detail })
    }
  },
  version: async function(_args, { insertResult }) {
    try {
      const message = execSync('/usr/local/bin/yt-dlp --version').toString()
      insertResult({ message })
    } catch(e) {
      insertResult({ status: 'error', message: e.message, error: e, code: e.code, detail: e.detail })
    }
  },
  update: async function(_args, { insertResult }) {
    try {
      const message = execSync('/usr/local/bin/yt-dlp -U').toString()
      insertResult({ message })
    } catch(e) {
      insertResult({ status: 'error', message: e.message, error: e, code: e.code, detail: e.detail })
    }
  },
}

module.exports = {
  handlers
}
