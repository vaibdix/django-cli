package main

import (
	"fmt"
	"os/exec"
)

func (m *Model) runDjangoMigrations(projectPath string) error {
	m.updateProgress("Running database migrations...")

	var pythonPath string
	var runCmd func(args ...string) *exec.Cmd

	if isUvAvailable() {
		pythonPath = "uv"
		runCmd = func(args ...string) *exec.Cmd {
			uvArgs := append([]string{"run", "python"}, args...)
			return exec.Command(pythonPath, uvArgs...)
		}
	} else {
		pythonPath = getPythonPath(projectPath)
		runCmd = func(args ...string) *exec.Cmd {
			return exec.Command(pythonPath, args...)
		}
	}

	// Run makemigrations
	cmd := runCmd("manage.py", "makemigrations")
	cmd.Dir = projectPath
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to create migrations: %v\nOutput: %s", err, string(output))
	}
	m.stepMessages = append(m.stepMessages, "✅ Created database migrations.")

	// Run migrate
	cmd = runCmd("manage.py", "migrate")
	cmd.Dir = projectPath
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to apply migrations: %v\nOutput: %s", err, string(output))
	}
	m.stepMessages = append(m.stepMessages, "✅ Applied database migrations.")

	// Create sample data if REST framework is enabled and we have an app
	if m.setupRestFramework && m.appName != "" {
		m.updateProgress("Creating sample data...")
		cmd = runCmd("manage.py", "create_sample_data")
		cmd.Dir = projectPath
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to create sample data: %v\nOutput: %s", err, string(output))
		}
		m.stepMessages = append(m.stepMessages, "✅ Created sample data for the API.")
	}

	return nil
}