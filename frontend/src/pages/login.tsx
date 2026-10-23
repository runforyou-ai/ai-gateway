import { useEffect, useState, type FormEvent } from "react"
import { Link, useNavigate } from "react-router"
import { useLogin } from "@refinedev/core"
import { installationStatus } from "@/api/generated/operations"
import { ApiError } from "@/api/client"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { InputPassword } from "@/components/refine-ui/form/input-password"

export function LoginPage() {
  const navigate = useNavigate()
  const { mutate: doLogin, isPending } = useLogin()
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")
  const [error, setError] = useState("")
  const [checking, setChecking] = useState(true)

  useEffect(() => {
    let cancelled = false
    installationStatus()
      .then((status) => {
        if (cancelled) return
        if (!status.installed) {
          navigate("/register", { replace: true })
          return
        }
        setChecking(false)
      })
      .catch(() => {
        if (!cancelled) setChecking(false)
      })
    return () => {
      cancelled = true
    }
  }, [navigate])

  function onSubmit(e: FormEvent) {
    e.preventDefault()
    setError("")
    doLogin(
      { email, password },
      {
        onError: (err) => {
          setError(err instanceof ApiError ? err.message : "登录失败")
        },
      },
    )
  }

  if (checking) {
    return (
      <div className="flex min-h-svh items-center justify-center p-8 text-sm text-muted-foreground">
        加载中…
      </div>
    )
  }

  return (
    <div className="mx-auto flex min-h-svh max-w-sm flex-col justify-center gap-6 p-8">
      <div>
        <h1 className="text-xl font-medium">登录</h1>
        <p className="mt-1 text-sm text-muted-foreground">使用邮箱与密码进入管理台</p>
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
            autoComplete="current-password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
        </div>
        {error ? <p className="text-sm text-destructive">{error}</p> : null}
        <Button type="submit" disabled={isPending}>
          {isPending ? "登录中…" : "登录"}
        </Button>
      </form>
      <p className="text-sm text-muted-foreground">
        还没有账号？{" "}
        <Link className="text-foreground underline" to="/register">
          注册
        </Link>
      </p>
    </div>
  )
}
