export type CallOptions = {
  method: string
  path: string
  body?: unknown
}

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
  const res = await fetch(opts.path, init)
  if (!res.ok) {
    const text = await res.text()
    throw new Error(text || res.statusText)
  }
  if (res.status === 204) {
    return undefined as T
  }
  return (await res.json()) as T
}
