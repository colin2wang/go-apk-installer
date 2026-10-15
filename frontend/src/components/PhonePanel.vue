<script setup lang="ts">
import { reactive, watch } from 'vue'
import { api } from '../api'
import type { CheckResult, Phone } from '../types'
import { t } from '../i18n'

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
  <section>
    <template v-if="props.phone">
      <!-- 分组：基本信息 -->
      <h3 class="group-title">{{ t('form.basic') }}</h3>
      <div class="grid2">
        <label class="field">
          <span class="f-lbl">{{ t('form.name') }} <i>{{ t('form.nameHint') }}</i></span>
          <input v-model="form.name" class="inp" :placeholder="t('form.namePlaceholder')" />
        </label>
        <label class="field">
          <span class="f-lbl">{{ t('form.serial') }} <i>{{ t('form.serialHint') }}</i></span>
          <input v-model="form.serial" class="inp mono" :placeholder="t('form.serialPlaceholder')" />
        </label>
      </div>

      <hr class="rule" />

      <!-- 分组：安装包 -->
      <h3 class="group-title">{{ t('form.package') }}</h3>
      <label class="field">
        <span class="f-lbl">{{ t('form.apkPath') }} <i>{{ t('form.apkPathHint') }}</i></span>
        <div class="row">
          <input v-model="form.apkPath" class="inp mono grow" :placeholder="t('form.apkPlaceholder')" />
          <button class="btn" @click="pickApk">
            <svg width="13" height="13" viewBox="0 0 13 13" fill="none" stroke="currentColor" stroke-width="1.3"><path d="M1.5 10.5V3a1 1 0 0 1 1-1h2.6l1.2 1.5h4.2a1 1 0 0 1 1 1v6a1 1 0 0 1-1 1h-8a1 1 0 0 1-1-1Z" /></svg>
            {{ t('form.browse') }}
          </button>
        </div>
      </label>

      <hr class="rule" />

      <!-- 分组：安装选项 -->
      <h3 class="group-title">{{ t('form.options') }}</h3>
      <div class="opts">
        <label class="chk"><input v-model="form.reinstall" type="checkbox" /> {{ t('form.reinstall') }} <i>{{ t('form.reinstallHint') }}</i></label>
        <label class="chk"><input v-model="form.allowDowngrade" type="checkbox" /> {{ t('form.downgrade') }} <i>-d</i></label>
        <label class="chk"><input v-model="form.grantPermissions" type="checkbox" /> {{ t('form.grantPerms') }} <i>-g</i></label>
        <label class="chk"><input v-model="form.autoLaunch" type="checkbox" /> {{ t('form.autoLaunch') }}</label>
      </div>
      <input
        v-if="form.autoLaunch"
        v-model="form.launchPackage"
        class="inp mono pkg"
        :placeholder="t('form.launchPackagePlaceholder')"
      />

      <!-- 备注：弱化（虚线边框） -->
      <textarea
        v-model="form.notes"
        class="inp ta note-ta"
        rows="2"
        :placeholder="t('form.notesPlaceholder')"
      />

      <!-- 操作条：右对齐，主按钮 + ghost 次按钮 -->
      <div class="actions">
        <button class="btn" :disabled="props.busyInstall" @click="submitSave">
          <svg width="13" height="13" viewBox="0 0 13 13" fill="none" stroke="currentColor" stroke-width="1.3"><path d="M2 2h7.5L11.5 4v7a.5.5 0 0 1-.5.5H2a.5.5 0 0 1-.5-.5v-8.5A.5.5 0 0 1 2 2Z" /><path d="M4 2v3h5V2M4 11V7.5h5V11" /></svg>
          {{ t('action.save') }}
        </button>
        <button class="btn" :disabled="props.busyCheck || props.busyInstall" @click="submitCheck">
          <svg width="13" height="13" viewBox="0 0 13 13" fill="none" stroke="currentColor" stroke-width="1.3"><path d="M4.5 8.5v2a2 2 0 0 0 4 0v-2M2 5.5h9M6.5 1.5v2M9 1.2v2.3" /></svg>
          {{ props.busyCheck ? t('action.checking') : t('action.check') }}
        </button>
        <span class="actions-sep" />
        <button class="btn primary big" :disabled="props.busyInstall" @click="submitInstall">
          <svg width="13" height="13" viewBox="0 0 13 13" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M6.5 1.5v7M3.5 5.5l3 3 3-3M2 11.5h9" /></svg>
          {{ props.busyInstall ? t('action.installing') : t('action.install') }}
        </button>
      </div>

      <div v-if="props.checkResult" class="result" :class="props.checkResult.ok ? 'ok' : 'bad'">
        {{ (props.checkResult.ok ? t('check.ok') : t('check.fail')) + props.checkResult.message }}
      </div>
      <div v-if="props.installMsg" class="result" :class="props.installMsg.ok ? 'ok' : 'bad'">
        {{ props.installMsg.text }}
      </div>
      <div class="hint">{{ t('form.hint') }}</div>
    </template>

    <div v-else class="empty big-empty">{{ t('form.empty') }}</div>
  </section>
</template>

<style scoped>
.note-ta {
  background: transparent;
  border-style: dashed;
  border-color: var(--border-default);
  color: var(--text-secondary);
  font-size: 12px;
  resize: none;
}
.actions { justify-content: flex-end; }
.actions-sep { width: 1px; height: 18px; background: var(--border-default); margin: 0 2px; }
</style>
