const net = require('net');

module.exports = async function liquidsoap(command, host = '172.17.0.1', port = 1234) {
  const client = new net.Socket();

  return new Promise((resolve, reject) => {
    client.connect(port, host, function() {
      client.write(command + '\n');
    })

    client.on('data', function(data) {
      resolve(data.toString())
      client.write('quit\n')
      client.destroy();
    });

    client.on('close', function() {
      console.log('liquidsoap telnet connection closed');
    });

    client.on('error', function(e) {
      console.error(e)
      reject(e)
    });
  })
}

