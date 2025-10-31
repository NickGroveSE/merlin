package main

import (
	"fmt"
	"image"
	"image/png"
	"os"

	// "path/filepath"
	// "time"
	"unsafe"

	"github.com/kbinani/screenshot"
	"golang.org/x/sys/windows"
)

func capture() (string, error) {
	windowTitle := "Overwatch" // 👈 change this to your target window title

	img, err := captureWindowByTitle(windowTitle)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}

	// --- Ensure temp folder exists ---
	outputDir := "temp"
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		fmt.Println("Failed to create temp folder:", err)
		return "", err
	}

	outputPath, filename := generatePath(outputDir, "window_capture_")

	file, err := os.Create(outputPath)
	if err != nil {
		fmt.Println("Failed to create file:", err)
		return "", err
	}
	defer file.Close()

	if err := png.Encode(file, img); err != nil {
		fmt.Println("Failed to encode PNG:", err)
		return "", err
	}

	fmt.Printf("✅ Saved %s successfully\n", outputPath)
	return filename, nil
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
		int(rect.Left)+width/3,  // Start from horizontal midpoint
		int(rect.Top),           // Start from top
		int(rect.Right)-width/3, // Go to right edge
		int(rect.Top)+height,    // Go to vertical midpoint
	)

	fmt.Printf("Capturing top-right of window %q bounds: %+v\n", title, topRightBounds)

	// --- Capture that screen region ---
	img, err := screenshot.CaptureRect(topRightBounds)
	if err != nil {
		return nil, fmt.Errorf("capture failed: %w", err)
	}
	return img, nil
}
