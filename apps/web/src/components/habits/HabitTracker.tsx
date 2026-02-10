import { useMemo, useState } from 'react'
import {
  CaretLeftIcon,
  CaretRightIcon,
  PencilIcon,
  PlusIcon,
} from '@phosphor-icons/react'
import { useNavigate } from '@tanstack/react-router'
import type { Habit, HabitLog, QuarterInfo } from '@/types/habits'
import { Calendar } from '@/components/Calendar'
import { HabitLogList } from '@/components/habits/HabitLogList'
import { useHabits } from '@/hooks/api/habits'
import { useCreateHabitLog, useHabitLogs } from '@/hooks/api/habitLogs'

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

function groupHabitLogsByDate(habitLogs: Array<HabitLog>) {
  const grouped = new Map<string, Array<HabitLog>>()
  for (const habitLog of habitLogs) {
    const existing = grouped.get(habitLog.date)
    if (existing) {
      existing.push(habitLog)
    } else {
      grouped.set(habitLog.date, [habitLog])
    }
  }
  return grouped
}

function buildCalendarEntries(
  habitLogs: Array<HabitLog>,
  habits: Array<Habit>,
) {
  const habitMap = new Map(habits.map((h) => [h.id, h]))
  const grouped = groupHabitLogsByDate(habitLogs)
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
  const createHabitLog = useCreateHabitLog()

  const [currentQuarter, setCurrentQuarter] = useState(() => {
    const now = new Date()
    const quarter = Math.ceil((now.getMonth() + 1) / 3)
    return getQuarter(now.getFullYear(), quarter)
  })
  const [logHabitId, setLogHabitId] = useState('')
  const [logNote, setLogNote] = useState('')
  const [logDate, setLogDate] = useState(() =>
    new Date().toISOString().slice(0, 10),
  )

  const entries = useMemo(
    () => buildCalendarEntries(habitLogs, habits),
    [habitLogs, habits],
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
    createHabitLog.mutate(
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
    <div className="min-h-[calc(100vh-88px)] flex flex-col p-4 sm:p-6">
      <div className="mx-auto w-full max-w-6xl space-y-6">
        {/* Header */}
        <div className="flex items-center justify-between gap-3">
          <div>
            <h1 className="text-2xl font-bold text-base-content">Habits</h1>
            <p className="text-sm text-base-content/60">
              Track your daily habits
            </p>
          </div>
          <div className="flex gap-2">
            <button
              className="btn btn-primary btn-sm sm:btn-md"
              onClick={() => navigate({ to: '/habits/manage' })}
            >
              <PencilIcon size={20} />
              <span className="hidden sm:inline">Edit habits</span>
            </button>
            <button
              className="btn btn-primary btn-sm sm:btn-md"
              onClick={() => navigate({ to: '/habits/create' })}
            >
              <PlusIcon size={20} />
              <span className="hidden sm:inline">New habit</span>
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
          <span className="input input-bordered flex w-full shrink-0 items-center sm:w-auto">
            {new Date(logDate + 'T00:00:00').toLocaleDateString(undefined, {
              weekday: 'short',
              month: 'short',
              day: 'numeric',
            })}
          </span>
          <button
            className="btn btn-primary shrink-0"
            disabled={logHabitId === '' || createHabitLog.isPending}
            onClick={handleLogSubmit}
          >
            {createHabitLog.isPending ? (
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
              Recent Logs
            </h2>
            <button
              className="btn btn-ghost btn-sm"
              onClick={() => navigate({ to: '/habits/logs' })}
            >
              View all logs
            </button>
          </div>
          <HabitLogList limit={10} showFilters={false} />
        </div>
      </div>
    </div>
  )
}
