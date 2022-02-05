const express = require('express')
const cors = require('cors')
const morgan = require('morgan')
const bodyParser = require('body-parser')
const jsonParser = bodyParser.json()
const PubSub = require('pubsub-js');
const { postgraphile } = require('postgraphile')
const MyPlugin = require('./schema-extension')

module.exports = function webserver({ pgClient, port }) {
  const app = express()

  // enable all cors requests
  app.use(cors())

  // add logging middleware
  app.use(morgan('combined'))

  // add postgraphile
  app.use(
    postgraphile(
      process.env.DATABASE_URL,
      "public",
      {
        watchPg: true,
        graphiql: true,
        enhanceGraphiql: true,
        ownerConnectionString: process.env.ROOT_DATABASE_URL, // to setup watch fixtures
        appendPlugins: [MyPlugin],
        dynamicJson: true,
      }
    )
  )

  app.get('/', (req, res) => {
    res.send('Hello World!\n')
  })

  app.get('/next/:station_slug', async (req, res) => {
    try {
      const { rows: [{ station_id }] } = await pgClient.query("select id as station_id from stations where slug = $1", [req.params.station_slug]);

      const { rows: [track] } = await pgClient.query('select id, filename from tracks where station_id = $1 order by bucket, fuzz, created_at limit 1', [station_id])

      if (!track) {
        console.error(`${req.params.station_slug}: no track found`)
        res.status(404).send('404')
        return
      } else {
        console.log(`/next/${req.params.station_slug}`, { track, station_id })
        await pgClient.query('insert into track_events (station_id, track_id, action) values ($1, $2, $3)', [station_id, track.id, 'queued'])
        res.send(`${track.filename}`)
      }
    } catch(e) {
      console.error(req.path, e)
      res.status(500).send('500')
    }
  })

  app.post('/now/:slug', jsonParser, async (req, res) => {
    console.log('SENT NOW', req.params, req.body.filename)

    const { slug } = req.params;
    const { filename } = req.body;

    try {
      const { rows } = await pgClient.query(`
select tracks.id as track_id, stations.id as station_id
 from tracks 
 join stations on tracks.station_id = stations.id
 where tracks.filename = $1 and stations.slug = $2
      `, [filename, slug])

      if (rows[0]) {
        await pgClient.query(`
  insert into track_events (
     station_id,
     track_id,
     action
  ) values ($1, $2, 'played')
        `, [rows[0].station_id, rows[0].track_id]);
      }
    } catch(e) {
      // TODO: handle better
      console.error(e);
    }

    PubSub.publish('NOW', { ...req.body, station: req.params.slug })
    res.sendStatus(200)
  });

  app.listen(port, () => {
    console.log(`App listening at http://localhost:${port}`)
  })

  return app
}
