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