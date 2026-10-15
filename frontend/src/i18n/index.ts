// 国际化：轻量实现（无运行时依赖），三份文案 + 响应式 t() + 语言切换
import { ref } from 'vue'
import { zhCN } from './zh-CN'
import { en } from './en'
import { ja } from './ja'

export type Locale = 'zh-CN' | 'en' | 'ja'

const messages: Record<Locale, Record<string, string>> = {
  'zh-CN': zhCN,
  en,
  ja,
}

const LOCALE_KEY = 'apk-installer.locale'

function detectLocale(): Locale {
  const saved = localStorage.getItem(LOCALE_KEY)
  if (saved && saved in messages) return saved as Locale
  const nav = navigator.language.toLowerCase()
  if (nav.startsWith('zh')) return 'zh-CN'
  if (nav.startsWith('ja')) return 'ja'
  return 'en'
}

export const locale = ref<Locale>(detectLocale())

type Vars = Record<string, string | number>

const interpolate = (s: string, vars?: Vars): string =>
  vars ? s.replace(/\{(\w+)\}/g, (m, k) => (k in vars ? String(vars[k]) : m)) : s

// t(key, vars?)：按 key 取当前语言文案，缺失时回退中文，再回退 key 本身。
// 普通函数（内部读 locale.value），脚本与模板中均可调用；模板内调用同样具备响应式。
export function t(key: string, vars?: Vars): string {
  const table = messages[locale.value]
  const fallback = messages['zh-CN']
  return interpolate(table[key] ?? fallback[key] ?? key, vars)
}

export function setLocale(l: Locale) {
  locale.value = l
  localStorage.setItem(LOCALE_KEY, l)
}

export const localeOptions: { value: Locale; label: string }[] = [
  { value: 'zh-CN', label: '中文' },
  { value: 'en', label: 'English' },
  { value: 'ja', label: '日本語' },
]
