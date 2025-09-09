package main

import (
	"go/build"
	"testing"
)

// RED: 必要な依存関係がimportできるかテスト
func TestDependencies_ShouldImportRequiredPackages(t *testing.T) {
	requiredPackages := []string{
		"github.com/rivo/tview",
		"github.com/spf13/cobra",
		"github.com/spf13/viper", 
		"github.com/google/uuid",
		"github.com/stretchr/testify/assert",
	}

	for _, pkg := range requiredPackages {
		if _, err := build.Import(pkg, ".", build.FindOnly); err != nil {
			t.Errorf("Required package %s cannot be imported: %v", pkg, err)
		}
	}
}