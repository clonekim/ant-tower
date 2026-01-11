package sysagent

import "golang.org/x/sys/windows"

func HideConsole() {
	hwnd, _, _ := procGetConsoleWindow.Call()
	if hwnd != 0 {
		procShowWindow.Call(hwnd, uintptr(windows.SW_HIDE)) // SW_HIDE = 0
	}
}
