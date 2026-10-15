<script setup lang="ts">
import { ref, computed, watch, nextTick } from 'vue'
import type { LogLine } from '../types'
import { t } from '../i18n'

const props = defineProps<{ logs: LogLine[] }>()
const emit = defineEmits<{ (e: 'clear'): void; (e: 'open-configs'): void }>()

const body = ref<HTMLElement | null>(null)
const collapsed = ref(false)

type LevelFilter = 'default' | 'info' | 'ok' | 'warn' | 'err' | 'debug'
const level = ref<LevelFilter>('default')

const filters: { value: LevelFilter; label: string }[] = [
  { value: 'default', label: 'ALL' },
  { value: 'info', label: 'INFO' },
  { value: 'ok', label: 'OK' },
  { value: 'warn', label: 'WARN' },
  { value: 'err', label: 'ERROR' },
  { value: 'debug', label: 'ADB' },
]

// debug 视图额外显示 adb 原始输出；其余视图只显示对应级别
const shown = computed<LogLine[]>(() => {
  if (level.value === 'default') return props.logs.filter(l => l.level !== 'adb')
  return props.logs.filter(l => l.level === level.value)
})

watch(
  () => [props.logs.length, shown.value.length] as const,
  async () => {
    await nextTick()
    if (body.value) body.value.scrollTop = body.value.scrollHeight
  },
)
</script>

<template>
  <section class="logconsole" :class="{ collapsed }">
    <!-- 头部：可折叠 + 文字级别过滤 -->
    <div class="log-head">
      <button class="h-title" @click="collapsed = !collapsed">
        <svg
          width="12" height="12" viewBox="0 0 12 12" fill="none"
          stroke="currentColor" stroke-width="2"
          :style="{ transform: collapsed ? 'rotate(-90deg)' : '', transition: 'transform .15s' }"
        ><path d="M2.5 4.5 6 8l3.5-3.5" /></svg>
        {{ t('log.title') }}
      </button>
      <span class="h-count">{{ shown.length }}</span>
      <span class="spacer" />
      <!-- 级别过滤：文字切换 -->
      <div class="log-filter">
        <button
          v-for="f in filters" :key="f.value"
          :class="{ active: level === f.value }"
          @click="level = f.value"
        >{{ f.label }}</button>
      </div>
      <button class="log-icon" :title="t('log.openConfigs')" @click="emit('open-configs')">
        <svg width="12" height="12" viewBox="0 0 13 13" fill="none" stroke="currentColor" stroke-width="1.3"><path d="M1.5 10.5V3a1 1 0 0 1 1-1h2.6l1.2 1.5h4.2a1 1 0 0 1 1 1v6a1 1 0 0 1-1 1h-8a1 1 0 0 1-1-1Z" /></svg>
      </button>
      <button class="log-icon danger" :title="t('log.clear')" @click="emit('clear')">
        <svg width="12" height="12" viewBox="0 0 12 12" fill="none" stroke="currentColor" stroke-width="1.3"><path d="M2 3h8M4.5 3V1.8h3V3M3.2 3l.5 7.2h4.6L8.8 3M5 5v3.4M7 5v3.4" /></svg>
      </button>
    </div>
    <!-- 日志内容：等宽字体 + 左侧色条 -->
    <div v-show="!collapsed" ref="body" class="log-body">
      <div v-for="(l, i) in shown" :key="i" class="log-line" :class="l.level">
        <span class="bar" />
        <span class="lt">{{ l.time }}</span><span class="lv">{{ l.level.toUpperCase() }}</span><span class="tx">{{ l.text }}</span>
      </div>
      <div v-if="shown.length === 0" class="muted pad">{{ t('log.empty') }}</div>
    </div>
  </section>
</template>
