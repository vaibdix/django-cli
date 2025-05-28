package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// getPythonPath returns the correct Python executable path within the venv for the current OS
func getPythonPath(projectPath string) string {
	if runtime.GOOS == "windows" {
		return filepath.Join(projectPath, ".venv", "Scripts", "python.exe")
	}
	return filepath.Join(projectPath, ".venv", "bin", "python")
}

// getPipPath returns the correct pip executable path within the venv for the current OS
func getPipPath(projectPath string) string {
	if runtime.GOOS == "windows" {
		return filepath.Join(projectPath, ".venv", "Scripts", "pip.exe")
	}
	return filepath.Join(projectPath, ".venv", "bin", "pip")
}

// isCommandAvailable checks if a command is available in the system PATH
func isCommandAvailable(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

// getPythonCommand returns the best available Python command (python3 or python) for creating the venv.
func getPythonCommand() string {
	if isCommandAvailable("python3") {
		return "python3"
	}
	if isCommandAvailable("python") {
		return "python"
	}
	return "" // Indicates no suitable Python command found
}

// addToListInSettingsPy attempts to add an item to a Python list (like INSTALLED_APPS or MIDDLEWARE)
// in a settings.py content string. This is a simplified helper.
func addToListInSettingsPy(settingsContent, listName, itemToAdd string) (string, error) {
	quotedItem := fmt.Sprintf("'%s'", strings.Trim(itemToAdd, "'\""))

	if strings.Contains(settingsContent, quotedItem) {
		return settingsContent, nil // Already exists
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

	// Find the position of the opening bracket '['
	actualListStartIndex := listStartIndex + strings.Index(settingsContent[listStartIndex:], "[")

	// Find the corresponding closing bracket ']' for this list
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

	// Determine indentation (simple: use 4 spaces from the line of the list marker)
    // Find the start of the line where listMarker is found
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


	// Check content inside the list just before the closing bracket
	contentBeforeClosingBracket := strings.TrimSpace(settingsContent[actualListStartIndex+1 : listEndIndex])

	var newEntry string
	if contentBeforeClosingBracket == "" { // Empty list
		newEntry = fmt.Sprintf("\n%s%s,\n%s", itemIndent, quotedItem, baseIndent)
	} else if strings.HasSuffix(contentBeforeClosingBracket, ",") { // List has items and ends with a comma
		newEntry = fmt.Sprintf("\n%s%s,", itemIndent, quotedItem)
	} else { // List has items but does not end with a comma
		newEntry = fmt.Sprintf(",\n%s%s,", itemIndent, quotedItem)
	}

	return settingsContent[:listEndIndex] + newEntry + settingsContent[listEndIndex:], nil
}