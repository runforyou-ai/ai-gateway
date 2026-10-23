import type { AuthProvider } from "@refinedev/core"
import { ApiError } from "@/api/client"
import { register as registerAccount } from "@/api/generated/operations"
import { login, logout, me } from "@/api/session"

export const authProvider: AuthProvider = {
  login: async ({ email, password, providerName }) => {
    if (providerName) {
      return {
        success: false,
        error: {
          name: "LoginError",
          message: "Social login is not supported",
        },
      }
    }
    await login({ email, password })
    return { success: true, redirectTo: "/" }
  },
  register: async ({ email, password, providerName }) => {
    if (providerName) {
      return {
        success: false,
        error: {
          name: "RegisterError",
          message: "Social sign-up is not supported",
        },
      }
    }
    await registerAccount({ email, password })
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
    const local = user.email.split("@")[0] || user.email
    return {
      id: user.id,
      name: user.email,
      email: user.email,
      role: user.role,
      fullName: user.email,
      firstName: local,
      lastName: "",
    }
  },
  onError: async (error) => {
    if (error instanceof ApiError && error.status === 401) {
      return { logout: true, redirectTo: "/login" }
    }
    return { error }
  },
}
