<script setup lang="ts">
import { ref } from 'vue'
import type { SDKStatus } from '../types'

const props = defineProps<{
  sdk: SDKStatus
  theme: string
  busy: boolean
}>()

const emit = defineEmits<{
  (e: 'pick-sdk'): void
  (e: 'apply-sdk', dir: string): void
  (e: 'toggle-theme'): void
}>()

const dirInput = ref('')
const editing = ref(false)

function startEdit() {
  dirInput.value = props.sdk.sdkDir
  editing.value = true
}

function apply() {
  emit('apply-sdk', dirInput.value.trim())
  editing.value = false
}
</script>

<template>
  <header class="topbar">
    <div class="brand">
      <span class="logo">📦</span>
      <div>
        <div class="title">APK 装机助手</div>
        <div class="sub">ADB · 一键安装 / 一机一配置</div>
      </div>
    </div>

    <div class="sdk-box">
      <span class="lbl">Android SDK</span>
      <template v-if="!editing">
        <span class="sdk-path" :title="props.sdk.sdkDir || '未设置'">{{ props.sdk.sdkDir || '未设置' }}</span>
        <button class="btn" :disabled="busy" @click="emit('pick-sdk')">选择目录</button>
        <button class="btn ghost" :disabled="busy" @click="startEdit">手动填写</button>
      </template>
      <template v-else>
        <input
          v-model="dirInput"
          class="inp sdk-input mono"
          placeholder="例如 C:\Users\<你>\AppData\Local\Android\Sdk"
          @keyup.enter="apply"
        />
        <button class="btn primary" @click="apply">应用</button>
        <button class="btn ghost" @click="editing = false">取消</button>
      </template>
      <span class="badge" :class="props.sdk.ready ? 'ok' : 'bad'">
        {{ props.sdk.ready ? `✓ adb 就绪（${props.sdk.version.replace('Android Debug Bridge version ', 'v')}）` : '✕ adb 未就绪' }}
      </span>
    </div>

    <button class="btn ghost theme-btn" title="切换主题" @click="emit('toggle-theme')">
      {{ props.theme === 'dark' ? '🌙' : '☀️' }}
    </button>
  </header>
</template>
