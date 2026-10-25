import type { AccessControlProvider } from "@refinedev/core"
import { authProvider } from "@/providers/auth"

/** 仅管理员可访问的资源。 */
const adminOnly = new Set(["users", "settings"])

type Identity = {
  role?: string
}

export const accessControlProvider: AccessControlProvider = {
  can: async ({ resource }) => {
    if (!resource || !adminOnly.has(resource)) {
      return { can: true }
    }
    try {
      const identity = (await authProvider.getIdentity?.({})) as
        | Identity
        | undefined
      return { can: identity?.role === "admin" }
    } catch {
      return { can: false }
    }
  },
}
