<script setup lang="ts">
import type { Phone } from '../types'

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
      <span>手机配置（configs/*.yml）</span>
      <button class="btn small primary" @click="emit('create')">＋ 新建</button>
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
          :title="connMap[p.name]?.text || '连接状态未知，可点「检查连接」'"
        />
        <div class="pi-text">
          <div class="pi-name">{{ p.name }}</div>
          <div class="pi-sub mono">{{ p.serial || 'serial 留空·自动匹配' }}</div>
        </div>
        <button class="icon-btn" title="删除该手机配置" @click.stop="emit('delete', p.name)">🗑</button>
      </div>
      <div v-if="phones.length === 0" class="empty">暂无手机配置<br>点击「＋ 新建」创建</div>
    </div>
  </aside>
</template>
