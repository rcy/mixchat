const assert = require('assert').strict;

const webserver = require('./webserver.js')
const ircBot = require('./ircbot.js')

const { LIBERA_PASSWORD, LIBERA_CHANNEL, LIBERA_NICK } = process.env;
assert(LIBERA_PASSWORD)
assert(LIBERA_CHANNEL)
assert(LIBERA_NICK)

ircBot('irc.libera.chat', LIBERA_NICK, {
  //debug: true,
  port: 6697,
  secure: true,
  channels: [LIBERA_CHANNEL],
  sasl: true,
  userName: LIBERA_NICK,
  password: LIBERA_PASSWORD,
})

webserver({ port: 3010 })
