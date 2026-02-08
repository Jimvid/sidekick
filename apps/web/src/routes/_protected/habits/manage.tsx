import { createFileRoute } from '@tanstack/react-router'
import { HabitList } from '@/components/habits/HabitList'

export const Route = createFileRoute('/_protected/habits/manage')({
  component: RouteComponent,
})

function RouteComponent() {
  return <HabitList />
}
