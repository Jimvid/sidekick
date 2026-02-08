interface BadgeProps {
  label: string
  count?: number
  color: string
}

export const Badge = ({ label, count, color }: BadgeProps) => {
  return (
    <div className="flex items-center gap-2 rounded-full border border-base-300 bg-base-100 px-3 py-1.5">
      <span
        className="h-2.5 w-2.5 rounded-full"
        style={{ backgroundColor: color }}
      />
      <span className="text-sm text-base-content">{label}</span>
      {count !== undefined && (
        <span className="text-xs text-base-content/50">{count}</span>
      )}
    </div>
  )
}
