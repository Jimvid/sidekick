import { createFileRoute, redirect } from '@tanstack/react-router'
import { LoginPage } from './login'

export const Route = createFileRoute('/')({
  beforeLoad: ({ context }) => {
    if (context.auth.isSignedIn) {
      throw redirect({ to: '/dashboard' })
    }
  },
  component: LoginPage,
})
