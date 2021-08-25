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

      console.log('liquidsoap data: <<<', text, '>>>')

      const lines = text.split('\r\n')
      for (let line of lines) {
        console.log('liquidsoap line', line);
        if (line === "END") {
          console.log('liquidsoap write quit');
          client.write('quit\n')
        } else {
          result += line
        }
      }
    });

    client.on('close', function() {
      console.log('liquidsoap telnet connection closed (resolving)');
      resolve(result.toString())
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
