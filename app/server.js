const { Client } = require('pg')
const assert = require('assert').strict;

const webserver = require('./webserver.js')
const ircBot = require('./ircbot.js')

const { LIBERA_PASSWORD, LIBERA_CHANNEL, LIBERA_NICK, DATABASE_URL } = process.env;
assert(LIBERA_PASSWORD)
assert(LIBERA_CHANNEL)
assert(LIBERA_NICK)
assert(DATABASE_URL)

async function start() {
  const pgClient = new Client({ connectionString: DATABASE_URL })
  await pgClient.connect()
  try {
    // make sure database connection is sound right away
    await pgClient.query('select * from events limit 1')
  } catch(e) {
    console.error(e)
    process.exit(1)
  }

  const { rows } = await pgClient.query('select channel from irc_channels')
  const channels = rows.map(row => row.channel)
  ircBot('irc.libera.chat', LIBERA_NICK, pgClient, {
    //debug: true,
    port: 6697,
    secure: true,
    channels: [...channels, '#fakeemb'],
    sasl: true,
    userName: LIBERA_NICK,
    password: LIBERA_PASSWORD,
  })

  const { makeWorkerUtils } = require("graphile-worker");
  makeWorkerUtils({
  }).then(workerUtils => {
    workerUtils.addJob('hello', { name: `****** APP BOOTED ${pgClient.database}` })
  })

  webserver({ pgClient, port: 3010 })
}

start()
