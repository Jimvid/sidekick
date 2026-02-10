import { useMemo, useState } from 'react'
import { ArrowLeftIcon, TrashIcon } from '@phosphor-icons/react'
import { useNavigate } from '@tanstack/react-router'
import { useHabits } from '@/hooks/api/habits'
import { useDeleteHabitLog, useHabitLogs } from '@/hooks/api/habitLogs'

export function HabitLogList() {
  const navigate = useNavigate()
  const { data: habits = [], isLoading: habitsLoading } = useHabits()
  const { data: logs = [], isLoading: logsLoading } = useHabitLogs()
  const deleteLog = useDeleteHabitLog()

  const [selectedHabitId, setSelectedHabitId] = useState('')
  const [dateFrom, setDateFrom] = useState('')
  const [dateTo, setDateTo] = useState('')

  const habitMap = useMemo(() => new Map(habits.map((h) => [h.id, h])), [habits])

  const filteredLogs = useMemo(() => {
    let result = logs.filter((log) => habitMap.has(log.habitId))

    if (selectedHabitId) {
      result = result.filter((log) => log.habitId === selectedHabitId)
    }
    if (dateFrom) {
      result = result.filter((log) => log.date >= dateFrom)
    }
    if (dateTo) {
      result = result.filter((log) => log.date <= dateTo)
    }

    return result.sort((a, b) => {
      if (a.date !== b.date) return b.date.localeCompare(a.date)
      return b.createdAt - a.createdAt
    })
  }, [logs, selectedHabitId, dateFrom, dateTo, habitMap])

  const isLoading = habitsLoading || logsLoading

  if (isLoading) {
    return (
      <div className="flex min-h-[calc(100vh-88px)] items-center justify-center">
        <span className="loading loading-spinner loading-lg" />
      </div>
    )
  }

  const hasFilters = selectedHabitId || dateFrom || dateTo

  return (
    <div className="min-h-[calc(100vh-88px)] flex flex-col p-6">
      <div className="mx-auto w-full max-w-6xl space-y-6">
        {/* Header */}
        <div className="flex items-center gap-3">
          <button
            className="btn btn-ghost btn-sm btn-circle"
            onClick={() => navigate({ to: '/habits' })}
          >
            <ArrowLeftIcon size={20} />
          </button>
          <div>
            <h1 className="text-2xl font-bold text-base-content">Habit Logs</h1>
            <p className="text-sm text-base-content/60">
              Browse and manage your logged habits
            </p>
          </div>
        </div>

        {/* Filter bar */}
        <div className="flex flex-wrap items-end justify-between gap-3">
          <div className="form-control w-full sm:w-48">
            <label className="label">
              <span className="label-text">Habit</span>
            </label>
            <select
              className="select select-bordered w-full"
              value={selectedHabitId}
              onChange={(e) => setSelectedHabitId(e.target.value)}
            >
              <option value="">All habits</option>
              {habits.map((h) => (
                <option key={h.id} value={h.id}>
                  {h.name}
                </option>
              ))}
            </select>
          </div>
          <div className="flex flex-wrap items-end gap-3">
            <button
              className={`btn btn-ghost btn-sm${hasFilters ? '' : ' invisible'}`}
              onClick={() => {
                setSelectedHabitId('')
                setDateFrom('')
                setDateTo('')
              }}
            >
              Clear filters
            </button>
            <div className="form-control">
              <label className="label">
                <span className="label-text">From</span>
              </label>
              <input
                type="date"
                className="input input-bordered"
                value={dateFrom}
                onChange={(e) => setDateFrom(e.target.value)}
              />
            </div>
            <div className="form-control">
              <label className="label">
                <span className="label-text">To</span>
              </label>
              <input
                type="date"
                className="input input-bordered"
                value={dateTo}
                onChange={(e) => setDateTo(e.target.value)}
              />
            </div>
          </div>
        </div>

        {/* Log list */}
        {filteredLogs.length > 0 ? (
          <div className="space-y-2">
            {filteredLogs.map((log) => {
              const habit = habitMap.get(log.habitId)
              return (
                <div
                  key={log.id}
                  className="flex items-center gap-3 rounded-lg border border-base-content/10 px-4 py-3"
                >
                  <span
                    className="h-2.5 w-2.5 shrink-0 rounded-full"
                    style={{
                      backgroundColor: habit?.color ?? '#888',
                    }}
                  />
                  <span className="text-sm font-medium text-base-content">
                    {habit?.name ?? 'Unknown'}
                  </span>
                  {log.note && (
                    <span className="truncate text-sm text-base-content/50">
                      {log.note}
                    </span>
                  )}
                  <span className="ml-auto shrink-0 text-xs text-base-content/40">
                    {new Date(log.date + 'T00:00:00').toLocaleDateString(
                      undefined,
                      { weekday: 'short', month: 'short', day: 'numeric' },
                    )}
                  </span>
                  <button
                    className="btn btn-ghost btn-sm btn-circle shrink-0 text-error"
                    disabled={deleteLog.isPending}
                    onClick={() => {
                      if (
                        window.confirm(
                          `Delete this ${habit?.name ?? 'habit'} log?`,
                        )
                      ) {
                        deleteLog.mutate(log.id)
                      }
                    }}
                  >
                    <TrashIcon size={16} />
                  </button>
                </div>
              )
            })}
          </div>
        ) : (
          <div className="rounded-lg border border-base-content/10 p-6 text-center">
            <p className="text-sm text-base-content/50">
              {hasFilters
                ? 'No logs match the current filters.'
                : 'No habit logs yet. Log a habit from the tracker to see it here.'}
            </p>
          </div>
        )}
      </div>
    </div>
  )
}
