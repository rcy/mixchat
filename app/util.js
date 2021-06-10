function formatDuration(s) {
  const minutes = Math.floor(s / 60)
  const seconds = s % 60
  return `${minutes}` + ':' + `0${seconds}`.slice(-2)
}

module.exports = {
  formatDuration
}
