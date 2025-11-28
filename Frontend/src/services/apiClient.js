let authToken = null

const sanitizeBaseUrl = (value) => {
  if (!value) return 'http://localhost:8080/api'
  return value.endsWith('/') ? value.slice(0, -1) : value
}

const API_BASE_URL = sanitizeBaseUrl(import.meta.env.VITE_API_BASE_URL)

const buildURL = (path) => {
  if (/^https?:\/\//i.test(path)) {
    return path
  }
  return `${API_BASE_URL}${path.startsWith('/') ? path : `/${path}`}`
}

const parseJSON = async (response) => {
  const contentType = response.headers.get('content-type')
  if (contentType && contentType.includes('application/json')) {
    return response.json()
  }
  return null
}

const request = async (path, options = {}) => {
  const { method = 'GET', body, headers = {}, ...rest } = options
  const init = {
    method,
    headers: {
      Accept: 'application/json',
      ...headers,
    },
    ...rest,
  }

  if (body !== undefined && body !== null) {
    if (body instanceof FormData) {
      init.body = body
    } else if (typeof body === 'string') {
      init.body = body
      if (!init.headers['Content-Type']) {
        init.headers['Content-Type'] = 'application/json'
      }
    } else {
      init.body = JSON.stringify(body)
      init.headers['Content-Type'] = 'application/json'
    }
  }

  if (authToken) {
    init.headers.Authorization = `Bearer ${authToken}`
  }

  let response
  try {
    response = await fetch(buildURL(path), init)
  } catch (networkError) {
    const error = new Error('No se pudo conectar con el backend.')
    error.details = networkError
    throw error
  }

  const payload = await parseJSON(response)

  if (!response.ok) {
    const errorMessage = payload?.error || payload?.message || 'Error al comunicarse con la API'
    const error = new Error(errorMessage)
    error.status = response.status
    error.details = payload?.details ?? payload
    error.code = payload?.code
    throw error
  }

  if (payload === null) {
    return null
  }

  return payload?.data ?? payload
}

const apiClient = {
  get: (path, options) => request(path, { ...options, method: 'GET' }),
  post: (path, body, options) => request(path, { ...options, method: 'POST', body }),
  put: (path, body, options) => request(path, { ...options, method: 'PUT', body }),
  delete: (path, options) => request(path, { ...options, method: 'DELETE' }),
}

export const setAuthToken = (token) => {
  authToken = token || null
}

export const getApiBaseUrl = () => API_BASE_URL

export default apiClient
