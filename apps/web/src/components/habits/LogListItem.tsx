import type { ReactNode } from 'react'
import type { Habit, HabitLog } from '@/types/habits'

interface LogListItemProps {
  habitLog: HabitLog
  habit?: Habit
  date?: string
  actions?: ReactNode
}

export function LogListItem({
  habitLog,
  habit,
  date,
  actions,
}: LogListItemProps) {
  return (
    <div className="flex items-center gap-3 rounded-lg border border-base-content/10 px-4 py-3">
      <span
        className="h-2.5 w-2.5 shrink-0 rounded-full"
        style={{ backgroundColor: habit?.color ?? '#888' }}
      />
      <div className="min-w-0 flex-1">
        <span className="text-sm font-medium text-base-content">
          {habit?.name ?? 'Unknown'}
        </span>
        {habitLog.note && (
          <p className="truncate text-sm text-base-content/50">
            {habitLog.note}
          </p>
        )}
      </div>
      {date && (
        <span className="shrink-0 text-xs text-base-content/40">{date}</span>
      )}
      {actions && (
        <div className="flex shrink-0 items-center gap-1">{actions}</div>
      )}
    </div>
  )
}
