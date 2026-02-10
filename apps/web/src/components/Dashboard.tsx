import { useMemo } from 'react'
import { useNavigate } from '@tanstack/react-router'
import { ListIcon, PencilIcon, PlusIcon } from '@phosphor-icons/react'
import { useHabits } from '@/hooks/api/habits'
import { useHabitLogs } from '@/hooks/api/habitLogs'
import { LogListItem } from '@/components/habits/LogListItem'

function getRelativeDate(dateStr: string): string {
  const today = new Date()
  const date = new Date(dateStr + 'T00:00:00')
  const todayStr = today.toISOString().slice(0, 10)
  const yesterdayDate = new Date(today)
  yesterdayDate.setDate(yesterdayDate.getDate() - 1)
  const yesterdayStr = yesterdayDate.toISOString().slice(0, 10)

  if (dateStr === todayStr) return 'Today'
  if (dateStr === yesterdayStr) return 'Yesterday'

  const diffMs = today.getTime() - date.getTime()
  const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24))
  if (diffDays < 7) return `${diffDays} days ago`
  if (diffDays < 30) {
    const weeks = Math.floor(diffDays / 7)
    return `${weeks} ${weeks === 1 ? 'week' : 'weeks'} ago`
  }
  return date.toLocaleDateString(undefined, { month: 'short', day: 'numeric' })
}

export function Dashboard() {
  const navigate = useNavigate()
  const { data: habits, isLoading: habitsLoading } = useHabits()
  const { data: habitLogs, isLoading: logsLoading } = useHabitLogs()

  const isLoading = habitsLoading || logsLoading

  const recentHabitLogs = useMemo(() => {
    if (!habitLogs || !habits) return []
    const habitMap = new Map(habits.map((h) => [h.id, h]))
    return [...habitLogs]
      .sort((a, b) => {
        if (a.date !== b.date) return b.date.localeCompare(a.date)
        return b.createdAt - a.createdAt
      })
      .slice(0, 5)
      .map((habitLog) => ({
        ...habitLog,
        habit: habitMap.get(habitLog.habitId),
      }))
  }, [habitLogs, habits])

  if (isLoading) {
    return (
      <div className="flex min-h-[calc(100vh-88px)] items-center justify-center">
        <span className="loading loading-spinner loading-lg" />
      </div>
    )
  }

  const hasLogs = habitLogs && habitLogs.length > 0

  return (
    <div className="min-h-[calc(100vh-88px)] flex flex-col p-4 sm:p-6">
      <div className="mx-auto w-full max-w-6xl space-y-6">
        {/* Header */}
        <div>
          <h1 className="text-2xl font-bold text-base-content">Dashboard</h1>
          <p className="text-sm text-base-content/60">Your habit overview</p>
        </div>

        {/* Quick Actions */}
        <div className="grid grid-cols-2 gap-3 sm:grid-cols-4">
          <button
            className="flex flex-col items-center gap-2 rounded-lg border border-base-content/10 p-3 text-center transition-colors hover:bg-base-200/60 sm:flex-row sm:items-center sm:gap-3 sm:p-4 sm:text-left"
            onClick={() => navigate({ to: '/habits/create' })}
          >
            <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-primary/10 text-primary">
              <PlusIcon size={20} />
            </div>
            <div>
              <p className="text-sm font-medium text-base-content">New habit</p>
              <p className="hidden text-xs text-base-content/50 sm:block">
                Create a habit to track
              </p>
            </div>
          </button>
          <button
            className="flex flex-col items-center gap-2 rounded-lg border border-base-content/10 p-3 text-center transition-colors hover:bg-base-200/60 sm:flex-row sm:items-center sm:gap-3 sm:p-4 sm:text-left"
            onClick={() => navigate({ to: '/habits/logs' })}
          >
            <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-secondary/10 text-secondary">
              <PlusIcon size={20} />
            </div>
            <div>
              <p className="text-sm font-medium text-base-content">Log habit</p>
              <p className="hidden text-xs text-base-content/50 sm:block">
                Record today's activity
              </p>
            </div>
          </button>
          <button
            className="flex flex-col items-center gap-2 rounded-lg border border-base-content/10 p-3 text-center transition-colors hover:bg-base-200/60 sm:flex-row sm:items-center sm:gap-3 sm:p-4 sm:text-left"
            onClick={() => navigate({ to: '/habits/logs' })}
          >
            <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-info/10 text-info">
              <ListIcon size={20} />
            </div>
            <div>
              <p className="text-sm font-medium text-base-content">View logs</p>
              <p className="hidden text-xs text-base-content/50 sm:block">
                Browse all habit logs
              </p>
            </div>
          </button>
          <button
            className="flex flex-col items-center gap-2 rounded-lg border border-base-content/10 p-3 text-center transition-colors hover:bg-base-200/60 sm:flex-row sm:items-center sm:gap-3 sm:p-4 sm:text-left"
            onClick={() => navigate({ to: '/habits/manage' })}
          >
            <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-accent/10 text-accent">
              <PencilIcon size={20} />
            </div>
            <div>
              <p className="text-sm font-medium text-base-content">
                Manage habits
              </p>
              <p className="hidden text-xs text-base-content/50 sm:block">
                Edit or delete habits
              </p>
            </div>
          </button>
        </div>

        {/* Recent Activity */}
        <section>
          <h2 className="mb-3 text-sm font-semibold text-base-content/60">
            Recent Activity
          </h2>
          {hasLogs ? (
            <div className="space-y-2">
              {recentHabitLogs.map((habitLog) => (
                <LogListItem
                  key={habitLog.id}
                  habitLog={habitLog}
                  habit={habitLog.habit}
                  date={getRelativeDate(habitLog.date)}
                />
              ))}
            </div>
          ) : (
            <div className="rounded-lg border border-base-content/10 p-6 text-center">
              <p className="text-sm text-base-content/50">
                No activity yet. Log a habit to see it here.
              </p>
            </div>
          )}
        </section>
      </div>
    </div>
  )
}
