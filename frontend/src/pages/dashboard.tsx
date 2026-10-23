import { useGetIdentity, useTranslation } from "@refinedev/core"
import { ListView, ListViewHeader } from "@/components/refine-ui/views/list-view"

type Identity = {
  email?: string
  role?: string
}

export function DashboardPage() {
  const { data: user } = useGetIdentity<Identity>()
  const { translate } = useTranslation()

  return (
    <ListView>
      <ListViewHeader title={translate("dashboard.title")} canCreate={false} />
      <div className="flex flex-col gap-2 text-sm text-muted-foreground">
        <p>{translate("dashboard.intro")}</p>
        {user?.email ? (
          <p>
            {translate("dashboard.currentUser", { email: user.email })}
            {user.role ? `（${user.role}）` : null}
          </p>
        ) : null}
      </div>
    </ListView>
  )
}
