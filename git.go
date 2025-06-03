package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func (m *Model) initializeGitRepository(projectPath string) error {
	if !m.initializeGit {
		return nil
	}

	gitCmd := exec.Command("git", "init")
	gitCmd.Dir = projectPath
	if output, err := gitCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to initialize Git repository: %v\nOutput: %s", err, string(output))
	}
	m.stepMessages = append(m.stepMessages, "✅ Git repository initialized.")

	gitignoreContent := `# Django
*.log
*.pot
*.pyc
__pycache__/
local_settings.py
db.sqlite3
db.sqlite3-journal
media

# Virtual environment
venv/
.venv/
env/
ENV/

# IDE
.vscode/
.idea/
*.swp
*.swo

# OS
.DS_Store
Thumbs.db
`
	if err := os.WriteFile(filepath.Join(projectPath, ".gitignore"), []byte(gitignoreContent), 0644); err != nil {
		return fmt.Errorf("failed to create .gitignore: %v", err)
	}
	m.stepMessages = append(m.stepMessages, "✅ .gitignore file created.")

	return nil
}