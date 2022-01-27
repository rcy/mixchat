const youtubedl = require('youtube-dl-exec')

module.exports = async function youtubeDownload(url) {
  const output = await youtubedl(url, {
    quiet: true,
    extractAudio: true,
    audioFormat: 'vorbis',
    //dumpSingleJson: true,
    // noWarnings: true,
    noCallHome: true,
    // noCheckCertificate: true,
    // preferFreeFormats: true,
    youtubeSkipDashManifest: true,
    //output: '/media/%(id)s.%(ext)s',
    //    referer: 'https://example.com',
    addMetadata: true,
    restrictFilenames: true,
    noPlaylist: true,
    maxDownloads: 1, // prevent channel downloads for now
    exec: "mv {} /media && echo {}", // output
  })
  console.log(output)
  return `/media/${output}`
}
