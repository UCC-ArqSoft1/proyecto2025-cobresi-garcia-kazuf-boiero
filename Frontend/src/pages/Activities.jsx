import { useMemo, useState } from 'react'
import { Link } from 'react-router-dom'
import Navbar from '../components/Navbar.jsx'
import { useActivities } from '../contexts/ActivitiesContext.jsx'
import ActivityImage from '../components/ActivityImage.jsx'

const DAY_LABELS = ['Domingo', 'Lunes', 'Martes', 'Miércoles', 'Jueves', 'Viernes', 'Sábado']

const ActivitiesPage = () => {
  const [query, setQuery] = useState('')
  const { activities, activitiesLoading, activitiesError, refreshActivities } = useActivities()

  const filteredActivities = useMemo(() => {
    const normalizedQuery = query.trim().toLowerCase()
    if (!normalizedQuery) return activities

    return activities.filter((activity) => {
      const haystack = `${activity.title} ${activity.category} ${activity.instructor} ${activity.description}`.toLowerCase()

      return haystack.includes(normalizedQuery)
    })
  }, [activities, query])

  return (
    <>
      <Navbar />

      <main className="activities-page">
        <header className="activities-hero">
          <h1>Actividades deportivas</h1>
          <p>Buscá por nombre, profesor, categoría o intensidad.</p>
        </header>

        <section className="activities-search">
          <input
            type="text"
            className="activities-search-input"
            placeholder="Buscar actividad..."
            value={query}
            onChange={(event) => setQuery(event.target.value)}
          />
        </section>

        <section className="activities-grid">
          {activitiesLoading ? (
            <p className="activities-empty">Cargando actividades…</p>
          ) : activitiesError ? (
            <div className="activities-empty">
              <p>No pudimos cargar las actividades.</p>
              <button type="button" className="btn-secondary" onClick={() => refreshActivities().catch(() => {})}>
                Reintentar
              </button>
            </div>
          ) : filteredActivities.length === 0 ? (
            <p className="activities-empty">No encontramos actividades con ese nombre o profesor.</p>
          ) : (
            filteredActivities.map((activity) => {
              const availableSlots =
                typeof activity.availableSlots === 'number' ? activity.availableSlots : activity.capacity
              return (
                <Link to={`/activities/${activity.id}`} className="activity-card" key={activity.id}>
                  <div className="activity-card-media">
                    <ActivityImage
                      src={activity.imageUrl}
                      alt={`Foto de ${activity.title}`}
                      className="activity-card-image"
                    />
                  </div>

                  <div className="activity-card-body">
                    <h2 className="activity-title">{activity.title}</h2>

                    <div className="activity-info">
                      <p className="activity-label">Instructor/a</p>
                      <p className="activity-list">{activity.instructor}</p>

                      <p className="activity-label">Horario</p>
                      <p className="activity-list">
                        {DAY_LABELS[activity.dayOfWeek]} · {activity.startTime} - {activity.endTime}
                      </p>
                    </div>

                    <p className="activity-capacity">
                      Cupos: {availableSlots}/{activity.capacity}
                    </p>
                  </div>
                </Link>
              )
            })
          )}
        </section>
      </main>
    </>
  )
}

export default ActivitiesPage
