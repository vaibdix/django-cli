package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"runtime"
)

func getPythonPath(projectPath string) string {
	if runtime.GOOS == "windows" {
		return filepath.Join(projectPath, ".venv", "Scripts", "python.exe")
	}
	return filepath.Join(projectPath, ".venv", "bin", "python")
}
func getPipPath(projectPath string) string {
	if runtime.GOOS == "windows" {
		return filepath.Join(projectPath, ".venv", "Scripts", "pip.exe")
	}
	return filepath.Join(projectPath, ".venv", "bin", "pip")
}
func isCommandAvailable(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}
func addToListInSettingsPy(settingsContent, listName, itemToAdd string) (string, error) {
	quotedItem := fmt.Sprintf("'%s'", strings.Trim(itemToAdd, "'\""))

	if strings.Contains(settingsContent, quotedItem) {
		return settingsContent, nil
	}

	listMarker := fmt.Sprintf("%s = [", listName)
	listStartIndex := strings.Index(settingsContent, listMarker)
	if listStartIndex == -1 {
		listMarker = fmt.Sprintf("%s=[", listName) // Try without space
		listStartIndex = strings.Index(settingsContent, listMarker)
		if listStartIndex == -1 {
			return settingsContent, fmt.Errorf("could not find list '%s' in settings", listName)
		}
	}

	actualListStartIndex := listStartIndex + strings.Index(settingsContent[listStartIndex:], "[")
	openBracketCount := 0
	listEndIndex := -1
	for i := actualListStartIndex; i < len(settingsContent); i++ {
		if settingsContent[i] == '[' {
			openBracketCount++
		} else if settingsContent[i] == ']' {
			openBracketCount--
			if openBracketCount == 0 {
				listEndIndex = i
				break
			}
		}
	}

	if listEndIndex == -1 {
		return settingsContent, fmt.Errorf("could not find closing bracket for list '%s'", listName)
	}
	lineStartForListMarker := 0
	if idx := strings.LastIndex(settingsContent[:listStartIndex], "\n"); idx != -1 {
		lineStartForListMarker = idx + 1
	}
	baseIndent := ""
	for _, r := range settingsContent[lineStartForListMarker:listStartIndex] {
		if r == ' ' || r == '\t' {
			baseIndent += string(r)
		} else {
			break
		}
	}
	itemIndent := baseIndent + "    "
	contentBeforeClosingBracket := strings.TrimSpace(settingsContent[actualListStartIndex+1 : listEndIndex])

	var newEntry string
	if contentBeforeClosingBracket == "" { // Empty list
		newEntry = fmt.Sprintf("\n%s%s,\n%s", itemIndent, quotedItem, baseIndent)
	} else if strings.HasSuffix(contentBeforeClosingBracket, ",") { // List has items and ends with a comma
		newEntry = fmt.Sprintf("\n%s%s,", itemIndent, quotedItem)
	} else {
		newEntry = fmt.Sprintf(",\n%s%s,", itemIndent, quotedItem)
	}

	return settingsContent[:listEndIndex] + newEntry + settingsContent[listEndIndex:], nil
}

func validateProjectName(name string) error {
	if name == "" {
		return fmt.Errorf("project name cannot be empty")
	}

	invalidChars := []string{"<", ">", ":", "\"", "|", "?", "*", " ", "/", "\\"}
	for _, char := range invalidChars {
		if strings.Contains(name, char) {
			return fmt.Errorf("project name cannot contain '%s'", char)
		}
	}

	if _, err := os.Stat(name); err == nil {
		return fmt.Errorf("directory '%s' already exists", name)
	}

	reserved := []string{"and", "as", "assert", "break", "class", "continue",
		"def", "del", "elif", "else", "except", "exec", "finally", "for",
		"from", "global", "if", "import", "in", "is", "lambda", "not",
		"or", "pass", "print", "raise", "return", "try", "while", "with", "yield"}

	for _, word := range reserved {
		if strings.ToLower(name) == word {
			return fmt.Errorf("'%s' is a Python reserved word and cannot be used as project name", name)
		}
	}

	return nil
}

func validateDjangoVersion(version string) error {
	if version == "" || version == "latest" {
		return nil
	}

	parts := strings.Split(version, ".")
	if len(parts) < 2 || len(parts) > 3 {
		return fmt.Errorf("invalid Django version format. Use format like '4.2.0' or '5.1'")
	}
	for _, p := range parts {
		for _, r := range p {
			if r < '0' || r > '9' {
				return fmt.Errorf("invalid Django version format. Use format like '4.2.0' or '5.1' (numeric parts only)")
			}
		}
	}

	return nil
}
