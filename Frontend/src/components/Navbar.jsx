import { useEffect, useMemo, useState } from 'react'
import { Link, useLocation, useNavigate } from 'react-router-dom'
import { useAuth } from '../contexts/AuthContext.jsx'

const Navbar = ({ links, cta }) => {
  const location = useLocation()
  const navigate = useNavigate()
  const { isAuthenticated, isAdmin, logout, user } = useAuth()
  const [isMenuOpen, setIsMenuOpen] = useState(false)

  const toggleMenu = () => setIsMenuOpen((prev) => !prev)
  const closeMenu = () => setIsMenuOpen(false)

  const computedLinks = useMemo(() => {
    if (links?.length) return links

    const baseLinks = [
      { label: 'Home', to: '/' },
      { label: 'About Us', to: '/', hash: '#about' },
      { label: 'Actividades', to: '/activities' },
    ]

    if (isAuthenticated) {
      baseLinks.push({ label: 'Mis actividades', to: '/mis-actividades' })
    }

    if (isAdmin) {
      baseLinks.push({ label: 'Administrar', to: '/admin/activities/new' })
    }

    return baseLinks
  }, [links, isAuthenticated, isAdmin])

  useEffect(() => {
    setIsMenuOpen(false)
  }, [location.pathname, location.hash])

  const handleLogout = () => {
    closeMenu()
    logout()
    navigate('/')
  }

  const ctaConfig =
    cta ??
    (isAuthenticated
      ? { label: 'Salir', onClick: handleLogout }
      : {
          to: '/login',
          label: 'Miembros',
        })

  const isLinkActive = (link) => {
    if (link.hash) {
      return location.pathname === link.to && location.hash === link.hash
    }

    return location.pathname === link.to
  }

  const renderCta = () => {
    if (ctaConfig.onClick) {
      return (
        <button
          type="button"
          className="btn-members"
          onClick={() => {
            closeMenu()
            ctaConfig.onClick()
          }}
        >
          {ctaConfig.label}
        </button>
      )
    }

    return (
      <Link to={ctaConfig.to} className="btn-members" onClick={closeMenu}>
        {ctaConfig.label}
      </Link>
    )
  }

  return (
    <header className="header">
      <nav className="navbar">
        <Link to="/" className="logo" aria-label="Ir al inicio">
          <span className="logo-top">PIPO&apos;S</span>
          <span className="logo-bottom">GYM</span>
        </Link>

        <button
          type="button"
          className={`menu-toggle ${isMenuOpen ? 'open' : ''}`}
          aria-label="Abrir menú de navegación"
          aria-expanded={isMenuOpen}
          aria-controls="nav-content"
          onClick={toggleMenu}
        >
          <span />
          <span />
          <span />
        </button>

        <div className={`nav-content ${isMenuOpen ? 'open' : ''}`} id="nav-content">
          <ul className="nav-links">
            {computedLinks.map((link) => {
              const target = link.hash ? { pathname: link.to, hash: link.hash } : link.to
              const className = isLinkActive(link) ? 'active' : ''

              return (
                <li key={link.label}>
                  <Link to={target} className={className} onClick={closeMenu}>
                    {link.label}
                  </Link>
                </li>
              )
            })}
          </ul>

          <div className="navbar-actions">
            {user ? <span className="user-pill">{user.role === 'admin' ? 'Admin' : 'Socio'}</span> : null}
            {renderCta()}
          </div>
        </div>
      </nav>
    </header>
  )
}

export default Navbar
