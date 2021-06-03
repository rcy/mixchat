const fetch = require('node-fetch')
const convert = require('xml-js');
async function fetchXspf() {
  const raw = await fetch('http://radio.nonzerosoftware.com:8000/emb.ogg.xspf')
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

  return match[1]
}

module.exports = {
  fetchXspf,
  countListeners,
}
