export type CallOptions = {
  method: string
  path: string
  body?: unknown
}

export type ApiErrorBody = {
  message: string
  fields?: Record<string, string>
}

export class ApiError extends Error {
  status: number
  fields?: Record<string, string>

  constructor(status: number, message: string, fields?: Record<string, string>) {
    super(message)
    this.name = "ApiError"
    this.status = status
    this.fields = fields
  }
}

const apiBase = (import.meta.env.VITE_API_URL as string | undefined)?.replace(/\/$/, "") ?? ""

export async function call<T>(opts: CallOptions, signal?: AbortSignal): Promise<T> {
  const init: RequestInit = {
    method: opts.method,
    headers: { Accept: "application/json" },
    signal,
  }
  if (opts.body !== undefined) {
    init.headers = { ...init.headers, "Content-Type": "application/json" }
    init.body = JSON.stringify(opts.body)
  }
  const res = await fetch(`${apiBase}${opts.path}`, init)
  if (!res.ok) {
    let message = res.statusText || "request failed"
    let fields: Record<string, string> | undefined
    try {
      const body = (await res.json()) as Partial<ApiErrorBody>
      if (typeof body.message === "string" && body.message) {
        message = body.message
      }
      if (body.fields && typeof body.fields === "object") {
        fields = body.fields
      }
    } catch {
      // ignore non-JSON error bodies
    }
    throw new ApiError(res.status, message, fields)
  }
  if (res.status === 204) {
    return undefined as T
  }
  return (await res.json()) as T
}
