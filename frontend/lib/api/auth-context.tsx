"use client"

import { createContext, useCallback, useContext, useEffect, useMemo, useState, type ReactNode } from "react"
import { login as apiLogin, logout as apiLogout, me, type AuthUser } from "./auth"
import { getAccessToken, getRefreshToken } from "./http"

interface AuthContextValue {
  user: AuthUser | null
  loading: boolean
  login: (email: string, password: string) => Promise<AuthUser>
  logout: () => Promise<void>
}

const AuthContext = createContext<AuthContextValue | undefined>(undefined)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<AuthUser | null>(null)
  // No token → not loading (guest); token → loading until /auth/me resolves.
  const [loading, setLoading] = useState(() => getAccessToken() !== null)

  useEffect(() => {
    if (!getAccessToken()) return
    let cancelled = false
    me()
      .then((data) => {
        if (!cancelled) setUser(data)
      })
      .catch(() => {
        if (!cancelled) setUser(null)
      })
      .finally(() => {
        if (!cancelled) setLoading(false)
      })
    return () => {
      cancelled = true
    }
  }, [])

  const login = useCallback(async (email: string, password: string) => {
    const result = await apiLogin(email, password)
    setUser(result.user)
    return result.user
  }, [])

  const logout = useCallback(async () => {
    await apiLogout(getRefreshToken())
    setUser(null)
  }, [])

  const value = useMemo(
    () => ({ user, loading, login, logout }),
    [user, loading, login, logout],
  )

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

export function useAuth(): AuthContextValue {
  const ctx = useContext(AuthContext)
  if (!ctx) {
    throw new Error("useAuth must be used within AuthProvider")
  }
  return ctx
}
