package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"apk-installer/internal/adb"
	"apk-installer/internal/config"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// SDKStatus SDK 与 adb 的就绪状态。
type SDKStatus struct {
	SDKDir  string `json:"sdkDir"`  // 当前配置的 SDK 目录
	ADBPath string `json:"adbPath"` // 实际使用的 adb 可执行文件
	Version string `json:"version"` // adb version 首行
	Ready   bool   `json:"ready"`   // true = adb 可执行且可用
}

// AppState 启动时的整体状态（配置 + adb 状态）。
type AppState struct {
	Config config.AppConfig `json:"config"`
	SDK    SDKStatus        `json:"sdk"`
}

// CheckResult 单台手机的连接检查结果。
type CheckResult struct {
	OK      bool            `json:"ok"`
	Message string          `json:"message"`
	Device  *adb.DeviceInfo `json:"device"`
}

// InstallOutcome 一次安装的结果。
type InstallOutcome struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	ElapsedMs int64  `json:"elapsedMs"`
}

// LogLine 推送到前端日志台的一行日志。
type LogLine struct {
	Time  string `json:"time"`
	Level string `json:"level"` // info / ok / warn / err / adb
	Text  string `json:"text"`
}

// App Wails 绑定入口。
type App struct {
	ctx context.Context
	cfg *config.Manager
}

func NewApp() *App { return &App{} }

// startup 由 Wails 调用。
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.cfg = config.NewManager(configsDir())
	if err := a.cfg.Ensure(); err != nil {
		a.log("err", "初始化配置目录失败: %v", err)
		return
	}
	a.log("info", "配置目录: %s", a.cfg.Dir())
	if st := a.sdkStatus(); st.Ready {
		a.log("ok", "adb 就绪: %s", st.Version)
	} else {
		a.log("warn", "adb 尚未就绪，请先在顶部选择 Android SDK 目录")
	}
}

// shutdown 由 Wails 调用。
func (a *App) shutdown(ctx context.Context) {}

// log 通过事件推送一行日志到前端。
func (a *App) log(level, format string, args ...any) {
	if a.ctx == nil {
		return
	}
	runtime.EventsEmit(a.ctx, "apk:log", LogLine{
		Time:  time.Now().Format("15:04:05"),
		Level: level,
		Text:  fmt.Sprintf(format, args...),
	})
}

// configsDir 配置目录固定在可执行文件同级 configs/；
// 目录不可写时（如 exe 放在只读位置）回退到工作目录。
func configsDir() string {
	if exe, err := os.Executable(); err == nil {
		d := filepath.Join(filepath.Dir(exe), "configs")
		if canWrite(d) {
			return d
		}
	}
	wd, _ := os.Getwd()
	return filepath.Join(wd, "configs")
}

func canWrite(dir string) bool {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return false
	}
	probe := filepath.Join(dir, ".write_probe")
	if err := os.WriteFile(probe, []byte("1"), 0o644); err != nil {
		return false
	}
	_ = os.Remove(probe)
	return true
}

// isAbsPath 兼容 Windows 盘符的绝对路径判断（Linux 容器下 filepath.IsAbs 不识别 C:\）。
func isAbsPath(p string) bool {
	if filepath.IsAbs(p) {
		return true
	}
	if len(p) >= 3 && p[1] == ':' && (p[2] == '\\' || p[2] == '/') {
		c := p[0]
		return (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z')
	}
	return false
}

// sdkStatus 探测当前 adb 状态（不产生 UI 副作用）。
func (a *App) sdkStatus() SDKStatus {
	app, err := a.cfg.LoadApp()
	if err != nil {
		return SDKStatus{}
	}
	st := SDKStatus{SDKDir: app.SDKDir}
	exe, err := adb.LocateADB(app.SDKDir)
	if err != nil {
		return st
	}
	st.ADBPath = exe
	v, err := adb.New(exe, nil).Version(a.ctx)
	if err != nil {
		return st
	}
	st.Version = v
	st.Ready = true
	return st
}

// resolveADB 取一个可用的 adb 实例；不可用返回带指引的错误。
func (a *App) resolveADB() (*adb.ADB, error) {
	app, err := a.cfg.LoadApp()
	if err != nil {
		return nil, err
	}
	exe, err := adb.LocateADB(app.SDKDir)
	if err != nil {
		return nil, err
	}
	return adb.New(exe, nil), nil
}

// ---------- 绑定方法（前端可调用） ----------

// GetState 返回全局配置 + adb 就绪状态。
func (a *App) GetState() (AppState, error) {
	cfg, err := a.cfg.LoadApp()
	if err != nil {
		cfg = config.DefaultApp()
	}
	return AppState{Config: cfg, SDK: a.sdkStatus()}, nil
}

// SelectSDKDir 弹出目录选择框，选中后校验 adb 并保存到 app.yml。
// 用户取消时返回当前状态。
func (a *App) SelectSDKDir() (SDKStatus, error) {
	dir, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "选择 Android SDK 目录（内含 platform-tools/adb.exe）",
	})
	if err != nil {
		return a.sdkStatus(), err
	}
	if dir == "" {
		return a.sdkStatus(), nil
	}
	return a.SetSDKDir(dir)
}

// SetSDKDir 设置并保存 SDK 目录（也接受用户手工输入后提交）。
func (a *App) SetSDKDir(dir string) (SDKStatus, error) {
	dir = strings.TrimSpace(dir)
	cfg, err := a.cfg.LoadApp()
	if err != nil {
		cfg = config.DefaultApp()
	}
	cfg.SDKDir = dir
	exe, locErr := adb.LocateADB(dir)
	if locErr != nil {
		// 目录无效也保存（用户可能先建目录后装 SDK），但不置 Ready
		if err := a.cfg.SaveApp(cfg); err != nil {
			return SDKStatus{}, err
		}
		a.log("warn", "%v", locErr)
		return SDKStatus{SDKDir: dir}, fmt.Errorf("%w", locErr)
	}
	if err := a.cfg.SaveApp(cfg); err != nil {
		return SDKStatus{}, err
	}
	v, verr := adb.New(exe, nil).Version(a.ctx)
	st := SDKStatus{SDKDir: dir, ADBPath: exe}
	if verr != nil {
		a.log("err", "%v", verr)
		return st, verr
	}
	st.Version = v
	st.Ready = true
	a.log("ok", "Android SDK 已设置: %s（%s）", dir, v)
	return st, nil
}

// SetTheme 保存界面主题（dark / light）到 app.yml。
func (a *App) SetTheme(theme string) error {
	if theme != "dark" && theme != "light" {
		theme = "dark"
	}
	cfg, err := a.cfg.LoadApp()
	if err != nil {
		cfg = config.DefaultApp()
	}
	cfg.Theme = theme
	return a.cfg.SaveApp(cfg)
}

// SelectAPKFile 弹出文件选择框（*.apk），取消返回空串。
func (a *App) SelectAPKFile() (string, error) {
	p, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "选择 APK 安装包",
		Filters: []runtime.FileFilter{
			{DisplayName: "APK 安装包 (*.apk)", Pattern: "*.apk"},
		},
	})
	if err != nil {
		return "", err
	}
	return p, nil
}

// ListPhones 读取 configs 下全部手机配置（一机一 yml）。
func (a *App) ListPhones() ([]config.Phone, error) {
	phones, warns, err := a.cfg.ListPhones()
	for _, w := range warns {
		a.log("warn", "%s", w)
	}
	if err != nil {
		return nil, err
	}
	return phones, nil
}

// SavePhone 保存手机配置；oldName 与名称不同时执行重命名。
func (a *App) SavePhone(oldName string, p config.Phone) error {
	if err := a.cfg.SavePhone(oldName, p); err != nil {
		a.log("err", "%v", err)
		return err
	}
	if oldName != "" && oldName != p.Name {
		a.log("ok", "手机「%s」已重命名为「%s」", oldName, p.Name)
	} else {
		a.log("ok", "手机「%s」配置已保存（configs/%s.yml）", p.Name, config.SanitizeName(p.Name))
	}
	return nil
}

// DeletePhone 删除手机配置文件。
func (a *App) DeletePhone(name string) error {
	if err := a.cfg.DeletePhone(name); err != nil {
		a.log("err", "删除「%s」失败: %v", name, err)
		return err
	}
	a.log("info", "已删除手机配置「%s」", name)
	return nil
}

// ListDevices 列出当前 adb 可见的设备。
func (a *App) ListDevices() ([]adb.DeviceInfo, error) {
	a2, err := a.resolveADB()
	if err != nil {
		a.log("warn", "%v", err)
		return nil, err
	}
	devices, err := a2.Devices(a.ctx)
	if err != nil {
		a.log("err", "%v", err)
		return nil, err
	}
	a.log("info", "adb devices：%d 台设备", len(devices))
	if devices == nil {
		devices = []adb.DeviceInfo{} // nil 切片会被序列化为 null，前端遍历会报错
	}
	return devices, nil
}

// CheckPhone 检查某台手机的连接状态（安装前检查的第一步，也可单独执行）。
func (a *App) CheckPhone(name string) (CheckResult, error) {
	p, err := a.cfg.LoadPhone(name)
	if err != nil {
		return CheckResult{OK: false, Message: err.Error()}, nil
	}
	a2, err := a.resolveADB()
	if err != nil {
		a.log("warn", "「%s」连接检查: %v", p.Name, err)
		return CheckResult{OK: false, Message: err.Error()}, nil
	}
	devices, err := a2.DevicesVerbose(a.ctx, func(line string) {
		a.log("adb", "%s", line)
	})
	if err != nil {
		return CheckResult{OK: false, Message: err.Error()}, nil
	}
	dev, err := adb.ResolveDevice(devices, p.Serial)
	if err != nil {
		a.log("warn", "「%s」连接检查未通过: %v", p.Name, err)
		return CheckResult{OK: false, Message: err.Error()}, nil
	}
	a.log("ok", "「%s」连接检查通过：%s", p.Name, dev.Describe())
	return CheckResult{OK: true, Message: dev.Describe(), Device: dev}, nil
}

// InstallAPK 安装指定手机配置的 APK。
// 严格执行两道前置检查（不通过绝不执行安装）：
//  1. adb 可用（SDK 目录已配置且 adb 可执行）；
//  2. adb 能连到目标手机（设备在线且已授权）。
func (a *App) InstallAPK(name string) (InstallOutcome, error) {
	start := time.Now()
	p, err := a.cfg.LoadPhone(name)
	if err != nil {
		return InstallOutcome{}, err
	}
	// ① 前置：apk 路径
	apk := strings.TrimSpace(p.APKPath)
	if apk == "" {
		a.log("err", "「%s」未配置 APK 路径，请先选择 .apk 安装包", p.Name)
		return InstallOutcome{}, fmt.Errorf("未配置 APK 路径，请先选择 .apk 安装包")
	}
	if !isAbsPath(apk) {
		a.log("err", "APK 路径必须是绝对路径: %s", apk)
		return InstallOutcome{}, fmt.Errorf("APK 路径必须是绝对路径: %s", apk)
	}
	if st, err := os.Stat(apk); err != nil || st.IsDir() {
		a.log("err", "APK 文件不存在: %s", apk)
		return InstallOutcome{}, fmt.Errorf("APK 文件不存在: %s", apk)
	}
	if ext := filepath.Ext(apk); !strings.EqualFold(ext, ".apk") {
		a.log("err", "仅支持 .apk 安装包，当前: %s", apk)
		return InstallOutcome{}, fmt.Errorf("仅支持 .apk 安装包，当前: %s", apk)
	}
	// ② 前置：adb 可用
	a2, err := a.resolveADB()
	if err != nil {
		a.log("err", "安装中止：%v", err)
		return InstallOutcome{}, err
	}
	a.log("info", "▶ 开始安装「%s」", p.Name)
	a.log("info", "① adb 就绪：%s", a2.Exe())
	// ③ 前置：连接检查
	devices, err := a2.DevicesVerbose(a.ctx, func(line string) {
		a.log("adb", "%s", line)
	})
	if err != nil {
		a.log("err", "安装中止：%v", err)
		return InstallOutcome{}, err
	}
	dev, err := adb.ResolveDevice(devices, p.Serial)
	if err != nil {
		a.log("err", "② 连接检查未通过，安装中止：%v", err)
		return InstallOutcome{}, fmt.Errorf("连接检查未通过：%w", err)
	}
	a.log("ok", "② 连接检查通过：%s", dev.Describe())
	// ④ 安装（流式输出日志）
	a.log("info", "③ 正在安装 %s ...", filepath.Base(apk))
	opts := adb.InstallOptions{
		Reinstall:      p.Reinstall,
		AllowDowngrade: p.AllowDowngrade,
		GrantPerms:     p.GrantPerms,
	}
	res, err := a2.Install(a.ctx, dev.Serial, apk, opts, func(line string) {
		if strings.TrimSpace(line) != "" {
			a.log("adb", "%s", line)
		}
	})
	elapsed := time.Since(start)
	if err != nil {
		a.log("err", "✘ 安装异常：%v", err)
		return InstallOutcome{Success: false, Message: err.Error(), ElapsedMs: elapsed.Milliseconds()}, nil
	}
	if !res.Success {
		a.log("err", "✘ 安装失败：%s（耗时 %.1f 秒）", res.Message, elapsed.Seconds())
		return InstallOutcome{Success: false, Message: res.Message, ElapsedMs: elapsed.Milliseconds()}, nil
	}
	a.log("ok", "✔ 安装成功（耗时 %.1f 秒）", elapsed.Seconds())
	// ⑤ 可选：自动启动
	if p.AutoLaunch && strings.TrimSpace(p.LaunchPackage) != "" {
		if msg, err := a2.LaunchApp(a.ctx, dev.Serial, strings.TrimSpace(p.LaunchPackage)); err != nil {
			a.log("warn", "自动启动失败：%v", err)
		} else {
			a.log("ok", "自动启动：%s（%s）", p.LaunchPackage, msg)
		}
	}
	return InstallOutcome{Success: true, Message: res.Message, ElapsedMs: elapsed.Milliseconds()}, nil
}

// OpenConfigsDir 在系统资源管理器中打开 configs 目录。
func (a *App) OpenConfigsDir() error {
	runtime.BrowserOpenURL(a.ctx, a.cfg.Dir())
	return nil
}
