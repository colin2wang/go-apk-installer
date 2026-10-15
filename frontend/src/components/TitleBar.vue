<script setup lang="ts">
import { WindowMinimise, WindowToggleMaximise, Quit } from '../../wailsjs/runtime/runtime'
import { t } from '../i18n'
</script>

<template>
  <header class="titlebar" style="--wails-draggable: drag">
    <!-- 左：几何 Logo + 标题 -->
    <div class="tb-brand">
      <div class="tb-logo"><i /></div>
      <span class="tb-title">APK Installer</span>
    </div>
    <div class="tb-spacer" />
    <!-- 右：窗口控制按钮 -->
    <div class="tb-controls" style="--wails-draggable: no-drag">
      <button class="tb-btn" :title="t('title.minimize')" @click="WindowMinimise()">
        <svg width="14" height="14" viewBox="0 0 14 14" stroke="currentColor" stroke-width="1.5"><path d="M3 7h8" /></svg>
      </button>
      <button class="tb-btn" :title="t('title.maximize')" @click="WindowToggleMaximise()">
        <svg width="11" height="11" viewBox="0 0 11 11" fill="none" stroke="currentColor" stroke-width="1.5"><rect x="1" y="1" width="9" height="9" rx="1" /></svg>
      </button>
      <button class="tb-btn tb-close" :title="t('title.close')" @click="Quit()">
        <svg width="15" height="15" viewBox="0 0 15 15" stroke="currentColor" stroke-width="1.5"><path d="M4 4l7 7M11 4l-7 7" /></svg>
      </button>
    </div>
  </header>
</template>

<style scoped>
.titlebar {
  height: 36px;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  padding-left: 14px;
  user-select: none;
  background: var(--bg-surface);
  border-bottom: 1px solid var(--border-subtle);
}
.tb-brand { display: flex; align-items: center; gap: 8px; pointer-events: none; }
.tb-logo {
  width: 16px; height: 16px; border-radius: 3px;
  background: var(--accent);
  display: flex; align-items: center; justify-content: center;
}
.tb-logo i { width: 6px; height: 6px; border-radius: 1px; background: var(--bg-base); }
.tb-title { font-size: 13px; font-weight: 500; color: var(--text-primary); }
.tb-spacer { flex: 1; }
.tb-controls { display: flex; height: 100%; }
.tb-btn {
  width: 46px; height: 100%;
  display: flex; align-items: center; justify-content: center;
  color: var(--text-muted);
  background: transparent; border: none; cursor: pointer;
  transition: background-color 0.15s, color 0.15s;
}
.tb-btn:hover { background: var(--bg-raised); color: var(--text-primary); }
.tb-close:hover { background: #c0392b; color: #fff; }
</style>
