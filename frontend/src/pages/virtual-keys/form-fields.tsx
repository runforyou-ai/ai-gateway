import { useState } from "react"
import { useTranslation } from "@refinedev/core"
import type { Control, UseFormWatch } from "react-hook-form"
import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form"
import { Input } from "@/components/ui/input"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import {
  BUDGET_PRESETS,
  BUDGET_RESET_OPTIONS,
  EXPIRES_IN_OPTIONS,
  budgetPresetFromAmount,
} from "./options"

const UNLIMITED = "__unlimited"
const CUSTOM = "__custom"

export type VirtualKeyFormValues = {
  name: string
  maxBudget: string
  expiresIn: string
  budgetReset: string
  disabled?: boolean
  prefix?: string
  spent?: string
  lastUsedAt?: string
}

type Props = {
  control: Control<VirtualKeyFormValues>
  watch: UseFormWatch<VirtualKeyFormValues>
  showMeta?: boolean
}

export function VirtualKeyFormFields({ control, watch, showMeta }: Props) {
  const { translate } = useTranslation()
  const maxBudget = watch("maxBudget") ?? ""
  const derivedPreset = budgetPresetFromAmount(maxBudget)
  const [customMode, setCustomMode] = useState(derivedPreset === "custom")
  const budgetPreset = customMode ? "custom" : derivedPreset

  return (
    <>
      <FormField
        control={control}
        name="name"
        render={({ field }) => (
          <FormItem>
            <FormLabel>{translate("virtualKeys.fields.name")}</FormLabel>
            <FormControl>
              <Input {...field} autoComplete="off" />
            </FormControl>
            <FormMessage />
          </FormItem>
        )}
      />
      {showMeta && (
        <>
          <FormField
            control={control}
            name="prefix"
            render={({ field }) => (
              <FormItem>
                <FormLabel>{translate("virtualKeys.fields.prefix")}</FormLabel>
                <FormControl>
                  <Input
                    value={`${field.value ?? ""}…`}
                    disabled
                    readOnly
                    className="font-mono"
                  />
                </FormControl>
              </FormItem>
            )}
          />
          <FormField
            control={control}
            name="spent"
            render={({ field }) => (
              <FormItem>
                <FormLabel>{translate("virtualKeys.fields.spent")}</FormLabel>
                <FormControl>
                  <Input value={field.value || "0"} disabled readOnly />
                </FormControl>
              </FormItem>
            )}
          />
          <FormField
            control={control}
            name="lastUsedAt"
            render={({ field }) => (
              <FormItem>
                <FormLabel>
                  {translate("virtualKeys.fields.lastUsedAt")}
                </FormLabel>
                <FormControl>
                  <Input
                    value={
                      field.value
                        ? new Date(field.value).toLocaleString()
                        : translate("virtualKeys.never")
                    }
                    disabled
                    readOnly
                  />
                </FormControl>
              </FormItem>
            )}
          />
        </>
      )}
      <FormField
        control={control}
        name="maxBudget"
        render={({ field }) => (
          <FormItem>
            <FormLabel>{translate("virtualKeys.fields.maxBudget")}</FormLabel>
            <Select
              value={
                budgetPreset === ""
                  ? UNLIMITED
                  : budgetPreset === "custom"
                    ? CUSTOM
                    : budgetPreset
              }
              onValueChange={(v) => {
                if (v === UNLIMITED) {
                  setCustomMode(false)
                  field.onChange("")
                  return
                }
                if (v === CUSTOM) {
                  setCustomMode(true)
                  if (derivedPreset !== "custom") field.onChange("")
                  return
                }
                setCustomMode(false)
                field.onChange(v)
              }}
            >
              <FormControl>
                <SelectTrigger className="w-full">
                  <SelectValue />
                </SelectTrigger>
              </FormControl>
              <SelectContent>
                {BUDGET_PRESETS.map((amount) => (
                  <SelectItem key={amount} value={amount}>
                    {translate("virtualKeys.budgetAmount", { amount })}
                  </SelectItem>
                ))}
                <SelectItem value={CUSTOM}>
                  {translate("virtualKeys.budgetCustom")}
                </SelectItem>
                <SelectItem value={UNLIMITED}>
                  {translate("virtualKeys.unlimited")}
                </SelectItem>
              </SelectContent>
            </Select>
            <FormMessage />
          </FormItem>
        )}
      />
      {budgetPreset === "custom" && (
        <FormField
          control={control}
          name="maxBudget"
          render={({ field }) => (
            <FormItem>
              <FormLabel>
                {translate("virtualKeys.fields.customBudget")}
              </FormLabel>
              <FormControl>
                <Input
                  {...field}
                  type="number"
                  min="0"
                  step="any"
                  autoComplete="off"
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
      )}
      <FormField
        control={control}
        name="expiresIn"
        render={({ field }) => (
          <FormItem>
            <FormLabel>{translate("virtualKeys.fields.expiresIn")}</FormLabel>
            <Select
              value={field.value === "" ? UNLIMITED : field.value}
              onValueChange={(v) => field.onChange(v === UNLIMITED ? "" : v)}
            >
              <FormControl>
                <SelectTrigger className="w-full">
                  <SelectValue />
                </SelectTrigger>
              </FormControl>
              <SelectContent>
                {EXPIRES_IN_OPTIONS.map((value) => (
                  <SelectItem key={value} value={value}>
                    {translate(`virtualKeys.expiresIn.${value}`)}
                  </SelectItem>
                ))}
                <SelectItem value={UNLIMITED}>
                  {translate("virtualKeys.expiresIn.unlimited")}
                </SelectItem>
              </SelectContent>
            </Select>
            <FormMessage />
          </FormItem>
        )}
      />
      <FormField
        control={control}
        name="budgetReset"
        render={({ field }) => (
          <FormItem>
            <FormLabel>{translate("virtualKeys.fields.budgetReset")}</FormLabel>
            <Select
              value={field.value === "" ? UNLIMITED : field.value}
              onValueChange={(v) => field.onChange(v === UNLIMITED ? "" : v)}
            >
              <FormControl>
                <SelectTrigger className="w-full">
                  <SelectValue />
                </SelectTrigger>
              </FormControl>
              <SelectContent>
                {BUDGET_RESET_OPTIONS.map((value) => (
                  <SelectItem
                    key={value || "na"}
                    value={value === "" ? UNLIMITED : value}
                  >
                    {translate(
                      value === ""
                        ? "virtualKeys.budgetReset.na"
                        : `virtualKeys.budgetReset.${value}`,
                    )}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            <FormMessage />
          </FormItem>
        )}
      />
    </>
  )
}
