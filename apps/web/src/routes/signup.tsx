import { SignUp } from '@clerk/clerk-react'
import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/signup')({
  component: SignupPage,
})

function SignupPage() {
  return (
    <section className="px-4 py-8 flex justify-center">
      <SignUp />
    </section>
  )
}
