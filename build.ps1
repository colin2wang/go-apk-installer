# ============================================================
# APK Installer Assistant - build & package script (build.ps1)
#
# Usage (run from the project root):
#   .\build.ps1              # full build: frontend + bindings + Windows exe
#   .\build.ps1 -Clean       # clean build\bin and frontend\dist before building
#   .\build.ps1 -Package     # after building, package release\APK-Installer-windows64.zip
#   .\build.ps1 -Clean -Package
#
# Dependencies:
#   · Go 1.21+        https://go.dev/dl/
#   · Node.js 16+     https://nodejs.org/
#   · pnpm 9+         https://pnpm.io/installation
#   · Wails CLI v2    go install github.com/wailsapp/wails/v2/cmd/wails@latest
#     (falls back to a plain `go build` when the Wails CLI is not installed;
#       the fallback relies on the pre-generated frontend\wailsjs bindings)
# ============================================================
param(
    [switch]$Clean,
    [switch]$Package
)

$ErrorActionPreference = "Stop"
$Root = $PSScriptRoot
$ExeName = "APK Installer.exe"

if ($Clean) {
    Write-Host "==> Cleaning old artifacts" -ForegroundColor Yellow
    if (Test-Path "$Root\build\bin")        { Remove-Item -Recurse -Force "$Root\build\bin" }
    if (Test-Path "$Root\frontend\dist")    { Remove-Item -Recurse -Force "$Root\frontend\dist" }
}

Write-Host "==> [1/3] Installing frontend dependencies (pnpm install)" -ForegroundColor Cyan
Push-Location "$Root\frontend"
pnpm install
if ($LASTEXITCODE -ne 0) { Pop-Location; throw "pnpm install failed" }

Write-Host "==> [2/3] Building frontend (vue-tsc type check + vite)" -ForegroundColor Cyan
pnpm run build
if ($LASTEXITCODE -ne 0) { Pop-Location; throw "frontend build failed (vue-tsc type errors)" }
Pop-Location

$wails = Get-Command wails -ErrorAction SilentlyContinue
if ($wails) {
    Write-Host "==> [3/3] Generating Wails bindings + compiling Windows exe" -ForegroundColor Cyan
    Push-Location $Root
    wails generate module
    if ($LASTEXITCODE -ne 0) { Pop-Location; throw "wails generate module failed" }
    wails build -platform windows/amd64
    if ($LASTEXITCODE -ne 0) { Pop-Location; throw "wails build failed" }
    Pop-Location
} else {
    Write-Host "==> [3/3] Wails CLI not found, falling back to plain go build" -ForegroundColor Yellow
    Push-Location $Root
    go build -tags "desktop,production" -ldflags "-w -s -H windowsgui" -o "build\bin\$ExeName" .
    if ($LASTEXITCODE -ne 0) { Pop-Location; throw "go build failed" }
    Pop-Location
}

$exe = Get-ChildItem "$Root\build\bin\*.exe" | Select-Object -First 1
if (-not $exe) { throw "build finished but no exe found" }
Write-Host ("[OK] Build completed: {0}  ({1:N0} KB)" -f $exe.FullName, ($exe.Length / 1KB)) -ForegroundColor Green

if ($Package) {
    Write-Host "==> Packaging release" -ForegroundColor Cyan
    $rel = "$Root\release"
    if (Test-Path $rel) { Remove-Item -Recurse -Force $rel }
    New-Item -ItemType Directory -Path "$rel\APK Installer" | Out-Null
    Copy-Item $exe.FullName "$rel\APK Installer\"
    Copy-Item "$Root\使用说明.md" "$rel\APK Installer\"
    Compress-Archive -Path "$rel\APK Installer" -DestinationPath "$rel\APK-Installer-windows64.zip" -Force
    Write-Host "[OK] Packaging completed: $rel\APK-Installer-windows64.zip" -ForegroundColor Green
}
