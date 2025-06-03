package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func (m *Model) setupTailwindCSS(projectPath string) error {
	if !m.setupTailwind {
		return nil
	}

	m.updateProgress("Setting up Tailwind CSS v4...")

	if !isCommandAvailable("npm") {
		m.stepMessages = append(m.stepMessages, "⚠️  Warning: npm not found. Please install Node.js to use Tailwind CSS.")
		return nil
	}

	cmd := exec.Command("npm", "init", "-y")
	cmd.Dir = projectPath
	if output, err := cmd.CombinedOutput(); err != nil {
		m.stepMessages = append(m.stepMessages, fmt.Sprintf("⚠️  Warning: Failed to initialize npm: %v\nOutput: %s", err, string(output)))
		return nil
	}
	m.stepMessages = append(m.stepMessages, "✅ npm initialized.")

	cmd = exec.Command("npm", "install", "tailwindcss", "@tailwindcss/cli")
	cmd.Dir = projectPath
	if output, err := cmd.CombinedOutput(); err != nil {
		m.stepMessages = append(m.stepMessages, fmt.Sprintf("⚠️  Warning: Failed to install Tailwind CSS: %v\nOutput: %s", err, string(output)))
		return nil
	}
	m.stepMessages = append(m.stepMessages, "✅ Tailwind CSS v4 installed.")

	staticSrcPath := filepath.Join(projectPath, "static", "src")
	staticDistPath := filepath.Join(projectPath, "static", "dist")
	if err := os.MkdirAll(staticSrcPath, 0755); err != nil {
		m.stepMessages = append(m.stepMessages, fmt.Sprintf("⚠️  Warning: Failed to create static/src directory: %v", err))
		return nil
	}
	if err := os.MkdirAll(staticDistPath, 0755); err != nil {
		m.stepMessages = append(m.stepMessages, fmt.Sprintf("⚠️  Warning: Failed to create static/dist directory: %v", err))
		return nil
	}
	m.stepMessages = append(m.stepMessages, "✅ Tailwind directory structure created.")
	tailwindCSS := `@import "tailwindcss";`
	if err := os.WriteFile(filepath.Join(staticSrcPath, "styles.css"), []byte(tailwindCSS), 0644); err != nil {
		m.stepMessages = append(m.stepMessages, fmt.Sprintf("⚠️  Warning: Failed to create styles.css: %v", err))
		return nil
	}
	m.stepMessages = append(m.stepMessages, "✅ Tailwind source CSS created.")
	if err := m.updatePackageJSONForTailwind(projectPath); err != nil {
		m.stepMessages = append(m.stepMessages, fmt.Sprintf("⚠️  Warning: Failed to update package.json: %v", err))
		return nil
	}
	if m.createTemplates {
		if err := m.updateBaseTemplateForTailwind(projectPath); err != nil {
			m.stepMessages = append(m.stepMessages, fmt.Sprintf("⚠️  Warning: Failed to update base.html: %v", err))
			return nil
		}
	}
	cmd = exec.Command("npm", "run", "build:css")
	cmd.Dir = projectPath
	if output, err := cmd.CombinedOutput(); err != nil {
		m.stepMessages = append(m.stepMessages, fmt.Sprintf("⚠️  Warning: Failed to build Tailwind CSS: %v\nOutput: %s", err, string(output)))
	} else {
		m.stepMessages = append(m.stepMessages, "✅ Tailwind CSS compiled successfully.")
		m.stepMessages = append(m.stepMessages, "💡 Run 'npm run watch:css' for development or 'npm run build:css' for production.")
	}
	return nil
}

func (m *Model) updatePackageJSONForTailwind(projectPath string) error {
	packageJSONPath := filepath.Join(projectPath, "package.json")
	packageData, err := os.ReadFile(packageJSONPath)
	if err != nil {
		return fmt.Errorf("failed to read package.json: %v", err)
	}
	var packageJSON map[string]interface{}
	if err := json.Unmarshal(packageData, &packageJSON); err != nil {
		return fmt.Errorf("failed to parse package.json: %v", err)
	}
	scripts := map[string]interface{}{
		"build:css": "npx tailwindcss -i ./static/src/styles.css -o ./static/dist/styles.css",
		"watch:css": "npx tailwindcss -i ./static/src/styles.css -o ./static/dist/styles.css --watch",
	}
	packageJSON["scripts"] = scripts
	updatedData, err := json.MarshalIndent(packageJSON, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal package.json: %v", err)
	}
	if err := os.WriteFile(packageJSONPath, updatedData, 0644); err != nil {
		return fmt.Errorf("failed to write package.json: %v", err)
	}
	m.stepMessages = append(m.stepMessages, "✅ package.json updated with Tailwind scripts.")
	return nil
}

func (m *Model) updateBaseTemplateForTailwind(projectPath string) error {
	baseTemplatePath := filepath.Join(projectPath, "templates", "base.html")
	baseContent, err := os.ReadFile(baseTemplatePath)
	if err != nil {
		return fmt.Errorf("failed to read base.html: %v", err)
	}
	updatedBaseContent := strings.Replace(string(baseContent),
		`<link rel="stylesheet" href="{% static 'css/style.css' %}">`,
		`<link rel="stylesheet" href="{% static 'dist/styles.css' %}">`, 1)
	if err := os.WriteFile(baseTemplatePath, []byte(updatedBaseContent), 0644); err != nil {
		return fmt.Errorf("failed to update base.html: %v", err)
	}

	m.stepMessages = append(m.stepMessages, "✅ base.html updated to use Tailwind CSS.")
	return nil
}
