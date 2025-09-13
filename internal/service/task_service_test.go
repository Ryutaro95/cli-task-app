package service

import (
	"context"
	"testing"

	"task-cli/internal/model"
	"task-cli/internal/validator"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository はRepositoryのモック実装
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Save(ctx context.Context, data *model.AppData) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

func (m *MockRepository) Load(ctx context.Context) (*model.AppData, error) {
	args := m.Called(ctx)
	return args.Get(0).(*model.AppData), args.Error(1)
}

func (m *MockRepository) CreateBackup(ctx context.Context, data *model.AppData) (string, error) {
	args := m.Called(ctx, data)
	return args.String(0), args.Error(1)
}

func (m *MockRepository) RestoreFromBackup(ctx context.Context, backupPath string) (*model.AppData, error) {
	args := m.Called(ctx, backupPath)
	return args.Get(0).(*model.AppData), args.Error(1)
}

// RED: TaskServiceのテスト
func TestTaskService_New_ShouldCreateService(t *testing.T) {
	// Given
	mockRepo := &MockRepository{}
	validator := validator.New()

	// When
	service := NewTaskService(mockRepo, validator)

	// Then
	assert.NotNil(t, service)
}

func TestTaskService_CreateTask_WithValidRequest_ShouldReturnTask(t *testing.T) {
	// Given
	mockRepo := &MockRepository{}
	validator := validator.New()
	service := NewTaskService(mockRepo, validator)
	ctx := context.Background()

	// 初期データの準備
	appData := model.NewAppData()
	mockRepo.On("Load", ctx).Return(appData, nil)
	mockRepo.On("Save", ctx, mock.AnythingOfType("*model.AppData")).Return(nil)

	request := CreateTaskRequest{
		Title:       "Test Task",
		Description: "Test Description",
		Priority:    model.PriorityMedium,
		Tags:        []string{"test"},
	}

	// When
	task, err := service.CreateTask(ctx, request)

	// Then
	assert.NoError(t, err)
	assert.NotNil(t, task)
	assert.Equal(t, "Test Task", task.Title)
	assert.Equal(t, "Test Description", task.Description)
	assert.Equal(t, model.PriorityMedium, task.Priority)
	assert.Equal(t, model.StatusTodo, task.Status)
	mockRepo.AssertExpectations(t)
}

func TestTaskService_CreateTask_WithInvalidRequest_ShouldReturnError(t *testing.T) {
	// Given
	mockRepo := &MockRepository{}
	validator := validator.New()
	service := NewTaskService(mockRepo, validator)
	ctx := context.Background()

	request := CreateTaskRequest{
		Title:       "", // 空のタイトル（無効）
		Description: "Test Description",
		Priority:    model.PriorityMedium,
		Tags:        []string{"test"},
	}

	// When
	task, err := service.CreateTask(ctx, request)

	// Then
	assert.Error(t, err)
	assert.Nil(t, task)
	assert.Contains(t, err.Error(), "title")
}

func TestTaskService_UpdateTask_WithValidData_ShouldUpdateTask(t *testing.T) {
	// Given
	mockRepo := &MockRepository{}
	validator := validator.New()
	service := NewTaskService(mockRepo, validator)
	ctx := context.Background()

	// 既存タスクのあるデータを準備
	appData := model.NewAppData()
	existingTask, _ := model.NewTask("Old Title", "Old Description", model.PriorityLow, []string{"old"})
	appData.AddTask(existingTask)

	mockRepo.On("Load", ctx).Return(appData, nil)
	mockRepo.On("Save", ctx, mock.AnythingOfType("*model.AppData")).Return(nil)

	request := UpdateTaskRequest{
		ID:          existingTask.ID,
		Title:       "Updated Title",
		Description: "Updated Description",
		Priority:    model.PriorityHigh,
		Status:      model.StatusInProgress,
		Tags:        []string{"updated"},
	}

	// When
	task, err := service.UpdateTask(ctx, request)

	// Then
	assert.NoError(t, err)
	assert.NotNil(t, task)
	assert.Equal(t, "Updated Title", task.Title)
	assert.Equal(t, "Updated Description", task.Description)
	assert.Equal(t, model.PriorityHigh, task.Priority)
	assert.Equal(t, model.StatusInProgress, task.Status)
	mockRepo.AssertExpectations(t)
}

func TestTaskService_UpdateTask_WithNonexistentID_ShouldReturnError(t *testing.T) {
	// Given
	mockRepo := &MockRepository{}
	validator := validator.New()
	service := NewTaskService(mockRepo, validator)
	ctx := context.Background()

	appData := model.NewAppData()
	mockRepo.On("Load", ctx).Return(appData, nil)

	request := UpdateTaskRequest{
		ID:          "non-existent-id",
		Title:       "Updated Title",
		Description: "Updated Description",
		Priority:    model.PriorityHigh,
		Status:      model.StatusInProgress,
		Tags:        []string{"updated"},
	}

	// When
	task, err := service.UpdateTask(ctx, request)

	// Then
	assert.Error(t, err)
	assert.Nil(t, task)
	assert.Contains(t, err.Error(), "not found")
	mockRepo.AssertExpectations(t)
}

func TestTaskService_DeleteTask_WithValidID_ShouldRemoveTask(t *testing.T) {
	// Given
	mockRepo := &MockRepository{}
	validator := validator.New()
	service := NewTaskService(mockRepo, validator)
	ctx := context.Background()

	// 既存タスクのあるデータを準備
	appData := model.NewAppData()
	existingTask, _ := model.NewTask("Task to Delete", "Description", model.PriorityMedium, []string{"delete"})
	appData.AddTask(existingTask)

	mockRepo.On("Load", ctx).Return(appData, nil)
	mockRepo.On("Save", ctx, mock.AnythingOfType("*model.AppData")).Return(nil)

	// When
	err := service.DeleteTask(ctx, existingTask.ID)

	// Then
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTaskService_ToggleTask_ShouldChangeStatus(t *testing.T) {
	// Given
	mockRepo := &MockRepository{}
	validator := validator.New()
	service := NewTaskService(mockRepo, validator)
	ctx := context.Background()

	// Todoステータスのタスクを準備
	appData := model.NewAppData()
	todoTask, _ := model.NewTask("Todo Task", "Description", model.PriorityMedium, []string{"toggle"})
	appData.AddTask(todoTask)

	mockRepo.On("Load", ctx).Return(appData, nil)
	mockRepo.On("Save", ctx, mock.AnythingOfType("*model.AppData")).Return(nil)

	// When
	task, err := service.ToggleTaskStatus(ctx, todoTask.ID)

	// Then
	assert.NoError(t, err)
	assert.NotNil(t, task)
	assert.Equal(t, model.StatusCompleted, task.Status)
	assert.NotNil(t, task.CompletedAt)
	mockRepo.AssertExpectations(t)
}

func TestTaskService_GetAllTasks_ShouldReturnTasks(t *testing.T) {
	// Given
	mockRepo := &MockRepository{}
	validator := validator.New()
	service := NewTaskService(mockRepo, validator)
	ctx := context.Background()

	appData := model.NewAppData()
	task1, _ := model.NewTask("Task 1", "Description 1", model.PriorityLow, []string{"test"})
	task2, _ := model.NewTask("Task 2", "Description 2", model.PriorityHigh, []string{"test"})
	appData.AddTask(task1)
	appData.AddTask(task2)

	mockRepo.On("Load", ctx).Return(appData, nil)

	// When
	tasks, err := service.GetAllTasks(ctx)

	// Then
	assert.NoError(t, err)
	assert.Len(t, tasks, 2)
	mockRepo.AssertExpectations(t)
}

func TestTaskService_SearchTasks_WithQuery_ShouldReturnMatchingTasks(t *testing.T) {
	// Given
	mockRepo := &MockRepository{}
	validator := validator.New()
	service := NewTaskService(mockRepo, validator)
	ctx := context.Background()

	appData := model.NewAppData()
	task1, _ := model.NewTask("Buy groceries", "Need to buy milk", model.PriorityMedium, []string{"shopping"})
	task2, _ := model.NewTask("Fix bug", "Fix login issue", model.PriorityHigh, []string{"dev"})
	appData.AddTask(task1)
	appData.AddTask(task2)

	mockRepo.On("Load", ctx).Return(appData, nil)

	// When
	results, err := service.SearchTasks(ctx, "buy")

	// Then
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "Buy groceries", results[0].Title)
	mockRepo.AssertExpectations(t)
}