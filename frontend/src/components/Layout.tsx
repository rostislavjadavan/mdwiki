import { Navbar } from './Navbar'

interface Props {
  children: React.ReactNode
}

export function Layout({ children }: Props) {
  return (
    <div className="min-h-screen bg-gh-bg text-gh-text">
      <Navbar />
      <main className="max-w-[1012px] mx-auto px-4 py-6">
        {children}
      </main>
    </div>
  )
}
