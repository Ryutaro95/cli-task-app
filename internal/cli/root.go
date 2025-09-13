package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"task-cli/internal/repository"
	"task-cli/internal/service"
	"task-cli/internal/ui"
	"task-cli/internal/validator"

	"github.com/spf13/cobra"
)

// Config はアプリケーションの設定
type Config struct {
	DataDir string
	Theme   string
}

// NewConfig は新しい設定を作成する
func NewConfig() *Config {
	homeDir, _ := os.UserHomeDir()
	defaultDataDir := filepath.Join(homeDir, ".task-cli")
	
	return &Config{
		DataDir: defaultDataDir,
		Theme:   "default",
	}
}

// SetDataDir はデータディレクトリを設定する
func (c *Config) SetDataDir(dataDir string) {
	c.DataDir = dataDir
}

// GetDataDir はデータディレクトリを取得する
func (c *Config) GetDataDir() string {
	return c.DataDir
}

// SetTheme はテーマを設定する
func (c *Config) SetTheme(theme string) {
	c.Theme = theme
}

// GetTheme はテーマを取得する
func (c *Config) GetTheme() string {
	return c.Theme
}

// Validate は設定を検証する
func (c *Config) Validate() error {
	switch c.Theme {
	case "default", "dark", "light":
		return nil
	default:
		return errors.New("invalid theme: must be one of 'default', 'dark', or 'light'")
	}
}

var (
	globalConfig = NewConfig()
)

// NewRootCommand は新しいルートコマンドを作成する
func NewRootCommand() *cobra.Command {
	// 各テスト用に新しい設定を作成
	config := NewConfig()
	
	rootCmd := &cobra.Command{
		Use:   "task-cli",
		Short: "Task management TUI application",
		Long: `A terminal-based task management application with a text user interface.
Manage your tasks efficiently with keyboard shortcuts and a clean interface.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAppWithConfig(config)
		},
		Version: "1.0.0",
	}

	// フラグを設定
	rootCmd.PersistentFlags().StringVar(&config.DataDir, "data-dir", config.DataDir,
		"Directory to store task data")
	rootCmd.PersistentFlags().StringVar(&config.Theme, "theme", config.Theme,
		"Theme to use (default, dark, light)")

	return rootCmd
}

// Execute はルートコマンドを実行する
func Execute() error {
	// メイン実行用は単一の設定を使用
	rootCmd := &cobra.Command{
		Use:   "task-cli",
		Short: "Task management TUI application",
		Long: `A terminal-based task management application with a text user interface.
Manage your tasks efficiently with keyboard shortcuts and a clean interface.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAppWithConfig(globalConfig)
		},
		Version: "1.0.0",
	}

	// フラグを設定
	rootCmd.PersistentFlags().StringVar(&globalConfig.DataDir, "data-dir", globalConfig.DataDir,
		"Directory to store task data")
	rootCmd.PersistentFlags().StringVar(&globalConfig.Theme, "theme", globalConfig.Theme,
		"Theme to use (default, dark, light)")

	return rootCmd.Execute()
}

// runAppWithConfig は指定された設定でアプリケーションを実行する
func runAppWithConfig(config *Config) error {
	// 設定を検証
	if err := config.Validate(); err != nil {
		return fmt.Errorf("configuration error: %w", err)
	}

	// コンポーネントを初期化
	repo := repository.NewFileRepository(config.DataDir)
	validator := validator.New()
	taskService := service.NewTaskService(repo, validator)
	stateManager := service.NewStateManager()
	
	// テーマを作成
	theme, err := createTheme(config.Theme)
	if err != nil {
		return fmt.Errorf("failed to create theme: %w", err)
	}

	// UIアプリケーションを作成
	app := ui.NewApp(taskService, stateManager, theme)

	// アプリケーションを初期化
	if err := app.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize application: %w", err)
	}

	// アプリケーションを実行
	return app.Run()
}

// createTheme は指定されたテーマ名からテーマを作成する
func createTheme(themeName string) (*ui.Theme, error) {
	switch themeName {
	case "default":
		return ui.NewTheme(), nil
	case "dark":
		return ui.NewDarkTheme(), nil
	case "light":
		return ui.NewLightTheme(), nil
	default:
		return nil, fmt.Errorf("unknown theme: %s", themeName)
	}
}