export const EXPIRES_IN_OPTIONS = [
  "1h",
  "1d",
  "3d",
  "7d",
  "30d",
  "90d",
  "180d",
  "365d",
] as const

export const BUDGET_PRESETS = [
  "10",
  "50",
  "100",
  "200",
  "300",
  "500",
  "1000",
] as const

export type BudgetPreset = (typeof BUDGET_PRESETS)[number] | "custom" | ""

export const BUDGET_RESET_OPTIONS = ["", "daily", "weekly", "monthly"] as const

export function budgetPresetFromAmount(amount: string): BudgetPreset {
  if (!amount) return ""
  if ((BUDGET_PRESETS as readonly string[]).includes(amount)) {
    return amount as BudgetPreset
  }
  return "custom"
}
