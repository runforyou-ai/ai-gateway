import { useMemo } from "react"
import { useTranslation } from "@refinedev/core"
import { useForm } from "@refinedev/react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { z } from "zod"
import {
  CreateView,
  CreateViewHeader,
} from "@/components/refine-ui/views/create-view"
import { InputPassword } from "@/components/refine-ui/form/input-password"
import { Button } from "@/components/ui/button"
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
import { Switch } from "@/components/ui/switch"

type FormValues = {
  name: string
  email: string
  password: string
  role: string
  disabled: boolean
}

export function UserCreatePage() {
  const { translate } = useTranslation()

  const schema = useMemo(
    () =>
      z.object({
        name: z
          .string()
          .trim()
          .min(1, translate("users.fieldErrors.name.invalid"))
          .max(64, translate("users.fieldErrors.name.invalid")),
        email: z
          .string()
          .trim()
          .email(translate("users.fieldErrors.email.invalid")),
        password: z
          .string()
          .min(8, translate("users.fieldErrors.password.invalid")),
        role: z.enum(["admin", "user"]),
        disabled: z.boolean(),
      }),
    [translate],
  )

  const {
    refineCore: { onFinish, formLoading },
    ...form
  } = useForm<FormValues>({
    resolver: zodResolver(schema) as never,
    defaultValues: {
      name: "",
      email: "",
      password: "",
      role: "user",
      disabled: false,
    },
    refineCoreProps: {
      resource: "users",
      action: "create",
      redirect: "list",
    },
  })

  return (
    <CreateView>
      <CreateViewHeader
        title={translate("users.createTitle")}
        withBackButton={false}
      />
      <Form {...form}>
        <form
          className="space-y-6 p-4"
          onSubmit={form.handleSubmit(onFinish)}
        >
          <FormField
            control={form.control}
            name="name"
            render={({ field }) => (
              <FormItem>
                <FormLabel>{translate("users.fields.name")}</FormLabel>
                <FormControl>
                  <Input {...field} autoComplete="off" />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
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
            name="password"
            render={({ field }) => (
              <FormItem>
                <FormLabel>{translate("users.fields.password")}</FormLabel>
                <FormControl>
                  <InputPassword {...field} autoComplete="new-password" />
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
          <FormField
            control={form.control}
            name="disabled"
            render={({ field }) => (
              <FormItem className="flex items-center justify-between gap-4">
                <FormLabel>{translate("users.fields.disabled")}</FormLabel>
                <FormControl>
                  <Switch
                    checked={field.value}
                    onCheckedChange={field.onChange}
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          <Button type="submit" disabled={formLoading}>
            {formLoading
              ? translate("buttons.save")
              : translate("buttons.create")}
          </Button>
        </form>
      </Form>
    </CreateView>
  )
}
