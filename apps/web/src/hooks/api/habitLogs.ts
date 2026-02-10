import { useAuth } from '@clerk/clerk-react'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import type { HabitLog, HabitLogReq } from '@/types/habits'
import { api } from '@/lib/api'
import { notifyError, notifySuccess } from '@/lib/notify'

const HABIT_LOGS_KEY = ['habit-logs']

export function useHabitLogs() {
  const { getToken } = useAuth()

  return useQuery({
    queryKey: HABIT_LOGS_KEY,
    queryFn: async () => {
      const token = await getToken()
      return api.get<Array<HabitLog>>('/habit-logs', { token })
    },
  })
}

export function useHabitLog(id: string) {
  const { getToken } = useAuth()

  return useQuery({
    queryKey: [...HABIT_LOGS_KEY, id],
    queryFn: async () => {
      const token = await getToken()
      return api.get<HabitLog>(`/habit-logs/${id}`, { token })
    },
  })
}

export function useCreateHabitLog() {
  const { getToken } = useAuth()
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (data: HabitLogReq) => {
      const token = await getToken()
      return api.post<HabitLog>('/habit-logs', { token }, data)
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: HABIT_LOGS_KEY })
      notifySuccess('Habit logged')
    },
    onError: (err) => {
      notifyError(err.message)
    },
  })
}

export function useUpdateHabitLog() {
  const { getToken } = useAuth()
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async ({ id, data }: { id: string; data: HabitLogReq }) => {
      const token = await getToken()
      return api.put<HabitLog>(`/habit-logs/${id}`, { token }, data)
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: HABIT_LOGS_KEY })
      notifySuccess('Log updated')
    },
    onError: (err) => {
      notifyError(err.message)
    },
  })
}

export function useDeleteHabitLog() {
  const { getToken } = useAuth()
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (id: string) => {
      const token = await getToken()
      return api.delete<{ message: string }>(`/habit-logs/${id}`, { token })
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: HABIT_LOGS_KEY })
      notifySuccess('Log deleted')
    },
    onError: (err) => {
      notifyError(err.message)
    },
  })
}
