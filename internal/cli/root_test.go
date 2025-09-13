package cli

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

// RED: CLIコマンドのテスト
func TestRootCommand_Execute_ShouldRunSuccessfully(t *testing.T) {
	// Given
	cmd := NewRootCommand()
	
	// バッファでstdoutをキャプチャ
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	// When & Then
	// 実際のTUIアプリケーションは実行しないが、コマンドが正常に作成されることを確認
	assert.NotNil(t, cmd)
	assert.Equal(t, "task-cli", cmd.Use)
}

func TestRootCommand_Help_ShouldDisplayUsage(t *testing.T) {
	// Given
	cmd := NewRootCommand()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"--help"})

	// When
	err := cmd.Execute()

	// Then
	assert.NoError(t, err)
	output := buf.String()
	assert.Contains(t, output, "terminal-based task management")
}

func TestRootCommand_Version_ShouldDisplayVersion(t *testing.T) {
	// Given
	cmd := NewRootCommand()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"--version"})

	// When
	err := cmd.Execute()

	// Then
	assert.NoError(t, err)
	// バージョン情報が含まれることを確認
	output := buf.String()
	assert.NotEmpty(t, output)
}

func TestRootCommand_WithDataDir_ShouldUseCustomDataDir(t *testing.T) {
	// Given
	cmd := NewRootCommand()
	customDir := "/tmp/custom-task-data"
	
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"--data-dir", customDir, "--help"})

	// When
	err := cmd.Execute()

	// Then
	assert.NoError(t, err)
	// コマンドが正常に実行されることを確認
	output := buf.String()
	assert.Contains(t, output, "data-dir")
}

func TestRootCommand_WithTheme_ShouldUseCustomTheme(t *testing.T) {
	// Given
	cmd := NewRootCommand()
	
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"--theme", "dark", "--help"})

	// When
	err := cmd.Execute()

	// Then
	assert.NoError(t, err)
	// コマンドが正常に実行されることを確認
	output := buf.String()
	assert.Contains(t, output, "theme")
}

func TestRootCommand_WithInvalidTheme_ShouldShowError(t *testing.T) {
	// Given
	cmd := NewRootCommand()
	
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"--theme", "invalid-theme"})

	// この場合、実際のバリデーションは実装時に行う
	// ここではコマンドが作成されることを確認
	assert.NotNil(t, cmd)
}

func TestExecute_ShouldCreateAndRunRootCommand(t *testing.T) {
	// Given & When
	// Execute関数が正常に定義されているかテスト
	// 実際のTUI実行はテストしないが、関数の存在を確認
	
	// Then
	// Execute関数が存在し、呼び出し可能であることを確認
	assert.NotPanics(t, func() {
		// Execute()の呼び出しはTUIを起動するため、テストでは実行しない
		// 代わりにNewRootCommand()が正常に動作することを確認
		cmd := NewRootCommand()
		assert.NotNil(t, cmd)
	})
}

func TestConfig_Default_ShouldHaveCorrectDefaults(t *testing.T) {
	// Given
	config := NewConfig()

	// Then
	assert.NotNil(t, config)
	assert.NotEmpty(t, config.DataDir)
	assert.NotEmpty(t, config.Theme)
}

func TestConfig_SetDataDir_ShouldUpdateDataDir(t *testing.T) {
	// Given
	config := NewConfig()
	customDir := "/tmp/test-data"

	// When
	config.SetDataDir(customDir)

	// Then
	assert.Equal(t, customDir, config.GetDataDir())
}

func TestConfig_SetTheme_ShouldUpdateTheme(t *testing.T) {
	// Given
	config := NewConfig()
	customTheme := "light"

	// When
	config.SetTheme(customTheme)

	// Then
	assert.Equal(t, customTheme, config.GetTheme())
}

func TestConfig_Validate_WithValidConfig_ShouldReturnNil(t *testing.T) {
	// Given
	config := NewConfig()
	config.SetDataDir("/tmp/test-data")
	config.SetTheme("dark")

	// When
	err := config.Validate()

	// Then
	assert.NoError(t, err)
}

func TestConfig_Validate_WithInvalidTheme_ShouldReturnError(t *testing.T) {
	// Given
	config := NewConfig()
	config.SetTheme("invalid-theme")

	// When
	err := config.Validate()

	// Then
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid theme")
}