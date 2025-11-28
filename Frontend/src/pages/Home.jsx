import { Link } from 'react-router-dom'
import Navbar from '../components/Navbar.jsx'

const HomePage = () => {
  return (
    <>
      <Navbar />

      <main>
        <section id="home" className="hero">
          <div className="hero-overlay" />

          <div className="hero-content">
            <h1 className="hero-title">
              <span className="hero-title-top">PIPO&apos;S</span>
              <span className="hero-title-bottom">GYM</span>
            </h1>
          </div>
        </section>

        <section id="about" className="about">
          <div className="about-card">
            <h2 className="about-title">¿Quiénes somos?</h2>

            <p className="about-text">
              En Pipo&apos;s Gym creemos que entrenar es mucho más que levantar peso: es crecer, superarse y construir
              hábitos que transforman tu vida. Nacimos con la idea de crear un espacio donde cada persona, sin importar
              su nivel, se sienta acompañada, motivada y capaz de alcanzar sus metas.
            </p>

            <p className="about-text">
              Nuestro equipo está formado por profesionales apasionados por el movimiento, la salud y el bienestar,
              listos para guiarte en cada paso. Acá vas a encontrar un ambiente cercano, energía real y una comunidad
              que te impulsa a ser tu mejor versión.
            </p>

            <p className="about-highlight">
              Pipo&apos;s Gym no es solo un gimnasio. Es tu lugar para empezar, mejorar y nunca frenar.
            </p>

            <div className="about-services">
              <article className="service-card service-trainer">
                <span className="service-label">Personal Trainer</span>
              </article>

              <article className="service-card service-kine">
                <span className="service-label">Kinesiólogo/a</span>
              </article>

              <article className="service-card service-7dias">
                <span className="service-label">7 Días</span>
              </article>
            </div>
          </div>
        </section>

        <section className="cta-section">
          <div className="cta-card">
            <p className="cta-quote">
              «Este es tu momento. Sumate al gimnasio y empezá a construir la mejor versión de vos mismo. Acá vas a
              encontrar energía, guía y un espacio donde cada día cuenta. Tu cambio empieza hoy.»
            </p>

            <p className="cta-sub">Empezar es el paso más difícil. El resto lo hacemos juntos.</p>

            <Link to="/activities" className="cta-button">
              DESCUBRÍ LAS ACTIVIDADES
            </Link>
          </div>
        </section>
      </main>
    </>
  )
}

export default HomePage
