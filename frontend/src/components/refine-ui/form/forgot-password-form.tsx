"use client"

import { useState } from "react"
import { ArrowLeft } from "lucide-react"
import {
  useForgotPassword,
  useLink,
  useNotification,
  useRefineOptions,
  useTranslation,
} from "@refinedev/core"
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
import { cn } from "@/lib/utils"

export const ForgotPasswordForm = () => {
  const [email, setEmail] = useState("")

  const Link = useLink()
  const { translate } = useTranslation()
  const { open } = useNotification()
  const { title } = useRefineOptions()
  const { mutate: forgotPassword } = useForgotPassword()

  const handleForgotPassword = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    forgotPassword(
      { email },
      {
        onSuccess: () => {
          open?.({
            type: "success",
            message: translate("pages.forgotPassword.successTitle"),
            description: translate("pages.forgotPassword.successDescription"),
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
            {translate("pages.forgotPassword.title")}
          </CardTitle>
          <CardDescription
            className={cn("text-muted-foreground", "font-medium")}
          >
            {translate("pages.forgotPassword.description")}
          </CardDescription>
        </CardHeader>

        <Separator />

        <CardContent className={cn("px-0")}>
          <form onSubmit={handleForgotPassword}>
            <div className={cn("flex", "flex-col", "gap-2")}>
              <Label htmlFor="email">
                {translate("pages.forgotPassword.fields.email")}
              </Label>
              <Input
                id="email"
                type="email"
                required
                value={email}
                onChange={(e) => setEmail(e.target.value)}
              />
            </div>
            <Button type="submit" size="lg" className={cn("w-full", "mt-6")}>
              {translate("pages.forgotPassword.buttons.submit")}
            </Button>
          </form>

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
              <span>{translate("pages.forgotPassword.buttons.back")}</span>
            </Link>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}

ForgotPasswordForm.displayName = "ForgotPasswordForm"
