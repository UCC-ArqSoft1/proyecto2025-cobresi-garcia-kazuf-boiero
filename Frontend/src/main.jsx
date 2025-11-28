import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './styles/style.css'
import App from './App.jsx'
import { AuthProvider } from './contexts/AuthContext.jsx'
import { ActivitiesProvider } from './contexts/ActivitiesContext.jsx'

createRoot(document.getElementById('root')).render(
  <StrictMode>
    <AuthProvider>
      <ActivitiesProvider>
        <App />
      </ActivitiesProvider>
    </AuthProvider>
  </StrictMode>,
)
