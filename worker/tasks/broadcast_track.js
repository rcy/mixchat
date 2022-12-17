const fetch = require('node-fetch')
const util = require('util');

module.exports = async ({ id }, helpers) => {
  const { rows } = await helpers.query("select * from tracks where id = $1 limit 1", [id])

  if (rows[0].station_id !== 1) {
    return
  }

  const purl = rows[0].metadata.native.vorbis.find(x => x.id === 'PURL').value
  const { artist, title } = rows[0].metadata.common

  console.log({ rows, artist, title, purl })

  const result = await fetch(process.env.FLAPPER_URL, {
    method: "POST",
    body: `${artist}, ${title} ${purl}`
  })

  if (result.status !== 200) {
    throw new Error(`status not ok ${result.status}`);
  }
}