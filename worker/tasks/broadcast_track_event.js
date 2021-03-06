const fetch = require('node-fetch')
const util = require('util');

module.exports = async ({ id }, helpers) => {
  const { rows } = await helpers.query("select * from track_events join tracks on tracks.id = track_events.track_id where track_events.id = $1 limit 1", [id])

  if (rows[0].action !== 'played') {
    return
  }

  if (rows[0].station_id !== 1) {
    return
  }

  const purl = rows[0].metadata.native.vorbis.find(x => x.id === 'PURL').value
  const { artist, title } = rows[0].metadata.common

  console.log({ rows, artist, title, purl })

  const result = await fetch("https://nichevomit.fly.dev/twtr/2C5hvW3Ty1MXW68mBwqdQaz0pkf", {
    method: "POST",
    body: `${artist}, ${title} ${purl}`
  })

  // TODO: 500s come back no matter how it failed... often it will be because of a duplicate track
  console.log({ result })
}
