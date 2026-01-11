package sysagent

import (
	"fmt"
	"syscall"

	"golang.org/x/sys/windows"
)

var (
	kernel32             = syscall.NewLazyDLL("kernel32.dll")
	user32               = syscall.NewLazyDLL("user32.dll")
	procGetConsoleWindow = kernel32.NewProc("GetConsoleWindow")
	procShowWindow       = user32.NewProc("ShowWindow")
)

// Native API
const (
	EWX_LOGOFF   = 0x00000000
	EWX_SHUTDOWN = 0x00000001
	EWX_REBOOT   = 0x00000002
	EWX_FORCE    = 0x00000004
)

func ControlPower(action string) error {
	var flags uint32

	switch action {
	case "shutdown":
		flags = EWX_SHUTDOWN | EWX_FORCE
	case "reboot":
		flags = EWX_REBOOT | EWX_FORCE
	case "logoff":
		flags = EWX_LOGOFF | EWX_FORCE
	default:
		return fmt.Errorf("unknown action: %s", action)
	}

	if err := enablePrivilege("SeShutdownPrivilege"); err != nil {
		return fmt.Errorf("failed to acquire privilege: %v", err)
	}

	return exitWindowsEx(flags)
}

// Interal Helpers
func exitWindowsEx(flags uint32) error {
	proc := user32.NewProc("ExitWindowsEx")
	ret, _, err := proc.Call(uintptr(flags), 0)
	if ret == 0 {
		return fmt.Errorf("ExitWindowsEx failed: %v", err)
	}
	return nil
}

func enablePrivilege(name string) error {
	var token windows.Token
	pHandle := windows.CurrentProcess()
	if err := windows.OpenProcessToken(pHandle, windows.TOKEN_ADJUST_PRIVILEGES|windows.TOKEN_QUERY, &token); err != nil {
		return err
	}
	defer token.Close()

	var luid windows.LUID
	if err := windows.LookupPrivilegeValue(nil, windows.StringToUTF16Ptr(name), &luid); err != nil {
		return err
	}

	newState := windows.Tokenprivileges{
		PrivilegeCount: 1,
		Privileges:     [1]windows.LUIDAndAttributes{{Luid: luid, Attributes: windows.SE_PRIVILEGE_ENABLED}},
	}
	return windows.AdjustTokenPrivileges(token, false, &newState, 0, nil, nil)
}
