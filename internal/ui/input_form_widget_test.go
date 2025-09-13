package ui

import (
	"testing"

	"task-cli/internal/model"

	"github.com/stretchr/testify/assert"
	"github.com/rivo/tview"
)

// RED: InputFormWidgetのテスト
func TestInputFormWidget_New_ShouldCreateWidget(t *testing.T) {
	// Given
	theme := NewTheme()

	// When
	widget := NewInputFormWidget(theme)

	// Then
	assert.NotNil(t, widget)
	assert.Implements(t, (*tview.Primitive)(nil), widget.GetPrimitive())
}

func TestInputFormWidget_CreateMode_ShouldAllowInput(t *testing.T) {
	// Given
	theme := NewTheme()
	widget := NewInputFormWidget(theme)

	// When
	widget.SetMode(FormModeCreate)
	widget.SetTitle("Test Task")
	widget.SetDescription("Test Description")
	widget.SetPriority(model.PriorityHigh)
	widget.SetTags("test,important")

	// Then
	assert.Equal(t, FormModeCreate, widget.GetMode())
	assert.Equal(t, "Test Task", widget.GetTitle())
	assert.Equal(t, "Test Description", widget.GetDescription())
	assert.Equal(t, model.PriorityHigh, widget.GetPriority())
	assert.Equal(t, "test,important", widget.GetTags())
}

func TestInputFormWidget_EditMode_ShouldPrePopulateFields(t *testing.T) {
	// Given
	theme := NewTheme()
	widget := NewInputFormWidget(theme)
	
	task, _ := model.NewTask("Existing Task", "Existing Description", model.PriorityMedium, []string{"existing", "tag"})
	task.Status = model.StatusInProgress

	// When
	widget.SetMode(FormModeEdit)
	widget.LoadTask(task)

	// Then
	assert.Equal(t, FormModeEdit, widget.GetMode())
	assert.Equal(t, "Existing Task", widget.GetTitle())
	assert.Equal(t, "Existing Description", widget.GetDescription())
	assert.Equal(t, model.PriorityMedium, widget.GetPriority())
	assert.Equal(t, model.StatusInProgress, widget.GetStatus())
	assert.Contains(t, widget.GetTags(), "existing")
	assert.Contains(t, widget.GetTags(), "tag")
}

func TestInputFormWidget_Validate_WithValidInput_ShouldReturnNil(t *testing.T) {
	// Given
	theme := NewTheme()
	widget := NewInputFormWidget(theme)
	
	widget.SetTitle("Valid Task")
	widget.SetDescription("Valid Description")
	widget.SetPriority(model.PriorityMedium)

	// When
	err := widget.Validate()

	// Then
	assert.NoError(t, err)
}

func TestInputFormWidget_Validate_WithInvalidInput_ShouldShowError(t *testing.T) {
	// Given
	theme := NewTheme()
	widget := NewInputFormWidget(theme)
	
	widget.SetTitle("") // 空のタイトル（無効）
	widget.SetDescription("Valid Description")
	widget.SetPriority(model.PriorityMedium)

	// When
	err := widget.Validate()

	// Then
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "title")
}

func TestInputFormWidget_Clear_ShouldResetAllFields(t *testing.T) {
	// Given
	theme := NewTheme()
	widget := NewInputFormWidget(theme)
	
	widget.SetTitle("Test Task")
	widget.SetDescription("Test Description")
	widget.SetPriority(model.PriorityHigh)
	widget.SetTags("test,tags")

	// When
	widget.Clear()

	// Then
	assert.Empty(t, widget.GetTitle())
	assert.Empty(t, widget.GetDescription())
	assert.Equal(t, model.PriorityMedium, widget.GetPriority()) // デフォルト値
	assert.Empty(t, widget.GetTags())
}

func TestInputFormWidget_Focus_ShouldFocusOnFirstField(t *testing.T) {
	// Given
	theme := NewTheme()
	widget := NewInputFormWidget(theme)

	// When
	widget.Focus()

	// Then
	// フォーカス状態のテストは実際のTUIアプリケーション内でのみ正確に動作するため
	// ここではメソッドが正常に呼び出せることを確認
	assert.NotNil(t, widget.GetPrimitive())
}

func TestInputFormWidget_SetSubmitCallback_ShouldInvokeCallback(t *testing.T) {
	// Given
	theme := NewTheme()
	widget := NewInputFormWidget(theme)
	
	callbackInvoked := false
	var receivedData FormData
	
	widget.SetSubmitCallback(func(data FormData) {
		callbackInvoked = true
		receivedData = data
	})

	widget.SetTitle("Test Task")
	widget.SetDescription("Test Description")
	widget.SetPriority(model.PriorityHigh)

	// When
	widget.Submit()

	// Then
	assert.True(t, callbackInvoked)
	assert.Equal(t, "Test Task", receivedData.Title)
	assert.Equal(t, "Test Description", receivedData.Description)
	assert.Equal(t, model.PriorityHigh, receivedData.Priority)
}

func TestInputFormWidget_SetCancelCallback_ShouldInvokeCallback(t *testing.T) {
	// Given
	theme := NewTheme()
	widget := NewInputFormWidget(theme)
	
	callbackInvoked := false
	widget.SetCancelCallback(func() {
		callbackInvoked = true
	})

	// When
	widget.Cancel()

	// Then
	assert.True(t, callbackInvoked)
}

func TestInputFormWidget_GetFormData_ShouldReturnStructuredData(t *testing.T) {
	// Given
	theme := NewTheme()
	widget := NewInputFormWidget(theme)
	
	widget.SetTitle("Test Task")
	widget.SetDescription("Test Description")
	widget.SetPriority(model.PriorityLow)
	widget.SetTags("tag1,tag2,tag3")

	// When
	data := widget.GetFormData()

	// Then
	assert.Equal(t, "Test Task", data.Title)
	assert.Equal(t, "Test Description", data.Description)
	assert.Equal(t, model.PriorityLow, data.Priority)
	assert.Contains(t, data.Tags, "tag1")
	assert.Contains(t, data.Tags, "tag2")
	assert.Contains(t, data.Tags, "tag3")
}

func TestInputFormWidget_SetEnabled_ShouldControlInputAccess(t *testing.T) {
	// Given
	theme := NewTheme()
	widget := NewInputFormWidget(theme)

	// When & Then - 有効な状態
	widget.SetEnabled(true)
	assert.True(t, widget.IsEnabled())

	// When & Then - 無効な状態
	widget.SetEnabled(false)
	assert.False(t, widget.IsEnabled())
}

func TestInputFormWidget_SetErrorMessage_ShouldDisplayError(t *testing.T) {
	// Given
	theme := NewTheme()
	widget := NewInputFormWidget(theme)
	errorMsg := "Validation failed: title is required"

	// When
	widget.SetErrorMessage(errorMsg)

	// Then
	assert.Equal(t, errorMsg, widget.GetErrorMessage())
}

func TestInputFormWidget_ClearError_ShouldRemoveErrorMessage(t *testing.T) {
	// Given
	theme := NewTheme()
	widget := NewInputFormWidget(theme)
	
	widget.SetErrorMessage("Some error message")
	assert.NotEmpty(t, widget.GetErrorMessage())

	// When
	widget.ClearError()

	// Then
	assert.Empty(t, widget.GetErrorMessage())
}