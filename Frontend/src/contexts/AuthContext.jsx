import { createContext, useContext, useEffect, useMemo, useState } from 'react'
import { login as loginRequest, register as registerRequest } from '../services/authService.js'
import { setAuthToken } from '../services/apiClient.js'

const AuthContext = createContext(null)

const storageKey = 'gad.auth'

export const AuthProvider = ({ children }) => {
  const [authState, setAuthState] = useState(() => {
    try {
      const stored = localStorage.getItem(storageKey)
      if (!stored) {
        return { user: null, token: null }
      }
      const parsed = JSON.parse(stored)
      return {
        user: parsed.user ?? null,
        token: parsed.token ?? null,
      }
    } catch {
      return { user: null, token: null }
    }
  })
  const [authReady, setAuthReady] = useState(false)

  useEffect(() => {
    setAuthToken(authState.token)
    setAuthReady(true)
  }, [authState.token])

  const persistAuthState = (state) => {
    setAuthState(state)
    // Keep apiClient token in sync immediately (avoid race on first request after login/logout)
    setAuthToken(state.token)
    if (state.user && state.token) {
      localStorage.setItem(storageKey, JSON.stringify(state))
    } else {
      localStorage.removeItem(storageKey)
    }
  }

  const login = async ({ email, password }) => {
    const data = await loginRequest({ email, password })
    persistAuthState({
      user: data.user,
      token: data.token,
    })
    return data.user
  }

  const register = async ({ name, email, password }) => {
    const user = await registerRequest({ name, email, password })
    return user
  }

  const logout = () => {
    persistAuthState({ user: null, token: null })
  }

  const value = useMemo(
    () => ({
      user: authState.user,
      token: authState.token,
      isAuthenticated: Boolean(authState.user && authState.token),
      isAdmin: authState.user?.role === 'admin',
      authReady,
      login,
      register,
      logout,
    }),
    [authState.user, authState.token, authReady],
  )

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

export const useAuth = () => {
  const context = useContext(AuthContext)
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}
