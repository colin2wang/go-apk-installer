package adb

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

const devicesSample = `* daemon not running; starting now at tcp:5037
* daemon started successfully
List of devices attached
emulator-5554          device product:sdk_gphone64_x86_64 model:Pixel_6 device:x86_64 transport_id:1
1A2B3C4D5E6F           unauthorized usb:1-1 product:raven model:Pixel_6_Pro transport_id:3
BADDEVICE              offline
`

type fakeRunner struct {
	versionOut string
	devicesOut string
	devicesErr error
	installOut []string // 逐行输出
	installErr error
	shellOut   string
	shellErr   error
	lastArgs   []string
}

func (f *fakeRunner) Run(ctx context.Context, exe string, args []string) (string, string, error) {
	f.lastArgs = args
	switch args[0] {
	case "version":
		return f.versionOut, "", nil
	case "devices":
		return f.devicesOut, "", f.devicesErr
	case "shell":
		return f.shellOut, "", f.shellErr
	}
	return "", "", nil
}

func (f *fakeRunner) RunStreaming(ctx context.Context, exe string, args []string, onLine func(string)) (string, error) {
	f.lastArgs = args
	var all strings.Builder
	for _, l := range f.installOut {
		all.WriteString(l + "\n")
		if onLine != nil {
			onLine(l)
		}
	}
	return all.String(), f.installErr
}

func TestParseDevices(t *testing.T) {
	got := parseDevices(devicesSample)
	if len(got) != 3 {
		t.Fatalf("期望解析出 3 台设备，实际 %d: %+v", len(got), got)
	}
	if got[0].Serial != "emulator-5554" || got[0].State != "device" {
		t.Errorf("第一台解析错误: %+v", got[0])
	}
	if got[0].Model != "Pixel 6" {
		t.Errorf("model 下划线应转空格，实际 %q", got[0].Model)
	}
	if got[1].State != "unauthorized" {
		t.Errorf("第二台状态错误: %+v", got[1])
	}
	if got[2].State != "offline" {
		t.Errorf("第三台状态错误: %+v", got[2])
	}
}

func TestLocateADB(t *testing.T) {
	dir := t.TempDir()
	if _, err := LocateADB(dir); err == nil {
		t.Fatal("空 SDK 目录且 PATH 无 adb，应报错")
	}
	pt := filepath.Join(dir, "platform-tools")
	os.MkdirAll(pt, 0o755)
	exe := filepath.Join(pt, "adb")
	os.WriteFile(exe, []byte("#!/bin/sh\n"), 0o755)
	got, err := LocateADB(dir)
	if err != nil || got != exe {
		t.Fatalf("应定位到 platform-tools/adb，got=%q err=%v", got, err)
	}
	// PATH 回退
	empty := t.TempDir()
	pathDir := t.TempDir()
	os.WriteFile(filepath.Join(pathDir, "adb"), []byte("#!/bin/sh\n"), 0o755)
	t.Setenv("PATH", pathDir)
	got, err = LocateADB(empty)
	if err != nil || got != filepath.Join(pathDir, "adb") {
		t.Fatalf("应回退 PATH 中的 adb，got=%q err=%v", got, err)
	}
}

func TestResolveDevice(t *testing.T) {
	devices := []DeviceInfo{
		{Serial: "A1", State: "device", Model: "Pixel 6"},
		{Serial: "B2", State: "unauthorized"},
		{Serial: "C3", State: "device"},
	}
	// 指定 serial，在线可用
	d, err := ResolveDevice(devices, "A1")
	if err != nil || d.Serial != "A1" {
		t.Fatalf("A1 应可用: %v", err)
	}
	// 指定 serial，未授权
	_, err = ResolveDevice(devices, "B2")
	if err == nil || !strings.Contains(err.Error(), "未授权") {
		t.Fatalf("B2 未授权应有明确提示: %v", err)
	}
	// 指定 serial，未连接
	_, err = ResolveDevice(devices, "ZZ")
	if err == nil || !strings.Contains(err.Error(), "未连接") {
		t.Fatalf("ZZ 应报未连接: %v", err)
	}
	// 空 serial + 多台可用 → 歧义
	_, err = ResolveDevice(devices, "")
	if err == nil || !strings.Contains(err.Error(), "A1、C3") {
		t.Fatalf("多台设备应提示歧义并列出序列号: %v", err)
	}
	// 空 serial + 0 台 → 未连接
	_, err = ResolveDevice(nil, "")
	if err == nil || !strings.Contains(err.Error(), "未检测到") {
		t.Fatalf("无设备应报未连接: %v", err)
	}
	// 空 serial + 恰一台可用 → 自动选中
	single := []DeviceInfo{{Serial: "ONLY", State: "device"}}
	d, err = ResolveDevice(single, "")
	if err != nil || d.Serial != "ONLY" {
		t.Fatalf("唯一可用设备应自动选中: %v", err)
	}
	// 只有不可用设备
	offline := []DeviceInfo{{Serial: "B2", State: "unauthorized"}}
	_, err = ResolveDevice(offline, "")
	if err == nil || !strings.Contains(err.Error(), "USB 调试授权") {
		t.Fatalf("仅离线设备应提示检查授权: %v", err)
	}
}

func TestInstallArgs(t *testing.T) {
	got := InstallArgs("A1", "D:\\apk\\demo.apk", InstallOptions{Reinstall: true, GrantPerms: true})
	want := []string{"-s", "A1", "install", "-r", "-g", "D:\\apk\\demo.apk"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("参数不一致: got=%v want=%v", got, want)
	}
	got = InstallArgs("A1", "x.apk", InstallOptions{AllowDowngrade: true})
	want = []string{"-s", "A1", "install", "-d", "x.apk"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("参数不一致: got=%v want=%v", got, want)
	}
}

func TestInstallSuccess(t *testing.T) {
	fr := &fakeRunner{installOut: []string{"Performing Streamed Install", "Success"}}
	a := New("adb", fr)
	var lines []string
	res, err := a.Install(context.Background(), "A1", "x.apk", InstallOptions{Reinstall: true}, func(l string) { lines = append(lines, l) })
	if err != nil {
		t.Fatal(err)
	}
	if !res.Success || res.Message != "安装成功" {
		t.Fatalf("应成功: %+v", res)
	}
	if len(lines) != 2 || lines[0] != "Performing Streamed Install" {
		t.Fatalf("流式回调不正确: %v", lines)
	}
	if !reflect.DeepEqual(fr.lastArgs, []string{"-s", "A1", "install", "-r", "x.apk"}) {
		t.Fatalf("实际执行参数错误: %v", fr.lastArgs)
	}
}

func TestInstallFailure(t *testing.T) {
	fr := &fakeRunner{installOut: []string{"Performing Streamed Install", "Failure [INSTALL_FAILED_UPDATE_INCOMPATIBLE]"}}
	a := New("adb", fr)
	res, err := a.Install(context.Background(), "A1", "x.apk", InstallOptions{}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if res.Success {
		t.Fatal("应失败")
	}
	if res.Message != "INSTALL_FAILED_UPDATE_INCOMPATIBLE" {
		t.Fatalf("应提取失败原因: %q", res.Message)
	}
}

func TestVersionAndLaunch(t *testing.T) {
	fr := &fakeRunner{
		versionOut: "Android Debug Bridge version 1.0.41\nVersion 34.0.5-10900879\n",
		shellOut:   "Events injected: 1\n",
	}
	a := New("adb", fr)
	v, err := a.Version(context.Background())
	if err != nil || !strings.Contains(v, "1.0.41") {
		t.Fatalf("版本解析错误: %q %v", v, err)
	}
	msg, err := a.LaunchApp(context.Background(), "A1", "com.example.app")
	if err != nil || msg == "" {
		t.Fatalf("启动应成功: %v %v", msg, err)
	}
	if !reflect.DeepEqual(fr.lastArgs, []string{"-s", "A1", "shell", "monkey", "-p", "com.example.app", "-c", "android.intent.category.LAUNCHER", "1"}) {
		t.Fatalf("monkey 参数错误: %v", fr.lastArgs)
	}
}
