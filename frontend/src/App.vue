<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import TopBar from './components/TopBar.vue'
import PhoneList from './components/PhoneList.vue'
import PhonePanel from './components/PhonePanel.vue'
import DevicePanel from './components/DevicePanel.vue'
import LogConsole from './components/LogConsole.vue'
import Modal from './components/Modal.vue'
import { api } from './api'
import { emptyPhone } from './types'
import type {
  CheckResult, DeviceInfo, LogLine, LogLevel, Phone, SDKStatus,
} from './types'

const sdk = ref<SDKStatus>({ sdkDir: '', adbPath: '', version: '', ready: false })
const theme = ref('dark')
const phones = ref<Phone[]>([])
const selected = ref('')
const devices = ref<DeviceInfo[]>([])
const logs = ref<LogLine[]>([])
const busyInstall = ref(false)
const busyCheck = ref(false)
const refreshing = ref(false)
const connMap = reactive<Record<string, { ok: boolean; text: string }>>({})
const checkResult = ref<CheckResult | null>(null)
const installMsg = ref<{ ok: boolean; text: string } | null>(null)
const toast = ref('')
const modal = ref<'' | 'create' | 'delete'>('')
const newName = ref('')
const pendingDelete = ref('')
let toastTimer: ReturnType<typeof setTimeout> | undefined

api.onLog(l => pushLogLine(l))

function pushLogLine(l: LogLine) {
  logs.value.push(l)
  if (logs.value.length > 600) logs.value.splice(0, logs.value.length - 600)
}

function log(level: LogLevel, text: string) {
  const d = new Date()
  const p = (n: number) => String(n).padStart(2, '0')
  pushLogLine({ time: `${p(d.getHours())}:${p(d.getMinutes())}:${p(d.getSeconds())}`, level, text })
}

function errMsg(e: unknown): string {
  return e instanceof Error ? e.message : String(e)
}

function showToast(t: string) {
  toast.value = t
  clearTimeout(toastTimer)
  toastTimer = setTimeout(() => (toast.value = ''), 2800)
}

const selectedPhone = computed(() => phones.value.find(p => p.name === selected.value) ?? null)

async function refreshPhones(select?: string) {
  phones.value = await api.ListPhones()
  if (select) selected.value = select
  else if (!phones.value.some(p => p.name === selected.value)) selected.value = phones.value[0]?.name ?? ''
}

async function refreshDevices(silent = false) {
  refreshing.value = true
  try {
    devices.value = (await api.ListDevices()) ?? []
    if (!silent) log('info', `设备列表已刷新：共 ${devices.value.length} 台`)
  } catch (e) {
    if (!silent) {
      log('err', `刷新设备失败：${errMsg(e)}`)
      showToast(errMsg(e))
    }
  } finally {
    refreshing.value = false
  }
}

async function init() {
  try {
    const st = await api.GetState()
    sdk.value = st.sdk
    theme.value = st.config.theme === 'light' ? 'light' : 'dark'
    document.documentElement.dataset['theme'] = theme.value
    await refreshPhones()
    await refreshDevices(true)
    log('info', `初始化完成：手机配置 ${phones.value.length} 台；adb ${st.sdk.ready ? `就绪（${st.sdk.version}）` : '未就绪，请在顶部选择 Android SDK 目录'}`)
  } catch (e) {
    log('err', `初始化失败：${errMsg(e)}`)
  }
}
onMounted(init)

// ---- SDK ----
async function pickSDK() {
  try {
    const st = await api.SelectSDKDir()
    sdk.value = st
    if (st.ready) await refreshDevices(true)
  } catch (e) {
    log('err', errMsg(e))
    showToast(errMsg(e))
  }
}

async function applySDK(dir: string) {
  if (!dir) {
    showToast('请先填写 SDK 目录路径')
    return
  }
  try {
    const st = await api.SetSDKDir(dir)
    sdk.value = st
    if (st.ready) await refreshDevices(true)
  } catch (e) {
    showToast(errMsg(e))
  }
}

// ---- 手机配置 ----
async function saveQuiet(p: Phone) {
  const oldName = selected.value
  await api.SavePhone(oldName, p)
  const renamed = oldName !== p.name
  if (renamed) delete connMap[oldName]
  await refreshPhones(p.name)
}

async function onSave(p: Phone) {
  try {
    await saveQuiet(p)
    showToast(`已保存 configs/${p.name}.yml`)
  } catch (e) {
    showToast(errMsg(e))
  }
}

async function onCheck(p: Phone) {
  busyCheck.value = true
  checkResult.value = null
  installMsg.value = null
  try {
    await saveQuiet(p)
    const r = await api.CheckPhone(p.name)
    checkResult.value = r
    connMap[p.name] = { ok: r.ok, text: r.message }
  } catch (e) {
    checkResult.value = { ok: false, message: errMsg(e), device: null }
    log('err', errMsg(e))
  } finally {
    busyCheck.value = false
  }
}

async function onInstall(p: Phone) {
  busyInstall.value = true
  installMsg.value = null
  checkResult.value = null
  try {
    await saveQuiet(p) // 所见即所装：安装前先落盘当前表单
    const res = await api.InstallAPK(p.name)
    installMsg.value = { ok: res.success, text: `${res.success ? '✔' : '✕'} ${res.message}` }
  } catch (e) {
    const t = errMsg(e)
    installMsg.value = { ok: false, text: `✕ ${t}` }
    log('err', t)
  } finally {
    busyInstall.value = false
  }
}

function confirmCreate() {
  const name = newName.value.trim()
  if (!name) {
    showToast('请输入手机名称')
    return
  }
  api.SavePhone('', emptyPhone(name))
    .then(async () => {
      await refreshPhones(name)
      log('ok', `已创建手机配置 configs/${name}.yml`)
      modal.value = ''
      newName.value = ''
    })
    .catch(e => showToast(errMsg(e)))
}

function confirmDelete() {
  const name = pendingDelete.value
  if (!name) return
  api.DeletePhone(name)
    .then(async () => {
      delete connMap[name]
      await refreshPhones()
      modal.value = ''
      log('info', `已删除手机配置 configs/${name}.yml`)
    })
    .catch(e => showToast(errMsg(e)))
}

async function toggleTheme() {
  theme.value = theme.value === 'dark' ? 'light' : 'dark'
  document.documentElement.dataset['theme'] = theme.value
  try {
    await api.SetTheme(theme.value)
  } catch {
    /* 主题保存失败不阻断 */
  }
}
</script>

<template>
  <div class="app">
    <TopBar :sdk="sdk" :theme="theme" :busy="busyInstall" @pick-sdk="pickSDK" @apply-sdk="applySDK" @toggle-theme="toggleTheme" />

    <div class="body">
      <PhoneList :phones="phones" :selected="selected" :conn-map="connMap" @select="selected = $event" @create="modal = 'create'" @delete="pendingDelete = $event; modal = 'delete'" />

      <main class="main">
        <DevicePanel :devices="devices" :refreshing="refreshing" :adb-ready="sdk.ready" @refresh="refreshDevices()" />
        <PhonePanel
          :phone="selectedPhone"
          :busy-install="busyInstall"
          :busy-check="busyCheck"
          :check-result="checkResult"
          :install-msg="installMsg"
          @save="onSave"
          @check="onCheck"
          @install="onInstall"
          @pick-apk="() => {}"
        />
      </main>
    </div>

    <LogConsole :logs="logs" @clear="logs = []" @open-configs="api.OpenConfigsDir().catch(() => {})" />

    <div class="toast" :class="{ show: toast }">{{ toast }}</div>

    <Modal :show="modal === 'create'" title="新建手机配置" @close="modal = ''" @confirm="confirmCreate">
      <p class="muted">将在 configs/ 目录下创建「<b>{{ newName.trim() || '?' }}</b>.yml」，一台手机一个配置文件。</p>
      <input v-model="newName" class="inp" placeholder="手机名称，如：办公测试机" @keyup.enter="confirmCreate" />
    </Modal>

    <Modal :show="modal === 'delete'" title="删除手机配置" danger @close="modal = ''" @confirm="confirmDelete">
      <p>确定删除「<b>{{ pendingDelete }}</b>」吗？</p>
      <p class="muted">将删除配置文件 configs/{{ pendingDelete }}.yml，此操作不可恢复。</p>
    </Modal>
  </div>
</template>
