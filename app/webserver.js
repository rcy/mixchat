const express = require('express')
const bodyParser = require('body-parser')
const jsonParser = bodyParser.json()
const PubSub = require('pubsub-js');

module.exports = function webserver({ pgClient, port }) {
  const app = express()

  app.get('/', (req, res) => {
    res.send('Hello World!\n')
  })

  app.get('/next', async (req, res) => {
    const { rows } = await pgClient.query('select id, filename from tracks order by bucket, fuzz, created_at limit 1')
    const track = rows[0]
    if (!track) {
      res.send('/media/thoop-RGcR9hVG4f4.ogg')
      return
    }
    await pgClient.query('insert into plays (track_id) values ($1)', [track.id])
    res.send(`${track.filename}\n`)
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
