import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/')({
  component: PlacesPage,
})

function PlacesPage() {
  return <p className="text-muted-foreground">No places yet.</p>
}
