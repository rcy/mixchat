const webserver = require('./webserver')
const ircBot = require('./ircbot')

const channel = "#emb-radio"
const nick = 'djfullmoon'

ircBot('irc.freenode.net', nick, {
  debug: true,
  port: 6667,
  channels: [channel],
  sasl: true,
  userName: nick,
  password: 'JJyf376fGgbPnfcz9',
})

ircBot('irc.libera.chat', nick, {
  debug: true,
  port: 6697,
  secure: true,
  channels: [channel],
  //  sasl: true,
  userName: nick,
  password: '6MdYaHHKBpSzEwLAa',
})

webserver({ port: 3010 })
