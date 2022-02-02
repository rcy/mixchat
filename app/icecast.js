const fetch = require('node-fetch')
const convert = require('xml-js');
const assert = require('assert').strict;

const ICECAST_URL = process.env['ICECAST_URL']
assert(ICECAST_URL)

async function fetchXspf(url) {
  console.log('fetchXspf', url)
  const raw = await fetch(url)
  const text = await raw.text()
  return JSON.parse(convert.xml2json(text, { arrayNotation: true }))
}

async function countListeners(station) {
  return await countFormatListeners(station, 'ogg') + await countFormatListeners(station, 'mp3')
}

async function countFormatListeners(station, format) {
  try {
    const url = `${ICECAST_URL}/${station}.${format}.xspf`
    const xspf = await fetchXspf(url)

    const match =
      xspf.elements[0].elements
          .find(e => e.name === 'trackList')
          .elements[0].elements
          .find(e => e.name === 'annotation')
          .elements[0].text
          .split('\n')
          .find(e => e.match('Current Listeners'))
          .match(/: (.+)$/)

    return +match[1]
  } catch(e) {
    return 0
  }
}

module.exports = {
  countListeners,
}
