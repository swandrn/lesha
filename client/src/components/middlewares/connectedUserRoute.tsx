import { Navigate } from 'react-router-dom'
import { useUser } from '../../hooks/useUser'

const ConnectedUserRoute = ({
  redirectPath = '/login',
  children
}: {
  redirectPath?: string
  children: React.ReactNode
}) => {
  const { user, loading } = useUser();

  if (loading) {
    return <div>Loading...</div> // or a spinner
  }

  if (!user) {
    return <Navigate to={redirectPath} replace />
  }

  return children;
}

export default ConnectedUserRoute;
