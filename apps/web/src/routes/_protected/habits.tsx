import { Outlet, createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/_protected/habits')({
  component: RouteComponent,
})

function RouteComponent() {
  return <Outlet />
}
