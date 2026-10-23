import { useGetIdentity } from "@refinedev/core"
import { ListView, ListViewHeader } from "@/components/refine-ui/views/list-view"

type Identity = {
  email?: string
  role?: string
}

export function DashboardPage() {
  const { data: user } = useGetIdentity<Identity>()

  return (
    <ListView>
      <ListViewHeader title="概览" canCreate={false} />
      <div className="flex flex-col gap-2 text-sm text-muted-foreground">
        <p>管理台骨架已接好（Refine Layout 01）。后续在此展示用量与状态。</p>
        {user?.email ? (
          <p>
            当前用户：{user.email}
            {user.role ? `（${user.role}）` : null}
          </p>
        ) : null}
      </div>
    </ListView>
  )
}
