package main

import (
	"fmt"
	"image"

	// "image/png"
	// "os"
	// "path/filepath"
	// "time"
	"unsafe"

	"github.com/kbinani/screenshot"
	"golang.org/x/sys/windows"
)

func captureHandler() *image.RGBA {
	windowTitle := "Overwatch" // 👈 change this to your target window title

	img, err := captureWindowByTitle(windowTitle)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}

	// // --- Ensure temp folder exists ---
	// outputDir := "temp"
	// if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
	// 	fmt.Println("Failed to create temp folder:", err)
	// 	return "failed"
	// }

	// timestamp := time.Now().Format("2006-01-02_15-04-05")
	// filename := fmt.Sprintf("window_capture_%s.png", timestamp)
	// outputPath := filepath.Join(outputDir, filename)

	// file, err := os.Create(outputPath)
	// if err != nil {
	// 	fmt.Println("Failed to create file:", err)
	// 	return "failed"
	// }
	// defer file.Close()

	// if err := png.Encode(file, img); err != nil {
	// 	fmt.Println("Failed to encode PNG:", err)
	// 	return "failed"
	// }

	// fmt.Printf("✅ Saved %s successfully\n", outputPath)
	return img
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

	bounds := image.Rect(int(rect.Left), int(rect.Top), int(rect.Right), int(rect.Bottom))
	fmt.Printf("Capturing window %q bounds: %+v\n", title, bounds)

	// --- Capture that screen region ---
	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		return nil, fmt.Errorf("capture failed: %w", err)
	}
	return img, nil
}
