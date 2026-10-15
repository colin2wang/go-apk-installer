// Package adb 封装对 Android SDK adb 命令行的调用：
// SDK 目录定位、版本探测、设备列表解析、连接判定与 APK 安装。
// 命令执行抽象为 Runner 接口，便于单元测试注入假输出。
package adb

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// installTimeout 单次 APK 安装的超时上限（大包安装留足余量）。
const installTimeout = 10 * time.Minute

// DeviceInfo 一台 adb 可见设备的信息（来自 adb devices -l）。
type DeviceInfo struct {
	Serial  string `json:"serial"`  // 序列号，如 emulator-5554 / 1A2B3C4D
	State   string `json:"state"`   // device / offline / unauthorized / recovery ...
	Model   string `json:"model"`   // 型号（model: 属性，下划线转空格）
	Product string `json:"product"` // 产品名（product: 属性）
}

// Ready 可安装 = 在线且已授权。
func (d DeviceInfo) Ready() bool { return d.State == "device" }

// Describe 一行人类可读摘要。
func (d DeviceInfo) Describe() string {
	desc := d.Serial + "（" + d.State + "）"
	if d.Model != "" {
		desc += " " + d.Model
	}
	return desc
}

// Runner 命令执行抽象。
type Runner interface {
	// Run 执行并等待结束，返回 stdout/stderr。
	Run(ctx context.Context, exe string, args []string) (string, string, error)
	// RunStreaming 逐行回调 stdout；返回完整 stdout 供结果判定。
	RunStreaming(ctx context.Context, exe string, args []string, onLine func(line string)) (string, error)
}

// ExecRunner 真实进程执行器。
type ExecRunner struct{}

func (ExecRunner) Run(ctx context.Context, exe string, args []string) (string, string, error) {
	cmd := exec.CommandContext(ctx, exe, args...)
	hideConsole(cmd)
	var out, errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	err := cmd.Run()
	return out.String(), errb.String(), err
}

func (ExecRunner) RunStreaming(ctx context.Context, exe string, args []string, onLine func(string)) (string, error) {
	cmd := exec.CommandContext(ctx, exe, args...)
	hideConsole(cmd)
	var out, errb bytes.Buffer
	cmd.Stderr = &errb
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}
	if err := cmd.Start(); err != nil {
		return "", err
	}
	reader := bufio.NewScanner(stdout)
	reader.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for reader.Scan() {
		line := strings.TrimRight(reader.Text(), "\r")
		out.WriteString(line)
		out.WriteString("\n")
		if onLine != nil {
			onLine(line)
		}
	}
	waitErr := cmd.Wait()
	if waitErr != nil {
		// adb 失败时 stderr 里常带具体原因（如 Failed to finalize session），一并返回
		if s := strings.TrimSpace(errb.String()); s != "" {
			out.WriteString(s)
			out.WriteString("\n")
		}
		out.WriteString(waitErr.Error())
		out.WriteString("\n")
	}
	return out.String(), waitErr
}

// ADB 面向一个已定位的 adb 可执行文件。
type ADB struct {
	exe string
	r   Runner
}

func New(exe string, r Runner) *ADB {
	if r == nil {
		r = ExecRunner{}
	}
	return &ADB{exe: exe, r: r}
}

// Exe 返回 adb 可执行文件路径。
func (a *ADB) Exe() string { return a.exe }

// LocateADB 在 Android SDK 目录中定位 adb：
// 依次尝试 <sdk>/platform-tools/adb(.exe)、<sdk>/adb(.exe)，
// 均未命中则回退系统 PATH 中的 adb。
func LocateADB(sdkDir string) (string, error) {
	var candidates []string
	if strings.TrimSpace(sdkDir) != "" {
		if os.PathSeparator == '\\' {
			candidates = []string{
				filepath.Join(sdkDir, "platform-tools", "adb.exe"),
				filepath.Join(sdkDir, "adb.exe"),
			}
		} else {
			candidates = []string{
				filepath.Join(sdkDir, "platform-tools", "adb"),
				filepath.Join(sdkDir, "adb"),
			}
		}
		for _, c := range candidates {
			if st, err := os.Stat(c); err == nil && !st.IsDir() {
				return c, nil
			}
		}
	}
	if p, err := exec.LookPath("adb"); err == nil {
		return p, nil
	}
	if strings.TrimSpace(sdkDir) != "" {
		return "", fmt.Errorf("在 Android SDK 目录中未找到 adb（已尝试 platform-tools 与根目录），且系统 PATH 中也没有 adb；请确认目录正确且已安装 platform-tools")
	}
	return "", fmt.Errorf("未配置 Android SDK 目录，且系统 PATH 中也没有 adb；请先在顶部选择 Android SDK 目录")
}

// Version 返回 adb 版本首行，如 "Android Debug Bridge version 1.0.41"。
func (a *ADB) Version(ctx context.Context) (string, error) {
	c, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	out, errOut, err := a.r.Run(c, a.exe, []string{"version"})
	if err != nil {
		return "", fmt.Errorf("执行 adb version 失败: %v %s", err, strings.TrimSpace(errOut))
	}
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			return line, nil
		}
	}
	return "", fmt.Errorf("adb version 无输出")
}

// parseDevices 解析 `adb devices -l` 输出。
func parseDevices(out string) []DeviceInfo {
	var list []DeviceInfo
	for _, raw := range strings.Split(out, "\n") {
		line := strings.TrimSpace(strings.TrimRight(raw, "\r"))
		if line == "" || strings.HasPrefix(line, "List of devices") || strings.Contains(line, "*") {
			continue
		}
		f := strings.Fields(line)
		if len(f) < 2 {
			continue
		}
		d := DeviceInfo{Serial: f[0], State: f[1]}
		for _, kv := range f[2:] {
			if strings.HasPrefix(kv, "model:") {
				d.Model = strings.ReplaceAll(strings.TrimPrefix(kv, "model:"), "_", " ")
			} else if strings.HasPrefix(kv, "product:") {
				d.Product = strings.TrimPrefix(kv, "product:")
			}
		}
		list = append(list, d)
	}
	return list
}

// Devices 返回当前 adb 可见的全部设备。
func (a *ADB) Devices(ctx context.Context) ([]DeviceInfo, error) {
	return a.DevicesVerbose(ctx, nil)
}

// DevicesVerbose 同 Devices，但把 `adb devices -l` 的每行原始输出回调给 onLine（可为 nil）。
func (a *ADB) DevicesVerbose(ctx context.Context, onLine func(string)) ([]DeviceInfo, error) {
	c, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	out, errOut, err := a.r.Run(c, a.exe, []string{"devices", "-l"})
	if err != nil {
		return nil, fmt.Errorf("执行 adb devices 失败: %v %s", err, strings.TrimSpace(errOut))
	}
	if onLine != nil {
		for _, raw := range strings.Split(out, "\n") {
			line := strings.TrimSpace(strings.TrimRight(raw, "\r"))
			if line != "" {
				onLine(line)
			}
		}
	}
	return parseDevices(out), nil
}

// ResolveDevice 依据手机配置解析目标设备（安装前连接检查的核心判定）。
//   - serial 非空：该序列号必须在线且已授权；
//   - serial 为空：恰好一台可用设备时自动使用；0 台报未连接，多台报歧义。
func ResolveDevice(devices []DeviceInfo, serial string) (*DeviceInfo, error) {
	serial = strings.TrimSpace(serial)
	if serial != "" {
		for i := range devices {
			if devices[i].Serial == serial {
				if !devices[i].Ready() {
					if devices[i].State == "unauthorized" {
						return nil, fmt.Errorf("设备 %s 已连接但未授权 USB 调试，请在手机屏幕上点击「允许 USB 调试」后重试", serial)
					}
					return nil, fmt.Errorf("设备 %s 当前状态为 %s，不可安装", serial, devices[i].State)
				}
				return &devices[i], nil
			}
		}
		return nil, fmt.Errorf("配置的设备 %s 未连接（当前在线 %d 台设备）", serial, len(devices))
	}
	var ready []DeviceInfo
	for _, d := range devices {
		if d.Ready() {
			ready = append(ready, d)
		}
	}
	switch len(ready) {
	case 1:
		return &ready[0], nil
	case 0:
		if len(devices) > 0 {
			var descs []string
			for _, d := range devices {
				descs = append(descs, d.Describe())
			}
			return nil, fmt.Errorf("检测到 %d 台设备但均不可安装（%s）；请检查手机上的 USB 调试授权", len(devices), strings.Join(descs, "、"))
		}
		return nil, fmt.Errorf("未检测到已连接的手机，请插上 USB 数据线并在手机上开启 USB 调试")
	default:
		var ss []string
		for _, d := range ready {
			ss = append(ss, d.Serial)
		}
		return nil, fmt.Errorf("检测到 %d 台在线设备（%s），请在该手机的配置文件中填写 serial 指定目标设备", len(ready), strings.Join(ss, "、"))
	}
}

// InstallOptions 安装参数（对应手机配置中的选项）。
type InstallOptions struct {
	Reinstall      bool // -r 覆盖安装保留数据
	AllowDowngrade bool // -d 允许降级
	GrantPerms     bool // -g 授予全部运行时权限
}

// InstallArgs 组装 adb install 参数（不含 adb 路径本身）。
func InstallArgs(serial, apk string, o InstallOptions) []string {
	args := []string{"-s", serial, "install"}
	if o.Reinstall {
		args = append(args, "-r")
	}
	if o.AllowDowngrade {
		args = append(args, "-d")
	}
	if o.GrantPerms {
		args = append(args, "-g")
	}
	return append(args, apk)
}

// InstallResult 一次安装的结果。
type InstallResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"` // 成功摘要或失败原因
}

// extractFailure 从输出中提取 Failure 原因（如 INSTALL_FAILED_UPDATE_INCOMPATIBLE）。
// 没有 Failure 行时返回空串，由调用方决定回退文案。
func extractFailure(out string) string {
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Failure") {
			if i := strings.Index(line, "["); i >= 0 && strings.HasSuffix(line, "]") {
				return line[i+1 : len(line)-1]
			}
			return strings.TrimSpace(strings.TrimPrefix(line, "Failure"))
		}
	}
	return ""
}

// Install 执行 `adb -s <serial> install ...`，逐行回调输出，判定 Success/Failure。
func (a *ADB) Install(ctx context.Context, serial, apk string, o InstallOptions, onLine func(string)) (InstallResult, error) {
	c, cancel := context.WithTimeout(ctx, installTimeout)
	defer cancel()
	args := InstallArgs(serial, apk, o)
	out, err := a.r.RunStreaming(c, a.exe, args, onLine)
	if strings.Contains(out, "Success") {
		return InstallResult{Success: true, Message: "安装成功"}, nil
	}
	if msg := extractFailure(out); msg != "" {
		return InstallResult{Success: false, Message: msg}, nil
	}
	if err != nil {
		msg := fmt.Sprintf("执行 adb install 失败: %v", err)
		// RunStreaming 失败时已把 stderr 并入 out，一并展示，避免只看到 exit status 1
		if extra := strings.TrimSpace(out); extra != "" {
			msg += "\n" + extra
		}
		return InstallResult{Success: false, Message: msg}, nil
	}
	return InstallResult{Success: false, Message: "安装失败（adb 无失败原因输出，可在日志台切到「调试」查看原始输出）"}, nil
}

// LaunchApp 通过 monkey 启动应用主界面（安装成功后的可选动作）。
func (a *ADB) LaunchApp(ctx context.Context, serial, pkg string) (string, error) {
	c, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	out, errOut, err := a.r.Run(c, a.exe, []string{"-s", serial, "shell", "monkey", "-p", pkg, "-c", "android.intent.category.LAUNCHER", "1"})
	if err != nil {
		return "", fmt.Errorf("启动应用失败: %v %s", err, strings.TrimSpace(errOut))
	}
	if strings.Contains(out, "Events injected: 0") || strings.Contains(errOut, "No activities found") {
		return "", fmt.Errorf("包 %s 未找到可启动的 Activity，请检查包名是否正确", pkg)
	}
	return "已向手机发送启动指令", nil
}
