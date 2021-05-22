const express = require('express')
const bodyParser = require('body-parser')
const jsonParser = bodyParser.json()

module.exports = function webserver({ port }) {
  const app = express()

  app.get('/', (req, res) => {
    res.send('Hello World!\n')
  })

  app.post('/youtube', jsonParser, async (req, res) => {
    try {
      const { filename, data } = requestYoutube(req.body.id)
      res.status(200).json({ filename, data })
    } catch(e) {
      console.error(e)
      res.sendStatus(500, { error: e.message })
    }
  })

  app.listen(port, () => {
    console.log(`App listening at http://localhost:${port}`)
  })

  return app
}

async function requestYoutube(url) {
  const filename = await download(url)
  const data = await liquidsoap(`request.push ${filename}`)
  return { filename, data }
}
