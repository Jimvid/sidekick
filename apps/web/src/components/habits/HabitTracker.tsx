import { useState } from 'react'
import {
  CaretLeftIcon,
  CaretRightIcon,
  PencilIcon,
  PlusIcon,
} from '@phosphor-icons/react'
import type { Habit, HabitLog, QuarterInfo } from '@/types/habits'
import { Calendar } from '@/components/Calendar'

const HABITS: Array<Habit> = [
  { id: 'reading', name: 'Reading', type: 'good', color: '#22c55e', count: 24 },
  {
    id: 'exercise',
    name: 'Exercise',
    type: 'good',
    color: '#3b82f6',
    count: 18,
  },
  {
    id: 'meditation',
    name: 'Meditation',
    type: 'good',
    color: '#a855f7',
    count: 15,
  },
  {
    id: 'water',
    name: 'Water Intake',
    type: 'good',
    color: '#14b8a6',
    count: 31,
  },
  {
    id: 'junkfood',
    name: 'Junk Food',
    type: 'bad',
    color: '#ef4444',
    count: 12,
  },
  { id: 'snacking', name: 'Snacking', type: 'bad', color: '#f97316', count: 8 },
  {
    id: 'latesleep',
    name: 'Late Sleep',
    type: 'bad',
    color: '#eab308',
    count: 5,
  },
]

const HABIT_LOGS: Array<HabitLog> = [
  { date: '2024-01-02', habitIds: ['reading', 'exercise'] },
  { date: '2024-01-05', habitIds: ['reading', 'meditation', 'water'] },
  { date: '2024-01-08', habitIds: ['exercise', 'water', 'junkfood'] },
  { date: '2024-01-12', habitIds: ['reading', 'meditation'] },
  { date: '2024-01-15', habitIds: ['water', 'exercise', 'reading'] },
  { date: '2024-01-18', habitIds: ['meditation', 'water', 'snacking'] },
  { date: '2024-01-22', habitIds: ['reading', 'junkfood'] },
  { date: '2024-01-25', habitIds: ['exercise', 'water', 'reading'] },
  { date: '2024-01-29', habitIds: ['reading', 'meditation', 'water'] },
  { date: '2024-01-30', habitIds: ['water', 'exercise'] },
  { date: '2024-02-01', habitIds: ['reading', 'water', 'meditation'] },
  { date: '2024-02-06', habitIds: ['exercise', 'reading', 'water'] },
  { date: '2024-02-08', habitIds: ['meditation', 'water'] },
  { date: '2024-02-12', habitIds: ['reading', 'exercise', 'snacking'] },
  { date: '2024-02-13', habitIds: ['water', 'meditation', 'reading'] },
  { date: '2024-02-19', habitIds: ['reading', 'water', 'junkfood'] },
  { date: '2024-02-20', habitIds: ['exercise', 'meditation', 'water'] },
  { date: '2024-02-22', habitIds: ['reading', 'latesleep'] },
  { date: '2024-02-26', habitIds: ['water', 'reading'] },
  { date: '2024-02-27', habitIds: ['meditation', 'exercise'] },
  { date: '2024-03-01', habitIds: ['reading', 'water'] },
  { date: '2024-03-04', habitIds: ['exercise', 'reading', 'meditation'] },
  { date: '2024-03-07', habitIds: ['water', 'reading', 'latesleep'] },
  { date: '2024-03-11', habitIds: ['reading', 'exercise', 'water'] },
  { date: '2024-03-14', habitIds: ['meditation', 'water', 'reading'] },
  { date: '2024-03-18', habitIds: ['water', 'exercise', 'junkfood'] },
  { date: '2024-03-19', habitIds: ['reading', 'meditation'] },
  { date: '2024-03-21', habitIds: ['water', 'reading', 'snacking'] },
  { date: '2024-03-25', habitIds: ['exercise', 'water', 'latesleep'] },
  { date: '2024-03-28', habitIds: ['reading', 'meditation', 'water'] },
  { date: '2024-03-29', habitIds: ['reading', 'water'] },
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
  const entries: Record<string, { colors: Array<string> }> = {}

  for (const log of logs) {
    entries[log.date] = {
      colors: log.habitIds
        .map((id) => habitMap.get(id)?.color)
        .filter((c): c is string => c !== undefined),
    }
  }

  return entries
}

export const HabitTracker = () => {
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
            <button className="btn btn-primary">
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
