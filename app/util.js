function formatDuration(s) {
  const minutes = Math.floor(s / 60)
  const seconds = s % 60
  return `0${minutes}`.slice(-2) + ':' + `0${seconds}`.slice(-2)
}

module.exports = {
  formatDuration
}
