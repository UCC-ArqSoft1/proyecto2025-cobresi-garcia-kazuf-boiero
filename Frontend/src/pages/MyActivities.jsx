import { Link } from 'react-router-dom'
import { useState } from 'react'
import Navbar from '../components/Navbar.jsx'
import Notification from '../components/Notification.jsx'
import { useActivities } from '../contexts/ActivitiesContext.jsx'

const DAY_LABELS = ['Domingo', 'Lunes', 'Martes', 'Miércoles', 'Jueves', 'Viernes', 'Sábado']

const MyActivitiesPage = () => {
  const { myActivities, myActivitiesLoading, refreshMyActivities, unenrollFromActivity } = useActivities()
  const [feedback, setFeedback] = useState({ type: null, message: '' })
  const [pendingActivityId, setPendingActivityId] = useState(null)

  const handleUnenroll = async (activityId) => {
    setPendingActivityId(activityId)
    try {
      await unenrollFromActivity(activityId)
      setFeedback({
        type: 'success',
        message: 'Cancelaste tu inscripción.',
      })
    } catch (error) {
      setFeedback({
        type: 'error',
        message: error.message ?? 'No pudimos cancelar tu inscripción.',
      })
    } finally {
      setPendingActivityId(null)
    }
  }

  return (
    <>
      <Navbar />

      <main className="subs-page">
        <h1 className="subs-title">Tus actividades</h1>

        <section className="subs-grid">
          {myActivitiesLoading ? (
            <p className="activities-empty">Cargando tus actividades…</p>
          ) : myActivities.length === 0 ? (
            <article className="subs-card not-found-card">
              <p className="subs-text">Aún no te inscribiste a ninguna actividad.</p>
              <Link to="/activities" className="btn-primary not-found-btn">
                Buscar actividades
              </Link>
            </article>
          ) : (
            myActivities.map((activity) => (
              <article className="subs-card" key={activity.id}>
                <h2 className="subs-activity">{activity.title}</h2>

                <div className="subs-info">
                  <p className="subs-label">Día</p>
                  <p className="subs-text">{DAY_LABELS[activity.dayOfWeek]}</p>

                  <p className="subs-label">Horario</p>
                  <p className="subs-text">
                    {activity.startTime} - {activity.endTime}
                  </p>
                </div>

                <div className="subs-card-actions">
                  <Link to={`/activities/${activity.id}`} className="btn-secondary subs-link">
                    Ver detalle
                  </Link>
                  <button
                    type="button"
                    className="btn-primary danger"
                    onClick={() => handleUnenroll(activity.id)}
                    disabled={pendingActivityId === activity.id}
                  >
                    {pendingActivityId === activity.id ? 'Procesando…' : 'Desinscribirme'}
                  </button>
                </div>
              </article>
            ))
          )}
        </section>

        <div className="subs-actions">
          <button type="button" className="btn-secondary" onClick={() => refreshMyActivities().catch(() => {})}>
            Actualizar listado
          </button>
        </div>

        <Notification
          type={feedback.type ?? 'success'}
          message={feedback.message}
          onClose={() => setFeedback({ type: null, message: '' })}
        />
      </main>
    </>
  )
}

export default MyActivitiesPage
