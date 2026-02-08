import { createFileRoute } from '@tanstack/react-router'
import { HabitForm } from '@/components/habits/HabitForm'

const MOCK_HABITS: Partial<Record<string, { id: string; name: string; description: string; color: string; count: number }>> = {
  reading: { id: 'reading', name: 'Reading', description: 'Read for at least 30 minutes', color: '#22c55e', count: 24 },
  exercise: { id: 'exercise', name: 'Exercise', description: 'Any physical activity', color: '#3b82f6', count: 18 },
  meditation: { id: 'meditation', name: 'Meditation', description: 'Mindfulness practice', color: '#a855f7', count: 15 },
  water: { id: 'water', name: 'Water Intake', description: 'Drink 8 glasses of water', color: '#14b8a6', count: 31 },
  junkfood: { id: 'junkfood', name: 'Junk Food', description: 'Track junk food consumption', color: '#ef4444', count: 12 },
  snacking: { id: 'snacking', name: 'Snacking', description: 'Track snacking between meals', color: '#f97316', count: 8 },
  latesleep: { id: 'latesleep', name: 'Late Sleep', description: 'Went to bed after midnight', color: '#eab308', count: 5 },
}

export const Route = createFileRoute('/_protected/habits/$habitId/edit')({
  component: RouteComponent,
})

function RouteComponent() {
  const { habitId } = Route.useParams()
  const habit = MOCK_HABITS[habitId]

  if (!habit) {
    return (
      <div className="flex min-h-[calc(100vh-88px)] items-center justify-center">
        <p className="text-base-content/60">Habit not found</p>
      </div>
    )
  }

  return <HabitForm habit={habit} />
}
