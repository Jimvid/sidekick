import { useState } from 'react'
import {
  CaretLeftIcon,
  CaretRightIcon,
  PencilIcon,
  PlusIcon,
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

  const entries = buildCalendarEntries(HABIT_LOGS, HABITS)

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
            <button className="btn btn-primary">
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
        <Calendar
          year={currentQuarter.year}
          months={currentQuarter.months}
          entries={entries}
        />
      </div>
    </div>
  )
}
