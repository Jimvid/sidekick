import { RouterProvider } from '@tanstack/react-router'
import { useAuth } from '@clerk/clerk-react'
import { QueryClientProvider } from '@tanstack/react-query'
import { router } from './router'
import { queryClient } from './queryClient'
import { Bounce, ToastContainer } from 'react-toastify'

export default function App() {
  const auth = useAuth()

  if (!auth.isLoaded) return null

  return (
    <QueryClientProvider client={queryClient}>
      <RouterProvider router={router} context={{ auth }} />
      <ToastContainer
        position="top-center"
        autoClose={5000}
        hideProgressBar={false}
        newestOnTop={false}
        closeOnClick={false}
        rtl={false}
        pauseOnFocusLoss
        draggable
        pauseOnHover
        theme="dark"
        transition={Bounce}
      />
    </QueryClientProvider>
  )
}
