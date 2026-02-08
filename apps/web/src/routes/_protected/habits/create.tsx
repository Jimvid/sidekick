import { createFileRoute } from '@tanstack/react-router'
import { HabitForm } from '@/components/habits/HabitForm'

export const Route = createFileRoute('/_protected/habits/create')({
  component: RouteComponent,
})

function RouteComponent() {
  return <HabitForm />
}
