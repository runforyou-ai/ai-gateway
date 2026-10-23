import { WelcomePage as RefineWelcomePage } from "@refinedev/core"
import { useTranslation } from "react-i18next"
import { Link } from "react-router"

/** 框架默认欢迎页；登录入口用 fixed + 更高 z-index，避免被欢迎页全屏层盖住。 */
export function WelcomePage() {
  const { t } = useTranslation()

  return (
    <>
      <Link
        to="/login"
        className="fixed right-4 top-4 z-50 inline-flex h-9 items-center rounded-md border border-white/30 bg-white/10 px-4 text-sm font-medium text-white backdrop-blur-sm transition-colors hover:bg-white/20 sm:right-6 sm:top-6"
      >
        {t("pages.welcome.login")}
      </Link>
      <RefineWelcomePage />
    </>
  )
}
