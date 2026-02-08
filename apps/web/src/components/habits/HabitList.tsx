import { useNavigate } from '@tanstack/react-router'
import { ArrowLeftIcon, PencilSimpleIcon, TrashIcon } from '@phosphor-icons/react'
import type { Habit } from '@/types/habits'

const HABITS: Array<Habit> = [
  { id: 'reading', name: 'Reading', description: 'Read for at least 30 minutes', color: '#22c55e', count: 24 },
  { id: 'exercise', name: 'Exercise', description: 'Any physical activity', color: '#3b82f6', count: 18 },
  { id: 'meditation', name: 'Meditation', description: 'Mindfulness practice', color: '#a855f7', count: 15 },
  { id: 'water', name: 'Water Intake', description: 'Drink 8 glasses of water', color: '#14b8a6', count: 31 },
  { id: 'junkfood', name: 'Junk Food', description: 'Track junk food consumption', color: '#ef4444', count: 12 },
  { id: 'snacking', name: 'Snacking', description: 'Track snacking between meals', color: '#f97316', count: 8 },
  { id: 'latesleep', name: 'Late Sleep', description: 'Went to bed after midnight', color: '#eab308', count: 5 },
]

export const HabitList = () => {
  const navigate = useNavigate()

  return (
    <div className="min-h-[calc(100vh-88px)] flex flex-col p-6">
      <div className="mx-auto w-full max-w-lg space-y-6">
        <div className="flex items-center gap-3">
          <button
            className="btn btn-ghost btn-sm btn-circle"
            onClick={() => navigate({ to: '/habits' })}
          >
            <ArrowLeftIcon size={20} />
          </button>
          <div>
            <h1 className="text-2xl font-bold text-base-content">
              Manage Habits
            </h1>
            <p className="text-sm text-base-content/60">
              Edit or remove your tracked habits
            </p>
          </div>
        </div>

        <div className="flex flex-col gap-2">
          {HABITS.map((habit) => (
            <div
              key={habit.id}
              className="flex items-center gap-3 rounded-lg border border-base-content/10 bg-base-100/50 p-4"
            >
              <span
                className="h-3 w-3 shrink-0 rounded-full"
                style={{ backgroundColor: habit.color }}
              />
              <div className="min-w-0 flex-1">
                <p className="font-medium text-base-content">{habit.name}</p>
                {habit.description && (
                  <p className="truncate text-sm text-base-content/50">
                    {habit.description}
                  </p>
                )}
              </div>
              <span className="text-xs text-base-content/40">
                {habit.count} logged
              </span>
              <button
                className="btn btn-ghost btn-sm btn-circle"
                onClick={() =>
                  navigate({ to: '/habits/$habitId/edit', params: { habitId: habit.id } })
                }
              >
                <PencilSimpleIcon size={18} />
              </button>
              <button
                className="btn btn-ghost btn-sm btn-circle text-error"
                onClick={() => alert('delete')}
              >
                <TrashIcon size={18} />
              </button>
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}
