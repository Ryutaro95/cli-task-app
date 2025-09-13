package validator

import (
	"testing"
	"time"

	"task-cli/internal/model"

	"github.com/stretchr/testify/assert"
)

// RED: Validatorのテスト
func TestValidator_New_ShouldCreateValidator(t *testing.T) {
	// When
	validator := New()

	// Then
	assert.NotNil(t, validator)
}

func TestValidator_ValidateTask_WithValidTask_ShouldReturnNil(t *testing.T) {
	// Given
	validator := New()
	task, err := model.NewTask("Valid Task", "Valid description", model.PriorityMedium, []string{"valid"})
	assert.NoError(t, err)

	// When
	result := validator.ValidateTask(task)

	// Then
	assert.NoError(t, result)
}

func TestValidator_ValidateTask_WithNilTask_ShouldReturnError(t *testing.T) {
	// Given
	validator := New()

	// When
	result := validator.ValidateTask(nil)

	// Then
	assert.Error(t, result)
	assert.Contains(t, result.Error(), "task cannot be nil")
}

func TestValidator_ValidateTask_WithEmptyTitle_ShouldReturnError(t *testing.T) {
	// Given
	validator := New()
	task := &model.Task{
		ID:          "test-id",
		Title:       "", // 空のタイトル
		Description: "Valid description",
		Status:      model.StatusTodo,
		Priority:    model.PriorityMedium,
		Tags:        []string{"test"},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// When
	result := validator.ValidateTask(task)

	// Then
	assert.Error(t, result)
	assert.Contains(t, result.Error(), "title")
}

func TestValidator_ValidateTask_WithTooLongTitle_ShouldReturnError(t *testing.T) {
	// Given
	validator := New()
	longTitle := make([]rune, 101) // 101文字のタイトル
	for i := range longTitle {
		longTitle[i] = 'a'
	}
	
	task := &model.Task{
		ID:          "test-id",
		Title:       string(longTitle),
		Description: "Valid description",
		Status:      model.StatusTodo,
		Priority:    model.PriorityMedium,
		Tags:        []string{"test"},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// When
	result := validator.ValidateTask(task)

	// Then
	assert.Error(t, result)
	assert.Contains(t, result.Error(), "title")
}

func TestValidator_ValidateTask_WithInvalidStatus_ShouldReturnError(t *testing.T) {
	// Given
	validator := New()
	task := &model.Task{
		ID:          "test-id",
		Title:       "Valid title",
		Description: "Valid description",
		Status:      model.Status("invalid_status"), // 無効なステータス
		Priority:    model.PriorityMedium,
		Tags:        []string{"test"},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// When
	result := validator.ValidateTask(task)

	// Then
	assert.Error(t, result)
	assert.Contains(t, result.Error(), "status")
}

func TestValidator_ValidateTask_WithInvalidPriority_ShouldReturnError(t *testing.T) {
	// Given
	validator := New()
	task := &model.Task{
		ID:          "test-id",
		Title:       "Valid title",
		Description: "Valid description",
		Status:      model.StatusTodo,
		Priority:    model.Priority("invalid_priority"), // 無効な優先度
		Tags:        []string{"test"},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// When
	result := validator.ValidateTask(task)

	// Then
	assert.Error(t, result)
	assert.Contains(t, result.Error(), "priority")
}

func TestValidator_ValidateTask_WithTooLongDescription_ShouldReturnError(t *testing.T) {
	// Given
	validator := New()
	longDescription := make([]rune, 501) // 501文字の説明
	for i := range longDescription {
		longDescription[i] = 'a'
	}
	
	task := &model.Task{
		ID:          "test-id",
		Title:       "Valid title",
		Description: string(longDescription),
		Status:      model.StatusTodo,
		Priority:    model.PriorityMedium,
		Tags:        []string{"test"},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// When
	result := validator.ValidateTask(task)

	// Then
	assert.Error(t, result)
	assert.Contains(t, result.Error(), "description")
}

func TestValidator_ValidateAppData_WithValidData_ShouldReturnNil(t *testing.T) {
	// Given
	validator := New()
	appData := model.NewAppData()
	task, _ := model.NewTask("Valid Task", "Valid description", model.PriorityMedium, []string{"valid"})
	appData.AddTask(task)

	// When
	result := validator.ValidateAppData(appData)

	// Then
	assert.NoError(t, result)
}

func TestValidator_ValidateAppData_WithNilData_ShouldReturnError(t *testing.T) {
	// Given
	validator := New()

	// When
	result := validator.ValidateAppData(nil)

	// Then
	assert.Error(t, result)
	assert.Contains(t, result.Error(), "appdata cannot be nil")
}

func TestValidator_ValidateAppData_WithInvalidTask_ShouldReturnError(t *testing.T) {
	// Given
	validator := New()
	appData := model.NewAppData()
	
	// 無効なタスクを直接追加
	invalidTask := &model.Task{
		ID:          "test-id",
		Title:       "", // 空のタイトル（無効）
		Description: "Valid description",
		Status:      model.StatusTodo,
		Priority:    model.PriorityMedium,
		Tags:        []string{"test"},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	appData.Tasks = append(appData.Tasks, invalidTask)

	// When
	result := validator.ValidateAppData(appData)

	// Then
	assert.Error(t, result)
	assert.Contains(t, result.Error(), "task")
}