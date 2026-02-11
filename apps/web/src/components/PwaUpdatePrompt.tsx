import { useEffect, useState } from 'react'

export const PwaUpdatePrompt = () => {
  const [needRefresh, setNeedRefresh] = useState(false)
  const [updateSW, setUpdateSW] = useState<(() => Promise<void>) | null>(null)

  useEffect(() => {
    if (!('serviceWorker' in navigator)) return

    import('virtual:pwa-register').then(({ registerSW }) => {
      const update = registerSW({
        onNeedRefresh() {
          setNeedRefresh(true)
        },
      })
      setUpdateSW(() => update)
    })
  }, [])

  const handleUpdate = () => {
    updateSW?.()
  }

  const handleDismiss = () => {
    setNeedRefresh(false)
  }

  if (!needRefresh) return null

  return (
    <div className="fixed bottom-20 left-0 right-0 z-50 flex justify-center px-4">
      <div className="alert alert-info shadow-lg max-w-md">
        <span>A new version is available.</span>
        <div className="flex gap-2">
          <button className="btn btn-sm btn-ghost" onClick={handleDismiss}>
            Later
          </button>
          <button className="btn btn-sm btn-primary" onClick={handleUpdate}>
            Update
          </button>
        </div>
      </div>
    </div>
  )
}
