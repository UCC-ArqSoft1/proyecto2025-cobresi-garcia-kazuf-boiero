import { useEffect, useMemo, useState } from 'react'
import placeholderImage from '../img/activity-placeholder.svg'

const sanitizeImageUrl = (rawUrl) => {
  if (!rawUrl) return ''
  const trimmed = rawUrl.trim()
  if (!trimmed) return ''

  if (trimmed.startsWith('//')) {
    if (typeof window !== 'undefined') {
      return `${window.location.protocol}${trimmed}`
    }
    return `https:${trimmed}`
  }

  if (trimmed.toLowerCase().startsWith('http://') && typeof window !== 'undefined' && window.location.protocol === 'https:') {
    return trimmed.replace(/^http:\/\//i, 'https://')
  }

  return trimmed
}

const ActivityImage = ({ src, alt, className }) => {
  const normalizedSrc = useMemo(() => sanitizeImageUrl(src), [src])
  const [currentSrc, setCurrentSrc] = useState(() => normalizedSrc || placeholderImage)

  useEffect(() => {
    setCurrentSrc(normalizedSrc || placeholderImage)
  }, [normalizedSrc])

  const handleError = () => {
    if (currentSrc !== placeholderImage) {
      setCurrentSrc(placeholderImage)
    }
  }

  return (
    <img
      src={currentSrc}
      alt={alt}
      className={className}
      onError={handleError}
      loading="lazy"
      referrerPolicy="no-referrer"
    />
  )
}

export default ActivityImage
