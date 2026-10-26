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
import type { CatalogModel } from "@/api/generated/contract"
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

const MODEL_TYPES = ["chat"] as const

type FormValues = {
  name: string
  displayName: string
  type: string
  disabled: boolean
}

export function ModelListPage() {
  const { translate } = useTranslation()

  const schema = useMemo(
    () =>
      z.object({
        name: z
          .string()
          .trim()
          .min(1, translate("models.fieldErrors.name.invalid"))
          .max(128, translate("models.fieldErrors.name.invalid"))
          .regex(
            /^[a-zA-Z0-9][a-zA-Z0-9._:/-]*$/,
            translate("models.fieldErrors.name.invalid"),
          ),
        displayName: z
          .string()
          .trim()
          .min(1, translate("models.fieldErrors.displayName.invalid"))
          .max(128, translate("models.fieldErrors.displayName.invalid")),
        type: z.enum(MODEL_TYPES),
        disabled: z.boolean(),
      }),
    [translate],
  )

  const createForm = useModalForm<CatalogModel, HttpError, FormValues>({
    resolver: zodResolver(schema) as never,
    defaultValues: {
      name: "",
      displayName: "",
      type: "chat",
      disabled: false,
    },
    refineCoreProps: {
      resource: "models",
      action: "create",
      redirect: false,
    },
  })

  const editForm = useModalForm<CatalogModel, HttpError, FormValues>({
    resolver: zodResolver(schema) as never,
    defaultValues: {
      name: "",
      displayName: "",
      type: "chat",
      disabled: false,
    },
    refineCoreProps: {
      resource: "models",
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

  const columns = useMemo<ColumnDef<CatalogModel>[]>(
    () => [
      {
        id: "name",
        accessorKey: "name",
        header: translate("models.fields.name"),
        cell: ({ getValue }) => (
          <code className="text-sm">{String(getValue() ?? "")}</code>
        ),
      },
      {
        id: "displayName",
        accessorKey: "displayName",
        header: translate("models.fields.displayName"),
      },
      {
        id: "type",
        accessorKey: "type",
        header: translate("models.fields.type"),
        cell: ({ getValue }) => {
          const t = String(getValue() ?? "")
          return translate(`models.types.${t}`, t)
        },
      },
      {
        id: "disabled",
        accessorKey: "disabled",
        header: translate("models.fields.disabled"),
        cell: ({ getValue }) =>
          getValue()
            ? translate("models.disabled.yes")
            : translate("models.disabled.no"),
      },
      {
        id: "actions",
        header: translate("models.fields.actions"),
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

  const table = useTable<CatalogModel>({
    columns,
    refineCoreProps: {
      resource: "models",
      pagination: { mode: "server" },
    },
  })

  return (
    <ListView>
      <ListViewHeader canCreate={false}>
        <Button type="button" onClick={() => createModal.show()}>
          {translate("buttons.create")}
        </Button>
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
            <DialogTitle>{translate("models.createTitle")}</DialogTitle>
          </DialogHeader>
          <Form {...createRhf}>
            <form
              className="space-y-4"
              onSubmit={createRhf.handleSubmit(createModal.submit)}
            >
              <ModelFormFields control={createRhf.control} />
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
            <DialogTitle>{translate("models.editTitle")}</DialogTitle>
          </DialogHeader>
          <Form {...editRhf}>
            <form
              className="space-y-4"
              onSubmit={editRhf.handleSubmit(editModal.submit)}
            >
              <ModelFormFields control={editRhf.control} />
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

function ModelFormFields({ control }: { control: Control<FormValues> }) {
  const { translate } = useTranslation()

  return (
    <>
      <FormField
        control={control}
        name="name"
        render={({ field }) => (
          <FormItem>
            <FormLabel>{translate("models.fields.name")}</FormLabel>
            <FormControl>
              <Input {...field} autoComplete="off" className="font-mono" />
            </FormControl>
            <FormDescription>{translate("models.nameHint")}</FormDescription>
            <FormMessage />
          </FormItem>
        )}
      />
      <FormField
        control={control}
        name="displayName"
        render={({ field }) => (
          <FormItem>
            <FormLabel>{translate("models.fields.displayName")}</FormLabel>
            <FormControl>
              <Input {...field} autoComplete="off" />
            </FormControl>
            <FormMessage />
          </FormItem>
        )}
      />
      <FormField
        control={control}
        name="type"
        render={({ field }) => (
          <FormItem>
            <FormLabel>{translate("models.fields.type")}</FormLabel>
            <Select value={field.value} onValueChange={field.onChange}>
              <FormControl>
                <SelectTrigger className="w-full">
                  <SelectValue />
                </SelectTrigger>
              </FormControl>
              <SelectContent>
                {MODEL_TYPES.map((t) => (
                  <SelectItem key={t} value={t}>
                    {translate(`models.types.${t}`)}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            <FormDescription>{translate("models.typeHint")}</FormDescription>
            <FormMessage />
          </FormItem>
        )}
      />
      <FormField
        control={control}
        name="disabled"
        render={({ field }) => (
          <FormItem className="flex items-center justify-between gap-4">
            <FormLabel>{translate("models.fields.disabled")}</FormLabel>
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
