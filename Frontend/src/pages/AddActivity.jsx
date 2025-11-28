import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import ActivityForm from '../components/ActivityForm.jsx'
import Navbar from '../components/Navbar.jsx'
import Notification from '../components/Notification.jsx'
import { useActivities } from '../contexts/ActivitiesContext.jsx'

const AddActivityPage = () => {
  const { createActivity } = useActivities()
  const navigate = useNavigate()
  const [feedback, setFeedback] = useState({ type: null, message: '' })
  const [redirectId, setRedirectId] = useState(null)

  useEffect(() => {
    return () => {
      if (redirectId) {
        clearTimeout(redirectId)
      }
    }
  }, [redirectId])

  const handleSubmit = async (payload) => {
    try {
      const created = await createActivity(payload)
      setFeedback({ type: 'success', message: 'Actividad creada correctamente. Redirigiendo…' })
      const timeout = setTimeout(() => navigate(`/activities/${created.id}`), 900)
      setRedirectId(timeout)
    } catch (error) {
      setFeedback({ type: 'error', message: error.message ?? 'No se pudo crear la actividad.' })
    }
  }

  return (
    <>
      <Navbar />

      <main className="login-page admin-form-page">
        <section className="login-layout single-column">
          <div className="login-panel">
            <h1 className="login-title">Nueva actividad</h1>
            <ActivityForm onSubmit={handleSubmit} submitLabel="Crear actividad" />
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

export default AddActivityPage
