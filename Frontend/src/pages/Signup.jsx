import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import Navbar from '../components/Navbar.jsx'
import Notification from '../components/Notification.jsx'
import { useAuth } from '../contexts/AuthContext.jsx'

const initialValues = {
  fullName: '',
  email: '',
  password: '',
  confirmPassword: '',
}

const SignupPage = () => {
  const navigate = useNavigate()
  const { login, register } = useAuth()
  const [feedback, setFeedback] = useState({ type: null, message: '' })
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [formValues, setFormValues] = useState(initialValues)

  const handleChange = (event) => {
    const { name, value } = event.target
    setFormValues((prev) => ({ ...prev, [name]: value }))
  }

  const handleSubmit = async (event) => {
    event.preventDefault()
    if (formValues.password !== formValues.confirmPassword) {
      setFeedback({ type: 'error', message: 'Las contraseñas no coinciden.' })
      return
    }

    setIsSubmitting(true)
    setFeedback({ type: null, message: '' })

    try {
      await register({
        name: formValues.fullName.trim(),
        email: formValues.email,
        password: formValues.password,
      })

      await login({
        email: formValues.email,
        password: formValues.password,
      })

      setFeedback({ type: 'success', message: 'Registro exitoso. Ya podés elegir tu actividad.' })
      navigate('/activities')
    } catch (error) {
      setFeedback({ type: 'error', message: error.message ?? 'No se pudo completar el registro.' })
    } finally {
      setIsSubmitting(false)
    }
  }

  const handleReset = () => setFormValues(initialValues)

  return (
    <>
      <Navbar />

      <main className="login-page">
        <section className="login-layout">
          <div className="login-panel">
            <h1 className="login-title">¡Hacete miembro!</h1>

            <form className="login-form" onSubmit={handleSubmit}>
              <div className="login-field">
                <input
                  name="fullName"
                  placeholder="Nombre y apellido"
                  value={formValues.fullName}
                  onChange={handleChange}
                  required
                />
              </div>

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

              <div className="login-field">
                <input
                  type="password"
                  name="confirmPassword"
                  placeholder="Confirmar contraseña"
                  value={formValues.confirmPassword}
                  onChange={handleChange}
                  required
                />
              </div>

              <div className="login-actions">
                <button type="reset" className="btn-secondary" onClick={handleReset}>
                  Limpiar
                </button>
                <button type="submit" className="btn-primary" disabled={isSubmitting}>
                  {isSubmitting ? 'Enviando…' : 'Registrarme'}
                </button>
              </div>
            </form>

            <p className="login-register">
              ¿Ya sos miembro?
              <Link to="/login" className="login-register-link">
                Iniciá sesión acá
              </Link>
            </p>
          </div>

          <div className="login-logo">
            <span className="login-logo-top">PIPO&apos;S</span>
            <span className="login-logo-bottom">GYM</span>
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

export default SignupPage
