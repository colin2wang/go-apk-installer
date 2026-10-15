import type {
  AppState, CheckResult, DeviceInfo, InstallOutcome, LogLine, Phone, SDKStatus,
} from '../types'

// Api 后端能力接口：桌面运行时由 Wails 绑定实现，浏览器/无头环境由 mock 实现
export interface Api {
  GetState(): Promise<AppState>
  SelectSDKDir(): Promise<SDKStatus>
  SetSDKDir(dir: string): Promise<SDKStatus>
  SelectAPKFile(): Promise<string>
  ListPhones(): Promise<Phone[]>
  SavePhone(oldName: string, p: Phone): Promise<void>
  DeletePhone(name: string): Promise<void>
  ListDevices(): Promise<DeviceInfo[]>
  CheckPhone(name: string): Promise<CheckResult>
  InstallAPK(name: string): Promise<InstallOutcome>
  OpenConfigsDir(): Promise<void>
  SetTheme(theme: string): Promise<void>
  onLog(cb: (l: LogLine) => void): void
}
