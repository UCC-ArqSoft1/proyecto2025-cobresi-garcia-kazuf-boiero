import { useEffect, useMemo, useState } from 'react'
import { Link, useLocation, useNavigate, useParams } from 'react-router-dom'
import Navbar from '../components/Navbar.jsx'
import Notification from '../components/Notification.jsx'
import { useActivities } from '../contexts/ActivitiesContext.jsx'
import { useAuth } from '../contexts/AuthContext.jsx'
import ActivityImage from '../components/ActivityImage.jsx'

const DAY_LABELS = ['Domingo', 'Lunes', 'Martes', 'Miércoles', 'Jueves', 'Viernes', 'Sábado']

const ActivityDetailPage = () => {
  const { activityId } = useParams()
  const navigate = useNavigate()
  const location = useLocation()
  const { loadActivityById, enrollInActivity, deleteActivity, unenrollFromActivity, myActivities } = useActivities()
  const { isAdmin, isAuthenticated } = useAuth()
  const [feedback, setFeedback] = useState({ type: null, message: '' })
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

  const alreadyEnrolled = useMemo(() => {
    if (!activity) return false
    return myActivities.some((item) => item.id === activity.id)
  }, [myActivities, activity])

  const availableSlots = useMemo(() => {
    if (!activity) return 0
    if (typeof activity.availableSlots === 'number') {
      return activity.availableSlots
    }
    return activity.capacity
  }, [activity])

  const isFull = activity ? availableSlots <= 0 : false

  const handleEnroll = async () => {
    if (!activity) return

    if (!isAuthenticated) {
      navigate('/login', { state: { from: location }, replace: false })
      return
    }

    try {
      const updated = await enrollInActivity(activity.id)
      if (updated) {
        setActivity(updated)
      }
      setFeedback({
        type: 'success',
        message: `Te inscribiste a ${activity.title}. ¡Nos vemos en la próxima clase!`,
      })
    } catch (error) {
      let message = error.message ?? 'No pudimos completar la inscripción.'
      if (error.code === 'SCHEDULE_CONFLICT') {
        message = 'Ya tenés una actividad inscripta que se solapa en día y horario.'
      } else if (error.code === 'NO_CAPACITY') {
        message = 'No quedan cupos disponibles para esta actividad.'
      }
      setFeedback({
        type: 'error',
        message,
      })
    }
  }

  const handleUnenroll = async () => {
    if (!activity) return

    try {
      const updated = await unenrollFromActivity(activity.id)
      if (updated) {
        setActivity(updated)
      }
      setFeedback({
        type: 'success',
        message: `Cancelaste tu inscripción a ${activity.title}.`,
      })
    } catch (error) {
      setFeedback({
        type: 'error',
        message: error.message ?? 'No pudimos procesar tu baja.',
      })
    }
  }

  const handleDelete = async () => {
    if (!activity) return
    const confirmed = window.confirm('¿Eliminar la actividad? Esta acción no se puede deshacer.')
    if (!confirmed) return

    try {
      await deleteActivity(activity.id)
      setFeedback({ type: 'success', message: 'Actividad eliminada. Redirigiendo…' })
      setTimeout(() => navigate('/activities', { replace: true }), 900)
    } catch (error) {
      setFeedback({ type: 'error', message: error.message ?? 'No se pudo eliminar la actividad.' })
    }
  }

  const renderAdminActions = () => {
    if (!isAdmin || !activity) return null

    return (
      <div className="detail-admin-actions">
        <Link to={`/admin/activities/${activity.id}/edit`} className="btn-secondary">
          Editar actividad
        </Link>
        <button type="button" className="btn-primary danger" onClick={handleDelete}>
          Eliminar
        </button>
      </div>
    )
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
          <h1 className="subs-title">La actividad que buscás no está disponible.</h1>
          <section className="subs-grid">
            <article className="subs-card not-found-card">
              <p className="subs-text">Podés volver al listado y explorar otras opciones.</p>
              <Link to="/activities" className="btn-primary not-found-btn">
                Volver a actividades
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

      <main className="activity-detail-page">
        <section className="activity-detail-card">
          <div className="activity-detail-media">
            <ActivityImage
              src={activity.imageUrl}
              alt={`Foto de ${activity.title}`}
              className="activity-detail-image"
            />
          </div>

          <div className="activity-detail-content">
            <p className="activity-pill">{activity.category}</p>
            <h1>{activity.title}</h1>
            <p className="detail-description">{activity.description}</p>

            <ul className="detail-list">
              <li>
                <span>Instructor/a:</span> {activity.instructor}
              </li>
              <li>
                <span>Día:</span> {DAY_LABELS[activity.dayOfWeek]}
              </li>
              <li>
                <span>Horario:</span> {activity.startTime} - {activity.endTime}
              </li>
              <li>
                <span>Cupos:</span> {availableSlots}/{activity.capacity}
              </li>
              {!activity.isActive ? (
                <li>
                  <span>Estado:</span> Inactiva
                </li>
              ) : null}
            </ul>

            <div className="detail-actions">
              <button
                type="button"
                className="btn-primary"
                onClick={handleEnroll}
                disabled={!activity.isActive || alreadyEnrolled || isFull}
              >
                {alreadyEnrolled ? 'Ya estás inscripto' : !activity.isActive ? 'Inactiva' : isFull ? 'Sin cupos' : 'Inscribirme'}
              </button>
              {alreadyEnrolled ? (
                <button type="button" className="btn-primary danger" onClick={handleUnenroll}>
                  Desinscribirme
                </button>
              ) : null}
              <Link to="/activities" className="btn-secondary">
                Volver al listado
              </Link>
            </div>

            <p className="detail-quota-message">
              {activity.isActive
                ? isFull
                  ? 'No quedan cupos disponibles.'
                  : `Quedan ${availableSlots} cupos disponibles.`
                : 'Esta actividad está inactiva.'}
            </p>

            {renderAdminActions()}
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

export default ActivityDetailPage
