import { createFileRoute } from '@tanstack/react-router'
import { HabitLogList } from '@/components/habits/HabitLogList'

export const Route = createFileRoute('/_protected/habits/logs')({
  component: RouteComponent,
})

function RouteComponent() {
  return <HabitLogList />
}
