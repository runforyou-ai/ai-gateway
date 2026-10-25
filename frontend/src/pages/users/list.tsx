import { useMemo } from "react"
import { useTranslation } from "@refinedev/core"
import { useTable } from "@refinedev/react-table"
import { type ColumnDef } from "@tanstack/react-table"
import { DataTable } from "@/components/refine-ui/data-table/data-table"
import { EditButton } from "@/components/refine-ui/buttons/edit"
import { DeleteButton } from "@/components/refine-ui/buttons/delete"
import {
  ListView,
  ListViewHeader,
} from "@/components/refine-ui/views/list-view"
import type { User } from "@/api/generated/contract"
import { InviteUserButton } from "./invite-button"

export function UserListPage() {
  const { translate } = useTranslation()

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
            <EditButton recordItemId={row.original.id} size="sm" />
            <DeleteButton recordItemId={row.original.id} size="sm" />
          </div>
        ),
      },
    ],
    [translate],
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
      <ListViewHeader>
        <InviteUserButton />
      </ListViewHeader>
      <DataTable table={table} />
    </ListView>
  )
}
