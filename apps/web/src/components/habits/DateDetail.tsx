import type { Habit, HabitLog } from '@/types/habits'

interface DateDetailProps {
  date: string
  logs: Array<HabitLog>
  habitMap: Map<string, Habit>
}

export const DateDetail = ({ date, logs, habitMap }: DateDetailProps) => {
  if (logs.length === 0) return null

  return (
    <div className="rounded-lg border border-base-content/10 bg-base-100/50 p-4">
      <h3 className="mb-3 text-sm font-semibold text-base-content">
        {new Date(date + 'T00:00:00').toLocaleDateString(undefined, {
          weekday: 'long',
          month: 'long',
          day: 'numeric',
        })}
      </h3>
      <div className="flex flex-col gap-2">
        {logs.map((log) => {
          const habit = habitMap.get(log.habitId)
          if (!habit) return null
          return (
            <div key={log.id} className="flex items-center gap-2">
              <span
                className="h-3 w-3 shrink-0 rounded-full"
                style={{ backgroundColor: habit.color }}
              />
              <span className="text-sm font-medium text-base-content">
                {habit.name}
              </span>
              {log.note && (
                <span className="text-sm text-base-content/60">
                  — {log.note}
                </span>
              )}
            </div>
          )
        })}
      </div>
    </div>
  )
}
