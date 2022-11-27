export default function TrackItem({ metadata }) {
  if (!metadata) {
    return null
  }
  const { common } = metadata
  const { artist, title } = common

  return [artist, title].join(' / ')
}
