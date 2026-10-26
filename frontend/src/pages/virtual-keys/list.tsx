import { useMemo, useState } from "react"
import { type HttpError, useTranslation } from "@refinedev/core"
import { useModalForm } from "@refinedev/react-hook-form"
import { useTable } from "@refinedev/react-table"
import { zodResolver } from "@hookform/resolvers/zod"
import { type ColumnDef } from "@tanstack/react-table"
import { z } from "zod"
import { DataTable } from "@/components/refine-ui/data-table/data-table"
import { DeleteButton } from "@/components/refine-ui/buttons/delete"
import { EditButton } from "@/components/refine-ui/buttons/edit"
import {
  ListView,
  ListViewHeader,
} from "@/components/refine-ui/views/list-view"
import type { VirtualKey, VirtualKeyCreated } from "@/api/generated/contract"
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
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form"
import { Input } from "@/components/ui/input"
import { Switch } from "@/components/ui/switch"
import {
  VirtualKeyFormFields,
  type VirtualKeyFormValues,
} from "./form-fields"

function formatTime(raw: string, neverLabel: string): string {
  if (!raw) return neverLabel
  const d = new Date(raw)
  return Number.isNaN(d.getTime()) ? raw : d.toLocaleString()
}

export function VirtualKeyListPage() {
  const { translate } = useTranslation()
  const [secret, setSecret] = useState("")
  const [copied, setCopied] = useState(false)

  const schema = useMemo(
    () =>
      z.object({
        name: z
          .string()
          .trim()
          .min(1, translate("virtualKeys.fieldErrors.name.invalid"))
          .max(64, translate("virtualKeys.fieldErrors.name.invalid")),
        maxBudget: z
          .string()
          .trim()
          .refine(
            (v) => v === "" || (!Number.isNaN(Number(v)) && Number(v) > 0),
            translate("virtualKeys.fieldErrors.maxBudget.invalid"),
          ),
        expiresIn: z.string(),
        budgetReset: z.string(),
        disabled: z.boolean().optional(),
        prefix: z.string().optional(),
        spent: z.string().optional(),
        lastUsedAt: z.string().optional(),
      }),
    [translate],
  )

  const createForm = useModalForm<
    VirtualKeyCreated,
    HttpError,
    VirtualKeyFormValues
  >({
    resolver: zodResolver(schema) as never,
    defaultValues: {
      name: "",
      maxBudget: "",
      expiresIn: "",
      budgetReset: "",
    },
    refineCoreProps: {
      resource: "virtual-keys",
      action: "create",
      redirect: false,
      onMutationSuccess: ({ data }) => {
        setSecret(data.secret ?? "")
        setCopied(false)
      },
    },
    modalProps: {
      autoSubmitClose: false,
    },
  })

  const editForm = useModalForm<VirtualKey, HttpError, VirtualKeyFormValues>({
    resolver: zodResolver(schema) as never,
    defaultValues: {
      name: "",
      maxBudget: "",
      expiresIn: "",
      budgetReset: "",
      disabled: false,
      prefix: "",
      spent: "0",
      lastUsedAt: "",
    },
    refineCoreProps: {
      resource: "virtual-keys",
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
    refineCore: { formLoading: editLoading, id: editId },
    ...editRhf
  } = editForm

  const columns = useMemo<ColumnDef<VirtualKey>[]>(
    () => [
      {
        id: "name",
        accessorKey: "name",
        header: translate("virtualKeys.fields.name"),
      },
      {
        id: "prefix",
        accessorKey: "prefix",
        header: translate("virtualKeys.fields.prefix"),
        cell: ({ getValue }) => (
          <code className="text-sm">{String(getValue() ?? "")}…</code>
        ),
      },
      {
        id: "maxBudget",
        accessorKey: "maxBudget",
        header: translate("virtualKeys.fields.maxBudget"),
        cell: ({ getValue }) => {
          const v = String(getValue() ?? "")
          return v
            ? translate("virtualKeys.budgetAmount", { amount: v })
            : translate("virtualKeys.unlimited")
        },
      },
      {
        id: "spent",
        accessorKey: "spent",
        header: translate("virtualKeys.fields.spent"),
      },
      {
        id: "budgetReset",
        accessorKey: "budgetReset",
        header: translate("virtualKeys.fields.budgetReset"),
        cell: ({ getValue }) => {
          const v = String(getValue() ?? "")
          return translate(
            v === ""
              ? "virtualKeys.budgetReset.na"
              : `virtualKeys.budgetReset.${v}`,
          )
        },
      },
      {
        id: "expiresAt",
        accessorKey: "expiresAt",
        header: translate("virtualKeys.fields.expiresAt"),
        cell: ({ getValue }) =>
          formatTime(
            String(getValue() ?? ""),
            translate("virtualKeys.never"),
          ),
      },
      {
        id: "disabled",
        accessorKey: "disabled",
        header: translate("virtualKeys.fields.disabled"),
        cell: ({ getValue }) =>
          getValue()
            ? translate("virtualKeys.disabled.yes")
            : translate("virtualKeys.disabled.no"),
      },
      {
        id: "createdAt",
        accessorKey: "createdAt",
        header: translate("virtualKeys.fields.createdAt"),
        cell: ({ getValue }) => formatTime(String(getValue() ?? ""), ""),
      },
      {
        id: "actions",
        header: translate("virtualKeys.fields.actions"),
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

  const table = useTable<VirtualKey>({
    columns,
    refineCoreProps: {
      resource: "virtual-keys",
      pagination: { mode: "server" },
    },
  })

  const closeCreate = () => {
    setSecret("")
    setCopied(false)
    createModal.close()
  }

  const copySecret = async () => {
    if (!secret) return
    try {
      await navigator.clipboard.writeText(secret)
      setCopied(true)
    } catch {
      setCopied(false)
    }
  }

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
          if (!open) closeCreate()
        }}
      >
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>
              {secret
                ? translate("virtualKeys.secretTitle")
                : translate("virtualKeys.createTitle")}
            </DialogTitle>
          </DialogHeader>
          {secret ? (
            <div className="space-y-4">
              <p className="text-sm text-muted-foreground">
                {translate("virtualKeys.secretHint")}
              </p>
              <div className="flex flex-col gap-2 sm:flex-row sm:items-center">
                <Input value={secret} readOnly className="font-mono text-sm" />
                <Button type="button" variant="outline" onClick={copySecret}>
                  {copied
                    ? translate("virtualKeys.copied")
                    : translate("virtualKeys.copy")}
                </Button>
              </div>
              <DialogFooter>
                <Button type="button" onClick={closeCreate}>
                  {translate("virtualKeys.done")}
                </Button>
              </DialogFooter>
            </div>
          ) : (
            <Form {...createRhf}>
              <form
                className="space-y-4"
                onSubmit={createRhf.handleSubmit(createModal.submit)}
              >
                <VirtualKeyFormFields
                  key={createModal.visible ? "create-open" : "create"}
                  control={createRhf.control}
                  watch={createRhf.watch}
                />
                <DialogFooter>
                  <Button type="submit" disabled={createLoading}>
                    {translate("buttons.create")}
                  </Button>
                </DialogFooter>
              </form>
            </Form>
          )}
        </DialogContent>
      </Dialog>

      <Dialog
        open={editModal.visible}
        onOpenChange={(open) => {
          if (!open) editModal.close()
        }}
      >
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>{translate("virtualKeys.editTitle")}</DialogTitle>
          </DialogHeader>
          <Form {...editRhf}>
            <form
              className="space-y-4"
              onSubmit={editRhf.handleSubmit(editModal.submit)}
            >
              <VirtualKeyFormFields
                key={
                  editModal.visible ? `edit-${String(editId ?? "")}` : "edit"
                }
                control={editRhf.control}
                watch={editRhf.watch}
                showMeta
              />
              <FormField
                control={editRhf.control}
                name="disabled"
                render={({ field }) => (
                  <FormItem className="flex items-center justify-between gap-4">
                    <FormLabel>
                      {translate("virtualKeys.fields.disabled")}
                    </FormLabel>
                    <FormControl>
                      <Switch
                        checked={Boolean(field.value)}
                        onCheckedChange={field.onChange}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
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
