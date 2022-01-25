const mm = require('music-metadata');
const util = require('util');

module.exports = async ({ id }, helpers) => {
  const { rows } = await helpers.query("select filename from tracks where id = $1", [id])

  const { filename } = rows[0]

  helpers.logger.info(`updating track metadata for id=${id} filename=${filename}`)

  const metadata = await mm.parseFile(filename);

  await helpers.query(`update tracks set metadata = $1 where id = $2`, [metadata, id]);
}
