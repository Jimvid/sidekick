import { useAuth } from '@clerk/clerk-react'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import type { Habit, HabitReq } from '@/types/habits'
import { api } from '@/lib/api'
import { notifyError, notifySuccess } from '@/lib/notify'

const HABITS_KEY = ['habits']

export function useHabits() {
  const { getToken } = useAuth()

  return useQuery({
    queryKey: HABITS_KEY,
    queryFn: async () => {
      const token = await getToken()
      return api.get<Array<Habit>>('/habits', { token })
    },
  })
}

export function useHabit(id: string) {
  const { getToken } = useAuth()

  return useQuery({
    queryKey: [...HABITS_KEY, id],
    queryFn: async () => {
      const token = await getToken()
      return api.get<Habit>(`/habits/${id}`, { token })
    },
  })
}

export function useCreateHabit() {
  const { getToken } = useAuth()
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (data: HabitReq) => {
      const token = await getToken()
      return api.post<Habit>('/habits', { token }, data)
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: HABITS_KEY })
      notifySuccess('Habit created')
    },
    onError: (err) => {
      notifyError(err.message)
    },
  })
}

export function useUpdateHabit() {
  const { getToken } = useAuth()
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async ({ id, data }: { id: string; data: HabitReq }) => {
      const token = await getToken()
      return api.put<Habit>(`/habits/${id}`, { token }, data)
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: HABITS_KEY })
      notifySuccess('Habit updated')
    },
    onError: (err) => {
      notifyError(err.message)
    },
  })
}

export function useDeleteHabit() {
  const { getToken } = useAuth()
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (id: string) => {
      const token = await getToken()
      return api.delete<{ message: string }>(`/habits/${id}`, { token })
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: HABITS_KEY })
      notifySuccess('Habit deleted')
    },
    onError: (err) => {
      notifyError(err.message)
    },
  })
}
