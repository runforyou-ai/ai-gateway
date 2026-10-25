import type { AuthProvider } from "@refinedev/core"
import { ApiError } from "@/api/client"
import { register as registerAccount } from "@/api/generated/operations"
import { requestPasswordReset, resetPassword } from "@/api/password"
import { login, logout, me } from "@/api/session"

export const authProvider: AuthProvider = {
  login: async ({ email, password, remember, providerName }) => {
    if (providerName) {
      return {
        success: false,
        error: {
          name: "LoginError",
          message: "Social login is not supported",
        },
      }
    }
    await login({ email, password, remember: Boolean(remember) })
    return { success: true, redirectTo: "/console" }
  },
  register: async ({ name, email, password, providerName }) => {
    if (providerName) {
      return {
        success: false,
        error: {
          name: "RegisterError",
          message: "Social sign-up is not supported",
        },
      }
    }
    await registerAccount({ name, email, password })
    await login({ email, password, remember: true })
    return { success: true, redirectTo: "/console" }
  },
  logout: async () => {
    try {
      await logout()
    } catch {
      // ignore
    }
    return { success: true, redirectTo: "/" }
  },
  forgotPassword: async ({ email }) => {
    await requestPasswordReset({ email: String(email ?? "") })
    return { success: true }
  },
  updatePassword: async ({ password }) => {
    const token =
      new URLSearchParams(window.location.search).get("token")?.trim() ?? ""
    if (!token) {
      return {
        success: false,
        error: {
          name: "ResetPasswordError",
          message: "missing token",
        },
      }
    }
    await resetPassword({ token, password: String(password ?? "") })
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
    return {
      id: user.id,
      name: user.name,
      email: user.email,
      role: user.role,
      fullName: user.name,
      firstName: user.name,
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
