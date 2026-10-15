package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestSanitizeName(t *testing.T) {
	cases := map[string]string{
		"  办公测试机  ":           "办公测试机",
		`a/b\c:d*e?f"g<h>i|j`: "a-b-c-d-e-f-g-h-i-j",
		"":                    "未命名手机",
		"   ":                 "未命名手机",
	}
	for in, want := range cases {
		if got := SanitizeName(in); got != want {
			t.Errorf("SanitizeName(%q)=%q want %q", in, got, want)
		}
	}
}

func TestYVQuoting(t *testing.T) {
	// Windows 路径必须加引号，且 YAML 可解析回原值
	p := `D:\apk\demo release.apk`
	rendered := "apk_path: " + yv(p) + "\n"
	var got struct {
		Path string `yaml:"apk_path"`
	}
	if err := yaml.Unmarshal([]byte(rendered), &got); err != nil {
		t.Fatal(err)
	}
	if got.Path != p {
		t.Fatalf("往返不一致: %q != %q", got.Path, p)
	}
	if strings.Contains(rendered, `apk_path: D:`) {
		t.Fatal("含冒号路径未加引号")
	}
}

func TestPhoneRoundtrip(t *testing.T) {
	p := Phone{
		Name:           "办公测试机",
		Serial:         "1A2B3C4D",
		APKPath:        `D:\apk\demo release v2.apk`,
		Reinstall:      true,
		AllowDowngrade: false,
		GrantPerms:     true,
		AutoLaunch:     true,
		LaunchPackage:  "com.example.app",
		Notes:          "第一行\n第二行: 带冒号 # 带井号",
	}
	data := RenderPhone(p)
	got, err := ParsePhone([]byte(data))
	if err != nil {
		t.Fatalf("解析失败: %v\n%s", err, data)
	}
	if got != p {
		t.Fatalf("往返不一致:\ngot:  %+v\nwant: %+v\n%s", got, p, data)
	}
	// 二次渲染必须与首次完全一致（注释稳定）
	if RenderPhone(got) != data {
		t.Fatalf("二次渲染与首次不一致:\n%s\n----\n%s", data, RenderPhone(got))
	}
}

func TestAppRoundtrip(t *testing.T) {
	c := AppConfig{SDKDir: `C:\Users\dev\AppData\Local\Android\Sdk`, Theme: "light"}
	data := RenderApp(c)
	got, err := ParseApp([]byte(data))
	if err != nil {
		t.Fatal(err)
	}
	if got != c {
		t.Fatalf("往返不一致: %+v != %+v", got, c)
	}
}

func TestManagerEnsureAndCRUD(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "configs")
	m := NewManager(dir)

	// Ensure：生成 app.yml + 示例手机
	if err := m.Ensure(); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dir, "app.yml")); err != nil {
		t.Fatal("app.yml 未生成")
	}
	phones, warns, err := m.ListPhones()
	if err != nil || len(warns) > 0 {
		t.Fatalf("ListPhones 异常: %v %v", err, warns)
	}
	if len(phones) != 1 || phones[0].Name != "示例手机" {
		t.Fatalf("示例手机未生成: %+v", phones)
	}
	// 再次 Ensure 幂等
	if err := m.Ensure(); err != nil {
		t.Fatal(err)
	}
	phones, _, _ = m.ListPhones()
	if len(phones) != 1 {
		t.Fatalf("Ensure 应幂等，实际 %d 台", len(phones))
	}

	// 新增
	p := DefaultPhone()
	p.Name = "车间平板"
	p.Serial = "emulator-5554"
	if err := m.SavePhone("", p); err != nil {
		t.Fatal(err)
	}
	// 同名新增应报错
	if err := m.SavePhone("", p); err == nil {
		t.Fatal("同名新增应报错")
	}
	// 重命名（旧文件应删除）
	p.Name = "车间平板02"
	if err := m.SavePhone("车间平板", p); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dir, "车间平板.yml")); !os.IsNotExist(err) {
		t.Fatal("旧文件未删除")
	}
	if _, err := os.Stat(filepath.Join(dir, "车间平板02.yml")); err != nil {
		t.Fatal("新文件未生成")
	}
	// 删除
	if err := m.DeletePhone("车间平板02"); err != nil {
		t.Fatal(err)
	}
	if err := m.DeletePhone("不存在"); err != nil {
		t.Fatal("删除不存在配置应幂等成功")
	}
}

func TestListPhonesSkipsBroken(t *testing.T) {
	dir := t.TempDir()
	m := NewManager(dir)
	os.WriteFile(filepath.Join(dir, "app.yml"), []byte(RenderApp(DefaultApp())), 0o644)
	os.WriteFile(filepath.Join(dir, "好机.yml"), []byte(RenderPhone(Phone{Name: "好机"})), 0o644)
	os.WriteFile(filepath.Join(dir, "坏机.yml"), []byte("name: [unclosed"), 0o644)
	os.WriteFile(filepath.Join(dir, "readme.txt"), []byte("skip me"), 0o644)

	phones, warns, err := m.ListPhones()
	if err != nil {
		t.Fatal(err)
	}
	if len(phones) != 1 || phones[0].Name != "好机" {
		t.Fatalf("应只解析出 1 台: %+v", phones)
	}
	if len(warns) != 1 || !strings.Contains(warns[0], "坏机.yml") {
		t.Fatalf("坏文件应有警告: %v", warns)
	}
}

func TestBrokenPhoneYAMLError(t *testing.T) {
	if _, err := ParsePhone([]byte("name: [unclosed")); err == nil {
		t.Fatal("应报错")
	}
	if _, err := ParsePhone([]byte("serial: xyz")); err == nil {
		t.Fatal("缺少 name 应报错")
	}
}
