import apiClient from './apiClient.js'

export const login = async ({ email, password }) =>
  apiClient.post('/auth/login', {
    email,
    password,
  })

export const register = async ({ name, email, password }) =>
  apiClient.post('/auth/register', {
    name,
    email,
    password,
  })
