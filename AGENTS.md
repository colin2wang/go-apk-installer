# AGENTS.md

Wails v2 + Vue 3 桌面应用：通过 adb 一键安装 APK 到手机（Windows）。

## 命令

- 构建打包：`.\build.ps1`（`-Clean` 清理旧产物，`-Package` 打包 release zip）
- Go 测试：`go test ./...`；静态检查：`go vet ./...`
- 前端（frontend/ 内）：`pnpm install`、`pnpm run build`（含 vue-tsc 类型检查）、`pnpm run dev`
- 包管理器统一用 **pnpm**，不要用 npm（wails.json、build.ps1 均已指向 pnpm）

## 结构

- `main.go` / `app.go`：Wails 入口与绑定方法（package main，必须在根目录）
- `internal/adb/`：adb 封装（定位、devices、install 流式输出）；`internal/config/`：YAML 配置（configs/ 目录一机一配置）
- `frontend/src/`：Vue 界面；`api/index.ts` 按 Wails 绑定分派，`api/mock.ts` 为浏览器无后端时的假实现，改后端方法时两者都要同步

## 项目特有注意事项

- Go module 名是 `apk-installer`（不是目录名 go-apk-installer）
- 子进程必须调 `hideConsole(cmd)`（internal/adb/hide_windows.go），否则 GUI 下每次调 adb 闪黑框
- 返回给前端的切片不能是 nil（JSON 会变 `null` 导致前端组件渲染崩溃），空列表返回 `[]adb.DeviceInfo{}`
- `exec.Cmd` 的 `Stdout` 与 `StdoutPipe()` 互斥，流式输出只走 pipe 并手动累积
- pnpm 11+ 拦截依赖构建脚本，esbuild 白名单在 `frontend/pnpm-workspace.yaml` 的 `allowBuilds`（`package.json` 的 `pnpm` 字段已失效）
- `-g`（授予运行时权限）在小米/MIUI 设备会被 SecurityException 拒绝，默认必须关闭
- 应用统一英文名 "APK Installer"；用户自行编译验证，改动完成后直接交付

## 维护规则

项目结构、构建/测试命令、架构边界、开发约定或本文件记录的其他事实发生变化时，必须在同一次改动中同步更新本文件。
