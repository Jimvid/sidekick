import { createFileRoute } from '@tanstack/react-router'
import { HabitTracker } from '@/components/habits/HabitTracker'

export const Route = createFileRoute('/_protected/habits/')({
  component: RouteComponent,
})

function RouteComponent() {
  return <HabitTracker />
}
