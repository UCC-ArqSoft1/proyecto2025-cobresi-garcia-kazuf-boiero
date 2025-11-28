import apiClient from './apiClient.js'

const toActivity = (payload) => ({
  id: payload.id,
  title: payload.title,
  description: payload.description,
  category: payload.category,
  dayOfWeek: payload.day_of_week,
  startTime: payload.start_time,
  endTime: payload.end_time,
  capacity: payload.capacity,
  instructor: payload.instructor,
  imageUrl: payload.image_url,
  isActive: payload.is_active,
  availableSlots: typeof payload.available_slots === 'number' ? payload.available_slots : null,
  enrolledCount: typeof payload.enrolled_count === 'number' ? payload.enrolled_count : null,
})

const toAdminPayload = (payload) => ({
  title: payload.title,
  description: payload.description,
  category: payload.category,
  day_of_week: Number(payload.dayOfWeek),
  start_time: payload.startTime,
  end_time: payload.endTime,
  capacity: Number(payload.capacity),
  instructor: payload.instructor,
  image_url: payload.imageUrl || '',
  is_active: payload.isActive ?? true,
})

const toMemberActivity = (payload) => ({
  id: payload.id,
  title: payload.title,
  description: payload.description,
  category: payload.category,
  dayOfWeek: payload.day_of_week,
  startTime: payload.start_time,
  endTime: payload.end_time,
  instructor: payload.instructor,
})

export const listActivities = async (filters = {}) => {
  const params = new URLSearchParams()
  if (filters.query) params.set('q', filters.query)
  if (filters.category) params.set('category', filters.category)
  if (typeof filters.day === 'number') params.set('day', filters.day)

  const query = params.toString()
  const data = await apiClient.get(query ? `/activities?${query}` : '/activities')
  return Array.isArray(data) ? data.map(toActivity) : []
}

export const getActivity = async (id) => {
  const data = await apiClient.get(`/activities/${id}`)
  return toActivity(data)
}

export const enrollInActivity = async (id) => apiClient.post(`/activities/${id}/enroll`)
export const unenrollFromActivity = async (id) => apiClient.delete(`/activities/${id}/enroll`)

export const getMyActivities = async () => {
  const data = await apiClient.get('/me/activities')
  return Array.isArray(data) ? data.map(toMemberActivity) : []
}

export const createActivity = async (payload) => {
  const data = await apiClient.post('/admin/activities', toAdminPayload(payload))
  return toActivity(data)
}

export const updateActivity = async (id, payload) => {
  const data = await apiClient.put(`/admin/activities/${id}`, toAdminPayload(payload))
  return toActivity(data)
}

export const deleteActivity = async (id) => apiClient.delete(`/admin/activities/${id}`)
