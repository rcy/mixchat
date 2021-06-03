const fs = require('fs');
const shuffle = require('lodash.shuffle')

let playlist = []
const requests = []

module.exports = {
  pushRequest(filename) {
    requests.push(filename)
  },
  getNext() {
    let next = requests.shift();
    
    if (!next) {
      next = nextFromPlaylist()
    }

    if (!next) {
      next = `./emergency.ogg`
    }

    console.log('next:', next)

    return next
  }
}

function nextFromPlaylist() {
  let next = playlist.shift()
  if (!next) {
    playlist = shuffle(loadPlaylist())
    next = playlist.shift()
  }
  return next
}

function loadPlaylist() {
  const files = fs.readdirSync('/media');
  return files.map(file => `/media/${file}`)
}