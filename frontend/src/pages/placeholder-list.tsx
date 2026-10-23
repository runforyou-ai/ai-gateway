import { ListView, ListViewHeader } from "@/components/refine-ui/views/list-view"

type PlaceholderListPageProps = {
  title?: string
  description?: string
}

/** 资源列表占位，待后端 CRUD 接入后再替换。 */
export function PlaceholderListPage({
  title,
  description = "后端尚未接入，此为菜单占位页。",
}: PlaceholderListPageProps) {
  return (
    <ListView>
      <ListViewHeader title={title} canCreate={false} />
      <p className="text-sm text-muted-foreground">{description}</p>
    </ListView>
  )
}
