import { useEffect, useState } from "react"
import { installationStatus } from "@/api/generated/operations"
import type { InstallationStatus } from "@/api/generated/contract"

/** 拉取公开安装状态；失败时 status 为 null，仍置 ready。 */
export function useInstallationStatus() {
  const [status, setStatus] = useState<InstallationStatus | null>(null)
  const [ready, setReady] = useState(false)

  useEffect(() => {
    let cancelled = false
    installationStatus()
      .then((next) => {
        if (cancelled) return
        setStatus(next)
        setReady(true)
      })
      .catch(() => {
        if (!cancelled) setReady(true)
      })
    return () => {
      cancelled = true
    }
  }, [])

  return { status, ready }
}
