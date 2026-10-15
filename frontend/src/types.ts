// 与后端 Go 结构体一一对应的类型定义（字段名 = JSON tag）
export interface Phone {
  name: string
  serial: string
  apkPath: string
  reinstall: boolean
  allowDowngrade: boolean
  grantPermissions: boolean
  autoLaunch: boolean
  launchPackage: string
  notes: string
}

export interface DeviceInfo {
  serial: string
  state: string
  model: string
  product: string
}

export interface SDKStatus {
  sdkDir: string
  adbPath: string
  version: string
  ready: boolean
}

export interface AppConfig {
  sdkDir: string
  theme: string
}

export interface AppState {
  config: AppConfig
  sdk: SDKStatus
}

export interface CheckResult {
  ok: boolean
  message: string
  device: DeviceInfo | null
}

export interface InstallOutcome {
  success: boolean
  message: string
  elapsedMs: number
}

export type LogLevel = 'info' | 'ok' | 'warn' | 'err' | 'adb'

export interface LogLine {
  time: string
  level: LogLevel
  text: string
}

export function emptyPhone(name: string): Phone {
  return {
    name,
    serial: '',
    apkPath: '',
    reinstall: true,
    allowDowngrade: false,
    grantPermissions: false,
    autoLaunch: false,
    launchPackage: '',
    notes: '',
  }
}
