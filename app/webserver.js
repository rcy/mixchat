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
      throw new Error('cannot find any track!')
      return
    }
    await pgClient.query('insert into plays (track_id, action) values ($1, $2)', [track.id, 'queued'])
    res.send(`${track.filename}\n`)
  })

  app.post('/now/:station', jsonParser, async (req, res) => {
    console.log('SENT NOW', req.body.filename)
    try {
      await pgClient.query('insert into plays (track_id, action) values ((select id from tracks where filename = $1), $2)', [req.body.filename, 'played'])
    } catch(e) {
      console.error(e)
    }
    PubSub.publish('NOW', req.body)
    res.sendStatus(200)
  });

  app.listen(port, () => {
    console.log(`App listening at http://localhost:${port}`)
  })

  return app
}
