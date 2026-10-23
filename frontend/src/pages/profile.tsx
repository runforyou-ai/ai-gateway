import { useMemo } from "react"
import { useGetIdentity, useTranslation } from "@refinedev/core"
import { useForm } from "@refinedev/react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { z } from "zod"
import { InputPassword } from "@/components/refine-ui/form/input-password"
import { LoadingOverlay } from "@/components/refine-ui/layout/loading-overlay"
import {
  ListView,
  ListViewHeader,
} from "@/components/refine-ui/views/list-view"
import { Button } from "@/components/ui/button"
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form"
import { Input } from "@/components/ui/input"

type ProfileFormValues = {
  name: string
  email: string
  role?: string
  currentPassword?: string
  newPassword?: string
  confirmPassword?: string
}

type Identity = {
  id?: number
}

export function ProfilePage() {
  const { translate } = useTranslation()
  const { data: identity, refetch: refetchIdentity } =
    useGetIdentity<Identity>()
  const userId = identity?.id

  const profileSchema = useMemo(
    () =>
      z
        .object({
          name: z
            .string()
            .trim()
            .min(1, translate("profile.fieldErrors.name.invalid"))
            .max(64, translate("profile.fieldErrors.name.invalid")),
          email: z
            .string()
            .trim()
            .email(translate("profile.fieldErrors.email.invalid")),
          role: z.string().optional(),
          currentPassword: z.string().optional(),
          newPassword: z.string().optional(),
          confirmPassword: z.string().optional(),
        })
        .superRefine((values, ctx) => {
          const changingPassword = Boolean(
            values.newPassword ||
              values.confirmPassword ||
              values.currentPassword
          )
          if (!changingPassword) return

          if (!values.newPassword) {
            ctx.addIssue({
              code: z.ZodIssueCode.custom,
              path: ["newPassword"],
              message: translate("profile.fieldErrors.newPassword.required"),
            })
          } else if (values.newPassword.length < 8) {
            ctx.addIssue({
              code: z.ZodIssueCode.custom,
              path: ["newPassword"],
              message: translate("profile.fieldErrors.newPassword.invalid"),
            })
          }

          if (!values.currentPassword) {
            ctx.addIssue({
              code: z.ZodIssueCode.custom,
              path: ["currentPassword"],
              message: translate(
                "profile.fieldErrors.currentPassword.required"
              ),
            })
          }

          if (values.newPassword !== values.confirmPassword) {
            ctx.addIssue({
              code: z.ZodIssueCode.custom,
              path: ["confirmPassword"],
              message: translate("profile.errors.passwordMismatch"),
            })
          }
        }),
    [translate]
  )

  const {
    refineCore: { onFinish, formLoading, query },
    ...form
  } = useForm<ProfileFormValues>({
    resolver: zodResolver(profileSchema) as never,
    defaultValues: {
      name: "",
      email: "",
      role: "",
      currentPassword: "",
      newPassword: "",
      confirmPassword: "",
    },
    refineCoreProps: {
      resource: "profile",
      action: "edit",
      id: userId,
      redirect: false,
      queryOptions: {
        enabled: userId != null,
      },
      successNotification: () => ({
        message: translate("profile.saved"),
        type: "success",
      }),
      onMutationSuccess: () => {
        void refetchIdentity()
      },
    },
  })

  return (
    <ListView>
      <ListViewHeader
        resource="profile"
        title={translate("profile.title")}
        canCreate={false}
      />
      <LoadingOverlay loading={formLoading || query?.isLoading || userId == null}>
        <Form {...form}>
          <form
            onSubmit={form.handleSubmit(async (data) => {
              await onFinish({
                name: data.name,
                email: data.email,
                currentPassword: data.newPassword
                  ? data.currentPassword
                  : undefined,
                newPassword: data.newPassword || undefined,
              })
              form.setValue("currentPassword", "")
              form.setValue("newPassword", "")
              form.setValue("confirmPassword", "")
            })}
            className="space-y-6 p-4"
          >
            <FormField
              control={form.control}
              name="name"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{translate("profile.fields.name")}</FormLabel>
                  <FormControl>
                    <Input
                      autoComplete="name"
                      disabled={formLoading}
                      {...field}
                    />
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
                  <FormLabel>{translate("profile.fields.email")}</FormLabel>
                  <FormControl>
                    <Input
                      type="email"
                      autoComplete="email"
                      disabled={formLoading}
                      {...field}
                    />
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
                  <FormLabel>{translate("profile.fields.role")}</FormLabel>
                  <FormControl>
                    <Input
                      value={
                        field.value ? translate(`roles.${field.value}`) : ""
                      }
                      disabled
                      readOnly
                    />
                  </FormControl>
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="currentPassword"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>
                    {translate("profile.fields.currentPassword")}
                  </FormLabel>
                  <FormControl>
                    <InputPassword
                      autoComplete="current-password"
                      disabled={formLoading}
                      {...field}
                    />
                  </FormControl>
                  <FormDescription>
                    {translate("profile.passwordHint")}
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="newPassword"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>
                    {translate("profile.fields.newPassword")}
                  </FormLabel>
                  <FormControl>
                    <InputPassword
                      autoComplete="new-password"
                      disabled={formLoading}
                      {...field}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="confirmPassword"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>
                    {translate("profile.fields.confirmPassword")}
                  </FormLabel>
                  <FormControl>
                    <InputPassword
                      autoComplete="new-password"
                      disabled={formLoading}
                      {...field}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <Button type="submit" disabled={formLoading}>
              {formLoading
                ? translate("profile.saving")
                : translate("buttons.save")}
            </Button>
          </form>
        </Form>
      </LoadingOverlay>
    </ListView>
  )
}
