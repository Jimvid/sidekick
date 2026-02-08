import { useState } from 'react'
import {
  CaretLeftIcon,
  CaretRightIcon,
  PencilIcon,
  PencilSimpleIcon,
  PlusIcon,
  TrashIcon,
} from '@phosphor-icons/react'
import { useNavigate } from '@tanstack/react-router'
import type { Habit, HabitLog, QuarterInfo } from '@/types/habits'
import { Calendar } from '@/components/Calendar'

const HABITS: Array<Habit> = [
  { id: 'reading', name: 'Reading', description: 'Read for at least 30 minutes', color: '#22c55e', count: 24 },
  { id: 'exercise', name: 'Exercise', description: 'Any physical activity', color: '#3b82f6', count: 18 },
  { id: 'meditation', name: 'Meditation', description: 'Mindfulness practice', color: '#a855f7', count: 15 },
  { id: 'water', name: 'Water Intake', description: 'Drink 8 glasses of water', color: '#14b8a6', count: 31 },
  { id: 'junkfood', name: 'Junk Food', description: 'Track junk food consumption', color: '#ef4444', count: 12 },
  { id: 'snacking', name: 'Snacking', description: 'Track snacking between meals', color: '#f97316', count: 8 },
  { id: 'latesleep', name: 'Late Sleep', description: 'Went to bed after midnight', color: '#eab308', count: 5 },
]

const HABIT_LOGS: Array<HabitLog> = [
  { date: '2024-01-02', entries: [{ habitId: 'reading' }, { habitId: 'exercise' }] },
  { date: '2024-01-05', entries: [{ habitId: 'reading' }, { habitId: 'meditation' }, { habitId: 'water' }] },
  { date: '2024-01-08', entries: [{ habitId: 'exercise' }, { habitId: 'water' }, { habitId: 'junkfood', note: 'Pizza for lunch' }] },
  { date: '2024-01-12', entries: [{ habitId: 'reading' }, { habitId: 'meditation' }] },
  { date: '2024-01-15', entries: [{ habitId: 'water' }, { habitId: 'exercise' }, { habitId: 'reading' }] },
  { date: '2024-01-18', entries: [{ habitId: 'meditation' }, { habitId: 'water' }, { habitId: 'snacking', note: 'Chips at 3pm' }] },
  { date: '2024-01-22', entries: [{ habitId: 'reading' }, { habitId: 'junkfood', note: 'Burger and fries' }] },
  { date: '2024-01-25', entries: [{ habitId: 'exercise' }, { habitId: 'water' }, { habitId: 'reading' }] },
  { date: '2024-01-29', entries: [{ habitId: 'reading' }, { habitId: 'meditation' }, { habitId: 'water' }] },
  { date: '2024-01-30', entries: [{ habitId: 'water' }, { habitId: 'exercise' }] },
  { date: '2024-02-01', entries: [{ habitId: 'reading' }, { habitId: 'water' }, { habitId: 'meditation' }] },
  { date: '2024-02-06', entries: [{ habitId: 'exercise' }, { habitId: 'reading' }, { habitId: 'water' }] },
  { date: '2024-02-08', entries: [{ habitId: 'meditation' }, { habitId: 'water' }] },
  { date: '2024-02-12', entries: [{ habitId: 'reading' }, { habitId: 'exercise' }, { habitId: 'snacking', note: 'Cookies' }] },
  { date: '2024-02-13', entries: [{ habitId: 'water' }, { habitId: 'meditation' }, { habitId: 'reading' }] },
  { date: '2024-02-19', entries: [{ habitId: 'reading' }, { habitId: 'water' }, { habitId: 'junkfood', note: 'Fried chicken' }] },
  { date: '2024-02-20', entries: [{ habitId: 'exercise' }, { habitId: 'meditation' }, { habitId: 'water' }] },
  { date: '2024-02-22', entries: [{ habitId: 'reading' }, { habitId: 'latesleep', note: '2am coding' }] },
  { date: '2024-02-26', entries: [{ habitId: 'water' }, { habitId: 'reading' }] },
  { date: '2024-02-27', entries: [{ habitId: 'meditation' }, { habitId: 'exercise' }] },
  { date: '2024-03-01', entries: [{ habitId: 'reading' }, { habitId: 'water' }] },
  { date: '2024-03-04', entries: [{ habitId: 'exercise' }, { habitId: 'reading' }, { habitId: 'meditation' }] },
  { date: '2024-03-07', entries: [{ habitId: 'water' }, { habitId: 'reading' }, { habitId: 'latesleep', note: 'Netflix binge' }] },
  { date: '2024-03-11', entries: [{ habitId: 'reading' }, { habitId: 'exercise' }, { habitId: 'water' }] },
  { date: '2024-03-14', entries: [{ habitId: 'meditation' }, { habitId: 'water' }, { habitId: 'reading' }] },
  { date: '2024-03-18', entries: [{ habitId: 'water' }, { habitId: 'exercise' }, { habitId: 'junkfood', note: 'Donuts at work' }] },
  { date: '2024-03-19', entries: [{ habitId: 'reading' }, { habitId: 'meditation' }] },
  { date: '2024-03-21', entries: [{ habitId: 'water' }, { habitId: 'reading' }, { habitId: 'snacking', note: 'Late night popcorn' }] },
  { date: '2024-03-25', entries: [{ habitId: 'exercise' }, { habitId: 'water' }, { habitId: 'latesleep' }] },
  { date: '2024-03-28', entries: [{ habitId: 'reading' }, { habitId: 'meditation' }, { habitId: 'water' }] },
  { date: '2024-03-29', entries: [{ habitId: 'reading' }, { habitId: 'water' }] },
]

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

function buildCalendarEntries(logs: Array<HabitLog>, habits: Array<Habit>) {
  const habitMap = new Map(habits.map((h) => [h.id, h]))
  const calendarEntries: Record<string, { colors: Array<string> }> = {}

  for (const log of logs) {
    calendarEntries[log.date] = {
      colors: log.entries
        .map((e) => habitMap.get(e.habitId)?.color)
        .filter((c): c is string => c !== undefined),
    }
  }

  return calendarEntries
}

export const HabitTracker = () => {
  const navigate = useNavigate()
  const [currentQuarter, setCurrentQuarter] = useState(() =>
    getQuarter(2024, 1),
  )
  const [logHabitId, setLogHabitId] = useState('')
  const [logNote, setLogNote] = useState('')
  const [logDate, setLogDate] = useState(() => new Date().toISOString().slice(0, 10))

  const entries = buildCalendarEntries(HABIT_LOGS, HABITS)
  const habitMap = new Map(HABITS.map((h) => [h.id, h]))
  const sortedLogs = [...HABIT_LOGS].sort((a, b) =>
    b.date.localeCompare(a.date),
  )
  const latestDate = sortedLogs[0]?.date ?? ''
  const weekAgo = new Date(latestDate + 'T00:00:00')
  weekAgo.setDate(weekAgo.getDate() - 7)
  const weekAgoStr = weekAgo.toISOString().slice(0, 10)
  const recentLogs = sortedLogs.filter((log) => log.date >= weekAgoStr)

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
              Current Streak:{' '}
              <strong className="text-base-content">7 days</strong>
            </span>
            <span>
              Total Logged:{' '}
              <strong className="text-base-content">113 habits</strong>
            </span>
          </div>
        </div>
        {/* Log Entry Form */}
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
              {HABITS.map((h) => (
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
            disabled={logHabitId === ''}
            onClick={() => {
              const habit = habitMap.get(logHabitId)
              alert(
                JSON.stringify(
                  { habitId: logHabitId, habit: habit?.name, note: logNote || undefined, date: logDate },
                  null,
                  2,
                ),
              )
              setLogHabitId('')
              setLogNote('')
            }}
          >
            <PlusIcon size={20} />
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
          <h2 className="mb-3 text-lg font-semibold text-base-content">
            Last 7 Days
          </h2>
          <div className="flex flex-col gap-4">
            {recentLogs.map((log) => (
              <div key={log.date}>
                <p className="mb-2 text-sm font-medium text-base-content/50">
                  {new Date(log.date + 'T00:00:00').toLocaleDateString(
                    undefined,
                    { weekday: 'long', month: 'short', day: 'numeric' },
                  )}
                </p>
                <div className="flex flex-col gap-1.5">
                  {log.entries.map((entry, i) => {
                    const habit = habitMap.get(entry.habitId)
                    if (!habit) return null
                    return (
                      <div
                        key={`${entry.habitId}-${i}`}
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
                          {entry.note && (
                            <p className="truncate text-sm text-base-content/40">
                              {entry.note}
                            </p>
                          )}
                        </div>
                        <button
                          className="btn btn-ghost btn-sm btn-circle shrink-0"
                          onClick={() =>
                            alert(
                              `Edit: ${habit.name} on ${log.date}${entry.note ? `\nNote: ${entry.note}` : ''}`,
                            )
                          }
                        >
                          <PencilSimpleIcon size={16} />
                        </button>
                        <button
                          className="btn btn-ghost btn-sm btn-circle shrink-0 text-error"
                          onClick={() => alert('delete')}
                        >
                          <TrashIcon size={16} />
                        </button>
                      </div>
                    )
                  })}
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  )
}
