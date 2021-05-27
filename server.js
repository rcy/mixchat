require('dotenv').config()
const assert = require('assert').strict;

const webserver = require('./webserver.js')
const ircBot = require('./ircbot.js')

const channel = "#emb-radio"
const nick = 'djfullmoon'

const { LIBERA_PASSWORD } = process.env;
assert(LIBERA_PASSWORD)

ircBot('irc.libera.chat', nick, {
  debug: true,
  port: 6697,
  secure: true,
  channels: [channel],
  sasl: true,
  userName: nick,
  password: LIBERA_PASSWORD,
})

webserver({ port: 3010 })
