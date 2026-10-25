import { useEffect } from "react"
import { useNavigate } from "react-router"
import { AuthPageLoading } from "@/components/auth-page-loading"
import { SignUpForm } from "@/components/refine-ui/form/sign-up-form"
import { useInstallationStatus } from "@/hooks/use-installation-status"

/** 首装或开放注册时可访问；否则跳转登录。 */
export function SignUpPage() {
  const navigate = useNavigate()
  const { status, ready } = useInstallationStatus()

  useEffect(() => {
    if (!ready) return
    if (status && status.installed && !status.registrationOpen) {
      navigate("/login", { replace: true })
    }
  }, [navigate, ready, status])

  if (!ready) {
    return <AuthPageLoading />
  }
  if (status && status.installed && !status.registrationOpen) {
    return <AuthPageLoading />
  }

  return <SignUpForm firstInstall={!status?.installed} />
}
