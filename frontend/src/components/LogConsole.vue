<script setup lang="ts">
import { ref, computed, watch, nextTick } from 'vue'
import type { LogLine } from '../types'

const props = defineProps<{ logs: LogLine[] }>()
const emit = defineEmits<{ (e: 'clear'): void; (e: 'open-configs'): void }>()

const body = ref<HTMLElement | null>(null)

type LevelFilter = 'default' | 'info' | 'ok' | 'warn' | 'err' | 'debug'
const level = ref<LevelFilter>('default')

const filters: { value: LevelFilter; label: string }[] = [
  { value: 'default', label: '全部' },
  { value: 'info', label: '信息' },
  { value: 'ok', label: '成功' },
  { value: 'warn', label: '警告' },
  { value: 'err', label: '错误' },
  { value: 'debug', label: '调试' },
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
  <section class="logconsole">
    <div class="card-head">
      <h3>运行日志</h3>
      <div class="spacer" />
      <select v-model="level" class="inp log-level">
        <option v-for="f in filters" :key="f.value" :value="f.value">{{ f.label }}</option>
      </select>
      <button class="btn small ghost" @click="emit('open-configs')">打开 configs 目录</button>
      <button class="btn small ghost" @click="emit('clear')">清空</button>
    </div>
    <div ref="body" class="log-body mono">
      <div v-for="(l, i) in shown" :key="i" class="log-line" :class="l.level">
        <span class="lt">{{ l.time }}</span><span class="lv">{{ l.level.toUpperCase().padEnd(4) }}</span><span class="tx">{{ l.text }}</span>
      </div>
      <div v-if="shown.length === 0" class="muted pad">暂无日志</div>
    </div>
  </section>
</template>
