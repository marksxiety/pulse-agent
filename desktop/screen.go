//go:build windows

package main

import "golang.org/x/sys/windows"

const (
	smCXScreen = 0
	smCYScreen = 1
)

var (
	user32               = windows.NewLazySystemDLL("user32.dll")
	procGetSystemMetrics = user32.NewProc("GetSystemMetrics")
)

// primaryScreenSize returns the size of the primary display in physical
// pixels. It returns zeroes if the call fails, which callers treat as "leave
// the window where the OS put it".
func primaryScreenSize() (int, int) {
	w, _, _ := procGetSystemMetrics.Call(uintptr(smCXScreen))
	h, _, _ := procGetSystemMetrics.Call(uintptr(smCYScreen))
	return int(w), int(h)
}
