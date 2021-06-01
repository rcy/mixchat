const net = require('net');

module.exports = async function liquidsoap(command, host = '172.17.0.1', port = 1234) {
  const client = new net.Socket();

  let result = '';

  return new Promise((resolve, reject) => {
    client.connect(port, host, function() {
      console.log('liquidsoap', command)
      client.write(command + '\n');
    })

    client.on('data', function(data) {
      const text = data.toString().trim()
      console.log('liquidsoap data: ', text)
      if (text === "END") {
        client.write('quit\n')
        resolve(result)
      } else {
        result += text
      }
    });

    client.on('close', function() {
      console.log('liquidsoap telnet connection closed');
      //resolve(result.toString())
    });

    client.on('error', function(e) {
      if (e.code == "ECONNRESET") {
        console.log('liquidsoap: connection reset')
      } else {
        console.error('liquidsoap error:', e)
        reject(e)
      }
    });
  })
}

