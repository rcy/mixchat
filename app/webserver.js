const express = require('express')
const morgan = require('morgan')
const bodyParser = require('body-parser')
const jsonParser = bodyParser.json()
const PubSub = require('pubsub-js');

module.exports = function webserver({ pgClient, port }) {
  const app = express()

  // add logging middleware
  app.use(morgan('combined'))

  app.get('/', (req, res) => {
    res.send('Hello World!\n')
  })

  app.get('/next/:station_slug', async (req, res) => {
    const { rows: [{ station_id }] } = await pgClient.query("select id as station_id from stations where slug = $1", [req.params.station_slug]);

    const { rows: [track] } = await pgClient.query('select id, filename from tracks where station_id = $1 order by bucket, fuzz, created_at limit 1', [station_id])

    if (!track) {
      console.error(`${req.params.station_slug}: no track found, sending ./404.ogg`)
      res.status(404).send('./404.ogg')
      return
    } else {
      console.log(`/next/${req.params.station_slug}`, { track, station_id })
      await pgClient.query('insert into track_events (station_id, track_id, action) values ($1, $2, $3)', [station_id, track.id, 'queued'])
      res.send(`${track.filename}\n`)
    }
  })

  app.post('/now/:station_slug', jsonParser, async (req, res) => {
    console.log('SENT NOW', req.params, req.body.filename)

    try {
      await pgClient.query(`
insert into track_events (
   station_id,
   track_id,
   action
) values (
   (select id from stations where slug = $1),
   (select id from tracks where filename = $2),
   'played'
)`, [req.params.station_slug, req.body.filename]);
    } catch(e) {
      // TODO: handle better
      console.error(e);
    }

    PubSub.publish('NOW', { ...req.body, station: req.params.station_slug })
    res.sendStatus(200)
  });

  app.listen(port, () => {
    console.log(`App listening at http://localhost:${port}`)
  })

  return app
}
