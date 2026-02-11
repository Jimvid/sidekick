import { useMemo, useState } from 'react'
import {
  ArrowLeftIcon,
  CheckIcon,
  PencilSimpleIcon,
  TrashIcon,
  XIcon,
} from '@phosphor-icons/react'
import { useNavigate } from '@tanstack/react-router'
import type { HabitLog } from '@/types/habits'
import { LogListItem } from '@/components/habits/LogListItem'
import { useHabits } from '@/hooks/api/habits'
import {
  useDeleteHabitLog,
  useHabitLogs,
  useUpdateHabitLog,
} from '@/hooks/api/habitLogs'

interface EditState {
  habitId: string
  date: string
  note: string
}

interface HabitLogListProps {
  limit?: number
  showFilters?: boolean
}

export function HabitLogList({ limit, showFilters = true }: HabitLogListProps) {
  const navigate = useNavigate()
  const { data: habits = [], isLoading: habitsLoading } = useHabits()
  const { data: habitLogs = [], isLoading: logsLoading } = useHabitLogs()
  const deleteHabitLog = useDeleteHabitLog()
  const updateHabitLog = useUpdateHabitLog()

  const [selectedHabitId, setSelectedHabitId] = useState('')
  const [dateFrom, setDateFrom] = useState('')
  const [dateTo, setDateTo] = useState('')
  const [editingId, setEditingId] = useState<string | null>(null)
  const [editState, setEditState] = useState<EditState>({
    habitId: '',
    date: '',
    note: '',
  })

  const habitMap = useMemo(
    () => new Map(habits.map((h) => [h.id, h])),
    [habits],
  )

  const filteredHabitLogs = useMemo(() => {
    let result = habitLogs.filter((habitLog) => habitMap.has(habitLog.habitId))

    if (showFilters) {
      if (selectedHabitId) {
        result = result.filter(
          (habitLog) => habitLog.habitId === selectedHabitId,
        )
      }
      if (dateFrom) {
        result = result.filter((habitLog) => habitLog.date >= dateFrom)
      }
      if (dateTo) {
        result = result.filter((habitLog) => habitLog.date <= dateTo)
      }
    }

    result.sort((a, b) => {
      if (a.date !== b.date) return b.date.localeCompare(a.date)
      return b.createdAt - a.createdAt
    })

    if (limit != null) {
      result = result.slice(0, limit)
    }

    return result
  }, [
    habitLogs,
    selectedHabitId,
    dateFrom,
    dateTo,
    habitMap,
    showFilters,
    limit,
  ])

  const startEditing = (habitLog: HabitLog) => {
    setEditingId(habitLog.id)
    setEditState({
      habitId: habitLog.habitId,
      date: habitLog.date,
      note: habitLog.note,
    })
  }

  const cancelEditing = () => {
    setEditingId(null)
  }

  const saveEdit = () => {
    if (!editingId || !editState.habitId || !editState.date) return
    updateHabitLog.mutate(
      {
        id: editingId,
        data: {
          habitId: editState.habitId,
          date: editState.date,
          note: editState.note,
        },
      },
      { onSuccess: () => setEditingId(null) },
    )
  }

  const isLoading = habitsLoading || logsLoading

  if (isLoading) {
    return (
      <div className="flex min-h-[calc(100vh-88px)] items-center justify-center">
        <span className="loading loading-spinner loading-lg" />
      </div>
    )
  }

  const hasFilters = selectedHabitId || dateFrom || dateTo

  const logList = (
    <>
      {filteredHabitLogs.length > 0 ? (
        <div className="space-y-2">
          {filteredHabitLogs.map((habitLog) => {
            const habit = habitMap.get(habitLog.habitId)

            if (editingId === habitLog.id) {
              const editHabit = habitMap.get(editState.habitId)
              return (
                <div
                  key={habitLog.id}
                  className="flex flex-wrap items-center gap-3 rounded-lg border border-primary/30 bg-base-200/30 px-4 py-3"
                >
                  <span
                    className="h-2.5 w-2.5 shrink-0 rounded-full"
                    style={{
                      backgroundColor: editHabit?.color ?? '#888',
                    }}
                  />
                  <select
                    className="select select-bordered select-sm"
                    value={editState.habitId}
                    onChange={(e) =>
                      setEditState((s) => ({ ...s, habitId: e.target.value }))
                    }
                  >
                    {habits.map((h) => (
                      <option key={h.id} value={h.id}>
                        {h.name}
                      </option>
                    ))}
                  </select>
                  <input
                    type="text"
                    className="input input-bordered input-sm flex-1"
                    placeholder="Note (optional)"
                    value={editState.note}
                    onChange={(e) =>
                      setEditState((s) => ({ ...s, note: e.target.value }))
                    }
                  />
                  <input
                    type="date"
                    className="input input-bordered input-sm"
                    value={editState.date}
                    onChange={(e) =>
                      setEditState((s) => ({ ...s, date: e.target.value }))
                    }
                  />
                  <button
                    className="btn btn-ghost btn-sm btn-circle shrink-0 text-success"
                    disabled={updateHabitLog.isPending}
                    onClick={saveEdit}
                  >
                    {updateHabitLog.isPending ? (
                      <span className="loading loading-spinner loading-xs" />
                    ) : (
                      <CheckIcon size={16} />
                    )}
                  </button>
                  <button
                    className="btn btn-ghost btn-sm btn-circle shrink-0"
                    onClick={cancelEditing}
                  >
                    <XIcon size={16} />
                  </button>
                </div>
              )
            }

            return (
              <LogListItem
                key={habitLog.id}
                habitLog={habitLog}
                habit={habit}
                date={new Date(habitLog.date + 'T00:00:00').toLocaleDateString(
                  undefined,
                  { weekday: 'short', month: 'short', day: 'numeric' },
                )}
                actions={
                  <>
                    <button
                      className="btn btn-ghost btn-sm btn-circle shrink-0"
                      onClick={() => startEditing(habitLog)}
                    >
                      <PencilSimpleIcon size={16} />
                    </button>
                    <button
                      className="btn btn-ghost btn-sm btn-circle shrink-0 text-error"
                      disabled={deleteHabitLog.isPending}
                      onClick={() => {
                        if (
                          window.confirm(
                            `Delete this ${habit?.name ?? 'habit'} log?`,
                          )
                        ) {
                          deleteHabitLog.mutate(habitLog.id)
                        }
                      }}
                    >
                      <TrashIcon size={16} />
                    </button>
                  </>
                }
              />
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
    </>
  )

  if (!showFilters) {
    return logList
  }

  return (
    <div className="min-h-[calc(100vh-88px)] flex flex-col p-4 sm:p-6">
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
        <div className="flex flex-col gap-3">
          <div className="flex items-end justify-between gap-3">
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
          </div>
          <div className="form-control w-full">
            <label className="label">
              <span className="label-text">From</span>
            </label>
            <input
              type="date"
              className="input input-bordered w-full"
              value={dateFrom}
              onChange={(e) => setDateFrom(e.target.value)}
            />
          </div>
          <div className="form-control w-full">
            <label className="label">
              <span className="label-text">To</span>
            </label>
            <input
              type="date"
              className="input input-bordered w-full"
              value={dateTo}
              onChange={(e) => setDateTo(e.target.value)}
            />
          </div>
        </div>

        {logList}
      </div>
    </div>
  )
}
