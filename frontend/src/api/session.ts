import { call } from "@/api/client"

export type SessionUser = {
  id: number
  name: string
  email: string
  role: string
}

export function login(
  input: { email: string; password: string; remember?: boolean },
  signal?: AbortSignal,
) {
  return call<SessionUser>({ method: "POST", path: "/api/login", body: input }, signal)
}

export function logout(signal?: AbortSignal) {
  return call<void>({ method: "POST", path: "/api/logout" }, signal)
}

export function me(signal?: AbortSignal) {
  return call<SessionUser>({ method: "GET", path: "/api/me" }, signal)
}

export type UpdateProfileInput = {
  name: string
  email: string
  currentPassword?: string
  newPassword?: string
}

export function updateProfile(input: UpdateProfileInput, signal?: AbortSignal) {
  return call<SessionUser>(
    {
      method: "PATCH",
      path: "/api/me",
      body: {
        name: input.name,
        email: input.email,
        currentPassword: input.currentPassword ?? "",
        newPassword: input.newPassword ?? "",
      },
    },
    signal,
  )
}
