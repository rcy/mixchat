const express = require('express')
const bodyParser = require('body-parser')
const jsonParser = bodyParser.json()
const { getNext } = require('./source.js')
const PubSub = require('pubsub-js');
const youtubeDownload = require('./youtube.js')
const { pushRequest } = require('./source.js')

module.exports = function webserver({ port }) {
  const app = express()

  app.get('/', (req, res) => {
    res.send('Hello World!\n')
  })

  app.post('/add', jsonParser, async (req, res) => {
    try {
      const { filename, data } = await youtubeDownload(req.body.url)
      pushRequest(filename)
      res.status(200).json({ filename, data })
    } catch(e) {
      console.error(e)
      res.send(400, { error: e.message })
    }
  })

  app.get('/next', async (req, res) => {
    const content = getNext()
    res.send(`${content}\n`)
  })

  app.post('/now', jsonParser, async (req, res) => {
    console.log('SENT NOW', req.body)
    PubSub.publish('NOW', req.body)
    res.sendStatus(200)
  });

  app.listen(port, () => {
    console.log(`App listening at http://localhost:${port}`)
  })

  return app
}

async function requestYoutube(url) {
  const filename = await download(url)
  const data = await liquidsoap(`request.push ${filename}`)
  return { filename, data }
}
