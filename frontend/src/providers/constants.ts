/** 与 `@/api/client` 一致：空则同域相对路径；开发可设 VITE_API_URL 直连后端。 */
const viteBase =
  (import.meta.env.VITE_API_URL as string | undefined)?.replace(/\/$/, "") ?? ""

/** Refine dataProvider 的 API 根路径（simple-rest 会拼 `/{resource}`）。 */
export const API_URL = `${viteBase}/api`
