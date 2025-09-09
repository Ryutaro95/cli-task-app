package main

import (
	"os"
	"path/filepath"
	"testing"
)

// RED: プロジェクト構造の検証テスト
func TestProjectStructure_ShouldHaveRequiredDirectories(t *testing.T) {
	requiredDirs := []string{
		"cmd",
		"internal",
		"internal/app",
		"internal/model",
		"internal/service",
		"internal/repository",
		"internal/ui",
		"internal/ui/widgets",
		"internal/validator",
		"pkg",
		"pkg/errors",
	}

	for _, dir := range requiredDirs {
		dirPath := filepath.Join(".", dir)
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			t.Errorf("Required directory %s does not exist", dir)
		}
	}
}

func TestProjectStructure_ShouldHaveMainFile(t *testing.T) {
	mainFile := "main.go"
	if _, err := os.Stat(mainFile); os.IsNotExist(err) {
		t.Errorf("Required file %s does not exist", mainFile)
	}
}

func TestProjectStructure_ShouldHaveGoMod(t *testing.T) {
	goModFile := "go.mod"
	if _, err := os.Stat(goModFile); os.IsNotExist(err) {
		t.Errorf("Required file %s does not exist", goModFile)
	}
}