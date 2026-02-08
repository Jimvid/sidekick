export type HabitType = 'good' | 'bad'

export interface Habit {
  id: string
  name: string
  type: HabitType
  color: string
  count: number
}

export interface HabitLog {
  date: string
  habitIds: Array<string>
}

export interface QuarterInfo {
  label: string
  year: number
  quarter: number
  months: Array<number>
}
