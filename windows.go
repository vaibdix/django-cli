package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func installOnWindows() error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("installation command is only supported on Windows")
	}

	// Get the executable path
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %v", err)
	}

	// Define installation directory
	localAppData := os.Getenv("LOCALAPPDATA")
	if localAppData == "" {
		return fmt.Errorf("LOCALAPPDATA environment variable not found")
	}

	installDir := filepath.Join(localAppData, "DjangoCLI")
	targetPath := filepath.Join(installDir, "django-cli.exe")

	// Create installation directory
	if err := os.MkdirAll(installDir, 0755); err != nil {
		return fmt.Errorf("failed to create installation directory: %v", err)
	}

	// Copy executable to installation directory
	input, err := os.ReadFile(exe)
	if err != nil {
		return fmt.Errorf("failed to read executable: %v", err)
	}

	if err := os.WriteFile(targetPath, input, 0755); err != nil {
		return fmt.Errorf("failed to copy executable: %v", err)
	}

	// Add to PATH using PowerShell
	cmd := exec.Command("powershell", "-Command", `
		$userPath = [Environment]::GetEnvironmentVariable("Path", "User")
		if ($userPath -notlike "*`+installDir+`*") {
			[Environment]::SetEnvironmentVariable("Path", "$userPath;`+installDir+`", "User")
		}
	`)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to update PATH: %v", err)
	}

	fmt.Println("\nDjango CLI installed successfully!")
	fmt.Println("You can now use 'django-cli' from any directory.")
	fmt.Printf("Installation directory: %s\n", installDir)
	fmt.Println("\nPlease restart your terminal for the PATH changes to take effect.")

	return nil
}
