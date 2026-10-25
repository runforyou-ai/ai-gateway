import { useMemo, useState } from "react"
import { useNotification, useTranslation } from "@refinedev/core"
import { zodResolver } from "@hookform/resolvers/zod"
import { useForm } from "react-hook-form"
import { z } from "zod"
import { inviteUser } from "@/api/generated/operations"
import { ApiError } from "@/api/client"
import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form"
import { Input } from "@/components/ui/input"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"

type FormValues = {
  email: string
  role: string
}

export function InviteUserButton() {
  const { translate } = useTranslation()
  const { open } = useNotification()
  const [openDialog, setOpenDialog] = useState(false)
  const [submitting, setSubmitting] = useState(false)

  const schema = useMemo(
    () =>
      z.object({
        email: z
          .string()
          .trim()
          .email(translate("users.fieldErrors.email.invalid")),
        role: z.enum(["admin", "user"]),
      }),
    [translate],
  )

  const form = useForm<FormValues>({
    resolver: zodResolver(schema) as never,
    defaultValues: { email: "", role: "user" },
  })

  const onSubmit = form.handleSubmit(async (values) => {
    setSubmitting(true)
    try {
      await inviteUser({ email: values.email, role: values.role })
      open?.({
        type: "success",
        message: translate("users.invite.success"),
      })
      form.reset({ email: "", role: "user" })
      setOpenDialog(false)
    } catch (error) {
      const message =
        error instanceof ApiError
          ? error.message
          : translate("users.invite.failed")
      open?.({ type: "error", message })
    } finally {
      setSubmitting(false)
    }
  })

  return (
    <Dialog open={openDialog} onOpenChange={setOpenDialog}>
      <DialogTrigger asChild>
        <Button variant="outline">{translate("users.invite.button")}</Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{translate("users.invite.title")}</DialogTitle>
          <DialogDescription>
            {translate("users.invite.description")}
          </DialogDescription>
        </DialogHeader>
        <Form {...form}>
          <form className="space-y-4" onSubmit={onSubmit}>
            <FormField
              control={form.control}
              name="email"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{translate("users.fields.email")}</FormLabel>
                  <FormControl>
                    <Input {...field} type="email" autoComplete="off" />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="role"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{translate("users.fields.role")}</FormLabel>
                  <Select value={field.value} onValueChange={field.onChange}>
                    <FormControl>
                      <SelectTrigger>
                        <SelectValue />
                      </SelectTrigger>
                    </FormControl>
                    <SelectContent>
                      <SelectItem value="admin">
                        {translate("roles.admin")}
                      </SelectItem>
                      <SelectItem value="user">
                        {translate("roles.user")}
                      </SelectItem>
                    </SelectContent>
                  </Select>
                  <FormMessage />
                </FormItem>
              )}
            />
            <DialogFooter>
              <Button type="submit" disabled={submitting}>
                {translate("users.invite.submit")}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  )
}
