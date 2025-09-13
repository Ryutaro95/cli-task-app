package ui

import (
	"fmt"

	"task-cli/internal/model"

	"github.com/gdamore/tcell/v2"
)

// Theme はUI要素の色とスタイルを定義する
type Theme struct {
	config ThemeConfig
}

// ThemeConfig はテーマの設定を定義
type ThemeConfig struct {
	Background     tcell.Color
	Foreground     tcell.Color
	Border         tcell.Color
	Highlight      tcell.Color
	Selection      tcell.Color
	StatusColors   map[model.Status]tcell.Color
	PriorityColors map[model.Priority]tcell.Color
}

// NewTheme はデフォルトテーマを作成する
func NewTheme() *Theme {
	return &Theme{
		config: getDefaultThemeConfig(),
	}
}

// NewCustomTheme はカスタムテーマを作成する
func NewCustomTheme(config ThemeConfig) *Theme {
	return &Theme{
		config: config,
	}
}

// NewDarkTheme はダークテーマを作成する
func NewDarkTheme() *Theme {
	return &Theme{
		config: getDarkThemeConfig(),
	}
}

// NewLightTheme はライトテーマを作成する
func NewLightTheme() *Theme {
	return &Theme{
		config: getLightThemeConfig(),
	}
}

// GetPriorityColor は優先度に対応する色を取得する
func (t *Theme) GetPriorityColor(priority model.Priority) tcell.Color {
	if color, exists := t.config.PriorityColors[priority]; exists {
		return color
	}
	return tcell.ColorWhite // デフォルトカラー
}

// GetStatusColor はステータスに対応する色を取得する
func (t *Theme) GetStatusColor(status model.Status) tcell.Color {
	if color, exists := t.config.StatusColors[status]; exists {
		return color
	}
	return tcell.ColorWhite // デフォルトカラー
}

// GetStatusStyle はステータスに対応するスタイルを取得する
func (t *Theme) GetStatusStyle(status model.Status) tcell.Style {
	color := t.GetStatusColor(status)
	return tcell.StyleDefault.Foreground(color).Background(t.config.Background)
}

// GetPriorityStyle は優先度に対応するスタイルを取得する
func (t *Theme) GetPriorityStyle(priority model.Priority) tcell.Style {
	color := t.GetPriorityColor(priority)
	return tcell.StyleDefault.Foreground(color).Background(t.config.Background)
}

// GetBackgroundColor は背景色を取得する
func (t *Theme) GetBackgroundColor() tcell.Color {
	return t.config.Background
}

// GetForegroundColor は前景色を取得する
func (t *Theme) GetForegroundColor() tcell.Color {
	return t.config.Foreground
}

// GetBorderColor はボーダー色を取得する
func (t *Theme) GetBorderColor() tcell.Color {
	return t.config.Border
}

// GetHighlightColor はハイライト色を取得する
func (t *Theme) GetHighlightColor() tcell.Color {
	return t.config.Highlight
}

// GetSelectionColor は選択色を取得する
func (t *Theme) GetSelectionColor() tcell.Color {
	return t.config.Selection
}

// FormatTaskText はタスクのテキストをフォーマットする
func (t *Theme) FormatTaskText(task *model.Task) string {
	statusSymbol := t.getStatusSymbol(task.Status)
	prioritySymbol := t.getPrioritySymbol(task.Priority)
	return fmt.Sprintf("%s %s %s", statusSymbol, prioritySymbol, task.Title)
}

// getStatusSymbol はステータスに対応するシンボルを取得する
func (t *Theme) getStatusSymbol(status model.Status) string {
	switch status {
	case model.StatusTodo:
		return "◯"
	case model.StatusInProgress:
		return "◐"
	case model.StatusCompleted:
		return "●"
	default:
		return "?"
	}
}

// getPrioritySymbol は優先度に対応するシンボルを取得する
func (t *Theme) getPrioritySymbol(priority model.Priority) string {
	switch priority {
	case model.PriorityHigh:
		return "🔴"
	case model.PriorityMedium:
		return "🟡"
	case model.PriorityLow:
		return "🟢"
	default:
		return "⚪"
	}
}

// getDefaultThemeConfig はデフォルトテーマの設定を取得する
func getDefaultThemeConfig() ThemeConfig {
	return ThemeConfig{
		Background: tcell.ColorBlack,
		Foreground: tcell.ColorWhite,
		Border:     tcell.ColorGray,
		Highlight:  tcell.ColorBlue,
		Selection:  tcell.ColorDarkBlue,
		StatusColors: map[model.Status]tcell.Color{
			model.StatusTodo:       tcell.ColorWhite,
			model.StatusInProgress: tcell.ColorBlue,
			model.StatusCompleted:  tcell.ColorGreen,
		},
		PriorityColors: map[model.Priority]tcell.Color{
			model.PriorityHigh:   tcell.ColorRed,
			model.PriorityMedium: tcell.ColorYellow,
			model.PriorityLow:    tcell.ColorGreen,
		},
	}
}

// getDarkThemeConfig はダークテーマの設定を取得する
func getDarkThemeConfig() ThemeConfig {
	return ThemeConfig{
		Background: tcell.ColorBlack,
		Foreground: tcell.ColorWhite,
		Border:     tcell.ColorDarkGray,
		Highlight:  tcell.ColorBlue,
		Selection:  tcell.ColorDarkBlue,
		StatusColors: map[model.Status]tcell.Color{
			model.StatusTodo:       tcell.ColorSilver,
			model.StatusInProgress: tcell.ColorLightBlue,
			model.StatusCompleted:  tcell.ColorLightGreen,
		},
		PriorityColors: map[model.Priority]tcell.Color{
			model.PriorityHigh:   tcell.ColorRed,
			model.PriorityMedium: tcell.ColorYellow,
			model.PriorityLow:    tcell.ColorGreen,
		},
	}
}

// getLightThemeConfig はライトテーマの設定を取得する
func getLightThemeConfig() ThemeConfig {
	return ThemeConfig{
		Background: tcell.ColorWhite,
		Foreground: tcell.ColorBlack,
		Border:     tcell.ColorGray,
		Highlight:  tcell.ColorBlue,
		Selection:  tcell.ColorLightBlue,
		StatusColors: map[model.Status]tcell.Color{
			model.StatusTodo:       tcell.ColorBlack,
			model.StatusInProgress: tcell.ColorDarkBlue,
			model.StatusCompleted:  tcell.ColorDarkGreen,
		},
		PriorityColors: map[model.Priority]tcell.Color{
			model.PriorityHigh:   tcell.ColorDarkRed,
			model.PriorityMedium: tcell.ColorOrange,
			model.PriorityLow:    tcell.ColorDarkGreen,
		},
	}
}