import { Outlet, createFileRoute, redirect } from '@tanstack/react-router'
import { MenuBar } from '../components/MenuBar'

export const Route = createFileRoute('/_protected')({
  beforeLoad: ({ context }) => {
    if (!context.auth.isSignedIn) {
      throw redirect({
        to: '/login',
      })
    }
  },
  component: () => (
    <main className="min-h-screen pb-20">
      <Outlet />
      <MenuBar />
    </main>
  ),
})
