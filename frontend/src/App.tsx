import { BrowserRouter, Routes, Route } from 'react-router-dom'
import { PageView } from './pages/PageView'
import { EditPage } from './pages/EditPage'
import { CreatePage } from './pages/CreatePage'
import { ListPages } from './pages/ListPages'
import { TrashList } from './pages/TrashList'
import { TrashPageView } from './pages/TrashPageView'
import { VersionsList } from './pages/VersionsList'
import { VersionView } from './pages/VersionView'

function NotFound() {
  return (
    <div className="min-h-screen bg-gh-bg flex items-center justify-center">
      <div className="text-center">
        <h1 className="text-4xl font-bold text-gh-text mb-2">404</h1>
        <p className="text-gh-muted mb-4">Page not found</p>
        <a href="/" className="text-gh-link hover:underline">Go home</a>
      </div>
    </div>
  )
}

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<PageView />} />
        <Route path="/list" element={<ListPages />} />
        <Route path="/create" element={<CreatePage />} />
        <Route path="/trash" element={<TrashList />} />
        <Route path="/trash/:page" element={<TrashPageView />} />
        <Route path="/edit/:page" element={<EditPage />} />
        <Route path="/:page/version" element={<VersionsList />} />
        <Route path="/:page/version/:ver" element={<VersionView />} />
        <Route path="/:page" element={<PageView />} />
        <Route path="*" element={<NotFound />} />
      </Routes>
    </BrowserRouter>
  )
}
