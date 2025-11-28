import { Link } from 'react-router-dom'
import Navbar from '../components/Navbar.jsx'

const NotFoundPage = () => {
  return (
    <>
      <Navbar />

      <main className="subs-page">
        <h1 className="subs-title">Ups, no encontramos la página que buscás.</h1>

        <section className="subs-grid">
          <article className="subs-card not-found-card">
            <p className="subs-text">Puede que el enlace esté roto o que la página se haya movido.</p>
            <Link to="/" className="btn-primary not-found-btn">
              Volver al inicio
            </Link>
          </article>
        </section>
      </main>
    </>
  )
}

export default NotFoundPage
