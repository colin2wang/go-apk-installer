import * as WailsApp from '../../wailsjs/go/main/App'
import { EventsOn } from '../../wailsjs/runtime/runtime'
import type { Api } from './types'
import { mockApi } from './mock'
import type { LogLine } from '../types'

// 桌面端（Wails WebView）注入了 window.go；纯浏览器/无头测试走 mock
const isDesktop = typeof window !== 'undefined' && (window as unknown as Record<string, unknown>)['go'] !== undefined

const realApi: Api = {
  GetState: () => WailsApp.GetState(),
  SelectSDKDir: () => WailsApp.SelectSDKDir(),
  SetSDKDir: (dir: string) => WailsApp.SetSDKDir(dir),
  SelectAPKFile: () => WailsApp.SelectAPKFile(),
  ListPhones: () => WailsApp.ListPhones(),
  SavePhone: (oldName: string, p: import('../types').Phone) => WailsApp.SavePhone(oldName, p),
  DeletePhone: (name: string) => WailsApp.DeletePhone(name),
  ListDevices: () => WailsApp.ListDevices(),
  CheckPhone: async (name: string) => {
    const r = await WailsApp.CheckPhone(name)
    // 生成的模型里指针字段是可选项，归一为 null
    return { ok: r.ok, message: r.message, device: r.device ?? null }
  },
  InstallAPK: (name: string) => WailsApp.InstallAPK(name),
  OpenConfigsDir: () => WailsApp.OpenConfigsDir(),
  SetTheme: (theme: string) => WailsApp.SetTheme(theme),
  onLog: (cb: (l: LogLine) => void) => {
    EventsOn('apk:log', (l: LogLine) => cb(l))
  },
}

export const api: Api = isDesktop ? realApi : mockApi
export type { Api }
