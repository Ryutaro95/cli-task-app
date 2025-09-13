package ui

import (
	"testing"

	"task-cli/internal/model"

	"github.com/stretchr/testify/assert"
	"github.com/gdamore/tcell/v2"
)

// RED: テーマシステムのテスト
func TestTheme_NewTheme_ShouldCreateDefaultTheme(t *testing.T) {
	// When
	theme := NewTheme()

	// Then
	assert.NotNil(t, theme)
}

func TestTheme_GetPriorityColor_ShouldReturnCorrectColor(t *testing.T) {
	// Given
	theme := NewTheme()

	tests := []struct {
		priority model.Priority
		expected tcell.Color
	}{
		{model.PriorityHigh, tcell.ColorRed},
		{model.PriorityMedium, tcell.ColorYellow},
		{model.PriorityLow, tcell.ColorGreen},
	}

	for _, tt := range tests {
		t.Run(string(tt.priority), func(t *testing.T) {
			// When
			color := theme.GetPriorityColor(tt.priority)

			// Then
			assert.Equal(t, tt.expected, color)
		})
	}
}

func TestTheme_GetStatusColor_ShouldReturnCorrectColor(t *testing.T) {
	// Given
	theme := NewTheme()

	tests := []struct {
		status   model.Status
		expected tcell.Color
	}{
		{model.StatusTodo, tcell.ColorWhite},
		{model.StatusInProgress, tcell.ColorBlue},
		{model.StatusCompleted, tcell.ColorGreen},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			// When
			color := theme.GetStatusColor(tt.status)

			// Then
			assert.Equal(t, tt.expected, color)
		})
	}
}

func TestTheme_GetStatusStyle_ShouldReturnCorrectStyle(t *testing.T) {
	// Given
	theme := NewTheme()

	// When
	todoStyle := theme.GetStatusStyle(model.StatusTodo)
	progressStyle := theme.GetStatusStyle(model.StatusInProgress)
	completedStyle := theme.GetStatusStyle(model.StatusCompleted)

	// Then
	assert.NotEqual(t, todoStyle, progressStyle)
	assert.NotEqual(t, todoStyle, completedStyle)
	assert.NotEqual(t, progressStyle, completedStyle)
}

func TestTheme_GetPriorityStyle_ShouldReturnCorrectStyle(t *testing.T) {
	// Given
	theme := NewTheme()

	// When
	highStyle := theme.GetPriorityStyle(model.PriorityHigh)
	mediumStyle := theme.GetPriorityStyle(model.PriorityMedium)
	lowStyle := theme.GetPriorityStyle(model.PriorityLow)

	// Then
	assert.NotEqual(t, highStyle, mediumStyle)
	assert.NotEqual(t, highStyle, lowStyle)
	assert.NotEqual(t, mediumStyle, lowStyle)
}

func TestTheme_GetUIColors_ShouldReturnPredefinedColors(t *testing.T) {
	// Given
	theme := NewTheme()

	// When & Then
	assert.NotEqual(t, tcell.ColorDefault, theme.GetBackgroundColor())
	assert.NotEqual(t, tcell.ColorDefault, theme.GetForegroundColor())
	assert.NotEqual(t, tcell.ColorDefault, theme.GetBorderColor())
	assert.NotEqual(t, tcell.ColorDefault, theme.GetHighlightColor())
	assert.NotEqual(t, tcell.ColorDefault, theme.GetSelectionColor())
}

func TestTheme_CreateCustomTheme_ShouldAllowCustomization(t *testing.T) {
	// Given
	customColors := ThemeConfig{
		Background:  tcell.ColorBlack,
		Foreground:  tcell.ColorWhite,
		Border:      tcell.ColorGray,
		Highlight:   tcell.ColorBlue,
		Selection:   tcell.ColorDarkBlue,
		StatusColors: map[model.Status]tcell.Color{
			model.StatusTodo:       tcell.ColorWhite,
			model.StatusInProgress: tcell.ColorYellow,
			model.StatusCompleted:  tcell.ColorGreen,
		},
		PriorityColors: map[model.Priority]tcell.Color{
			model.PriorityHigh:   tcell.ColorRed,
			model.PriorityMedium: tcell.ColorYellow,
			model.PriorityLow:    tcell.ColorGreen,
		},
	}

	// When
	theme := NewCustomTheme(customColors)

	// Then
	assert.Equal(t, tcell.ColorBlack, theme.GetBackgroundColor())
	assert.Equal(t, tcell.ColorWhite, theme.GetForegroundColor())
	assert.Equal(t, tcell.ColorGray, theme.GetBorderColor())
	assert.Equal(t, tcell.ColorBlue, theme.GetHighlightColor())
	assert.Equal(t, tcell.ColorDarkBlue, theme.GetSelectionColor())
}

func TestTheme_GetInvalidStatus_ShouldReturnDefaultColor(t *testing.T) {
	// Given
	theme := NewTheme()
	invalidStatus := model.Status("invalid")

	// When
	color := theme.GetStatusColor(invalidStatus)

	// Then
	assert.Equal(t, tcell.ColorWhite, color) // デフォルトカラー
}

func TestTheme_GetInvalidPriority_ShouldReturnDefaultColor(t *testing.T) {
	// Given
	theme := NewTheme()
	invalidPriority := model.Priority("invalid")

	// When
	color := theme.GetPriorityColor(invalidPriority)

	// Then
	assert.Equal(t, tcell.ColorWhite, color) // デフォルトカラー
}

func TestTheme_GetFormattedTaskText_ShouldApplyCorrectFormatting(t *testing.T) {
	// Given
	theme := NewTheme()
	task := &model.Task{
		ID:       "1",
		Title:    "Test Task",
		Priority: model.PriorityHigh,
		Status:   model.StatusTodo,
	}

	// When
	formattedText := theme.FormatTaskText(task)

	// Then
	assert.Contains(t, formattedText, "Test Task")
	assert.NotEmpty(t, formattedText)
}

func TestTheme_DarkTheme_ShouldUseDarkColors(t *testing.T) {
	// When
	darkTheme := NewDarkTheme()

	// Then
	assert.Equal(t, tcell.ColorBlack, darkTheme.GetBackgroundColor())
	assert.NotNil(t, darkTheme)
}

func TestTheme_LightTheme_ShouldUseLightColors(t *testing.T) {
	// When
	lightTheme := NewLightTheme()

	// Then
	assert.Equal(t, tcell.ColorWhite, lightTheme.GetBackgroundColor())
	assert.NotNil(t, lightTheme)
}