import { render, screen } from '@testing-library/react'
import { RouterProvider, createMemoryHistory } from '@tanstack/react-router'
import { describe, expect, it } from 'vitest'
import { buildRouter } from '@/app/router'

describe('App shell', () => {
  it('shows the app name on the home page', async () => {
    const router = buildRouter({ history: createMemoryHistory({ initialEntries: ['/'] }) })
    render(<RouterProvider router={router} />)
    expect(await screen.findByRole('heading', { name: 'Auchon' })).toBeInTheDocument()
  })
})
