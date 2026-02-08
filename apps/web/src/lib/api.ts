import { config } from '@/config/config'

type RequestOptions = {
  token: string | null
}

async function request<T>(
  method: string,
  path: string,
  opts: RequestOptions,
  body?: unknown,
): Promise<T> {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
  }

  if (opts.token) {
    headers['Authorization'] = `Bearer ${opts.token}`
  }

  const res = await fetch(`${config.apiUrl}${path}`, {
    method,
    headers,
    body: body != null ? JSON.stringify(body) : undefined,
  })

  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: res.statusText }))
    throw new Error(err.error ?? 'Request failed')
  }

  return res.json() as Promise<T>
}

export const api = {
  get: <T>(path: string, opts: RequestOptions) =>
    request<T>('GET', path, opts),
  post: <T>(path: string, opts: RequestOptions, body: unknown) =>
    request<T>('POST', path, opts, body),
  put: <T>(path: string, opts: RequestOptions, body: unknown) =>
    request<T>('PUT', path, opts, body),
  delete: <T>(path: string, opts: RequestOptions) =>
    request<T>('DELETE', path, opts),
}
