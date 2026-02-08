export interface Habit {
  id: string
  name: string
  description: string
  color: string
  createdAt: number
  updatedAt: number
}

export interface HabitReq {
  name: string
  description: string
  color: string
}

export interface HabitLogEntry {
  habitId: string
  note?: string
}

export interface HabitLog {
  date: string
  entries: Array<HabitLogEntry>
}

export interface QuarterInfo {
  label: string
  year: number
  quarter: number
  months: Array<number>
}
