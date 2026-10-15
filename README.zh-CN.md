# APK Installer（APK 装机助手）

[English](README.md) | 简体中文

基于 Wails v2 的桌面应用（Go + Vue 3）：通过 adb 一键安装 APK 到手机，一台手机一个配置文件。

![Platform](https://img.shields.io/badge/platform-Windows-blue) ![Go](https://img.shields.io/badge/Go-1.23+-00ADD8) ![Wails](https://img.shields.io/badge/Wails-v2-red)

## 功能特性

- **一键安装**——选好手机配置和 APK 即可安装；执行前自动做两道前置检查（adb 就绪、设备在线且已授权）。
- **一机一配置**——每台手机对应 `configs/` 下的一个 YAML 文件（名称、serial、APK 路径、安装参数）。serial 留空时，只连接一台手机会自动匹配。
- **安装参数**——覆盖安装（`-r`）、允许降级（`-d`）、授予运行时权限（`-g`，可选；小米/MIUI 会拒绝）、安装后自动启动。
- **实时日志台**——安装过程逐行输出；级别切到「调试」可查看 adb 原始输出。
- **设备面板**——实时显示已连接设备状态（已授权 / 未授权 / 离线）。

## 从源码构建

依赖：[Go 1.21+](https://go.dev/dl/)、[Node.js 16+](https://nodejs.org/)、[pnpm 9+](https://pnpm.io/installation)、[Wails CLI v2](https://wails.io/docs/gettingstarted/installation)（`go install github.com/wailsapp/wails/v2/cmd/wails@latest`）。

```powershell
.\build.ps1              # 完整构建 → build\bin\APK Installer.exe
.\build.ps1 -Clean       # 先清理旧产物
.\build.ps1 -Package     # 构建并打包 release\APK-Installer-windows64.zip
```

未安装 Wails CLI 时自动回退为 `go build` 直接编译（依赖仓库内已生成的 `frontend/wailsjs` 绑定）。

运行测试：`go test ./...`

## 使用说明

1. 运行 `APK Installer.exe`；首次启动会在 exe 旁创建 `configs\` 目录（含全局 `app.yml` 与示例手机配置）。
2. 在顶栏设置 **Android SDK 目录**（含 `platform-tools\adb.exe` 的目录）；徽标变绿即 adb 就绪。
3. 选中（或新建）一台手机，填写 APK 路径，点 **检查并安装**。安装过程与 adb 输出实时滚动到下方日志台。

## 项目结构

```
main.go / app.go        Wails 入口与绑定
internal/adb/           adb 封装（定位、设备列表、流式安装）
internal/config/        YAML 配置管理（configs/ 一机一文件）
frontend/               Vue 3 + Vite 界面（pnpm）
```

## 许可证

MIT
