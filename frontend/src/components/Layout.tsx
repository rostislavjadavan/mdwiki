import { Navbar } from './Navbar'
import { useTheme } from '../useTheme'

interface Props {
  children: React.ReactNode
}

export function Layout({ children }: Props) {
  const { theme, setTheme } = useTheme()

  return (
    <div className="min-h-screen bg-gh-bg text-gh-text">
      <Navbar theme={theme} setTheme={setTheme} />
      <main className="max-w-[1012px] mx-auto px-4 py-6">
        {children}
      </main>
    </div>
  )
}
