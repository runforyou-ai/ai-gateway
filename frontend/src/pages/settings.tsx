import { useTranslation } from "@refinedev/core"
import { ListView, ListViewHeader } from "@/components/refine-ui/views/list-view"

export function SettingsPage() {
  const { translate } = useTranslation()

  return (
    <ListView>
      <ListViewHeader title={translate("settings.title")} canCreate={false} />
      <p className="text-sm text-muted-foreground">
        {translate("settings.placeholder")}
      </p>
    </ListView>
  )
}
