import { call } from "@/api/client"

export function requestPasswordReset(
  input: { email: string },
  signal?: AbortSignal,
) {
  return call<void>(
    { method: "POST", path: "/api/forgot-password", body: input },
    signal,
  )
}

export function resetPassword(
  input: { token: string; password: string },
  signal?: AbortSignal,
) {
  return call<void>(
    { method: "POST", path: "/api/reset-password", body: input },
    signal,
  )
}
