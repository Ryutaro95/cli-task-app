package model

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// RED: Taskモデルのテスト
func TestTask_NewTask_WithValidData_ShouldCreateTask(t *testing.T) {
	// Given
	title := "Test Task"
	description := "Test Description"
	priority := PriorityHigh
	tags := []string{"test", "development"}

	// When
	task, err := NewTask(title, description, priority, tags)

	// Then
	assert.NoError(t, err)
	assert.NotEmpty(t, task.ID)
	assert.Equal(t, title, task.Title)
	assert.Equal(t, description, task.Description)
	assert.Equal(t, StatusTodo, task.Status)
	assert.Equal(t, priority, task.Priority)
	assert.Equal(t, tags, task.Tags)
	assert.False(t, task.CreatedAt.IsZero())
	assert.False(t, task.UpdatedAt.IsZero())
	assert.Nil(t, task.CompletedAt)
}

func TestTask_Validate_WithEmptyTitle_ShouldReturnError(t *testing.T) {
	// Given
	task := Task{
		ID:          "test-id",
		Title:       "", // 空のタイトル
		Description: "Description",
		Status:      StatusTodo,
		Priority:    PriorityMedium,
		Tags:        []string{"test"},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// When
	err := task.Validate()

	// Then
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "title")
}

func TestTask_Validate_WithTooLongTitle_ShouldReturnError(t *testing.T) {
	// Given
	longTitle := string(make([]byte, 101)) // 101文字のタイトル
	task := Task{
		ID:          "test-id",
		Title:       longTitle,
		Description: "Description",
		Status:      StatusTodo,
		Priority:    PriorityMedium,
		Tags:        []string{"test"},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// When
	err := task.Validate()

	// Then
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "title")
}

func TestTask_ToJSON_ShouldReturnValidJSON(t *testing.T) {
	// Given
	task := Task{
		ID:          "test-id",
		Title:       "Test Task",
		Description: "Test Description",
		Status:      StatusTodo,
		Priority:    PriorityHigh,
		Tags:        []string{"test"},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// When
	jsonData, err := json.Marshal(task)

	// Then
	assert.NoError(t, err)
	assert.Contains(t, string(jsonData), "test-id")
	assert.Contains(t, string(jsonData), "Test Task")
}

func TestTask_IsCompleted_WhenStatusCompleted_ShouldReturnTrue(t *testing.T) {
	// Given
	task := Task{Status: StatusCompleted}

	// When
	result := task.IsCompleted()

	// Then
	assert.True(t, result)
}

func TestTask_IsCompleted_WhenStatusTodo_ShouldReturnFalse(t *testing.T) {
	// Given
	task := Task{Status: StatusTodo}

	// When
	result := task.IsCompleted()

	// Then
	assert.False(t, result)
}

// Status型のテスト
func TestStatus_String_ShouldReturnCorrectValue(t *testing.T) {
	tests := []struct {
		status   Status
		expected string
	}{
		{StatusTodo, "todo"},
		{StatusInProgress, "in_progress"},
		{StatusCompleted, "completed"},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			// When
			result := tt.status.String()

			// Then
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestStatus_IsValid_WithValidStatus_ShouldReturnTrue(t *testing.T) {
	validStatuses := []Status{StatusTodo, StatusInProgress, StatusCompleted}

	for _, status := range validStatuses {
		t.Run(string(status), func(t *testing.T) {
			// When
			result := status.IsValid()

			// Then
			assert.True(t, result)
		})
	}
}

func TestStatus_IsValid_WithInvalidStatus_ShouldReturnFalse(t *testing.T) {
	invalidStatuses := []Status{"invalid", "unknown", ""}

	for _, status := range invalidStatuses {
		t.Run(string(status), func(t *testing.T) {
			// When
			result := status.IsValid()

			// Then
			assert.False(t, result)
		})
	}
}

// Priority型のテスト
func TestPriority_String_ShouldReturnCorrectValue(t *testing.T) {
	tests := []struct {
		priority Priority
		expected string
	}{
		{PriorityLow, "low"},
		{PriorityMedium, "medium"},
		{PriorityHigh, "high"},
	}

	for _, tt := range tests {
		t.Run(string(tt.priority), func(t *testing.T) {
			// When
			result := tt.priority.String()

			// Then
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPriority_IsValid_WithValidPriority_ShouldReturnTrue(t *testing.T) {
	validPriorities := []Priority{PriorityLow, PriorityMedium, PriorityHigh}

	for _, priority := range validPriorities {
		t.Run(string(priority), func(t *testing.T) {
			// When
			result := priority.IsValid()

			// Then
			assert.True(t, result)
		})
	}
}

func TestPriority_IsValid_WithInvalidPriority_ShouldReturnFalse(t *testing.T) {
	invalidPriorities := []Priority{"invalid", "unknown", ""}

	for _, priority := range invalidPriorities {
		t.Run(string(priority), func(t *testing.T) {
			// When
			result := priority.IsValid()

			// Then
			assert.False(t, result)
		})
	}
}