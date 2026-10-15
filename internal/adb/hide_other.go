//go:build !windows

package adb

import "os/exec"

func hideConsole(cmd *exec.Cmd) {}
