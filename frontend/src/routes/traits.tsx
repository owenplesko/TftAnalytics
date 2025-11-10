import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/traits')({
  component: RouteComponent,
})

function RouteComponent() {
  return <div>Hello "/traits"!</div>
}
