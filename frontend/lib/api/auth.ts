import { http, setTokens, clearTokens, getAccessToken } from "./http"

export interface AuthUser {
  id: string
  tenantId: string | null
  name: string
  email: string
  avatar: string
  role: string
  isAdmin: boolean
}

export interface LoginResult {
  accessToken: string
  refreshToken: string
  user: AuthUser
}

export async function login(email: string, password: string): Promise<LoginResult> {
  const data = await http.post<{
    access_token: string
    refresh_token: string
    user: {
      id: string
      tenant_id: string | null
      name: string
      email: string
      avatar: string
      role: string
      is_admin: boolean
    }
  }>("/auth/login", { email, password })

  setTokens(data.access_token, data.refresh_token)

  return {
    accessToken: data.access_token,
    refreshToken: data.refresh_token,
    user: {
      id: data.user.id,
      tenantId: data.user.tenant_id,
      name: data.user.name,
      email: data.user.email,
      avatar: data.user.avatar,
      role: data.user.role,
      isAdmin: data.user.is_admin,
    },
  }
}

export async function logout(refreshToken?: string | null): Promise<void> {
  try {
    if (refreshToken) {
      await http.post("/auth/logout", { refresh_token: refreshToken })
    }
  } catch {
    // best effort — always clear locally
  }
  clearTokens()
}

export async function me(): Promise<AuthUser> {
  const data = await http.get<{
    user: {
      id: string
      tenant_id: string | null
      name: string
      email: string
      avatar: string
      role: string
      is_admin: boolean
    }
  }>("/auth/me")

  return {
    id: data.user.id,
    tenantId: data.user.tenant_id,
    name: data.user.name,
    email: data.user.email,
    avatar: data.user.avatar,
    role: data.user.role,
    isAdmin: data.user.is_admin,
  }
}

export function hasToken(): boolean {
  return getAccessToken() !== null
}
