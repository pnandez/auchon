import { Outlet, createRootRoute } from '@tanstack/react-router'

export const Route = createRootRoute({
  component: RootLayout,
})

function RootLayout() {
  return (
    <div className="min-h-dvh bg-background text-foreground">
      <header className="border-b px-4 py-3">
        <h1 className="text-xl font-semibold">Auchon</h1>
      </header>
      <main className="p-4">
        <Outlet />
      </main>
    </div>
  )
}
