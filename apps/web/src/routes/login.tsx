import { SignIn } from '@clerk/clerk-react'
import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/login')({
  component: RouteComponent,
})

function RouteComponent() {
  return (
    <section className="px-4 py-8 flex justify-center">
      <SignIn signUpUrl="/signup" />
    </section>
  )
}
