import i18n from "i18next"
import { initReactI18next } from "react-i18next"
import LanguageDetector from "i18next-browser-languagedetector"
import en from "@/locales/en/common.json"
import zh from "@/locales/zh/common.json"

void i18n
  .use(LanguageDetector)
  .use(initReactI18next)
  .init({
    resources: {
      en: { common: en },
      zh: { common: zh },
    },
    ns: ["common"],
    defaultNS: "common",
    fallbackLng: "en",
    supportedLngs: ["en", "zh"],
    nonExplicitSupportedLngs: true,
    interpolation: { escapeValue: false },
    detection: {
      // 手动切换写入 localStorage；无记录时再按浏览器语言
      order: ["localStorage", "navigator"],
      caches: ["localStorage"],
      lookupLocalStorage: "ai-gateway.lng",
    },
  })

export default i18n
