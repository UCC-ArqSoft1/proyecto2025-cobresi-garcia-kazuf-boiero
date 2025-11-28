import { useState } from 'react'

const DAY_OPTIONS = [
  { value: 0, label: 'Domingo' },
  { value: 1, label: 'Lunes' },
  { value: 2, label: 'Martes' },
  { value: 3, label: 'Miércoles' },
  { value: 4, label: 'Jueves' },
  { value: 5, label: 'Viernes' },
  { value: 6, label: 'Sábado' },
]

const defaultValues = {
  title: '',
  description: '',
  category: '',
  dayOfWeek: 1,
  startTime: '08:00',
  endTime: '09:00',
  capacity: 10,
  instructor: '',
  imageUrl: '',
  isActive: true,
}

const ActivityForm = ({ initialValues = {}, onSubmit, submitLabel = 'Guardar' }) => {
  const [formValues, setFormValues] = useState({
    ...defaultValues,
    ...initialValues,
  })

  const handleChange = (event) => {
    const { name, value, type, checked } = event.target
    setFormValues((prev) => ({
      ...prev,
      [name]:
        type === 'number'
          ? Number(value)
          : type === 'checkbox'
            ? checked
            : name === 'dayOfWeek'
              ? Number(value)
              : value,
    }))
  }

  const handleSubmit = (event) => {
    event.preventDefault()

    const payload = {
      title: formValues.title.trim(),
      description: formValues.description.trim(),
      category: formValues.category.trim(),
      dayOfWeek: Number(formValues.dayOfWeek),
      startTime: formValues.startTime,
      endTime: formValues.endTime,
      capacity: Number(formValues.capacity),
      instructor: formValues.instructor.trim(),
      imageUrl: formValues.imageUrl.trim(),
      isActive: Boolean(formValues.isActive),
    }

    onSubmit(payload)
  }

  return (
    <form className="login-form activity-form" onSubmit={handleSubmit}>
      <div className="login-field">
        <input name="title" value={formValues.title} onChange={handleChange} placeholder="Título" required />
      </div>
      <div className="login-field">
        <textarea
          name="description"
          value={formValues.description}
          onChange={handleChange}
          placeholder="Descripción detallada"
          rows={4}
          required
        />
      </div>
      <div className="login-field">
        <input name="category" value={formValues.category} onChange={handleChange} placeholder="Categoría" required />
      </div>
      <div className="login-field">
        <select
          id="dayOfWeek"
          name="dayOfWeek"
          value={formValues.dayOfWeek}
          onChange={handleChange}
          aria-label="Día de la semana"
        >
          {DAY_OPTIONS.map((option) => (
            <option key={option.value} value={option.value}>
              {option.label}
            </option>
          ))}
        </select>
      </div>
      <div className="login-field two-columns">
        <input
          type="time"
          name="startTime"
          value={formValues.startTime}
          onChange={handleChange}
          placeholder="Inicio"
          required
        />
        <input
          type="time"
          name="endTime"
          value={formValues.endTime}
          onChange={handleChange}
          placeholder="Fin"
          required
        />
      </div>
      <div className="login-field">
        <input
          type="number"
          min="1"
          name="capacity"
          value={formValues.capacity}
          onChange={handleChange}
          placeholder="Cupos"
          required
        />
      </div>
      <div className="login-field">
        <input
          name="instructor"
          value={formValues.instructor}
          onChange={handleChange}
          placeholder="Instructor/a"
          required
        />
      </div>
      <div className="login-field">
        <input
          name="imageUrl"
          value={formValues.imageUrl}
          onChange={handleChange}
          placeholder="URL de imagen (opcional)"
        />
      </div>
      <div className="login-field checkbox-field">
        <label>
          <input type="checkbox" name="isActive" checked={Boolean(formValues.isActive)} onChange={handleChange} />
          <span>Actividad visible y disponible</span>
        </label>
      </div>

      <div className="login-actions">
        <button type="submit" className="btn-primary">
          {submitLabel}
        </button>
      </div>
    </form>
  )
}

export default ActivityForm
