const API_URL = (process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8000").replace(/\/$/, "") + "/api/v2"

export const ACCESS_TOKEN_KEY = "laku_access_token"
export const REFRESH_TOKEN_KEY = "laku_refresh_token"

export class ApiError extends Error {
  status: number
  code: number

  constructor(status: number, message: string, code = status) {
    super(message)
    this.status = status
    this.code = code
  }
}

export function getAccessToken(): string | null {
  if (typeof window === "undefined") return null
  return window.localStorage.getItem(ACCESS_TOKEN_KEY)
}

export function getRefreshToken(): string | null {
  if (typeof window === "undefined") return null
  return window.localStorage.getItem(REFRESH_TOKEN_KEY)
}

export function setTokens(access: string, refresh: string) {
  window.localStorage.setItem(ACCESS_TOKEN_KEY, access)
  window.localStorage.setItem(REFRESH_TOKEN_KEY, refresh)
}

export function clearTokens() {
  window.localStorage.removeItem(ACCESS_TOKEN_KEY)
  window.localStorage.removeItem(REFRESH_TOKEN_KEY)
}

// redirectToLogin sends the user to /login when the session cannot be refreshed.
export function redirectToLogin() {
  clearTokens()
  if (typeof window !== "undefined" && !window.location.pathname.startsWith("/login")) {
    window.location.href = "/login"
  }
}

async function tryRefreshToken(): Promise<boolean> {
  const refresh = getRefreshToken()
  if (!refresh) return false

  try {
    const res = await fetch(`${API_URL}/auth/refresh`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ refresh_token: refresh }),
    })
    const body = await res.json().catch(() => null)
    if (!res.ok || !body?.data?.access_token) return false
    setTokens(body.data.access_token, body.data.refresh_token)
    return true
  } catch {
    return false
  }
}

interface RequestOptions extends Omit<RequestInit, "body"> {
  body?: unknown
}

async function request<T>(path: string, options: RequestOptions = {}, retry = true): Promise<T> {
  const headers: Record<string, string> = {
    ...(options.body !== undefined ? { "Content-Type": "application/json" } : {}),
    ...((options.headers as Record<string, string>) ?? {}),
  }

  const token = getAccessToken()
  if (token) {
    headers["Authorization"] = `Bearer ${token}`
  }

  const res = await fetch(`${API_URL}${path}`, {
    ...options,
    headers,
    body: options.body !== undefined ? JSON.stringify(options.body) : undefined,
  })

  const body = await res.json().catch(() => null)

  if (res.status === 401 && retry && !path.startsWith("/auth/")) {
    const refreshed = await tryRefreshToken()
    if (refreshed) {
      return request<T>(path, options, false)
    }
    redirectToLogin()
    throw new ApiError(401, "Sesi berakhir, silakan login kembali")
  }

  if (!res.ok) {
    throw new ApiError(res.status, body?.msg ?? "Terjadi kesalahan pada server", body?.code)
  }

  return (body?.data ?? body) as T
}

export const http = {
  get: <T>(path: string) => request<T>(path),
  post: <T>(path: string, body?: unknown) => request<T>(path, { method: "POST", body }),
  put: <T>(path: string, body?: unknown) => request<T>(path, { method: "PUT", body }),
  patch: <T>(path: string, body?: unknown) => request<T>(path, { method: "PATCH", body }),
  delete: <T>(path: string) => request<T>(path, { method: "DELETE" }),
}

export interface Paginated<T> {
  data: T[]
  count: number | null
  pagination: {
    total: number
    page: number
    limit: number
    total_page: number
  }
}
