const fetch = require('node-fetch')
const convert = require('xml-js');
const assert = require('assert').strict;

const ICECAST_URL = process.env['ICECAST_URL']
assert(ICECAST_URL)

async function fetchXspf(station) {
  const url = `${ICECAST_URL}/${station}.ogg.xspf`
  console.log('fetchXspf', url)
  const raw = await fetch(url)
  const text = await raw.text()
  return JSON.parse(convert.xml2json(text, { arrayNotation: true }))
}

async function countListeners(station) {
  const xspf = await fetchXspf(station)

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
}

module.exports = {
  countListeners,
}
