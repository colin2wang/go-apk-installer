<script setup lang="ts">
import type { DeviceInfo } from '../types'

defineProps<{
  devices: DeviceInfo[]
  refreshing: boolean
  adbReady: boolean
}>()

const emit = defineEmits<{ (e: 'refresh'): void }>()

function stateText(s: string): string {
  switch (s) {
    case 'device': return '在线·已授权'
    case 'unauthorized': return '未授权调试'
    case 'offline': return '离线'
    case 'recovery': return 'Recovery'
    case 'sideload': return 'Sideload'
    default: return s
  }
}
</script>

<template>
  <section class="card">
    <div class="card-head">
      <h3>当前 adb 设备</h3>
      <button class="btn small" :disabled="refreshing || !adbReady" @click="emit('refresh')">
        {{ refreshing ? '刷新中…' : '刷新设备' }}
      </button>
    </div>
    <div class="dev-chips">
      <template v-if="devices.length > 0">
        <span v-for="d in devices" :key="d.serial" class="chip" :class="d.state">
          <span class="dot" :class="d.state === 'device' ? 'ok' : d.state === 'unauthorized' ? 'warn' : 'bad'" />
          <span class="mono">{{ d.serial }}</span>
          <span class="chip-state">{{ stateText(d.state) }}</span>
          <span v-if="d.model" class="chip-model">{{ d.model }}</span>
        </span>
      </template>
      <span v-else class="muted">{{ adbReady ? '未检测到设备，插上手机后点「刷新设备」' : 'adb 未就绪，先在顶部选择 Android SDK 目录' }}</span>
    </div>
  </section>
</template>
