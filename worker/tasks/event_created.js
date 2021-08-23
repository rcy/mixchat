const youtubeDownload = require('../youtube.js')

module.exports = async ({ id }, helpers) => {
  const { rows } = await helpers.query("select * from events where id = $1", [id])

  const event = rows[0]

  helpers.logger.info(JSON.stringify(event))

  if (event.name === 'IRC_COMMAND') {
    const args = [...event.data.tokens];
    const command = args.shift().toLowerCase()

    // ================ dispatch to handlers here
    if (command === 'add') {
      await insertResult({ status: 'adding' })

      const url = args.join(' ')
      console.log(`ripping ${url}`)

      let filename;

      // download
      try {
        filename = await youtubeDownload(url)
        await insertResult({ status: 'added', filename })
      } catch(e) {
        console.error(e)
        await insertResult({ status: 'error', error: e.message })
      }

      // add track to db
      try {
        const { rows } = await helpers.query("insert into tracks (filename, event_id) values ($1::text, $2::integer) returning id", [filename, event.id]);
        await insertResult({ status: 'queued', filename, id: rows[0].id })
      } catch(e) {
        console.error(e)
        await insertResult({ status: 'error', message: e.message, code: e.code, detail: e.detail })
      }
    }
  }

  async function insertResult(data) {
    return await helpers.query("insert into results (event_id, name, data) values ($1::integer, $2::text, $3::jsonb) returning id", [event.id, 'IRC_RESPONSE', data])
  }
};
