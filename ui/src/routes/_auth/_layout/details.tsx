import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/_auth/_layout/details')({
  component: RouteComponent,
})

function RouteComponent() {
  return <div>Hello "/_auth/_layout/details"!</div>
}
