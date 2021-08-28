const liquidsoap = require('../liquidsoap.js')
const youtubeDownload = require('../youtube.js')

module.exports = async ({ id }, helpers) => {
  const { rows } = await helpers.query("select * from events where id = $1", [id])

  const event = rows[0]

  helpers.logger.info(JSON.stringify(event))

  if (event.name === 'IRC_COMMAND') {
    const args = [...event.data.tokens];
    const command = args.shift().toLowerCase()

    // ================ dispatch to handlers here
    const handler = handlers[command]
    if (handler) {
      await handler(args, { event, insertResult, helpers })
    } else {
      await insertResult({ message: `Bad command !${command}. Type !help` })
    }
  }

  async function insertResult(data) {
    return await helpers.query("insert into results (event_id, name, data) values ($1::integer, $2::text, $3::jsonb) returning id", [event.id, 'IRC_RESPONSE', data])
  }
};

const handlers = {
  '?': async function(args, { helpers, insertResult }) {
    const { rows } = await helpers.query("select * from results where id = $1", [args[0]])
    if (rows[0]) {
      insertResult({ message: JSON.stringify(rows[0].data) })
    }
  },
  echo: async function(args, { insertResult }) {
    insertResult({ message: `${args.join(' ')}` })
  },
  skip: async function(args, { helpers, insertResult }) {
    try {
      const result = await liquidsoap("dynlist.skip")
      const { rows } = await helpers.query("insert into plays (action, track_id) values ('skipped', (select track_id from plays where action = 'played' order by created_at desc limit 1)) returning *")
      await insertResult({ result, rows })
    } catch(e) {
      await insertResult({ status: 'error', error: e })
    }
  },
  now: async function(args, { helpers, insertResult }) {
    const { rows } = await helpers.query("select tracks.id, filename, plays.created_at as started_at from plays join tracks on track_id = tracks.id where plays.action = 'played' order by plays.created_at DESC limit 1");
    const track = rows[0]
    if (track) {
      const message = `${track.id} ${track.filename}`
      await insertResult({ message })
    }
  },
  add: async function(args, { event, helpers, insertResult }) {
    //await insertResult({ status: 'adding' })

    const url = args.join(' ')
    console.log(`ripping ${url}`)

    let filename;

    // download
    try {
      filename = await youtubeDownload(url)
      //await insertResult({ status: 'added', filename })
    } catch(e) {
      console.error(e)
      await insertResult({ status: 'error', error: e, message: e.stderr.split('. ')[0] })
      return
    }

    // add track to db
    try {
      const { rows } = await helpers.query("insert into tracks (filename, event_id) values ($1::text, $2::integer) returning id", [filename, event.id]);
      const track_id = rows[0].id
      await insertResult({ filename, track_id, message: `queued track ${track_id} ${filename}` })
    } catch(e) {
      console.error(e)
      await insertResult({ status: 'error', message: e.message, error: e, code: e.code, detail: e.detail })
    }
  }
}
