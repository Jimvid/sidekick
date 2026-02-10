import { useNavigate } from '@tanstack/react-router'
import {
  ArrowLeftIcon,
  PencilSimpleIcon,
  TrashIcon,
} from '@phosphor-icons/react'
import { useDeleteHabit, useHabits } from '@/hooks/api/habits'

export const HabitList = () => {
  const navigate = useNavigate()
  const { data: habits = [], isLoading } = useHabits()
  const deleteHabit = useDeleteHabit()

  if (isLoading) {
    return (
      <div className="flex min-h-[calc(100vh-88px)] items-center justify-center">
        <span className="loading loading-spinner loading-lg" />
      </div>
    )
  }

  return (
    <div className="min-h-[calc(100vh-88px)] flex flex-col p-6">
      <div className="mx-auto w-full max-w-6xl space-y-6">
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
          {habits.map((habit) => (
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
              <button
                className="btn btn-ghost btn-sm btn-circle"
                onClick={() =>
                  navigate({
                    to: '/habits/$habitId/edit',
                    params: { habitId: habit.id },
                  })
                }
              >
                <PencilSimpleIcon size={18} />
              </button>
              <button
                className="btn btn-ghost btn-sm btn-circle text-error"
                disabled={deleteHabit.isPending}
                onClick={() => {
                  if (window.confirm(`Delete "${habit.name}"?`)) {
                    deleteHabit.mutate(habit.id)
                  }
                }}
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
