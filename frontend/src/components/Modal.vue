<script setup lang="ts">
import { t } from '../i18n'
defineProps<{ show: boolean; title: string; danger?: boolean }>()
const emit = defineEmits<{ (e: 'close'): void; (e: 'confirm'): void }>()
</script>

<template>
  <teleport to="body">
    <div v-if="show" class="modal-mask" @click.self="emit('close')">
      <div class="modal" :class="{ danger }">
        <div class="modal-title">{{ title }}</div>
        <div class="modal-body"><slot /></div>
        <div class="modal-foot">
          <button class="btn ghost" @click="emit('close')">{{ t('common.cancel') }}</button>
          <button class="btn" :class="danger ? 'danger' : 'primary'" @click="emit('confirm')">{{ t('common.confirm') }}</button>
        </div>
      </div>
    </div>
  </teleport>
</template>
