import { useMemo } from "react"
import { type HttpError, useTranslation } from "@refinedev/core"
import { useModalForm } from "@refinedev/react-hook-form"
import { useTable } from "@refinedev/react-table"
import { zodResolver } from "@hookform/resolvers/zod"
import type { Control } from "react-hook-form"
import { type ColumnDef } from "@tanstack/react-table"
import { z } from "zod"
import { DataTable } from "@/components/refine-ui/data-table/data-table"
import { DeleteButton } from "@/components/refine-ui/buttons/delete"
import { EditButton } from "@/components/refine-ui/buttons/edit"
import {
  ListView,
  ListViewHeader,
} from "@/components/refine-ui/views/list-view"
import { InputPassword } from "@/components/refine-ui/form/input-password"
import type { User } from "@/api/generated/contract"
import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
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
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { Switch } from "@/components/ui/switch"
import { InviteUserButton } from "./invite-button"

type FormValues = {
  name: string
  email: string
  password: string
  role: string
  disabled: boolean
}

export function UserListPage() {
  const { translate } = useTranslation()

  const createSchema = useMemo(
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

  const editSchema = useMemo(
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
          .refine(
            (v) => !v || v.length >= 8,
            translate("users.fieldErrors.password.invalid"),
          ),
        role: z.enum(["admin", "user"]),
        disabled: z.boolean(),
      }),
    [translate],
  )

  const createForm = useModalForm<User, HttpError, FormValues>({
    resolver: zodResolver(createSchema) as never,
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
      redirect: false,
    },
  })

  const editForm = useModalForm<User, HttpError, FormValues>({
    resolver: zodResolver(editSchema) as never,
    defaultValues: {
      name: "",
      email: "",
      password: "",
      role: "user",
      disabled: false,
    },
    refineCoreProps: {
      resource: "users",
      action: "edit",
      redirect: false,
    },
  })

  const {
    modal: createModal,
    refineCore: { formLoading: createLoading },
    ...createRhf
  } = createForm

  const {
    modal: editModal,
    refineCore: { formLoading: editLoading },
    ...editRhf
  } = editForm

  const columns = useMemo<ColumnDef<User>[]>(
    () => [
      {
        id: "id",
        accessorKey: "id",
        header: translate("users.fields.id"),
        size: 80,
      },
      {
        id: "name",
        accessorKey: "name",
        header: translate("users.fields.name"),
      },
      {
        id: "email",
        accessorKey: "email",
        header: translate("users.fields.email"),
      },
      {
        id: "role",
        accessorKey: "role",
        header: translate("users.fields.role"),
        cell: ({ getValue }) => {
          const role = String(getValue() ?? "")
          return translate(`roles.${role}`, role)
        },
      },
      {
        id: "disabled",
        accessorKey: "disabled",
        header: translate("users.fields.disabled"),
        cell: ({ getValue }) =>
          getValue()
            ? translate("users.disabled.yes")
            : translate("users.disabled.no"),
      },
      {
        id: "actions",
        header: translate("users.fields.actions"),
        size: 120,
        cell: ({ row }) => (
          <div className="flex items-center gap-2">
            <EditButton
              recordItemId={row.original.id}
              size="sm"
              onClick={() => editModal.show(row.original.id)}
            />
            <DeleteButton recordItemId={row.original.id} size="sm" />
          </div>
        ),
      },
    ],
    [translate, editModal],
  )

  const table = useTable<User>({
    columns,
    refineCoreProps: {
      resource: "users",
      pagination: { mode: "server" },
    },
  })

  return (
    <ListView>
      <ListViewHeader canCreate={false}>
        <div className="flex items-center gap-2">
          <InviteUserButton />
          <Button type="button" onClick={() => createModal.show()}>
            {translate("buttons.create")}
          </Button>
        </div>
      </ListViewHeader>
      <DataTable table={table} />

      <Dialog
        open={createModal.visible}
        onOpenChange={(open) => {
          if (!open) createModal.close()
        }}
      >
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{translate("users.createTitle")}</DialogTitle>
          </DialogHeader>
          <Form {...createRhf}>
            <form
              className="space-y-4"
              onSubmit={createRhf.handleSubmit(createModal.submit)}
            >
              <UserFormFields control={createRhf.control} passwordRequired />
              <DialogFooter>
                <Button type="submit" disabled={createLoading}>
                  {translate("buttons.create")}
                </Button>
              </DialogFooter>
            </form>
          </Form>
        </DialogContent>
      </Dialog>

      <Dialog
        open={editModal.visible}
        onOpenChange={(open) => {
          if (!open) editModal.close()
        }}
      >
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{translate("users.editTitle")}</DialogTitle>
          </DialogHeader>
          <Form {...editRhf}>
            <form
              className="space-y-4"
              onSubmit={editRhf.handleSubmit((values) =>
                editModal.submit({
                  ...values,
                  password: values.password || "",
                }),
              )}
            >
              <UserFormFields control={editRhf.control} />
              <DialogFooter>
                <Button type="submit" disabled={editLoading}>
                  {translate("buttons.save")}
                </Button>
              </DialogFooter>
            </form>
          </Form>
        </DialogContent>
      </Dialog>
    </ListView>
  )
}

function UserFormFields({
  control,
  passwordRequired = false,
}: {
  control: Control<FormValues>
  passwordRequired?: boolean
}) {
  const { translate } = useTranslation()

  return (
    <>
      <FormField
        control={control}
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
        control={control}
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
        control={control}
        name="password"
        render={({ field }) => (
          <FormItem>
            <FormLabel>{translate("users.fields.password")}</FormLabel>
            <FormControl>
              <InputPassword
                {...field}
                value={field.value ?? ""}
                autoComplete="new-password"
              />
            </FormControl>
            {!passwordRequired ? (
              <FormDescription>
                {translate("users.passwordHint")}
              </FormDescription>
            ) : null}
            <FormMessage />
          </FormItem>
        )}
      />
      <FormField
        control={control}
        name="role"
        render={({ field }) => (
          <FormItem>
            <FormLabel>{translate("users.fields.role")}</FormLabel>
            <Select value={field.value} onValueChange={field.onChange}>
              <FormControl>
                <SelectTrigger className="w-full">
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
        control={control}
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
    </>
  )
}
