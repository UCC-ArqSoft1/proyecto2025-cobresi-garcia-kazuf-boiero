import { useEffect, useState } from 'react'
import { Link, useNavigate, useParams } from 'react-router-dom'
import ActivityForm from '../components/ActivityForm.jsx'
import Navbar from '../components/Navbar.jsx'
import Notification from '../components/Notification.jsx'
import { useActivities } from '../contexts/ActivitiesContext.jsx'

const EditActivityPage = () => {
  const { activityId } = useParams()
  const navigate = useNavigate()
  const { loadActivityById, updateActivity } = useActivities()
  const [feedback, setFeedback] = useState({ type: null, message: '' })
  const [redirectId, setRedirectId] = useState(null)
  const [activity, setActivity] = useState(null)
  const [isLoading, setIsLoading] = useState(true)
  const [loadError, setLoadError] = useState(null)

  useEffect(() => {
    setIsLoading(true)
    setLoadError(null)
    loadActivityById(activityId)
      .then((data) => setActivity(data))
      .catch((error) => setLoadError(error))
      .finally(() => setIsLoading(false))
  }, [activityId, loadActivityById])

  useEffect(() => {
    return () => {
      if (redirectId) {
        clearTimeout(redirectId)
      }
    }
  }, [redirectId])

  const handleSubmit = async (payload) => {
    if (!activity) return
    try {
      await updateActivity(activity.id, payload)
      setFeedback({ type: 'success', message: 'Actividad actualizada. Redirigiendo…' })
      const timeout = setTimeout(() => navigate(`/activities/${activity.id}`), 900)
      setRedirectId(timeout)
    } catch (error) {
      setFeedback({ type: 'error', message: error.message ?? 'No se pudo actualizar la actividad.' })
    }
  }

  if (isLoading) {
    return (
      <>
        <Navbar />
        <main className="subs-page">
          <h1 className="subs-title">Cargando actividad…</h1>
        </main>
      </>
    )
  }

  if (loadError || !activity) {
    return (
      <>
        <Navbar />
        <main className="subs-page">
          <h1 className="subs-title">No encontramos la actividad.</h1>
          <section className="subs-grid">
            <article className="subs-card not-found-card">
              <Link to="/activities" className="btn-primary not-found-btn">
                Volver
              </Link>
            </article>
          </section>
        </main>
      </>
    )
  }

  return (
    <>
      <Navbar />

      <main className="login-page admin-form-page">
        <section className="login-layout single-column">
          <div className="login-panel">
            <h1 className="login-title">Editar {activity.title}</h1>
            <ActivityForm initialValues={activity} onSubmit={handleSubmit} submitLabel="Guardar cambios" />
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

export default EditActivityPage
