<script setup lang="ts">
import { ref } from 'vue'
import type { SDKStatus } from '../types'
import { t, locale, setLocale, localeOptions } from '../i18n'

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
    <div class="sdk-box">
      <span class="lbl">{{ t('sdk.label') }}</span>
      <template v-if="!editing">
        <span class="sdk-path" :title="props.sdk.sdkDir || t('sdk.notSet')">{{ props.sdk.sdkDir || t('sdk.notSet') }}</span>
        <button class="btn small" :disabled="busy" @click="emit('pick-sdk')">{{ t('sdk.pick') }}</button>
        <button class="btn small ghost" :disabled="busy" @click="startEdit">{{ t('sdk.manual') }}</button>
      </template>
      <template v-else>
        <input
          v-model="dirInput"
          class="inp sdk-input mono"
          :placeholder="t('sdk.placeholder')"
          @keyup.enter="apply"
        />
        <button class="btn small primary" @click="apply">{{ t('sdk.apply') }}</button>
        <button class="btn small ghost" @click="editing = false">{{ t('common.cancel') }}</button>
      </template>
      <span class="badge" :class="props.sdk.ready ? 'ok' : 'bad'">
        {{ props.sdk.ready
          ? t('sdk.ready', { version: props.sdk.version.replace('Android Debug Bridge version ', 'v') })
          : t('sdk.notReady') }}
      </span>
    </div>

    <!-- 语言切换 -->
    <select v-model="locale" class="inp lang-select" title="Language" @change="setLocale(locale)">
      <option v-for="o in localeOptions" :key="o.value" :value="o.value">{{ o.label }}</option>
    </select>

    <button class="btn ghost theme-btn" :title="t('sdk.toggleTheme')" @click="emit('toggle-theme')">
      <svg v-if="props.theme === 'dark'" width="14" height="14" viewBox="0 0 14 14" fill="none" stroke="currentColor" stroke-width="1.3"><path d="M12 8.5A5.5 5.5 0 0 1 5.5 2 5.5 5.5 0 1 0 12 8.5Z" /></svg>
      <svg v-else width="14" height="14" viewBox="0 0 14 14" fill="none" stroke="currentColor" stroke-width="1.3"><circle cx="7" cy="7" r="3" /><path d="M7 1v1.5M7 11.5V13M1 7h1.5M11.5 7H13M2.8 2.8l1 1M10.2 10.2l1 1M11.2 2.8l-1 1M3.8 10.2l-1 1" /></svg>
    </button>
  </header>
</template>

<style scoped>
.lang-select {
  width: auto;
  min-width: 92px;
  cursor: pointer;
}
</style>
