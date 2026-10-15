import type { Api } from './types'
import type {
  AppState, CheckResult, DeviceInfo, InstallOutcome, LogLine, Phone, SDKStatus,
} from '../types'

// 浏览器 / 无头测试环境的内存 Mock：与真实绑定同签名。
// 安装流程输出真实 adb install 的典型日志行，便于验证前端渲染。

const sleep = (ms: number) => new Promise<void>(r => setTimeout(r, ms))

function now(): string {
  const d = new Date()
  const p = (n: number) => String(n).padStart(2, '0')
  return `${p(d.getHours())}:${p(d.getMinutes())}:${p(d.getSeconds())}`
}

const DEFAULT_SDK = 'C:\\Users\\demo\\AppData\\Local\\Android\\Sdk'

let appCfg = { sdkDir: '', theme: 'dark' }
let phones: Phone[] = [
  {
    name: '示例手机',
    serial: '',
    apkPath: '',
    reinstall: true,
    allowDowngrade: false,
    grantPermissions: false,
    autoLaunch: false,
    launchPackage: '',
    notes: '示例配置：填好 serial 与 apk_path 即可一键安装。',
  },
  {
    name: '办公测试机',
    serial: 'emulator-5554',
    apkPath: 'D:\\apk\\demo-release.apk',
    reinstall: true,
    allowDowngrade: false,
    grantPermissions: false,
    autoLaunch: false,
    launchPackage: 'com.example.demo',
    notes: '',
  },
]
let devices: DeviceInfo[] = [
  { serial: 'emulator-5554', state: 'device', model: 'Pixel 6', product: 'sdk_gphone64' },
  { serial: '1A2B3C4D', state: 'unauthorized', model: 'Pixel 6 Pro', product: 'raven' },
]
const logCbs: ((l: LogLine) => void)[] = []

function log(level: LogLine['level'], text: string) {
  const l: LogLine = { time: now(), level, text }
  logCbs.forEach(cb => cb(l))
}

function sdkStatus(): SDKStatus {
  if (!appCfg.sdkDir) return { sdkDir: '', adbPath: '', version: '', ready: false }
  return {
    sdkDir: appCfg.sdkDir,
    adbPath: appCfg.sdkDir + '\\platform-tools\\adb.exe',
    version: 'Android Debug Bridge version 1.0.41',
    ready: true,
  }
}

export const mockApi: Api = {
  async GetState(): Promise<AppState> {
    return { config: { ...appCfg }, sdk: sdkStatus() }
  },

  async SelectSDKDir(): Promise<SDKStatus> {
    appCfg.sdkDir = DEFAULT_SDK
    log('ok', `Android SDK 已设置: ${DEFAULT_SDK}（Android Debug Bridge version 1.0.41）`)
    return sdkStatus()
  },

  async SetSDKDir(dir: string): Promise<SDKStatus> {
    const d = dir.trim()
    if (!d) {
      const msg = '未配置 Android SDK 目录，且系统 PATH 中也没有 adb；请先在顶部选择 Android SDK 目录'
      log('warn', msg)
      throw new Error(msg)
    }
    if (d.toLowerCase().includes('invalid')) {
      const msg = '在 Android SDK 目录中未找到 adb（已尝试 platform-tools 与根目录），且系统 PATH 中也没有 adb'
      log('warn', msg)
      throw new Error(msg)
    }
    appCfg.sdkDir = d
    log('ok', `Android SDK 已设置: ${d}（Android Debug Bridge version 1.0.41）`)
    return sdkStatus()
  },

  async SelectAPKFile(): Promise<string> {
    return 'D:\\apk\\new-demo-signed.apk'
  },

  async ListPhones(): Promise<Phone[]> {
    return JSON.parse(JSON.stringify(phones)) as Phone[]
  },

  async SavePhone(oldName: string, p: Phone): Promise<void> {
    const name = p.name.trim()
    if (!name) throw new Error('手机名称不能为空')
    const exists = (n: string) => phones.some(x => x.name === n)
    // 新建重名 / 重命名到已存在名称 → 拒绝
    if ((oldName === '' || oldName !== name) && exists(name)) {
      throw new Error(`已存在同名手机配置「${name}」，请换一个名称`)
    }
    const copy = JSON.parse(JSON.stringify(p)) as Phone
    const i = phones.findIndex(x => x.name === (oldName || name))
    if (i >= 0) phones[i] = copy
    else phones.push(copy)
    log('ok', oldName && oldName !== name
      ? `手机「${oldName}」已重命名为「${name}」`
      : `手机「${name}」配置已保存（configs/${name}.yml）`)
  },

  async DeletePhone(name: string): Promise<void> {
    phones = phones.filter(x => x.name !== name)
    log('info', `已删除手机配置「${name}」`)
  },

  async ListDevices(): Promise<DeviceInfo[]> {
    log('info', `adb devices：${devices.length} 台设备`)
    return JSON.parse(JSON.stringify(devices)) as DeviceInfo[]
  },

  async CheckPhone(name: string): Promise<CheckResult> {
    const p = phones.find(x => x.name === name)
    if (!p) return { ok: false, message: `手机「${name}」的配置文件不存在`, device: null }
    if (p.serial && p.serial === '1A2B3C4D') {
      const msg = '设备 1A2B3C4D 已连接但未授权 USB 调试，请在手机屏幕上点击「允许 USB 调试」后重试'
      log('warn', `「${p.name}」连接检查未通过: ${msg}`)
      return { ok: false, message: msg, device: null }
    }
    const ready = devices.filter(d => d.state === 'device')
    if (p.serial) {
      const d = devices.find(x => x.serial === p.serial)
      if (!d) {
        const msg = `配置的设备 ${p.serial} 未连接（当前在线 ${devices.length} 台设备）`
        log('warn', `「${p.name}」连接检查未通过: ${msg}`)
        return { ok: false, message: msg, device: null }
      }
      if (d.state !== 'device') {
        const msg = `设备 ${d.serial} 当前状态为 ${d.state}，不可安装`
        log('warn', `「${p.name}」连接检查未通过: ${msg}`)
        return { ok: false, message: msg, device: null }
      }
      log('ok', `「${p.name}」连接检查通过：${d.serial}（${d.state}） ${d.model}`)
      return { ok: true, message: `${d.serial}（${d.state}） ${d.model}`, device: d }
    }
    if (ready.length === 1) {
      const d = ready[0]
      log('ok', `「${p.name}」连接检查通过：${d.serial}（${d.state}） ${d.model}`)
      return { ok: true, message: `${d.serial}（${d.state}） ${d.model}`, device: d }
    }
    if (ready.length === 0) {
      const msg = '未检测到已连接的手机，请插上 USB 数据线并在手机上开启 USB 调试'
      log('warn', `「${p.name}」连接检查未通过: ${msg}`)
      return { ok: false, message: msg, device: null }
    }
    const msg = `检测到 ${ready.length} 台在线设备（${ready.map(d => d.serial).join('、')}），请在该手机的配置文件中填写 serial 指定目标设备`
    log('warn', `「${p.name}」连接检查未通过: ${msg}`)
    return { ok: false, message: msg, device: null }
  },

  async InstallAPK(name: string): Promise<InstallOutcome> {
    const start = Date.now()
    const p = phones.find(x => x.name === name)
    if (!p) throw new Error(`手机「${name}」的配置文件不存在`)
    // 与真实后端一致的两道前置检查
    if (!p.apkPath.trim()) {
      log('err', `「${p.name}」未配置 APK 路径，请先选择 .apk 安装包`)
      throw new Error('未配置 APK 路径，请先选择 .apk 安装包')
    }
    if (!appCfg.sdkDir) {
      const msg = '未配置 Android SDK 目录，且系统 PATH 中也没有 adb；请先在顶部选择 Android SDK 目录'
      log('err', `安装中止：${msg}`)
      throw new Error(msg)
    }
    const check = await mockApi.CheckPhone(name)
    if (!check.ok) {
      log('err', `② 连接检查未通过，安装中止：${check.message}`)
      throw new Error(`连接检查未通过：${check.message}`)
    }
    log('info', `▶ 开始安装「${p.name}」`)
    log('info', `① adb 就绪：${appCfg.sdkDir}\\platform-tools\\adb.exe`)
    log('ok', `② 连接检查通过：${check.message}`)
    log('info', `③ 正在安装 ${p.apkPath.split('\\').pop()} ...`)
    const lines = [
      'Performing Streamed Install',
      'Success',
    ]
    for (const line of lines) {
      await sleep(90)
      log('adb', line)
    }
    await sleep(80)
    log('ok', `✔ 安装成功（耗时 ${((Date.now() - start) / 1000).toFixed(1)} 秒）`)
    if (p.autoLaunch && p.launchPackage.trim()) {
      log('ok', `自动启动：${p.launchPackage}（已向手机发送启动指令）`)
    }
    return { success: true, message: '安装成功', elapsedMs: Date.now() - start }
  },

  async OpenConfigsDir(): Promise<void> {
    log('info', '（Mock）打开配置目录 configs/')
  },

  async SetTheme(theme: string): Promise<void> {
    appCfg.theme = theme === 'light' ? 'light' : 'dark'
  },

  onLog(cb: (l: LogLine) => void): void {
    logCbs.push(cb)
  },
}
