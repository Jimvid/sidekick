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

export interface HabitLog {
  id: string
  habitId: string
  date: string
  note: string
  createdAt: number
  updatedAt: number
}

export interface HabitLogReq {
  habitId: string
  date: string
  note: string
}

export interface QuarterInfo {
  label: string
  year: number
  quarter: number
  months: Array<number>
}
