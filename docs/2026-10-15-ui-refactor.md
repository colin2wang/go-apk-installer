# APK 装机助手 · 界面重构方案
**总体思路：换一套"有色温"的设计令牌 → 去卡片化重排主界面 → 用等宽字体重塑日志区 → 加自定义标题栏和微交互。** 下面是针对 Wails + Vue 3 + TS 的完整可运行代码。
## 一、方案速览
| 步骤 | 改动 | 解决的问题 | 工作量 |
|---|---|---|---|
| 1 | 定义 Design Tokens（CSS 变量） | 去掉 Tailwind 默认蓝 + 中性灰 | 10 分钟 |
| 2 | 自定义标题栏（Wails frameless） | 去 Windows 默认标题栏，整体感更强 | 15 分钟 |
| 3 | 主界面重排（App.vue） | 拆卡片，用分隔线 + 分组标题 | 30 分钟 |
| 4 | 日志面板重构 | 等宽字体、色条级别标识 | 20 分钟 |
| 5 | 基础组件（按钮/输入框/Toast） | 去 emoji 图标、去绿色横幅 | 30 分钟 |
---
## 二、安装依赖
```bash
# frontend 目录下
npm i -D tailwindcss postcss autoprefixer
npx tailwindcss init -p
# 图标库（lucide-vue-next 已废弃，改用 @lucide/vue）
npm i @lucide/vue
# 等宽字体（本地加载，避免打包网络字体）
# 从 https://www.jetbrains.com/lp/mono/ 下载 woff2 放到 src/assets/fonts/
```
---
## 三、Design Tokens（`src/styles/tokens.css`）
**这一步决定"去 AI 味"的 80%。** 放弃 `#0f172a` 系中性灰，改用带暖调的琥珀主色。
```css
/* src/styles/tokens.css */
:root {
  /* ===== 背景层次：暖调暗色，而非中性灰蓝 ===== */
  --bg-base:    #14120f;   /* 最底层，带一点暖褐 */
  --bg-surface: #1c1a17;   /* 面板层 */
  --bg-raised:  #232019;   /* 浮起元素（hover、下拉） */
  --bg-input:   #1a1815;   /* 输入框 */
  /* ===== 边框：低对比，弱化"框感" ===== */
  --border-subtle:  #2a2620;
  --border-default: #35302a;
  --border-strong:  #47403a;
  /* ===== 文字 ===== */
  --text-primary:   #ede8e2;
  --text-secondary: #a8a29a;
  --text-muted:     #6e6861;
  --text-inverse:   #14120f;
  /* ===== 主色：琥珀 ===== */
  --accent:        #f5a623;
  --accent-hover:  #ffb84d;
  --accent-active: #e09400;
  --accent-soft:   rgba(245, 166, 35, 0.12);
  --accent-ring:   rgba(245, 166, 35, 0.3);
  /* ===== 语义色（低调，饱和度低） ===== */
  --success: #7cb87c;
  --danger:  #e07a6a;
  --warning: #d4a857;
  --info:    #6a9bd4;
  /* ===== 日志级别色 ===== */
  --log-info:    #8a9aa8;
  --log-warn:    #d4a857;
  --log-error:   #e07a6a;
  --log-success: #7cb87c;
  /* ===== 字体 ===== */
  --font-sans: -apple-system, BlinkMacSystemFont, "Segoe UI",
               "PingFang SC", "Microsoft YaHei", sans-serif;
  --font-mono: "JetBrains Mono", "Sarasa Mono SC", Consolas,
               "SF Mono", monospace;
  /* ===== 圆角：分级，不要全场 8px ===== */
  --radius-sm:   4px;
  --radius-md:   6px;
  --radius-lg:   10px;
  /* ===== 间距基准（收紧 15%） ===== */
  --space-xs: 4px;
  --space-sm: 8px;
  --space-md: 12px;
  --space-lg: 16px;
  --space-xl: 24px;
}
```
### Tailwind 映射（`tailwind.config.js`）
```js
/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{vue,ts}'],
  theme: {
    extend: {
      colors: {
        bg:       'var(--bg-base)',
        surface:  'var(--bg-surface)',
        raised:   'var(--bg-raised)',
        input:    'var(--bg-input)',
        'border-subtle':   'var(--border-subtle)',
        'border-default':  'var(--border-default)',
        'border-strong':   'var(--border-strong)',
        'text-primary':    'var(--text-primary)',
        'text-secondary':  'var(--text-secondary)',
        'text-muted':      'var(--text-muted)',
        accent: {
          DEFAULT: 'var(--accent)',
          hover:   'var(--accent-hover)',
          active:  'var(--accent-active)',
          soft:    'var(--accent-soft)',
        },
      },
      fontFamily: {
        sans: ['var(--font-sans)'],
        mono: ['var(--font-mono)'],
      },
      borderRadius: {
        sm: 'var(--radius-sm)',
        md: 'var(--radius-md)',
        lg: 'var(--radius-lg)',
      },
    },
  },
  plugins: [],
}
```
### 主样式（`src/styles/main.css`）
```css
@import './tokens.css';
@tailwind base;
@tailwind components;
@tailwind utilities;
@layer base {
  html, body, #app {
    height: 100%;
    margin: 0;
    background: var(--bg-base);
    color: var(--text-primary);
    font-family: var(--font-sans);
    font-size: 13px;              /* 整体字号降 1px，信息密度提升 */
    -webkit-font-smoothing: antialiased;
    overflow: hidden;
  }
  /* 全局 focus：用主色 ring，不用浏览器默认蓝框 */
  *:focus-visible {
    outline: 2px solid var(--accent-ring);
    outline-offset: 1px;
  }
  /* 滚动条 */
  ::-webkit-scrollbar { width: 10px; height: 10px; }
  ::-webkit-scrollbar-track { background: transparent; }
  ::-webkit-scrollbar-thumb {
    background: var(--border-default);
    border-radius: 5px;
    border: 2px solid var(--bg-base);
  }
  ::-webkit-scrollbar-thumb:hover { background: var(--border-strong); }
}
```
---
## 四、Wails 无边框窗口配置
**自定义标题栏是最容易"出质感"的改动**，让界面看起来像原生工具而非网页。
### Go 侧（`main.go`）
```go
package main
import (
	"embed"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)
//go:embed all:frontend/dist
var assets embed.FS
func main() {
	app := NewApp()
	err := wails.Run(&options.App{
		Title:     "APK 装机助手",
		Width:     1100,
		Height:    720,
		MinWidth:  900,
		MinHeight: 600,
		AssetServer: &assetserver.Options{Assets: assets},
		Frameless:  true,              // 关键：无边框
		BackgroundColour: &options.RGBA{R: 20, G: 18, B: 15, A: 255},
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			DisableWindowIcon:    false,
		},
		OnStartup:  app.startup,
		Bind:       []interface{}{app},
	})
	if err != nil {
		println("Error:", err.Error())
	}
}
```
拖拽通过 CSS 变量 `--wails-draggable: drag` 实现，任何设置了该属性的元素都可作为拖拽区。
---
## 五、自定义标题栏（`src/components/TitleBar.vue`）
```vue
<script setup lang="ts">
import { Minus, Square, X } from '@lucide/vue'
import { WindowMinimise, Quit } from '../../wailsjs/runtime/runtime'
// 切换最大化（Wails 运行时）
async function toggleMaximize() {
  // 需在 wailsJS 中绑定 IsMinimised/ToggleMaximise，此处简化
  // 可通过 window.runtime 或绑定方法调用
  // @ts-ignore
  window.runtime.WindowToggleMaximise()
}
</script>
<template>
  <header
    class="h-9 flex items-center justify-between pl-4 pr-0 select-none
           border-b border-[color:var(--border-subtle)] bg-[color:var(--bg-surface)]"
    style="--wails-draggable: drag"
  >
    <!-- 左：Logo + 标题 -->
    <div class="flex items-center gap-2 pointer-events-none">
      <!-- 简单几何 logo，不用 emoji -->
      <div class="w-4 h-4 rounded-[3px] bg-[color:var(--accent)]
                  flex items-center justify-center">
        <div class="w-1.5 h-1.5 bg-[color:var(--bg-base)] rounded-[1px]"></div>
      </div>
      <span class="text-[13px] font-medium text-[color:var(--text-primary)]">
        APK 装机助手
      </span>
    </div>
    <!-- 右：窗口控制按钮 -->
    <div class="flex h-full" style="--wails-draggable: no-drag">
      <button
        class="w-12 h-full flex items-center justify-center
               text-[color:var(--text-muted)] hover:bg-[color:var(--bg-raised)]
               hover:text-[color:var(--text-primary)] transition-colors"
        @click="WindowMinimise()"
      >
        <Minus :size="14" :stroke-width="1.5" />
      </button>
      <button
        class="w-12 h-full flex items-center justify-center
               text-[color:var(--text-muted)] hover:bg-[color:var(--bg-raised)]
               hover:text-[color:var(--text-primary)] transition-colors"
        @click="toggleMaximize"
      >
        <Square :size="11" :stroke-width="1.5" />
      </button>
      <button
        class="w-12 h-full flex items-center justify-center
               text-[color:var(--text-muted)] hover:bg-[#c0392b]
               hover:text-white transition-colors"
        @click="Quit()"
      >
        <X :size="15" :stroke-width="1.5" />
      </button>
    </div>
  </header>
</template>
```
---
## 六、主布局（`src/App.vue`）— 去卡片化重排
**核心思路：主表单不套框，直接铺在背景上，用分组标题 + 细分隔线分区。**
```vue
<script setup lang="ts">
import { ref, computed } from 'vue'
import TitleBar from './components/TitleBar.vue'
import DeviceSidebar from './components/DeviceSidebar.vue'
import InstallForm from './components/InstallForm.vue'
import LogPanel from './components/LogPanel.vue'
import StatusBar from './components/StatusBar.vue'
import ToastHost from './components/ToastHost.vue'
const activeDevice = ref<string>('')
</script>
<template>
  <div class="h-screen flex flex-col bg-[color:var(--bg-base)]">
    <TitleBar />
    <!-- 主体三段：侧栏 / 主内容 / 状态栏 -->
    <div class="flex-1 flex min-h-0">
      <DeviceSidebar v-model:active="activeDevice" />
      <main class="flex-1 flex flex-col min-w-0 overflow-hidden">
        <div class="flex-1 overflow-y-auto px-6 py-5">
          <InstallForm :device-serial="activeDevice" />
        </div>
        <!-- 底部固定操作条 -->
        <StatusBar />
      </main>
    </div>
    <!-- 日志区（可折叠） -->
    <LogPanel />
    <ToastHost />
  </div>
</template>
```
### 设备侧栏（`DeviceSidebar.vue`）
```vue
<script setup lang="ts">
import { ref } from 'vue'
import { Smartphone, Plus } from '@lucide/vue'
interface Device {
  name: string
  serial: string
  online: boolean
}
const devices = ref<Device[]>([
  { name: '示例手机', serial: 'auto', online: false },
])
const active = defineModel<string>('active', { required: true })
</script>
<template>
  <aside class="w-52 shrink-0 border-r border-[color:var(--border-subtle)]
                bg-[color:var(--bg-surface)] flex flex-col">
    <!-- 分组标题：小、灰、加字距，不用粗体大标题 -->
    <div class="px-4 pt-4 pb-2 flex items-center justify-between">
      <span class="text-[11px] font-medium tracking-wider
                   text-[color:var(--text-muted)] uppercase">
        设备配置
      </span>
      <button
        class="w-5 h-5 rounded flex items-center justify-center
               text-[color:var(--text-muted)]
               hover:bg-[color:var(--bg-raised)] hover:text-[color:var(--accent)]
               transition-colors"
      >
        <Plus :size="13" :stroke-width="2" />
      </button>
    </div>
    <nav class="flex-1 overflow-y-auto px-2 pb-2">
      <button
        v-for="d in devices" :key="d.name"
        class="w-full text-left px-3 py-2 rounded-[var(--radius-md)]
               transition-colors group"
        :class="active === d.name
          ? 'bg-[color:var(--accent-soft)] text-[color:var(--accent)]'
          : 'text-[color:var(--text-secondary)] hover:bg-[color:var(--bg-raised)]'"
        @click="active = d.name"
      >
        <div class="flex items-center gap-2">
          <!-- 在线状态点 -->
          <span class="w-1.5 h-1.5 rounded-full shrink-0"
                :class="d.online
                  ? 'bg-[color:var(--success)]'
                  : 'bg-[color:var(--border-strong)]'" />
          <span class="text-[13px] truncate">{{ d.name }}</span>
        </div>
        <!-- serial 用等宽小字 -->
        <div class="pl-3.5 mt-0.5 font-mono text-[10px]
                    text-[color:var(--text-muted)] truncate">
          {{ d.serial }}
        </div>
      </button>
    </nav>
  </aside>
</template>
```
### 安装表单（`InstallForm.vue`）— 无卡片设计
```vue
<script setup lang="ts">
import { ref } from 'vue'
import { FolderOpen, HelpCircle } from '@lucide/vue'
defineProps<{ deviceSerial: string }>()
const form = ref({
  name: '示例手机',
  serial: '',
  apkPath: '',
  overwrite: true,
  allowDowngrade: false,
  grantPermissions: false,
  autoLaunch: false,
  note: '',
})
</script>
<template>
  <div class="max-w-2xl space-y-6">
    <!-- 分组：基本信息 -->
    <section>
      <h3 class="text-[11px] font-medium tracking-wider uppercase
                 text-[color:var(--text-muted)] mb-3">
        基本信息
      </h3>
      <div class="grid grid-cols-2 gap-4">
        <div>
          <label class="block text-[12px] text-[color:var(--text-secondary)] mb-1.5">
            配置名称
          </label>
          <input
            v-model="form.name"
            class="w-full h-8 px-3 rounded-[var(--radius-md)]
                   bg-[color:var(--bg-input)]
                   border border-[color:var(--border-default)]
                   text-[13px] text-[color:var(--text-primary)]
                   placeholder:text-[color:var(--text-muted)]
                   focus:border-[color:var(--accent)] focus:outline-none
                   transition-colors"
          />
        </div>
        <div>
          <label class="flex items-center gap-1 text-[12px]
                        text-[color:var(--text-secondary)] mb-1.5">
            设备 Serial
            <HelpCircle :size="12" :stroke-width="1.5"
                        class="text-[color:var(--text-muted)]" />
          </label>
          <input
            v-model="form.serial"
            placeholder="留空自动匹配"
            class="w-full h-8 px-3 rounded-[var(--radius-md)]
                   bg-[color:var(--bg-input)] font-mono text-[12px]
                   border border-[color:var(--border-default)]
                   text-[color:var(--text-primary)]
                   placeholder:text-[color:var(--text-muted)]
                   focus:border-[color:var(--accent)] focus:outline-none
                   transition-colors"
          />
        </div>
      </div>
    </section>
    <!-- 分组：APK 路径 -->
    <section>
      <h3 class="text-[11px] font-medium tracking-wider uppercase
                 text-[color:var(--text-muted)] mb-3">
        安装包
      </h3>
      <div class="flex gap-2">
        <input
          v-model="form.apkPath"
          placeholder="选择 .apk 文件"
          class="flex-1 h-8 px-3 rounded-[var(--radius-md)]
                 bg-[color:var(--bg-input)] font-mono text-[12px]
                 border border-[color:var(--border-default)]
                 text-[color:var(--text-primary)]
                 placeholder:text-[color:var(--text-muted)]
                 focus:border-[color:var(--accent)] focus:outline-none
                 transition-colors"
        />
        <button
          class="h-8 px-3 rounded-[var(--radius-md)]
                 border border-[color:var(--border-default)]
                 bg-[color:var(--bg-raised)]
                 text-[color:var(--text-secondary)]
                 hover:border-[color:var(--border-strong)]
                 hover:text-[color:var(--text-primary)]
                 transition-colors flex items-center gap-1.5"
        >
          <FolderOpen :size="13" :stroke-width="1.5" />
          浏览
        </button>
      </div>
    </section>
    <!-- 分组：安装选项（一行内联，不各占一块） -->
    <section>
      <h3 class="text-[11px] font-medium tracking-wider uppercase
                 text-[color:var(--text-muted)] mb-3">
        安装选项
      </h3>
      <div class="flex flex-wrap gap-x-5 gap-y-2">
        <label class="flex items-center gap-2 cursor-pointer group">
          <input type="checkbox" v-model="form.overwrite"
                 class="peer sr-only" />
          <span class="w-3.5 h-3.5 rounded-[3px] border
                       border-[color:var(--border-strong)]
                       peer-checked:bg-[color:var(--accent)]
                       peer-checked:border-[color:var(--accent)]
                       flex items-center justify-center transition-colors">
            <svg v-if="form.overwrite" width="10" height="10" viewBox="0 0 10 10"
                 fill="none" stroke="var(--text-inverse)" stroke-width="2">
              <path d="M2 5l2.5 2.5L8 3" stroke-linecap="round"/>
            </svg>
          </span>
          <span class="text-[12px] text-[color:var(--text-secondary)]
                       group-hover:text-[color:var(--text-primary)]
                       transition-colors">
            覆盖安装 <span class="text-[color:var(--text-muted)]">-r</span>
          </span>
        </label>
        <label class="flex items-center gap-2 cursor-pointer group">
          <input type="checkbox" v-model="form.allowDowngrade" class="peer sr-only" />
          <span class="w-3.5 h-3.5 rounded-[3px] border
                       border-[color:var(--border-strong)]
                       peer-checked:bg-[color:var(--accent)]
                       peer-checked:border-[color:var(--accent)]
                       transition-colors"></span>
          <span class="text-[12px] text-[color:var(--text-secondary)]
                       group-hover:text-[color:var(--text-primary)]">
            允许降级 <span class="text-[color:var(--text-muted)]">-d</span>
          </span>
        </label>
        <!-- 其余两个 checkbox 同上 -->
      </div>
    </section>
    <!-- 备注：弱化 -->
    <section>
      <textarea
        v-model="form.note"
        rows="2"
        placeholder="备注（可选）"
        class="w-full px-3 py-2 rounded-[var(--radius-md)] resize-none
               bg-transparent border border-dashed
               border-[color:var(--border-default)]
               text-[12px] text-[color:var(--text-secondary)]
               placeholder:text-[color:var(--text-muted)]
               focus:border-[color:var(--accent)] focus:outline-none
               transition-colors"
      />
    </section>
  </div>
</template>
```
### 底部操作条（`StatusBar.vue`）
```vue
<script setup lang="ts">
import {ref} from 'vue'
import { Zap, Plug, Save } from '@lucide/vue'
import { useToast } from '../composables/useToast'
const toast = useToast()
const installing = ref(false)
async function handleInstall() {
  installing.value = true
  try {
    // 调用 Go 方法
    // await InstallApk(...)
    toast.success('安装成功')
  } catch (e: any) {
    toast.error(e.message || '安装失败')
  } finally {
    installing.value = false
  }
}
</script>
<template>
  <!-- 右对齐主按钮，次按钮弱化为 ghost -->
  <footer class="shrink-0 h-12 px-6 flex items-center justify-end gap-2
                 border-t border-[color:var(--border-subtle)]
                 bg-[color:var(--bg-surface)]">
    <button
      class="h-8 px-3 rounded-[var(--radius-md)]
             text-[12px] text-[color:var(--text-secondary)]
             hover:bg-[color:var(--bg-raised)]
             hover:text-[color:var(--text-primary)]
             transition-colors flex items-center gap-1.5"
    >
      <Plug :size="13" :stroke-width="1.5" />
      检查连接
    </button>
    <button
      class="h-8 px-3 rounded-[var(--radius-md)]
             text-[12px] text-[color:var(--text-secondary)]
             hover:bg-[color:var(--bg-raised)]
             hover:text-[color:var(--text-primary)]
             transition-colors flex items-center gap-1.5"
    >
      <Save :size="13" :stroke-width="1.5" />
      保存配置
    </button>
    <div class="w-px h-5 bg-[color:var(--border-default)] mx-1"></div>
    <button
      :disabled="installing"
      class="h-8 px-4 rounded-[var(--radius-md)]
             bg-[color:var(--accent)]
             text-[color:var(--text-inverse)] text-[12px] font-medium
             hover:bg-[color:var(--accent-hover)]
             active:scale-[0.98] disabled:opacity-50
             disabled:cursor-not-allowed transition-all
             flex items-center gap-1.5"
      @click="handleInstall"
    >
      <Zap :size="13" :stroke-width="1.5"
           :class="installing && 'animate-pulse'" />
      {{ installing ? '安装中…' : '检查并安装' }}
    </button>
  </footer>
</template>
```
---
## 七、日志面板（`LogPanel.vue`）— 质感关键
```vue
<script setup lang="ts">
import { ref, computed, nextTick, watch } from 'vue'
import { ChevronDown, Trash2, Folder } from '@lucide/vue'
type Level = 'INFO' | 'WARN' | 'ERROR' | 'SUCCESS'
interface LogEntry {
  time: string
  level: Level
  msg: string
}
const logs = ref<LogEntry[]>([
  { time: '16:22:24', level: 'INFO', msg: 'adb devices: 0 台设备' },
  { time: '16:22:24', level: 'INFO', msg: '设备列表已刷新' },
  // ...
])
const filter = ref<'all' | Level>('all')
const collapsed = ref(false)
const filtered = computed(() =>
  filter.value === 'all' ? logs.value : logs.value.filter(l => l.level === filter.value)
)
// 色条颜色映射（不用彩色 tag）
const levelColor: Record<Level, string> = {
  INFO:    'var(--log-info)',
  WARN:    'var(--log-warn)',
  ERROR:   'var(--log-error)',
  SUCCESS: 'var(--log-success)',
}
watch(() => logs.value.length, async () => {
  await nextTick()
  scrollEl?.scrollTo({ top: scrollEl.scrollHeight })
})
const scrollEl = ref<HTMLElement>()
</script>
<template>
  <section
    class="shrink-0 border-t border-[color:var(--border-subtle)]
           bg-[color:var(--bg-surface)] flex flex-col"
    :style="{ height: collapsed ? '32px' : '180px' }"
  >
    <!-- 头部：可折叠 -->
    <div class="h-8 shrink-0 px-4 flex items-center gap-3
                border-b border-[color:var(--border-subtle)]">
      <button
        class="flex items-center gap-1.5 text-[11px] font-medium
               tracking-wider uppercase text-[color:var(--text-muted)]
               hover:text-[color:var(--text-secondary)] transition-colors"
        @click="collapsed = !collapsed"
      >
        <ChevronDown :size="12" :stroke-width="2"
                     class="transition-transform"
                     :class="collapsed && '-rotate-90'" />
        运行日志
      </button>
      <span class="font-mono text-[10px] text-[color:var(--text-muted)]">
        {{ logs.length }}
      </span>
      <div class="flex-1"></div>
      <!-- 级别过滤：文字切换，不用彩色 tag -->
      <div class="flex gap-1">
        <button
          v-for="f in ['all', 'INFO', 'WARN', 'ERROR'] as const"
          :key="f"
          class="px-2 h-5 rounded-[3px] text-[10px] font-mono
                 transition-colors"
          :class="filter === f
            ? 'bg-[color:var(--bg-raised)] text-[color:var(--text-primary)]'
            : 'text-[color:var(--text-muted)] hover:text-[color:var(--text-secondary)]'"
          @click="filter = f"
        >
          {{ f.toUpperCase() }}
        </button>
      </div>
      <button class="w-5 h-5 rounded flex items-center justify-center
                     text-[color:var(--text-muted)]
                     hover:text-[color:var(--text-primary)] transition-colors">
        <Folder :size="12" :stroke-width="1.5" />
      </button>
      <button class="w-5 h-5 rounded flex items-center justify-center
                     text-[color:var(--text-muted)]
                     hover:text-[color:var(--danger)] transition-colors">
        <Trash2 :size="12" :stroke-width="1.5" />
      </button>
    </div>
    <!-- 日志内容 -->
    <div
      v-show="!collapsed"
      ref="scrollEl"
      class="flex-1 overflow-y-auto px-4 py-2 font-mono text-[11px] leading-relaxed"
    >
      <div
        v-for="(log, i) in filtered" :key="i"
        class="flex gap-2 py-0.5 hover:bg-[color:var(--bg-raised)]
               hover:mx-[-8px] hover:px-2 rounded-[3px] transition-colors
               cursor-text group"
      >
        <!-- 左侧色条替代彩色 tag -->
        <span class="w-[2px] shrink-0 rounded-full"
              :style="{ background: levelColor[log.level], opacity: 0.6 }" />
        <span class="text-[color:var(--text-muted)] shrink-0
                     tabular-nums select-none">
          {{ log.time }}
        </span>
        <span class="shrink-0 w-14 select-none"
              :style="{ color: levelColor[log.level], opacity: 0.8 }">
          {{ log.level }}
        </span>
        <span class="text-[color:var(--text-secondary)] break-all">
          {{ log.msg }}
        </span>
      </div>
    </div>
  </section>
</template>
```
---
## 八、Toast 通知（替换绿色横幅）
### `src/composables/useToast.ts`
```ts
import { reactive } from 'vue'
export interface Toast {
  id: number
  type: 'success' | 'error' | 'info'
  msg: string
}
const state = reactive<{ toasts: Toast[] }>({ toasts: [] })
let uid = 0
export function useToast() {
  function push(type: Toast['type'], msg: string, duration = 3000) {
    const id = ++uid
    state.toasts.push({ id, type, msg })
    setTimeout(() => {
      const i = state.toasts.findIndex(t => t.id === id)
      if (i > -1) state.toasts.splice(i, 1)
    }, duration)
  }
  return {
    toasts: state.toasts,
    success: (m: string) => push('success', m),
    error:   (m: string) => push('error', m),
    info:    (m: string) => push('info', m),
  }
}
```
### `src/components/ToastHost.vue`
```vue
<script setup lang="ts">
import { CheckCircle2, AlertCircle, Info } from '@lucide/vue'
import { useToast } from '../composables/useToast'
const { toasts } = useToast()
const icon = {
  success: CheckCircle2,
  error: AlertCircle,
  info: Info,
}
const color = {
  success: 'var(--success)',
  error: 'var(--danger)',
  info: 'var(--info)',
}
</script>
<template>
  <!-- 右上角，不占版面 -->
  <div class="fixed top-12 right-4 z-50 flex flex-col gap-2 pointer-events-none">
    <TransitionGroup name="toast">
      <div
        v-for="t in toasts" :key="t.id"
        class="pointer-events-auto flex items-center gap-2
               px-3 py-2 rounded-[var(--radius-md)]
               bg-[color:var(--bg-raised)]
               border border-[color:var(--border-default)]
               shadow-lg min-w-[200px]"
      >
        <component
          :is="icon[t.type]"
          :size="14" :stroke-width="1.5"
          :style="{ color: color[t.type] }"
        />
        <span class="text-[12px] text-[color:var(--text-primary)]">{{ t.msg }}</span>
      </div>
    </TransitionGroup>
  </div>
</template>
<style scoped>
.toast-enter-active, .toast-leave-active {
  transition: all 0.2s ease;
}
.toast-enter-from { opacity: 0; transform: translateY(-8px); }
.toast-leave-to   { opacity: 0; transform: translateX(20px); }
</style>
```
---
## 九、微交互细节清单（收尾打磨）
| 项目 | 实现要点 |
|---|---|
| **主按钮反馈** | `active:scale-[0.98]` + hover 提亮 8%（`--accent-hover`） |
| **输入框 focus** | `focus:border-[color:var(--accent)]`，去掉浏览器默认蓝框 |
| **日志新行淡入** | `<TransitionGroup>` + `opacity` 200ms |
| **侧栏项激活** | 用 `--accent-soft`（12% 透明主色）背景，不用满色块 |
| **tabular-nums** | 时间戳/序号加 `tabular-nums`，数字对齐 |
| **select-none** | 状态栏、日志时间戳、图标都加 `select-none`，避免误选 |
| **过渡统一** | 所有 `transition-colors` 默认 150ms，不过度动画 |
---
## 十、可选进阶（进一步提升"原生工具感"）
- **Wails v3 alpha 的原生非客户区**支持 Windows 11 Snap Layouts，自定义标题栏可以获得原生窗口贴靠体验
- **透明亚克力效果**：Windows 下开启 `WindowIsTranslucent: true` + 前端 `backdrop-blur`，但会增加 GPU 占用，酌情使用
- **主题切换**：在 `tokens.css` 里写两套 `:root[data-theme='light']` 变量，右上角主题按钮通过 `document.documentElement.dataset.theme` 切换即可
## 关键文件清单
```
frontend/src/
├── styles/
│   ├── tokens.css      ← 设计令牌（第一步要做）
│   └── main.css
├── components/
│   ├── TitleBar.vue
│   ├── DeviceSidebar.vue
│   ├── InstallForm.vue
│   ├── StatusBar.vue
│   ├── LogPanel.vue
│   └── ToastHost.vue
├── composables/
│   └── useToast.ts
└── App.vue
```
建议先落地 `tokens.css` + `TitleBar.vue`（共约 25 分钟），这两项完成后界面的"气质"就会彻底改变，剩下的组件重构可以分批进行。Lucide 图标库在 Vue 3 中的组件化用法支持完整的 props 响应式与 tree-shaking，`@lucide/vue` 是当前官方维护的包名。
