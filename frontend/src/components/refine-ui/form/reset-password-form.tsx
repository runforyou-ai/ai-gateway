"use client"

import { useMemo, useState } from "react"
import { ArrowLeft } from "lucide-react"
import {
  useLink,
  useNotification,
  useRefineOptions,
  useTranslation,
  useUpdatePassword,
} from "@refinedev/core"
import { Button } from "@/components/ui/button"
import { Label } from "@/components/ui/label"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { Separator } from "@/components/ui/separator"
import { InputPassword } from "@/components/refine-ui/form/input-password"
import { cn } from "@/lib/utils"

export function ResetPasswordForm() {
  const [password, setPassword] = useState("")
  const [confirmPassword, setConfirmPassword] = useState("")
  const [mismatch, setMismatch] = useState(false)

  const Link = useLink()
  const { translate } = useTranslation()
  const { open } = useNotification()
  const { title } = useRefineOptions()
  const { mutate: updatePassword } = useUpdatePassword()

  const token = useMemo(() => {
    const params = new URLSearchParams(window.location.search)
    return params.get("token")?.trim() ?? ""
  }, [])

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    if (password !== confirmPassword) {
      setMismatch(true)
      return
    }
    setMismatch(false)
    updatePassword(
      { password },
      {
        onSuccess: () => {
          open?.({
            type: "success",
            message: translate("pages.resetPassword.successTitle"),
            description: translate("pages.resetPassword.successDescription"),
          })
        },
      },
    )
  }

  return (
    <div
      className={cn(
        "flex",
        "flex-col",
        "items-center",
        "justify-center",
        "px-6",
        "py-8",
        "min-h-svh",
      )}
    >
      <div className={cn("flex", "items-center", "justify-center")}>
        {title.icon && (
          <div
            className={cn("text-foreground", "[&>svg]:w-12", "[&>svg]:h-12")}
          >
            {title.icon}
          </div>
        )}
      </div>

      <Card className={cn("sm:w-[456px]", "p-12", "mt-6")}>
        <CardHeader className={cn("px-0")}>
          <CardTitle
            className={cn(
              "text-blue-600",
              "dark:text-blue-400",
              "text-3xl",
              "font-semibold",
            )}
          >
            {translate("pages.resetPassword.title")}
          </CardTitle>
          <CardDescription
            className={cn("text-muted-foreground", "font-medium")}
          >
            {token
              ? translate("pages.resetPassword.description")
              : translate("pages.resetPassword.missingToken")}
          </CardDescription>
        </CardHeader>

        <Separator />

        <CardContent className={cn("px-0")}>
          {token ? (
            <form onSubmit={handleSubmit}>
              <div className={cn("flex", "flex-col", "gap-2")}>
                <Label htmlFor="password">
                  {translate("pages.resetPassword.fields.password")}
                </Label>
                <InputPassword
                  id="password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  required
                  minLength={8}
                  autoComplete="new-password"
                />
              </div>
              <div className={cn("flex", "flex-col", "gap-2", "mt-6")}>
                <Label htmlFor="confirmPassword">
                  {translate("pages.resetPassword.fields.confirmPassword")}
                </Label>
                <InputPassword
                  id="confirmPassword"
                  value={confirmPassword}
                  onChange={(e) => setConfirmPassword(e.target.value)}
                  required
                  minLength={8}
                  autoComplete="new-password"
                />
                {mismatch ? (
                  <p className="text-sm text-destructive">
                    {translate("pages.resetPassword.errors.passwordMismatch")}
                  </p>
                ) : null}
              </div>
              <Button type="submit" size="lg" className={cn("w-full", "mt-6")}>
                {translate("pages.resetPassword.buttons.submit")}
              </Button>
            </form>
          ) : null}

          <div className={cn("mt-8")}>
            <Link
              to="/login"
              className={cn(
                "inline-flex",
                "items-center",
                "gap-2",
                "text-sm",
                "text-muted-foreground",
                "hover:text-foreground",
                "transition-colors",
              )}
            >
              <ArrowLeft className={cn("w-4", "h-4")} />
              <span>{translate("pages.resetPassword.buttons.back")}</span>
            </Link>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}

ResetPasswordForm.displayName = "ResetPasswordForm"
