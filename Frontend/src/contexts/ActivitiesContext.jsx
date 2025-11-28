import { createContext, useCallback, useContext, useEffect, useMemo, useState } from 'react'
import { useAuth } from './AuthContext.jsx'
import {
  createActivity as createActivityRequest,
  deleteActivity as deleteActivityRequest,
  enrollInActivity as enrollInActivityRequest,
  getActivity as getActivityRequest,
  getMyActivities as getMyActivitiesRequest,
  listActivities,
  unenrollFromActivity as unenrollFromActivityRequest,
  updateActivity as updateActivityRequest,
} from '../services/activitiesService.js'

const ActivitiesContext = createContext(null)

export const ActivitiesProvider = ({ children }) => {
  const { token, user, isAuthenticated, authReady } = useAuth()
  const [activities, setActivities] = useState([])
  const [activitiesLoading, setActivitiesLoading] = useState(true)
  const [activitiesError, setActivitiesError] = useState(null)
  const [myActivities, setMyActivities] = useState([])
  const [myActivitiesLoading, setMyActivitiesLoading] = useState(false)

  const fetchActivities = useCallback(async (filters) => {
    setActivitiesLoading(true)
    setActivitiesError(null)
    try {
      const data = await listActivities(filters)
      setActivities(data)
      return data
    } catch (error) {
      setActivitiesError(error)
      throw error
    } finally {
      setActivitiesLoading(false)
    }
  }, [])

  const refreshMyActivities = useCallback(async () => {
    if (!token) {
      setMyActivities([])
      return []
    }

    setMyActivitiesLoading(true)
    try {
      const data = await getMyActivitiesRequest()
      setMyActivities(data)
      return data
    } finally {
      setMyActivitiesLoading(false)
    }
  }, [token])

  useEffect(() => {
    fetchActivities().catch(() => {})
  }, [fetchActivities])

  useEffect(() => {
    // Esperar a que AuthContext haya finalizado la inicialización
    if (!authReady) return

    if (!isAuthenticated || !token || !user) {
      setMyActivities([])
      return
    }

    refreshMyActivities().catch(() => {})
  }, [authReady, isAuthenticated, token, user, refreshMyActivities])

  const createActivity = useCallback(async (payload) => {
    const created = await createActivityRequest(payload)
    setActivities((prev) => [created, ...prev])
    return created
  }, [])

  const updateActivity = useCallback(async (id, payload) => {
    const updated = await updateActivityRequest(id, payload)
    setActivities((prev) => prev.map((activity) => (activity.id === updated.id ? updated : activity)))
    return updated
  }, [])

  const deleteActivity = useCallback(async (id) => {
    await deleteActivityRequest(id)
    setActivities((prev) => prev.filter((activity) => activity.id !== Number(id)))
  }, [])

  const loadActivityById = useCallback(
    async (id, options = {}) => {
      const numericId = Number(id)
      if (!options.force) {
        const cached = activities.find((activity) => activity.id === numericId)
        if (cached) {
          return cached
        }
      }

      const remote = await getActivityRequest(id)
      setActivities((prev) => {
        const exists = prev.some((activity) => activity.id === remote.id)
        if (!exists) {
          return [...prev, remote]
        }
        return prev.map((activity) => (activity.id === remote.id ? remote : activity))
      })
      return remote
    },
    [activities],
  )

  const enrollInActivity = useCallback(
    async (activityId) => {
      await enrollInActivityRequest(activityId)
      const updated = await loadActivityById(activityId, { force: true }).catch(() => null)
      await refreshMyActivities()
      return updated
    },
    [loadActivityById, refreshMyActivities],
  )

  const unenrollFromActivity = useCallback(
    async (activityId) => {
      await unenrollFromActivityRequest(activityId)
      const updated = await loadActivityById(activityId, { force: true }).catch(() => null)
      await refreshMyActivities()
      return updated
    },
    [loadActivityById, refreshMyActivities],
  )

  const value = useMemo(
    () => ({
      activities,
      activitiesLoading,
      activitiesError,
      refreshActivities: fetchActivities,
      createActivity,
      updateActivity,
      deleteActivity,
      enrollInActivity,
      unenrollFromActivity,
      loadActivityById,
      myActivities,
      myActivitiesLoading,
      refreshMyActivities,
    }),
    [
      activities,
      activitiesLoading,
      activitiesError,
      fetchActivities,
      createActivity,
      updateActivity,
      deleteActivity,
      enrollInActivity,
      unenrollFromActivity,
      loadActivityById,
      myActivities,
      myActivitiesLoading,
      refreshMyActivities,
    ],
  )

  return <ActivitiesContext.Provider value={value}>{children}</ActivitiesContext.Provider>
}

export const useActivities = () => {
  const context = useContext(ActivitiesContext)
  if (!context) {
    throw new Error('useActivities must be used within an ActivitiesProvider')
  }
  return context
}
