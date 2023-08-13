const { create } = require('youtube-dl-exec')
const youtubedl = create('/usr/local/bin/yt-dlp')

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
  const basename = output.split('/').pop()
  return `/media/${basename}`
}
