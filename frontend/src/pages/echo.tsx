import { useState, type KeyboardEvent } from "react"
import { echo } from "@/api/generated/operations"
import { Input } from "@/components/ui/input"

export function EchoPage() {
  const [text, setText] = useState("")
  const [message, setMessage] = useState("")
  const [error, setError] = useState("")
  const [loading, setLoading] = useState(false)

  async function submit() {
    setLoading(true)
    setError("")
    try {
      const result = await echo({ text })
      setMessage(result.message)
    } catch (e) {
      setMessage("")
      setError(e instanceof Error ? e.message : "request failed")
    } finally {
      setLoading(false)
    }
  }

  function onKeyDown(e: KeyboardEvent<HTMLInputElement>) {
    if (e.key !== "Enter" || loading) {
      return
    }
    e.preventDefault()
    void submit()
  }

  return (
    <div className="mx-auto flex min-h-svh max-w-lg flex-col justify-center gap-4 p-8">
      <h1 className="text-xl font-medium">Echo</h1>
      <Input
        value={text}
        onChange={(e) => setText(e.target.value)}
        onKeyDown={onKeyDown}
        placeholder="输入任意字符串后回车"
        disabled={loading}
        autoFocus
      />
      {message ? <p className="text-sm text-foreground">{message}</p> : null}
      {error ? <p className="text-sm text-destructive">{error}</p> : null}
    </div>
  )
}
