"use client"

import { useEffect, useMemo, useState } from "react"
import { ArrowLeft } from "lucide-react"
import {
  useLink,
  useNotification,
  useRefineOptions,
  useTranslation,
} from "@refinedev/core"
import { acceptInvite, lookupInvite } from "@/api/generated/operations"
import { ApiError } from "@/api/client"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
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

export function AcceptInviteForm() {
  const [name, setName] = useState("")
  const [password, setPassword] = useState("")
  const [confirmPassword, setConfirmPassword] = useState("")
  const [mismatch, setMismatch] = useState(false)
  const [email, setEmail] = useState("")
  const [lookupError, setLookupError] = useState(false)
  const [loading, setLoading] = useState(true)
  const [submitting, setSubmitting] = useState(false)
  const [done, setDone] = useState(false)

  const Link = useLink()
  const { translate } = useTranslation()
  const { open } = useNotification()
  const { title } = useRefineOptions()

  const token = useMemo(() => {
    const params = new URLSearchParams(window.location.search)
    return params.get("token")?.trim() ?? ""
  }, [])

  useEffect(() => {
    let cancelled = false
    ;(async () => {
      if (!token) {
        setLookupError(true)
        setLoading(false)
        return
      }
      try {
        const invite = await lookupInvite({ token })
        if (cancelled) return
        setEmail(invite.email)
      } catch {
        if (cancelled) return
        setLookupError(true)
      } finally {
        if (!cancelled) setLoading(false)
      }
    })()
    return () => {
      cancelled = true
    }
  }, [token])

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    if (password !== confirmPassword) {
      setMismatch(true)
      return
    }
    setMismatch(false)
    setSubmitting(true)
    try {
      await acceptInvite({ token, name, password })
      setDone(true)
      open?.({
        type: "success",
        message: translate("pages.acceptInvite.successTitle"),
        description: translate("pages.acceptInvite.successDescription"),
      })
    } catch (error) {
      const message =
        error instanceof ApiError
          ? error.message
          : translate("pages.acceptInvite.failed")
      open?.({ type: "error", message })
    } finally {
      setSubmitting(false)
    }
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
            {translate("pages.acceptInvite.title")}
          </CardTitle>
          <CardDescription
            className={cn("text-muted-foreground", "font-medium")}
          >
            {loading
              ? translate("pages.acceptInvite.loading")
              : lookupError
                ? translate("pages.acceptInvite.missingToken")
                : done
                  ? translate("pages.acceptInvite.successDescription")
                  : translate("pages.acceptInvite.description")}
          </CardDescription>
        </CardHeader>

        <Separator />

        <CardContent className={cn("px-0")}>
          {!loading && !lookupError && !done ? (
            <form onSubmit={handleSubmit}>
              <div className={cn("flex", "flex-col", "gap-2")}>
                <Label htmlFor="email">
                  {translate("pages.acceptInvite.fields.email")}
                </Label>
                <Input id="email" value={email} disabled readOnly />
              </div>
              <div className={cn("flex", "flex-col", "gap-2", "mt-6")}>
                <Label htmlFor="name">
                  {translate("pages.acceptInvite.fields.name")}
                </Label>
                <Input
                  id="name"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  required
                  autoComplete="name"
                />
              </div>
              <div className={cn("flex", "flex-col", "gap-2", "mt-6")}>
                <Label htmlFor="password">
                  {translate("pages.acceptInvite.fields.password")}
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
                  {translate("pages.acceptInvite.fields.confirmPassword")}
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
                    {translate("pages.acceptInvite.errors.passwordMismatch")}
                  </p>
                ) : null}
              </div>
              <Button
                type="submit"
                size="lg"
                className={cn("w-full", "mt-6")}
                disabled={submitting}
              >
                {translate("pages.acceptInvite.buttons.submit")}
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
              <span>{translate("pages.acceptInvite.buttons.back")}</span>
            </Link>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}

AcceptInviteForm.displayName = "AcceptInviteForm"
