import { useState } from 'react'
import { Link, useLocation, useNavigate } from 'react-router-dom'
import Navbar from '../components/Navbar.jsx'
import Notification from '../components/Notification.jsx'
import { useAuth } from '../contexts/AuthContext.jsx'

const LoginPage = () => {
  const navigate = useNavigate()
  const location = useLocation()
  const { login } = useAuth()
  const [feedback, setFeedback] = useState({ type: null, message: '' })
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [formValues, setFormValues] = useState({ email: '', password: '' })

  const handleChange = (event) => {
    const { name, value } = event.target
    setFormValues((prev) => ({ ...prev, [name]: value }))
  }

  const handleSubmit = async (event) => {
    event.preventDefault()

    if (!formValues.email || !formValues.password) return
    setIsSubmitting(true)
    setFeedback({ type: null, message: '' })

    try {
      await login({
        email: formValues.email,
        password: formValues.password,
      })

      setFeedback({
        type: 'success',
        message: 'Bienvenido/a. Token generado correctamente.',
      })

      const redirectTo = location.state?.from?.pathname ?? '/'
      navigate(redirectTo, { replace: true })
    } catch (error) {
      setFeedback({
        type: 'error',
        message: error.message ?? 'No se pudo iniciar sesión.',
      })
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <>
      <Navbar />

      <main className="login-page">
        <section className="login-layout">
          <div className="login-logo">
            <span className="login-logo-top">PIPO&apos;S</span>
            <span className="login-logo-bottom">GYM</span>
          </div>

          <div className="login-panel">
            <h1 className="login-title">Log In</h1>
            <p className="login-helper">
              Para pruebas podés usar <strong>admin@example.com</strong> o <strong>socia@example.com</strong> con la
              contraseña <em>contra123</em>.
            </p>

            <form className="login-form" onSubmit={handleSubmit}>
              <div className="login-field">
                <input
                  type="email"
                  name="email"
                  placeholder="Email"
                  value={formValues.email}
                  onChange={handleChange}
                  required
                />
              </div>

              <div className="login-field">
                <input
                  type="password"
                  name="password"
                  placeholder="Contraseña"
                  value={formValues.password}
                  onChange={handleChange}
                  required
                />
              </div>

              <div className="login-actions">
                <button type="reset" className="btn-secondary" onClick={() => setFormValues({ email: '', password: '' })}>
                  Limpiar
                </button>
                <button type="submit" className="btn-primary" disabled={isSubmitting}>
                  {isSubmitting ? 'Ingresando…' : 'Ingresar'}
                </button>
              </div>
            </form>

            <p className="login-register">
              ¿No sos miembro?
              <Link to="/signup" className="login-register-link">
                Registrate acá
              </Link>
            </p>
          </div>
        </section>

        <Notification
          type={feedback.type ?? 'success'}
          message={feedback.message}
          onClose={() => setFeedback({ type: null, message: '' })}
        />
      </main>
    </>
  )
}

export default LoginPage
