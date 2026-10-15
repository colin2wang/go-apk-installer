# APK Installer

English | [简体中文](README.zh-CN.md)

A Wails v2 desktop app (Go + Vue 3) that installs APK files to Android phones via adb — one config file per phone, one-click install.

![Platform](https://img.shields.io/badge/platform-Windows-blue) ![Go](https://img.shields.io/badge/Go-1.23+-00ADD8) ![Wails](https://img.shields.io/badge/Wails-v2-red)

## Features

- **One-click install** — pick a phone config and an APK, click install; the full two-step pre-check (adb ready + device connected/authorized) runs before anything is pushed to the device.
- **Per-phone configs** — each phone is a single YAML file under `configs/` (name, serial, APK path, install flags). Serial left blank auto-matches the only connected device.
- **Install options** — reinstall (`-r`), allow downgrade (`-d`), grant runtime permissions (`-g`, opt-in; rejected by MIUI), auto-launch after install.
- **Live log console** — install output streams line by line; switch the level filter to **Debug** to see raw adb output.
- **Device panel** — shows connected devices with state (authorized / unauthorized / offline).

## Build from source

Prerequisites: [Go 1.21+](https://go.dev/dl/), [Node.js 16+](https://nodejs.org/), [pnpm 9+](https://pnpm.io/installation), [Wails CLI v2](https://wails.io/docs/gettingstarted/installation) (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`).

```powershell
.\build.ps1              # full build → build\bin\APK Installer.exe
.\build.ps1 -Clean       # clean old artifacts first
.\build.ps1 -Package     # build and package release\APK-Installer-windows64.zip
```

Without the Wails CLI the script falls back to a plain `go build`, relying on the pre-generated `frontend/wailsjs` bindings.

Run tests: `go test ./...`

## Usage

1. Run `APK Installer.exe`; on first launch it creates a `configs\` folder next to the exe (global `app.yml` + a sample phone config).
2. Set the **Android SDK directory** in the top bar (the folder containing `platform-tools\adb.exe`); the badge turns green when adb is ready.
3. Select a phone (or create one), fill in the APK path, click **Check & Install**. Progress and adb output scroll into the log console below.

## Project layout

```
main.go / app.go        Wails entry & bindings
internal/adb/           adb wrapper (locate, devices, streaming install)
internal/config/        YAML config manager (one file per phone in configs/)
frontend/               Vue 3 + Vite UI (pnpm)
```

## License

MIT
