<script setup lang="ts">
import type { Phone } from '../types'
import { t } from '../i18n'

defineProps<{
  phones: Phone[]
  selected: string
  connMap: Record<string, { ok: boolean; text: string }>
}>()

const emit = defineEmits<{
  (e: 'select', name: string): void
  (e: 'create'): void
  (e: 'delete', name: string): void
}>()
</script>

<template>
  <aside class="phone-list">
    <div class="side-head">
      <span>{{ t('sidebar.title') }} · CONFIGS</span>
      <button class="icon-btn" :title="t('sidebar.create')" @click="emit('create')">
        <svg width="13" height="13" viewBox="0 0 13 13" stroke="currentColor" stroke-width="1.8"><path d="M6.5 2v9M2 6.5h9" /></svg>
      </button>
    </div>
    <div class="side-body">
      <div
        v-for="p in phones"
        :key="p.name"
        class="phone-item"
        :class="{ active: p.name === selected }"
        @click="emit('select', p.name)"
      >
        <span
          class="dot"
          :class="connMap[p.name] ? (connMap[p.name].ok ? 'ok' : 'bad') : 'unknown'"
          :title="connMap[p.name]?.text || t('sidebar.connUnknown')"
        />
        <div class="pi-text">
          <div class="pi-name">{{ p.name }}</div>
          <div class="pi-sub">{{ p.serial || t('sidebar.serialAuto') }}</div>
        </div>
        <button class="icon-btn" :title="t('sidebar.delete')" @click.stop="emit('delete', p.name)">
          <svg width="12" height="12" viewBox="0 0 12 12" fill="none" stroke="currentColor" stroke-width="1.3"><path d="M2 3h8M4.5 3V1.8h3V3M3.2 3l.5 7.2h4.6L8.8 3M5 5v3.4M7 5v3.4" /></svg>
        </button>
      </div>
      <div v-if="phones.length === 0" class="empty">{{ t('sidebar.empty') }}<br>{{ t('sidebar.emptyHint') }}</div>
    </div>
  </aside>
</template>
