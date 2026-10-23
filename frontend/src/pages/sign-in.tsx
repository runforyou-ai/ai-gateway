import { useEffect, useState } from "react"
import { useNavigate } from "react-router"
import { installationStatus } from "@/api/generated/operations"
import { SignInForm } from "@/components/refine-ui/form/sign-in-form"

/** 官方 SignInForm 入口；无用户时跳转注册。 */
export function SignInPage() {
  const navigate = useNavigate()
  const [ready, setReady] = useState(false)

  useEffect(() => {
    let cancelled = false
    installationStatus()
      .then((status) => {
        if (cancelled) return
        if (!status.installed) {
          navigate("/register", { replace: true })
          return
        }
        setReady(true)
      })
      .catch(() => {
        if (!cancelled) setReady(true)
      })
    return () => {
      cancelled = true
    }
  }, [navigate])

  if (!ready) {
    return (
      <div className="flex min-h-svh items-center justify-center p-8 text-sm text-muted-foreground">
        Loading…
      </div>
    )
  }

  return <SignInForm />
}
