const liquidsoap = require('../liquidsoap.js')

module.exports = async ({ id }, helpers) => {
  const { rows } = await helpers.query("select * from stations where id = $1", [id])

  const station = rows[0]

  helpers.logger.info(JSON.stringify(station))

  // telnet into liquidsoap and make the station
  const result = await liquidsoap(`meta.make_station ${station.slug}`)

  helpers.logger.info(result)
}
