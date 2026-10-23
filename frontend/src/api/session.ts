import { call } from "@/api/client"

export type SessionUser = {
  id: number
  email: string
  role: string
}

export function login(input: { email: string; password: string }, signal?: AbortSignal) {
  return call<SessionUser>({ method: "POST", path: "/api/login", body: input }, signal)
}

export function logout(signal?: AbortSignal) {
  return call<void>({ method: "POST", path: "/api/logout" }, signal)
}

export function me(signal?: AbortSignal) {
  return call<SessionUser>({ method: "GET", path: "/api/me" }, signal)
}
