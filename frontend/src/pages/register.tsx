import { useState, type FormEvent } from "react"
import { Link, useNavigate } from "react-router"
import { useLogin } from "@refinedev/core"
import { register } from "@/api/generated/operations"
import { ApiError } from "@/api/client"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { InputPassword } from "@/components/refine-ui/form/input-password"

export function RegisterPage() {
  const navigate = useNavigate()
  const { mutate: doLogin } = useLogin()
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")
  const [error, setError] = useState("")
  const [loading, setLoading] = useState(false)

  async function onSubmit(e: FormEvent) {
    e.preventDefault()
    setError("")
    setLoading(true)
    try {
      await register({ email, password })
      doLogin(
        { email, password },
        {
          onSuccess: () => navigate("/", { replace: true }),
          onError: (err) => {
            setError(err instanceof ApiError ? err.message : "注册成功但登录失败，请手动登录")
            navigate("/login", { replace: true })
          },
        },
      )
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "注册失败")
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="mx-auto flex min-h-svh max-w-sm flex-col justify-center gap-6 p-8">
      <div>
        <h1 className="text-xl font-medium">注册</h1>
        <p className="mt-1 text-sm text-muted-foreground">
          第一个注册的用户将成为管理员
        </p>
      </div>
      <form className="flex flex-col gap-4" onSubmit={onSubmit}>
        <div className="flex flex-col gap-2">
          <Label htmlFor="email">邮箱</Label>
          <Input
            id="email"
            type="email"
            autoComplete="username"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
            autoFocus
          />
        </div>
        <div className="flex flex-col gap-2">
          <Label htmlFor="password">密码</Label>
          <InputPassword
            id="password"
            autoComplete="new-password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
            minLength={8}
          />
          <p className="text-xs text-muted-foreground">至少 8 位</p>
        </div>
        {error ? <p className="text-sm text-destructive">{error}</p> : null}
        <Button type="submit" disabled={loading}>
          {loading ? "提交中…" : "注册"}
        </Button>
      </form>
      <p className="text-sm text-muted-foreground">
        已有账号？{" "}
        <Link className="text-foreground underline" to="/login">
          登录
        </Link>
      </p>
    </div>
  )
}
