<script setup lang="ts">
import { reactive, watch } from 'vue'
import { api } from '../api'
import type { CheckResult, Phone } from '../types'

const props = defineProps<{
  phone: Phone | null
  busyInstall: boolean
  busyCheck: boolean
  checkResult: CheckResult | null
  installMsg: { ok: boolean; text: string } | null
}>()

const emit = defineEmits<{
  (e: 'save', p: Phone): void
  (e: 'check', p: Phone): void
  (e: 'install', p: Phone): void
}>()

const form = reactive<Phone>({
  name: '', serial: '', apkPath: '',
  reinstall: true, allowDowngrade: false, grantPermissions: false,
  autoLaunch: false, launchPackage: '', notes: '',
})

let loadedFor = ''
watch(
  () => props.phone,
  p => {
    if (!p) return
    if (p.name === loadedFor) return
    loadedFor = p.name
    Object.assign(form, JSON.parse(JSON.stringify(p)))
  },
  { immediate: true, deep: true },
)

// 切换到未加载过的手机时重置（处理同组件切换）
watch(
  () => props.phone?.name,
  (n, o) => {
    if (n !== undefined && n !== o && props.phone && n !== loadedFor) {
      loadedFor = n
      Object.assign(form, JSON.parse(JSON.stringify(props.phone)))
    }
  },
)

function submitSave() {
  if (props.phone) emit('save', { ...form })
}
function submitCheck() {
  if (props.phone) emit('check', { ...form })
}
function submitInstall() {
  if (props.phone) emit('install', { ...form })
}
async function pickApk() {
  try {
    const path = await api.SelectAPKFile()
    if (path) form.apkPath = path
  } catch {
    /* 选择取消或失败时保持原值 */
  }
}
</script>

<template>
  <section class="card">
    <div class="card-head">
      <h3>手机配置与安装</h3>
      <span v-if="props.phone" class="muted small-note">配置文件：configs/{{ props.phone.name }}.yml</span>
    </div>

    <template v-if="props.phone">
      <div class="grid2">
        <label class="field">
          <span class="f-lbl">手机名称 <i>（= yml 文件名）</i></span>
          <input v-model="form.name" class="inp" placeholder="如：办公测试机" />
        </label>
        <label class="field">
          <span class="f-lbl">设备 serial <i>（adb devices 第一列，留空=自动匹配唯一在线设备）</i></span>
          <input v-model="form.serial" class="inp mono" placeholder="如 emulator-5554 / 1A2B3C4D" />
        </label>
      </div>

      <label class="field">
        <span class="f-lbl">APK 安装包路径 <i>（绝对路径，.apk）</i></span>
        <div class="row">
          <input v-model="form.apkPath" class="inp mono grow" placeholder="D:\apk\demo-release.apk" />
          <button class="btn" @click="pickApk">浏览…</button>
        </div>
      </label>

      <div class="opts">
        <label class="chk"><input v-model="form.reinstall" type="checkbox" /> 覆盖安装 <i>-r 保留数据升级</i></label>
        <label class="chk"><input v-model="form.allowDowngrade" type="checkbox" /> 允许降级 <i>-d</i></label>
        <label class="chk"><input v-model="form.grantPermissions" type="checkbox" /> 授予运行时权限 <i>-g</i></label>
        <label class="chk"><input v-model="form.autoLaunch" type="checkbox" /> 安装后自动启动</label>
        <input
          v-if="form.autoLaunch"
          v-model="form.launchPackage"
          class="inp mono pkg"
          placeholder="包名，如 com.example.app"
        />
      </div>

      <label class="field">
        <span class="f-lbl">备注</span>
        <textarea v-model="form.notes" class="inp ta" rows="2" placeholder="给这台手机留点说明…" />
      </label>

      <div class="actions">
        <button class="btn primary big" :disabled="props.busyInstall" @click="submitInstall">
          {{ props.busyInstall ? '⏳ 检查并安装中…' : '⬇ 检查并安装' }}
        </button>
        <button class="btn" :disabled="props.busyCheck || props.busyInstall" @click="submitCheck">
          {{ props.busyCheck ? '检查中…' : '🔌 检查连接' }}
        </button>
        <button class="btn ghost" :disabled="props.busyInstall" @click="submitSave">💾 保存配置</button>
      </div>

      <div v-if="props.checkResult" class="result" :class="props.checkResult.ok ? 'ok' : 'bad'">
        {{ props.checkResult.ok ? '✓ 连接正常：' : '✕ 未通过：' }}{{ props.checkResult.message }}
      </div>
      <div v-if="props.installMsg" class="result" :class="props.installMsg.ok ? 'ok' : 'bad'">
        {{ props.installMsg.text }}
      </div>
      <div class="hint">
        安装前会自动执行两道检查：① adb 是否可用 → ② adb 能否连到该手机；任一不通过都不会执行安装。
      </div>
    </template>

    <div v-else class="empty big-empty">左侧选择或新建一台手机开始</div>
  </section>
</template>
