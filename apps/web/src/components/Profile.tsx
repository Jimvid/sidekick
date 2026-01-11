import { useAuth, useUser } from '@clerk/clerk-react'
import { CalendarIcon, EnvelopeIcon, SignOutIcon, UserIcon } from '@phosphor-icons/react'

export const Profile = () => {
  const { user } = useUser()
  const { signOut } = useAuth()

  const handleSignOut = () => {
    signOut()
  }

  if (!user) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <span className="loading loading-spinner loading-lg"></span>
      </div>
    )
  }

  return (
    <div className="min-h-[calc(100vh-88px)] flex flex-col p-6">
      <div className="flex flex-col w-full flex-1 max-w-2xl mx-auto space-y-6">
        {/* Header */}
        <div className="text-center">
          <h1 className="text-2xl font-bold text-base-content mb-2">Profile</h1>
          <p className="text-base-content/70">Manage your account settings</p>
        </div>

        <div className="card w-full bg-base-200 shadow-lg">
          <div className="card-body">
            <div className="flex items-center gap-4">
              <div className="avatar">
                <div className="w-full rounded-full ring ring-primary ring-offset-base-100 ring-offset-2">
                  {user.imageUrl ? (
                    <img
                      className="max-w-24"
                      src={user.imageUrl}
                      alt="Profile"
                    />
                  ) : (
                    <div className="bg-base-300 flex items-center justify-center">
                      <UserIcon size={40} className="text-base-content/60" />
                    </div>
                  )}
                </div>
              </div>

              <div className="flex-1 space-y-2">
                <h2 className="card-title text-base-content">
                  {user.fullName || user.firstName || 'Anonymous User'}
                </h2>

                {user.primaryEmailAddress && (
                  <div className="flex items-center gap-2 text-base-content/70">
                    <EnvelopeIcon size={16} />
                    <span className="text-sm">
                      {user.primaryEmailAddress.emailAddress}
                    </span>
                  </div>
                )}

                {user.createdAt && (
                  <div className="flex items-center gap-2 text-base-content/70">
                    <CalendarIcon size={16} />
                    <span className="text-sm">
                      Member since{' '}
                      {new Date(user.createdAt).toLocaleDateString()}
                    </span>
                  </div>
                )}
              </div>
            </div>
          </div>
        </div>
        <button onClick={handleSignOut} className="btn btn-primary w-full">
          <SignOutIcon size={20} />
          Sign Out
        </button>
      </div>
    </div>
  )
}
