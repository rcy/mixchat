expressApp()

ircConnect('irc.freenode.net', 'djfullmoon', {
  channels: ["#emb-radio"],
  debug: true,
  sasl: true,
  userName: 'djfullmoon',
  password: 'JJyf376fGgbPnfcz9',
})

ircConnect('irc.libera.chat', 'djfullmoon', {
  port: 6697,
  secure: true,
  channels: ["#emb-radio-debug"],
  debug: true,
  //  sasl: true,
  userName: 'djfullmoon',
  password: '6MdYaHHKBpSzEwLAa',
})
