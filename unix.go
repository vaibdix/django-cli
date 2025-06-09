package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func installOnUnix() error {
	if runtime.GOOS == "windows" {
		return fmt.Errorf("this installation method is only for Linux and macOS")
	}

	// Get the executable path
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %v", err)
	}

	// Define installation directory
	var installDir string
	if runtime.GOOS == "darwin" {
		installDir = "/usr/local/bin"
	} else {
		// Linux
		installDir = "/usr/local/bin"
	}

	// Create target path
	targetPath := filepath.Join(installDir, "django-cli")

	// Copy executable to installation directory
	input, err := os.ReadFile(exe)
	if err != nil {
		return fmt.Errorf("failed to read executable: %v", err)
	}

	if err := os.WriteFile(targetPath, input, 0755); err != nil {
		return fmt.Errorf("failed to copy executable: %v", err)
	}

	fmt.Println("\nDjango CLI installed successfully!")
	fmt.Println("You can now use 'django-cli' from any directory.")
	fmt.Printf("Installation directory: %s\n", installDir)

	return nil
}
