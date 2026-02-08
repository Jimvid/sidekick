import { useForm } from '@tanstack/react-form'
import { useNavigate } from '@tanstack/react-router'
import { ArrowLeftIcon } from '@phosphor-icons/react'
import type { Habit } from '@/types/habits'
import { useCreateHabit, useUpdateHabit } from '@/hooks/api/habits'

const COLOR_OPTIONS = [
  { label: 'Green', value: '#22c55e' },
  { label: 'Blue', value: '#3b82f6' },
  { label: 'Purple', value: '#a855f7' },
  { label: 'Teal', value: '#14b8a6' },
  { label: 'Red', value: '#ef4444' },
  { label: 'Orange', value: '#f97316' },
  { label: 'Yellow', value: '#eab308' },
  { label: 'Pink', value: '#ec4899' },
]

interface HabitFormProps {
  habit?: Habit
}

export const HabitForm = ({ habit }: HabitFormProps) => {
  const navigate = useNavigate()
  const isEditing = habit != null
  const createHabit = useCreateHabit()
  const updateHabit = useUpdateHabit()
  const isPending = createHabit.isPending || updateHabit.isPending

  const form = useForm({
    defaultValues: {
      name: habit?.name ?? '',
      description: habit?.description ?? '',
      color: habit?.color ?? '#22c55e',
    },
    onSubmit: ({ value }) => {
      const data = {
        name: value.name,
        description: value.description,
        color: value.color,
      }

      if (isEditing) {
        updateHabit.mutate(
          { id: habit.id, data },
          { onSuccess: () => navigate({ to: '/habits' }) },
        )
      } else {
        createHabit.mutate(data, {
          onSuccess: () => navigate({ to: '/habits' }),
        })
      }
    },
  })

  return (
    <div className="min-h-[calc(100vh-88px)] flex flex-col p-6">
      <div className="mx-auto w-full max-w-lg space-y-6">
        <div className="flex items-center gap-3">
          <button
            className="btn btn-ghost btn-sm btn-circle"
            onClick={() => navigate({ to: '/habits' })}
          >
            <ArrowLeftIcon size={20} />
          </button>
          <div>
            <h1 className="text-2xl font-bold text-base-content">
              {isEditing ? 'Edit Habit' : 'New Habit'}
            </h1>
            <p className="text-sm text-base-content/60">
              {isEditing
                ? 'Update your habit details'
                : 'Create a new habit to track'}
            </p>
          </div>
        </div>

        <form
          onSubmit={(e) => {
            e.preventDefault()
            e.stopPropagation()
            form.handleSubmit()
          }}
          className="space-y-6"
        >
          <form.Field
            name="name"
            validators={{
              onChange: ({ value }) =>
                value.trim().length === 0 ? 'Name is required' : undefined,
            }}
          >
            {(field) => (
              <label className="form-control w-full">
                <div className="label">
                  <span className="label-text">Name</span>
                </div>
                <input
                  type="text"
                  placeholder="e.g. Working Out"
                  className="input input-bordered w-full"
                  value={field.state.value}
                  onChange={(e) => field.handleChange(e.target.value)}
                  onBlur={field.handleBlur}
                />
                {field.state.meta.errors.length > 0 && (
                  <div className="label">
                    <span className="label-text-alt text-error">
                      {field.state.meta.errors.join(', ')}
                    </span>
                  </div>
                )}
              </label>
            )}
          </form.Field>

          <form.Field name="description">
            {(field) => (
              <label className="form-control w-full">
                <div className="label">
                  <span className="label-text">Description</span>
                </div>
                <textarea
                  placeholder="e.g. Track what I eat so I can be more mindful"
                  className="textarea textarea-bordered w-full"
                  rows={3}
                  value={field.state.value}
                  onChange={(e) => field.handleChange(e.target.value)}
                  onBlur={field.handleBlur}
                />
              </label>
            )}
          </form.Field>

          <form.Field name="color">
            {(field) => (
              <label className="form-control w-full">
                <div className="label">
                  <span className="label-text">Color</span>
                </div>
                <div className="flex flex-wrap items-center gap-2">
                  {COLOR_OPTIONS.map((opt) => (
                    <button
                      key={opt.value}
                      type="button"
                      onClick={() => field.handleChange(opt.value)}
                      className={`flex h-10 w-10 items-center justify-center rounded-full border-2 transition-transform hover:scale-110 ${
                        field.state.value === opt.value
                          ? 'border-base-content scale-110'
                          : 'border-transparent'
                      }`}
                    >
                      <span
                        className="h-7 w-7 rounded-full"
                        style={{ backgroundColor: opt.value }}
                      />
                    </button>
                  ))}
                  <label className="flex h-10 w-10 cursor-pointer items-center justify-center rounded-full border-2 border-dashed border-base-content/30 transition-transform hover:scale-110">
                    <input
                      type="color"
                      className="sr-only"
                      value={field.state.value}
                      onChange={(e) => field.handleChange(e.target.value)}
                    />
                    <span
                      className="h-7 w-7 rounded-full"
                      style={{ backgroundColor: field.state.value }}
                    />
                  </label>
                </div>
              </label>
            )}
          </form.Field>

          <div className="flex gap-2 pt-4">
            <button
              type="button"
              className="btn btn-ghost flex-1 border border-base-content/10"
              onClick={() => navigate({ to: '/habits' })}
            >
              Cancel
            </button>
            <form.Subscribe selector={(state) => [state.canSubmit]}>
              {([canSubmit]) => (
                <button
                  type="submit"
                  className="btn btn-primary flex-1"
                  disabled={!canSubmit || isPending}
                >
                  {isPending
                    ? <span className="loading loading-spinner loading-sm" />
                    : isEditing ? 'Save Changes' : 'Create Habit'}
                </button>
              )}
            </form.Subscribe>
          </div>
        </form>
      </div>
    </div>
  )
}
