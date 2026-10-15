<script setup lang="ts">
import type { DeviceInfo } from '../types'
import { t } from '../i18n'

defineProps<{
  devices: DeviceInfo[]
  refreshing: boolean
  adbReady: boolean
}>()

const emit = defineEmits<{ (e: 'refresh'): void }>()

function stateText(s: string): string {
  switch (s) {
    case 'device': return t('state.device')
    case 'unauthorized': return t('state.unauthorized')
    case 'offline': return t('state.offline')
    case 'recovery': return t('state.recovery')
    case 'sideload': return t('state.sideload')
    default: return s
  }
}
</script>

<template>
  <section>
    <div class="dev-head">
      <h3 class="group-title">{{ t('devices.title') }}</h3>
      <span class="spacer" />
      <button class="btn small" :disabled="refreshing || !adbReady" @click="emit('refresh')">
        <svg width="12" height="12" viewBox="0 0 12 12" fill="none" stroke="currentColor" stroke-width="1.4"><path d="M10.5 6a4.5 4.5 0 1 1-1.3-3.2M9.5 1v2.3H7.2" /></svg>
        {{ refreshing ? t('devices.refreshing') : t('devices.refresh') }}
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
      <span v-else class="muted">{{ adbReady ? t('devices.none') : t('devices.adbNotReady') }}</span>
    </div>
    <hr class="rule" />
  </section>
</template>

<style scoped>
.dev-head { display: flex; align-items: center; gap: 10px; }
.dev-head .group-title { margin-bottom: 0; }
.spacer { flex: 1; }
</style>
