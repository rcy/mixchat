const requests = []

module.exports = {
  pushRequest(filename) {
    requests.push(filename)
  },
  getNext() {
    let next = requests.shift();
    
    if (!next) {
      next = `./emergency.ogg`
    }

    console.log('next:', next)

    return next
  }
}
