export default function TrackItem({ metadata }) {
  if (!metadata) {
    return null
  }
  const { common } = metadata
  const { artist, title, year } = common

  return [artist, title].join(' / ')
}
