import { useEffect, useState } from 'react'
import { CaretLeftIcon, CaretRightIcon } from '@phosphor-icons/react'

const MONTH_NAMES = [
  'January',
  'February',
  'March',
  'April',
  'May',
  'June',
  'July',
  'August',
  'September',
  'October',
  'November',
  'December',
]

const DAY_HEADERS = ['Su', 'Mo', 'Tu', 'We', 'Th', 'Fr', 'Sa']

interface CalendarEntry {
  colors: Array<string>
}

interface CalendarProps {
  year: number
  months: Array<number>
  entries: Record<string, CalendarEntry>
  selectedDate?: string
  onDateSelect?: (date: string) => void
}

function getMonthGrid(year: number, month: number) {
  const firstDay = new Date(year, month, 1).getDay()
  const daysInMonth = new Date(year, month + 1, 0).getDate()
  const weeks: Array<Array<number | null>> = []
  let currentWeek: Array<number | null> = Array(firstDay).fill(null)

  for (let day = 1; day <= daysInMonth; day++) {
    currentWeek.push(day)
    if (currentWeek.length === 7) {
      weeks.push(currentWeek)
      currentWeek = []
    }
  }

  if (currentWeek.length > 0) {
    while (currentWeek.length < 7) {
      currentWeek.push(null)
    }
    weeks.push(currentWeek)
  }

  return weeks
}

function formatDateKey(year: number, month: number, day: number) {
  return `${year}-${String(month + 1).padStart(2, '0')}-${String(day).padStart(2, '0')}`
}

function MonthGrid({
  year,
  month,
  entries,
  hideTitle,
  selectedDate,
  onDateSelect,
}: {
  year: number
  month: number
  entries: Record<string, CalendarEntry>
  hideTitle?: boolean
  selectedDate?: string
  onDateSelect?: (date: string) => void
}) {
  const weeks = getMonthGrid(year, month)

  return (
    <div className="rounded-lg border border-base-content/10 p-4">
      {!hideTitle && (
        <h3 className="mb-3 text-sm font-semibold text-base-content">
          {MONTH_NAMES[month]} {year}
        </h3>
      )}
      <div className="grid grid-cols-7 gap-1">
        {DAY_HEADERS.map((header) => (
          <div
            key={header}
            className="py-1 text-center text-xs font-medium text-base-content/40"
          >
            {header}
          </div>
        ))}
        {weeks.flat().map((day, i) => {
          if (day === null) {
            return <div key={`empty-${i}`} className="aspect-square" />
          }

          const dateKey = formatDateKey(year, month, day)
          const entry = entries[dateKey] as CalendarEntry | undefined
          const hasEntries = entry != null && entry.colors.length > 0
          const isSelected = dateKey === selectedDate

          return (
            <div
              key={dateKey}
              onClick={() => onDateSelect?.(dateKey)}
              className={`flex aspect-square cursor-pointer flex-col items-center justify-center gap-0.5 rounded border transition-all hover:scale-110 hover:bg-base-200/60 ${
                isSelected
                  ? 'border-primary ring-1 ring-primary'
                  : 'border-base-content/10'
              } ${hasEntries ? 'bg-base-200' : 'bg-base-100/50'}`}
            >
              <span className="text-xs text-base-content/70">{day}</span>
              {hasEntries && (
                <div className="flex gap-0.5">
                  {entry.colors.slice(0, 4).map((color, ci) => (
                    <span
                      key={ci}
                      className="h-1.5 w-1.5 rounded-full"
                      style={{ backgroundColor: color }}
                    />
                  ))}
                </div>
              )}
            </div>
          )
        })}
      </div>
    </div>
  )
}

export const Calendar = ({
  year,
  months,
  entries,
  selectedDate,
  onDateSelect,
}: CalendarProps) => {
  const [mobileIndex, setMobileIndex] = useState(() => {
    const currentMonth = new Date().getMonth()
    const idx = months.indexOf(currentMonth)
    return idx >= 0 ? idx : 0
  })

  useEffect(() => {
    const currentMonth = new Date().getMonth()
    const idx = months.indexOf(currentMonth)
    setMobileIndex(idx >= 0 ? idx : 0)
  }, [year, months[0]])

  return (
    <>
      {/* Mobile: single month with navigation */}
      <div className="lg:hidden">
        <div className="mb-3 flex items-center justify-between">
          <button
            className="btn btn-ghost btn-sm btn-circle"
            disabled={mobileIndex === 0}
            onClick={() => setMobileIndex((i) => i - 1)}
          >
            <CaretLeftIcon size={18} />
          </button>
          <span className="text-sm font-semibold text-base-content">
            {MONTH_NAMES[months[mobileIndex]]} {year}
          </span>
          <button
            className="btn btn-ghost btn-sm btn-circle"
            disabled={mobileIndex === months.length - 1}
            onClick={() => setMobileIndex((i) => i + 1)}
          >
            <CaretRightIcon size={18} />
          </button>
        </div>
        <MonthGrid
          year={year}
          month={months[mobileIndex]}
          entries={entries}
          hideTitle
          selectedDate={selectedDate}
          onDateSelect={onDateSelect}
        />
      </div>

      {/* Desktop: all months */}
      <div className="hidden gap-6 lg:grid lg:grid-cols-3">
        {months.map((month) => (
          <MonthGrid
            key={month}
            year={year}
            month={month}
            entries={entries}
            selectedDate={selectedDate}
            onDateSelect={onDateSelect}
          />
        ))}
      </div>
    </>
  )
}
