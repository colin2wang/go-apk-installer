// Package config 负责configs 目录下的全部 YAML 配置：
//   - app.yml        全局配置（Android SDK 目录、主题）
//   - <手机名>.yml   一台手机一个配置文件（一机一 yml）
//
// 写出采用内置渲染器（带中文注释，方便用户手工编辑），
// 读取采用 yaml.v3 解析；两者的一致性由往返单测保证。
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// AppConfig 全局配置（configs/app.yml）。
type AppConfig struct {
	SDKDir string `json:"sdkDir" yaml:"sdk_dir"` // Android SDK 目录
	Theme  string `json:"theme" yaml:"theme"`    // 界面主题 dark / light
}

// Phone 单台手机配置（configs/<手机名>.yml）。
type Phone struct {
	Name           string `json:"name" yaml:"name"`                          // 手机名称（=文件名）
	Serial         string `json:"serial" yaml:"serial"`                      // 设备序列号，留空自动匹配唯一在线设备
	APKPath        string `json:"apkPath" yaml:"apk_path"`                   // APK 安装包路径
	Reinstall      bool   `json:"reinstall" yaml:"reinstall"`                // -r 覆盖安装
	AllowDowngrade bool   `json:"allowDowngrade" yaml:"allow_downgrade"`     // -d 允许降级
	GrantPerms     bool   `json:"grantPermissions" yaml:"grant_permissions"` // -g 授予运行时权限
	AutoLaunch     bool   `json:"autoLaunch" yaml:"auto_launch"`             // 安装后自动启动
	LaunchPackage  string `json:"launchPackage" yaml:"launch_package"`       // 启动的包名
	Notes          string `json:"notes" yaml:"notes"`                        // 备注
}

// DefaultApp 默认全局配置。
func DefaultApp() AppConfig { return AppConfig{SDKDir: "", Theme: "dark"} }

// DefaultPhone 首次运行生成的示例手机配置。
func DefaultPhone() Phone {
	return Phone{
		Name:           "示例手机",
		Reinstall:      true,
		AllowDowngrade: false,
		GrantPerms:     false,
		AutoLaunch:     false,
		LaunchPackage:  "",
		Notes:          "示例配置：填好 serial 与 apk_path 即可一键安装；serial 留空时，只连接一台手机会自动匹配。",
	}
}

// SanitizeName 将手机名称清理为合法文件名（剔除 Windows 非法字符）。
func SanitizeName(s string) string {
	s = strings.TrimSpace(s)
	var b strings.Builder
	for _, r := range s {
		switch r {
		case '\\', '/', ':', '*', '?', '"', '<', '>', '|':
			b.WriteRune('-')
		default:
			if r >= 32 {
				b.WriteRune(r)
			}
		}
	}
	out := strings.TrimSpace(b.String())
	if out == "" {
		out = "未命名手机"
	}
	if len([]rune(out)) > 60 {
		out = string([]rune(out)[:60])
	}
	return out
}

// yv 将字符串渲染为安全的 YAML 标量（需要时加双引号）。
func yv(s string) string {
	if s == "" {
		return `""`
	}
	if strings.ContainsAny(s, ":#{}[],&*?|>'\"%@`\\\n\t\r") ||
		strings.TrimSpace(s) != s ||
		strings.HasPrefix(s, "-") {
		return strconv.Quote(s)
	}
	switch strings.ToLower(s) {
	case "true", "false", "null", "yes", "no", "on", "off", "~":
		return strconv.Quote(s)
	}
	return s
}

// RenderApp 渲染 app.yml（带中文注释）。
func RenderApp(c AppConfig) string {
	if c.Theme == "" {
		c.Theme = "dark"
	}
	var b strings.Builder
	b.WriteString("# ============================================================\n")
	b.WriteString("# APK 装机助手 · 全局配置（app.yml）\n")
	b.WriteString("# -----------------------------------------------------------\n")
	b.WriteString("# sdk_dir : Android SDK 目录（内含 platform-tools/adb.exe）\n")
	b.WriteString("#           建议直接在软件界面顶部「选择 SDK 目录」自动填写\n")
	b.WriteString("# theme   : 界面主题 dark / light\n")
	b.WriteString("# ============================================================\n\n")
	b.WriteString("# Android SDK 目录\n")
	b.WriteString("sdk_dir: " + yv(c.SDKDir) + "\n\n")
	b.WriteString("# 界面主题\n")
	b.WriteString("theme: " + yv(c.Theme) + "\n")
	return b.String()
}

// RenderPhone 渲染单台手机的 yml（带中文注释）。
func RenderPhone(p Phone) string {
	var b strings.Builder
	b.WriteString("# ============================================================\n")
	b.WriteString("# 手机配置文件（一台手机一个 yml，文件名 = 手机名称）\n")
	b.WriteString("# 本文件由「APK 装机助手」维护，也可用文本编辑器手工修改\n")
	b.WriteString("# -----------------------------------------------------------\n")
	b.WriteString("# 安装动作 = adb -s <serial> install [-r] [-d] [-g] <apk_path>\n")
	b.WriteString("# ============================================================\n\n")
	b.WriteString("# 手机名称（须与文件名保持一致）\n")
	b.WriteString("name: " + yv(p.Name) + "\n\n")
	b.WriteString("# 设备序列号（即 adb devices 列表第一列，如 emulator-5554 / 1A2B3C4D）\n")
	b.WriteString("# 留空 = 自动匹配：仅连接一台手机时自动使用该手机\n")
	b.WriteString("serial: " + yv(p.Serial) + "\n\n")
	b.WriteString("# APK 安装包绝对路径（.apk 文件）\n")
	b.WriteString("apk_path: " + yv(p.APKPath) + "\n\n")
	b.WriteString("# 覆盖安装：已安装时保留应用数据升级（adb install -r）\n")
	b.WriteString("reinstall: " + strconv.FormatBool(p.Reinstall) + "\n\n")
	b.WriteString("# 允许降级安装：安装低于已装版本号的 APK（adb install -d）\n")
	b.WriteString("allow_downgrade: " + strconv.FormatBool(p.AllowDowngrade) + "\n\n")
	b.WriteString("# 安装时授予全部运行时权限（adb install -g）\n")
	b.WriteString("grant_permissions: " + strconv.FormatBool(p.GrantPerms) + "\n\n")
	b.WriteString("# 安装成功后自动启动应用\n")
	b.WriteString("auto_launch: " + strconv.FormatBool(p.AutoLaunch) + "\n\n")
	b.WriteString("# 自动启动的包名（auto_launch 为 true 时必填，如 com.example.app）\n")
	b.WriteString("launch_package: " + yv(p.LaunchPackage) + "\n\n")
	b.WriteString("# 备注\n")
	b.WriteString("notes: " + yv(p.Notes) + "\n")
	return b.String()
}

// ParseApp 解析 app.yml。
func ParseApp(data []byte) (AppConfig, error) {
	var c AppConfig
	if err := yaml.Unmarshal(data, &c); err != nil {
		return c, fmt.Errorf("解析 app.yml 失败: %w", err)
	}
	if c.Theme == "" {
		c.Theme = "dark"
	}
	return c, nil
}

// ParsePhone 解析单台手机的 yml。
func ParsePhone(data []byte) (Phone, error) {
	var p Phone
	if err := yaml.Unmarshal(data, &p); err != nil {
		return p, fmt.Errorf("解析手机配置失败: %w", err)
	}
	p.Name = strings.TrimSpace(p.Name)
	if p.Name == "" {
		return p, fmt.Errorf("手机配置缺少 name 字段")
	}
	return p, nil
}

// Manager 管理一个 configs 目录。
type Manager struct{ dir string }

func NewManager(dir string) *Manager { return &Manager{dir: dir} }

func (m *Manager) Dir() string { return m.dir }

func (m *Manager) appPath() string { return filepath.Join(m.dir, "app.yml") }
func (m *Manager) phonePath(n string) string {
	return filepath.Join(m.dir, SanitizeName(n)+".yml")
}

// Ensure 创建 configs 目录与缺失的默认配置（app.yml、首台示例手机）。
func (m *Manager) Ensure() error {
	if err := os.MkdirAll(m.dir, 0o755); err != nil {
		return fmt.Errorf("创建 configs 目录失败: %w", err)
	}
	if _, err := os.Stat(m.appPath()); os.IsNotExist(err) {
		if err := m.SaveApp(DefaultApp()); err != nil {
			return err
		}
	}
	phones, _, err := m.ListPhones()
	if err != nil {
		return err
	}
	if len(phones) == 0 {
		return m.SavePhone("", DefaultPhone())
	}
	return nil
}

// LoadApp 读取全局配置；文件缺失时返回默认值。
func (m *Manager) LoadApp() (AppConfig, error) {
	data, err := os.ReadFile(m.appPath())
	if os.IsNotExist(err) {
		return DefaultApp(), nil
	}
	if err != nil {
		return DefaultApp(), err
	}
	return ParseApp(data)
}

// SaveApp 写回全局配置。
func (m *Manager) SaveApp(c AppConfig) error {
	return os.WriteFile(m.appPath(), []byte(RenderApp(c)), 0o644)
}

// ListPhones 列出全部手机配置（按名称排序，跳过 app.yml）。
// 个别文件解析失败不阻断整体，原因汇总在返回的 error（可为空）。
func (m *Manager) ListPhones() ([]Phone, []string, error) {
	entries, err := os.ReadDir(m.dir)
	if os.IsNotExist(err) {
		return nil, nil, nil
	}
	if err != nil {
		return nil, nil, err
	}
	var phones []Phone
	var warns []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if name == "app.yml" || name == "app.yaml" {
			continue
		}
		if ext := strings.ToLower(filepath.Ext(name)); ext != ".yml" && ext != ".yaml" {
			continue
		}
		data, err := os.ReadFile(filepath.Join(m.dir, name))
		if err != nil {
			warns = append(warns, "读取 "+name+" 失败: "+err.Error())
			continue
		}
		p, err := ParsePhone(data)
		if err != nil {
			warns = append(warns, name+" 被跳过: "+err.Error())
			continue
		}
		phones = append(phones, p)
	}
	sort.Slice(phones, func(i, j int) bool { return phones[i].Name < phones[j].Name })
	return phones, warns, nil
}

// LoadPhone 按名称读取手机配置。
func (m *Manager) LoadPhone(name string) (Phone, error) {
	data, err := os.ReadFile(m.phonePath(name))
	if err != nil {
		return Phone{}, fmt.Errorf("手机「%s」的配置文件不存在或不可读", name)
	}
	return ParsePhone(data)
}

// SavePhone 写入手机配置；oldName 非空且与现名称不同时执行重命名（删除旧文件）。
func (m *Manager) SavePhone(oldName string, p Phone) error {
	p.Name = strings.TrimSpace(p.Name)
	if p.Name == "" {
		return fmt.Errorf("手机名称不能为空")
	}
	target := m.phonePath(p.Name)
	if oldName == "" {
		if _, err := os.Stat(target); err == nil {
			return fmt.Errorf("已存在同名手机配置「%s」，请换一个名称", p.Name)
		}
	}
	if err := os.WriteFile(target, []byte(RenderPhone(p)), 0o644); err != nil {
		return fmt.Errorf("写入配置失败: %w", err)
	}
	if oldName != "" {
		oldTarget := m.phonePath(oldName)
		if oldTarget != target {
			if _, err := os.Stat(oldTarget); err == nil {
				if err := os.Remove(oldTarget); err != nil {
					return fmt.Errorf("重命名后删除旧文件失败: %w", err)
				}
			}
		}
	}
	return nil
}

// DeletePhone 删除手机配置（文件不存在视为成功）。
func (m *Manager) DeletePhone(name string) error {
	err := os.Remove(m.phonePath(name))
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
