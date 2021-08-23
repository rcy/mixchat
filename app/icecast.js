const fetch = require('node-fetch')
const convert = require('xml-js');
const assert = require('assert').strict;

const ICECAST_URL = process.env['ICECAST_URL']
assert(ICECAST_URL)

async function fetchXspf() {
  const raw = await fetch(`${ICECAST_URL}/emb.ogg.xspf`)
  const text = await raw.text()
  return JSON.parse(convert.xml2json(text, { arrayNotation: true }))
}

async function countListeners() {
  const xspf = await fetchXspf()

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
  fetchXspf,
  countListeners,
}
