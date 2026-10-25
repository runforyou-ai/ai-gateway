import { useEffect } from "react"
import { useNavigate } from "react-router"
import { AuthPageLoading } from "@/components/auth-page-loading"
import { SignInForm } from "@/components/refine-ui/form/sign-in-form"
import { useInstallationStatus } from "@/hooks/use-installation-status"

/** 官方 SignInForm 入口；无用户时跳转注册。 */
export function SignInPage() {
  const navigate = useNavigate()
  const { status, ready } = useInstallationStatus()

  useEffect(() => {
    if (!ready || !status) return
    if (!status.installed) {
      navigate("/register", { replace: true })
    }
  }, [navigate, ready, status])

  if (!ready || !status || !status.installed) {
    return <AuthPageLoading />
  }

  return <SignInForm showSignUp={status.registrationOpen} />
}
