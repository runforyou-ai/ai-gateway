import type { AuthProvider } from "@refinedev/core"
import { ApiError } from "@/api/client"
import { login, logout, me } from "@/api/session"

export const authProvider: AuthProvider = {
  login: async ({ email, password }) => {
    await login({ email, password })
    return { success: true, redirectTo: "/" }
  },
  logout: async () => {
    try {
      await logout()
    } catch {
      // ignore
    }
    return { success: true, redirectTo: "/login" }
  },
  check: async () => {
    try {
      await me()
      return { authenticated: true }
    } catch {
      return { authenticated: false, redirectTo: "/login", logout: true }
    }
  },
  getIdentity: async () => {
    const user = await me()
    return { id: user.id, name: user.email, role: user.role }
  },
  onError: async (error) => {
    if (error instanceof ApiError && error.status === 401) {
      return { logout: true, redirectTo: "/login" }
    }
    return { error }
  },
}
