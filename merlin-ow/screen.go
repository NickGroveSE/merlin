package main

import (
	"fmt"
	"image"
	"unsafe"

	"github.com/kbinani/screenshot"
	"golang.org/x/sys/windows"
)

func capture() (*image.RGBA, error) {
	windowTitle := "Overwatch" // 👈 change this to your target window title

	img, err := captureWindowByTitle(windowTitle)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	return img, nil
}

func captureWindowByTitle(title string) (*image.RGBA, error) {
	// --- Load user32.dll ---
	user32 := windows.NewLazySystemDLL("user32.dll")
	procFindWindowW := user32.NewProc("FindWindowW")
	procGetWindowRect := user32.NewProc("GetWindowRect")

	// --- Find the window by title ---
	titlePtr, _ := windows.UTF16PtrFromString(title)
	hwnd, _, _ := procFindWindowW.Call(0, uintptr(unsafe.Pointer(titlePtr)))
	if hwnd == 0 {
		return nil, fmt.Errorf("window with title %q not found", title)
	}

	// --- Get window rect ---
	var rect struct {
		Left, Top, Right, Bottom int32
	}
	ret, _, _ := procGetWindowRect.Call(hwnd, uintptr(unsafe.Pointer(&rect)))
	if ret == 0 {
		return nil, fmt.Errorf("failed to get window rect")
	}

	// --- Calculate top-right quadrant ---
	width := int(rect.Right - rect.Left)
	height := int(rect.Bottom - rect.Top)

	// Top-right quarter of the window
	topRightBounds := image.Rect(
		int(rect.Left),        // Start from horizontal midpoint
		int(rect.Top),         // Start from top
		int(rect.Right)+width, // Go to right edge
		int(rect.Top)+height,  // Go to vertical midpoint
	)

	fmt.Printf("Capturing top-right of window %q bounds: %+v\n", title, topRightBounds)

	// --- Capture that screen region ---
	img, err := screenshot.CaptureRect(topRightBounds)
	if err != nil {
		return nil, fmt.Errorf("capture failed: %w", err)
	}
	return img, nil
}
