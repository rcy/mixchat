const requests = []

let count = 0;

module.exports = {
  pushRequest(file) {
    requests.push(file)
  },
  getNext() {
    let next = requests.shift();
    
//    let rnd = Math.floor(Math.random(3))

    if (!next) {
      next = `/media/${count++ % 5}.ogg`
    }

    console.log('next:', next)

    return next
  }
}
