import { useEffect } from 'react'
import { useLocation } from 'react-router-dom'

const ScrollToHash = () => {
  const location = useLocation()

  useEffect(() => {
    if (location.hash) {
      const elementId = location.hash.replace('#', '')
      const target = document.getElementById(elementId)

      if (target) {
        target.scrollIntoView({ behavior: 'smooth' })
        return
      }
    }

    window.scrollTo(0, 0)
  }, [location.pathname, location.hash])

  return null
}

export default ScrollToHash
