package adb

import (
	"os/exec"
	"syscall"
)

// hideConsole prevents child processes (e.g. adb.exe) from flashing a
// console window when spawned from the GUI app on Windows.
func hideConsole(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: 0x08000000, // CREATE_NO_WINDOW
	}
}
