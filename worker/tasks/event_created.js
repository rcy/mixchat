const { execSync } = require('child_process')
const { handlers } = require('../lib/handlers.js')

module.exports = async ({ id }, helpers) => {
  const { rows } = await helpers.query("select * from events where id = $1", [id])

  const event = rows[0]

  helpers.logger.info(JSON.stringify(event))

  let insertResult;

  if (event.name === 'IRC_COMMAND') {
    insertResult = async (data) => {
      return await helpers.query("insert into results (station_id, event_id, name, data) values ($1::integer, $2::integer, $3::text, $4::jsonb) returning id", [event.station_id, event.id, 'IRC_RESPONSE', data])
    }
  } else if (event.name === 'WEB_COMMAND') {
    insertResult = async (data) => {
      return await helpers.query("insert into results (station_id, event_id, name, data) values ($1::integer, $2::integer, $3::text, $4::jsonb) returning id", [event.station_id, event.id, 'WEB_RESPONSE', data])
    }
  } else {
    helpers.logger.error(`dropping unhandled event.name: ${event.name}`)
    return
  }
  await processCommand({ event, helpers, insertResult })
};

async function processCommand({ event, helpers, insertResult }) {
  const args = [...event.data.tokens];
  const command = args.shift().toLowerCase()

  const handler = handlers[command]
  if (handler) {
    //await insertResult({ message: `Read command !${command}...` })
    await handler(args, { event, insertResult, helpers })
    //await insertResult({ message: `Read command !${command}...done` })
  } else {
    await insertResult({ message: `Bad command !${command}. Type !help` })
  }
}
