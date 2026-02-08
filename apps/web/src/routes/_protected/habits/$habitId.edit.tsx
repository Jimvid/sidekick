import { createFileRoute } from '@tanstack/react-router'
import { HabitForm } from '@/components/habits/HabitForm'
import { useHabit } from '@/hooks/api/habits'

export const Route = createFileRoute('/_protected/habits/$habitId/edit')({
  component: RouteComponent,
})

function RouteComponent() {
  const { habitId } = Route.useParams()
  const { data: habit, isLoading, isError } = useHabit(habitId)

  if (isLoading) {
    return (
      <div className="flex min-h-[calc(100vh-88px)] items-center justify-center">
        <span className="loading loading-spinner loading-lg" />
      </div>
    )
  }

  if (isError || !habit) {
    return (
      <div className="flex min-h-[calc(100vh-88px)] items-center justify-center">
        <p className="text-base-content/60">Habit not found</p>
      </div>
    )
  }

  return <HabitForm habit={habit} />
}
