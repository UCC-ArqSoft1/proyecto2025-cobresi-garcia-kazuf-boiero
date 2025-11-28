const Notification = ({ type = 'success', message, onClose }) => {
  if (!message) return null

  return (
    <div className={`notification notification-${type}`}>
      <p>{message}</p>
      {onClose ? (
        <button type="button" className="notification-close" onClick={onClose} aria-label="Cerrar">
          ×
        </button>
      ) : null}
    </div>
  )
}

export default Notification
