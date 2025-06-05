package main

import (
	"fmt"
	"os/exec"
)

func (m *Model) runDjangoMigrations(projectPath string) error {
	m.updateProgress("Running database migrations...")
	pythonPath := getPythonPath(projectPath)

	// Run makemigrations
	cmd := exec.Command(pythonPath, "manage.py", "makemigrations")
	cmd.Dir = projectPath
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to create migrations: %v\nOutput: %s", err, string(output))
	}
	m.stepMessages = append(m.stepMessages, "✅ Created database migrations.")

	// Run migrate
	cmd = exec.Command(pythonPath, "manage.py", "migrate")
	cmd.Dir = projectPath
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to apply migrations: %v\nOutput: %s", err, string(output))
	}
	m.stepMessages = append(m.stepMessages, "✅ Applied database migrations.")

	// Create sample data if REST framework is enabled and we have an app
	if m.setupRestFramework && m.appName != "" {
		m.updateProgress("Creating sample data...")
		cmd = exec.Command(pythonPath, "manage.py", "create_sample_data")
		cmd.Dir = projectPath
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to create sample data: %v\nOutput: %s", err, string(output))
		}
		m.stepMessages = append(m.stepMessages, "✅ Created sample data for the API.")
	}

	return nil
}
