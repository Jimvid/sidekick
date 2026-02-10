import { useMemo, useState } from 'react'
import {
  CaretLeftIcon,
  CaretRightIcon,
  PencilIcon,
  PlusIcon,
  TrashIcon,
} from '@phosphor-icons/react'
import { useNavigate } from '@tanstack/react-router'
import type { Habit, HabitLog, QuarterInfo } from '@/types/habits'
import { Calendar } from '@/components/Calendar'
import { useHabits } from '@/hooks/api/habits'
import { useCreateHabitLog, useDeleteHabitLog, useHabitLogs } from '@/hooks/api/habitLogs'

function getQuarter(year: number, quarter: number): QuarterInfo {
  const months = [
    (quarter - 1) * 3,
    (quarter - 1) * 3 + 1,
    (quarter - 1) * 3 + 2,
  ]
  return {
    label: `Q${quarter} ${year}`,
    year,
    quarter,
    months,
  }
}

function groupLogsByDate(logs: Array<HabitLog>) {
  const grouped = new Map<string, Array<HabitLog>>()
  for (const log of logs) {
    const existing = grouped.get(log.date)
    if (existing) {
      existing.push(log)
    } else {
      grouped.set(log.date, [log])
    }
  }
  return grouped
}

function buildCalendarEntries(logs: Array<HabitLog>, habits: Array<Habit>) {
  const habitMap = new Map(habits.map((h) => [h.id, h]))
  const grouped = groupLogsByDate(logs)
  const calendarEntries: Record<string, { colors: Array<string> }> = {}

  for (const [date, dateLogs] of grouped) {
    calendarEntries[date] = {
      colors: dateLogs
        .map((l) => habitMap.get(l.habitId)?.color)
        .filter((c): c is string => c !== undefined),
    }
  }

  return calendarEntries
}

export const HabitTracker = () => {
  const navigate = useNavigate()
  const { data: habits = [], isLoading: habitsLoading } = useHabits()
  const { data: habitLogs = [], isLoading: logsLoading } = useHabitLogs()
  const createLog = useCreateHabitLog()
  const deleteLog = useDeleteHabitLog()

  const [currentQuarter, setCurrentQuarter] = useState(() => {
    const now = new Date()
    const quarter = Math.ceil((now.getMonth() + 1) / 3)
    return getQuarter(now.getFullYear(), quarter)
  })
  const [logHabitId, setLogHabitId] = useState('')
  const [logNote, setLogNote] = useState('')
  const [logDate, setLogDate] = useState(() => new Date().toISOString().slice(0, 10))

  const habitMap = useMemo(() => new Map(habits.map((h) => [h.id, h])), [habits])
  const entries = useMemo(() => buildCalendarEntries(habitLogs, habits), [habitLogs, habits])

  const recentLogs = useMemo(() => {
    const sorted = [...habitLogs].sort((a, b) => b.date.localeCompare(a.date))
    const latestDate = sorted[0]?.date ?? ''
    if (!latestDate) return new Map<string, Array<HabitLog>>()
    const weekAgo = new Date(latestDate + 'T00:00:00')
    weekAgo.setDate(weekAgo.getDate() - 7)
    const weekAgoStr = weekAgo.toISOString().slice(0, 10)
    const recent = sorted.filter((log) => log.date >= weekAgoStr)
    return groupLogsByDate(recent)
  }, [habitLogs])

  const sortedRecentDates = useMemo(
    () => [...recentLogs.keys()].sort((a, b) => b.localeCompare(a)),
    [recentLogs],
  )

  const navigateQuarter = (direction: -1 | 1) => {
    setCurrentQuarter((prev) => {
      let newQuarter = prev.quarter + direction
      let newYear = prev.year
      if (newQuarter < 1) {
        newQuarter = 4
        newYear--
      } else if (newQuarter > 4) {
        newQuarter = 1
        newYear++
      }
      return getQuarter(newYear, newQuarter)
    })
  }

  const goToToday = () => {
    const now = new Date()
    const quarter = Math.ceil((now.getMonth() + 1) / 3)
    setCurrentQuarter(getQuarter(now.getFullYear(), quarter))
  }

  const handleLogSubmit = () => {
    createLog.mutate(
      { habitId: logHabitId, date: logDate, note: logNote },
      {
        onSuccess: () => {
          setLogHabitId('')
          setLogNote('')
        },
      },
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

  return (
    <div className="min-h-[calc(100vh-88px)] flex flex-col p-6">
      <div className="mx-auto w-full max-w-6xl space-y-6">
        {/* Header */}
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-2xl font-bold text-base-content">
              Habit Tracker
            </h1>
            <p className="text-sm text-base-content/60">
              Track your daily habits
            </p>
          </div>
          <div className="flex gap-2">
            <button
              className="btn btn-primary"
              onClick={() => navigate({ to: '/habits/manage' })}
            >
              <PencilIcon size={20} />
              Edit habits
            </button>
            <button
              className="btn btn-primary"
              onClick={() => navigate({ to: '/habits/create' })}
            >
              <PlusIcon size={20} />
              New habit
            </button>
          </div>
        </div>

        {/* Navigation */}
        <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
          <div className="flex items-center gap-2">
            <span className="text-lg font-semibold text-base-content">
              {currentQuarter.label}
            </span>
            <button
              className="btn btn-ghost btn-sm btn-circle"
              onClick={() => navigateQuarter(-1)}
            >
              <CaretLeftIcon size={18} />
            </button>
            <button
              className="btn btn-ghost btn-sm btn-circle"
              onClick={() => navigateQuarter(1)}
            >
              <CaretRightIcon size={18} />
            </button>
            <button
              className="btn btn-ghost btn-sm text-primary"
              onClick={goToToday}
            >
              Today
            </button>
          </div>
          <div className="flex items-center gap-4 text-sm text-base-content/60">
            <span>
              Total Logged:{' '}
              <strong className="text-base-content">{habitLogs.length}</strong>
            </span>
          </div>
        </div>
        {/* Log Entry Form */}
        <div className="flex flex-col gap-3 rounded-lg border border-base-content/10 bg-base-100/50 p-4 sm:flex-row sm:items-end">
          <div className="form-control flex-1">
            <label className="label">
              <span className="label-text">Habit</span>
            </label>
            <select
              className="select select-bordered w-full"
              value={logHabitId}
              onChange={(e) => setLogHabitId(e.target.value)}
            >
              <option value="">Select a habit...</option>
              {habits.map((h) => (
                <option key={h.id} value={h.id}>
                  {h.name}
                </option>
              ))}
            </select>
          </div>
          <div className="form-control flex-1">
            <label className="label">
              <span className="label-text">Note (optional)</span>
            </label>
            <input
              type="text"
              className="input input-bordered w-full"
              placeholder="e.g. 30 min run"
              value={logNote}
              onChange={(e) => setLogNote(e.target.value)}
            />
          </div>
          <span className="input input-bordered flex w-auto min-w-36 shrink-0 items-center">
            {new Date(logDate + 'T00:00:00').toLocaleDateString(undefined, {
              weekday: 'short',
              month: 'short',
              day: 'numeric',
            })}
          </span>
          <button
            className="btn btn-primary shrink-0"
            disabled={logHabitId === '' || createLog.isPending}
            onClick={handleLogSubmit}
          >
            {createLog.isPending ? (
              <span className="loading loading-spinner loading-sm" />
            ) : (
              <PlusIcon size={20} />
            )}
            Log
          </button>
        </div>

        <Calendar
          year={currentQuarter.year}
          months={currentQuarter.months}
          entries={entries}
          selectedDate={logDate}
          onDateSelect={setLogDate}
        />

        {/* Recent Log Entries */}
        <div>
          <div className="mb-3 flex items-center justify-between">
            <h2 className="text-lg font-semibold text-base-content">
              Last 7 Days
            </h2>
            <button
              className="btn btn-ghost btn-sm"
              onClick={() => navigate({ to: '/habits/logs' })}
            >
              View all logs
            </button>
          </div>
          <div className="flex flex-col gap-4">
            {sortedRecentDates.map((date) => {
              const dateLogs = recentLogs.get(date) ?? []
              return (
                <div key={date}>
                  <p className="mb-2 text-sm font-medium text-base-content/50">
                    {new Date(date + 'T00:00:00').toLocaleDateString(
                      undefined,
                      { weekday: 'long', month: 'short', day: 'numeric' },
                    )}
                  </p>
                  <div className="flex flex-col gap-1.5">
                    {dateLogs.map((log) => {
                      const habit = habitMap.get(log.habitId)
                      if (!habit) return null
                      return (
                        <div
                          key={log.id}
                          className="flex items-center gap-3 rounded-lg border border-base-content/10 bg-base-100/50 px-4 py-3"
                        >
                          <span
                            className="h-3 w-3 shrink-0 rounded-full"
                            style={{ backgroundColor: habit.color }}
                          />
                          <div className="min-w-0 flex-1">
                            <span className="text-sm font-medium text-base-content">
                              {habit.name}
                            </span>
                            {log.note && (
                              <p className="truncate text-sm text-base-content/40">
                                {log.note}
                              </p>
                            )}
                          </div>
                          <button
                            className="btn btn-ghost btn-sm btn-circle shrink-0 text-error"
                            disabled={deleteLog.isPending}
                            onClick={() => {
                              if (window.confirm(`Delete this ${habit.name} log?`)) {
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
                </div>
              )
            })}
          </div>
        </div>
      </div>
    </div>
  )
}
