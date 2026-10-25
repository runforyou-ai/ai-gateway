import { useEffect, useMemo, useState } from "react"
import { useNotification, useTranslation } from "@refinedev/core"
import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { z } from "zod"
import {
  getSiteSettings,
  updateSiteSettings,
} from "@/api/generated/operations"
import { ApiError } from "@/api/client"
import {
  ListView,
  ListViewHeader,
} from "@/components/refine-ui/views/list-view"
import { LoadingOverlay } from "@/components/refine-ui/layout/loading-overlay"
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
import { Switch } from "@/components/ui/switch"

type FormValues = {
  name: string
  registrationOpen: boolean
}

export function SettingsPage() {
  const { translate } = useTranslation()
  const { open } = useNotification()
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)

  const schema = useMemo(
    () =>
      z.object({
        name: z
          .string()
          .trim()
          .min(1, translate("settings.fieldErrors.name.invalid"))
          .max(64, translate("settings.fieldErrors.name.invalid")),
        registrationOpen: z.boolean(),
      }),
    [translate],
  )

  const form = useForm<FormValues>({
    resolver: zodResolver(schema) as never,
    defaultValues: { name: "", registrationOpen: false },
  })

  useEffect(() => {
    let cancelled = false
    ;(async () => {
      try {
        const site = await getSiteSettings()
        if (cancelled) return
        form.reset({
          name: site.name,
          registrationOpen: site.registrationOpen,
        })
      } catch (error) {
        if (cancelled) return
        const message =
          error instanceof ApiError
            ? error.message
            : translate("settings.loadFailed")
        open?.({ type: "error", message })
      } finally {
        if (!cancelled) setLoading(false)
      }
    })()
    return () => {
      cancelled = true
    }
  }, [form, open, translate])

  const onSubmit = form.handleSubmit(async (values) => {
    setSaving(true)
    try {
      const site = await updateSiteSettings({
        name: values.name,
        registrationOpen: values.registrationOpen,
      })
      form.reset({
        name: site.name,
        registrationOpen: site.registrationOpen,
      })
      open?.({
        type: "success",
        message: translate("settings.saved"),
      })
    } catch (error) {
      const message =
        error instanceof ApiError
          ? error.message
          : translate("settings.saveFailed")
      open?.({ type: "error", message })
    } finally {
      setSaving(false)
    }
  })

  return (
    <ListView>
      <ListViewHeader title={translate("settings.title")} canCreate={false} />
      <LoadingOverlay loading={loading}>
        <Form {...form}>
          <form className="space-y-6 p-4" onSubmit={onSubmit}>
            <FormField
              control={form.control}
              name="name"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{translate("settings.fields.name")}</FormLabel>
                  <FormControl>
                    <Input {...field} autoComplete="off" />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="registrationOpen"
              render={({ field }) => (
                <FormItem className="flex items-center justify-between gap-4">
                  <div className="space-y-1">
                    <FormLabel>
                      {translate("settings.fields.registrationOpen")}
                    </FormLabel>
                    <FormDescription>
                      {translate("settings.registrationOpenHint")}
                    </FormDescription>
                  </div>
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
            <Button type="submit" disabled={saving || loading}>
              {translate("buttons.save")}
            </Button>
          </form>
        </Form>
      </LoadingOverlay>
    </ListView>
  )
}
