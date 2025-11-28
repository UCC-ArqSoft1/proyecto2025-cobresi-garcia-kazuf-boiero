import { BrowserRouter, Route, Routes } from 'react-router-dom'
import ScrollToHash from './components/ScrollToHash.jsx'
import ProtectedRoute from './components/ProtectedRoute.jsx'
import AddActivityPage from './pages/AddActivity.jsx'
import ActivityDetailPage from './pages/ActivityDetail.jsx'
import ActivitiesPage from './pages/Activities.jsx'
import EditActivityPage from './pages/EditActivity.jsx'
import HomePage from './pages/Home.jsx'
import LoginPage from './pages/Login.jsx'
import MyActivitiesPage from './pages/MyActivities.jsx'
import NotFoundPage from './pages/NotFound.jsx'
import SignupPage from './pages/Signup.jsx'

function App() {
  return (
    <BrowserRouter>
      <ScrollToHash />
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/activities" element={<ActivitiesPage />} />
        <Route path="/activities/:activityId" element={<ActivityDetailPage />} />
        <Route path="/login" element={<LoginPage />} />
        <Route path="/signup" element={<SignupPage />} />
        <Route element={<ProtectedRoute />}>
          <Route path="/mis-actividades" element={<MyActivitiesPage />} />
        </Route>
        <Route element={<ProtectedRoute requireAdmin />}>
          <Route path="/admin/activities/new" element={<AddActivityPage />} />
          <Route path="/admin/activities/:activityId/edit" element={<EditActivityPage />} />
        </Route>
        <Route path="*" element={<NotFoundPage />} />
      </Routes>
    </BrowserRouter>
  )
}

export default App
