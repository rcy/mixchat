export default function MetadataLink({ metadata }) {
  const link = metadata?.native?.vorbis?.find(x => x.id === "PURL")?.value
  if (!link) {
    return null
  }

  const url = new URL(link)

  // www.youtube.com -> youtube, etc
  const host = url.host.replace(/^(.+\.)?(.+)\..+$/, '$2')

  const shortcode = ({
    bandcamp: 'bc',
    reddit: 're',
    soundcloud: 'sc',
    tiktok: 'tt',
    twitter: 'tw',
    youtube: 'yt',
    vimeo: 'vm'
  })[host] || '??';

  return (
    <a
      href={link}
      target="_blank"
    >{shortcode}</a>
  )
}
