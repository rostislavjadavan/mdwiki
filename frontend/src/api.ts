export interface Page {
  filename: string
  name: string
  content: string
  rawContent: string
}

export interface PageInfo {
  filename: string
  modTime: string
  version: number
}

export interface SearchResult {
  query: string
  filenames: SearchHit[]
  pageContent: SearchHit[]
}

export interface SearchHit {
  filename: string
  modTime: string
  score: number
  preview: string
}

async function request<T>(url: string, options?: RequestInit): Promise<T> {
  const res = await fetch(url, options)
  if (!res.ok) {
    const err = await res.json().catch(() => ({ message: res.statusText }))
    throw new Error(err.message || res.statusText)
  }
  if (res.status === 204) return undefined as T
  return res.json()
}

export const api = {
  listPages: () => request<PageInfo[]>('/api/pages'),
  getPage: (filename: string) => request<Page>(`/api/pages/${filename}`),
  createPage: (filename: string) =>
    request<Page>('/api/pages', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ filename }),
    }),
  updatePage: (filename: string, data: { filename: string; content: string }) =>
    request<Page>(`/api/pages/${filename}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    }),
  deletePage: (filename: string) =>
    request<void>(`/api/pages/${filename}`, { method: 'DELETE' }),
  search: (query: string) =>
    request<SearchResult>(`/api/search?query=${encodeURIComponent(query)}`),
  listTrash: () => request<PageInfo[]>('/api/trash'),
  getTrashPage: (filename: string) => request<Page>(`/api/trash/${filename}`),
  restoreTrash: (filename: string) =>
    request<void>(`/api/trash/${filename}/restore`, { method: 'POST' }),
  emptyTrash: () => request<void>('/api/trash', { method: 'DELETE' }),
  listVersions: (filename: string) =>
    request<PageInfo[]>(`/api/pages/${filename}/versions`),
  getVersion: (ver: string) => request<Page>(`/api/versions/${ver}`),
  restoreVersion: (ver: string) =>
    request<void>(`/api/versions/${ver}/restore`, { method: 'POST' }),
}
