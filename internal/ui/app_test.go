package ui

import (
	"context"
	"testing"
	"time"

	"task-cli/internal/model"
	"task-cli/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTaskService はTaskServiceのモック実装
type MockTaskService struct {
	mock.Mock
}

func (m *MockTaskService) CreateTask(ctx context.Context, request service.CreateTaskRequest) (*model.Task, error) {
	args := m.Called(ctx, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Task), args.Error(1)
}

func (m *MockTaskService) UpdateTask(ctx context.Context, request service.UpdateTaskRequest) (*model.Task, error) {
	args := m.Called(ctx, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Task), args.Error(1)
}

func (m *MockTaskService) DeleteTask(ctx context.Context, taskID string) error {
	args := m.Called(ctx, taskID)
	return args.Error(0)
}

func (m *MockTaskService) ToggleTaskStatus(ctx context.Context, taskID string) (*model.Task, error) {
	args := m.Called(ctx, taskID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Task), args.Error(1)
}

func (m *MockTaskService) GetAllTasks(ctx context.Context) ([]*model.Task, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*model.Task), args.Error(1)
}

func (m *MockTaskService) SearchTasks(ctx context.Context, query string) ([]*model.Task, error) {
	args := m.Called(ctx, query)
	return args.Get(0).([]*model.Task), args.Error(1)
}

func (m *MockTaskService) GetTasksByStatus(ctx context.Context, status model.Status) ([]*model.Task, error) {
	args := m.Called(ctx, status)
	return args.Get(0).([]*model.Task), args.Error(1)
}

func (m *MockTaskService) GetTasksByPriority(ctx context.Context, priority model.Priority) ([]*model.Task, error) {
	args := m.Called(ctx, priority)
	return args.Get(0).([]*model.Task), args.Error(1)
}

// RED: メインAppのテスト
func TestApp_New_ShouldCreateApp(t *testing.T) {
	// Given
	mockTaskService := &MockTaskService{}
	stateManager := service.NewStateManager()
	theme := NewTheme()

	// When
	app := NewApp(mockTaskService, stateManager, theme)

	// Then
	assert.NotNil(t, app)
}

func TestApp_Initialize_ShouldSetupComponents(t *testing.T) {
	// Given
	mockTaskService := &MockTaskService{}
	stateManager := service.NewStateManager()
	theme := NewTheme()
	app := NewApp(mockTaskService, stateManager, theme)

	// データの準備
	testTasks := []*model.Task{
		{ID: "1", Title: "Test Task 1", Priority: model.PriorityHigh},
		{ID: "2", Title: "Test Task 2", Priority: model.PriorityMedium},
	}
	mockTaskService.On("GetAllTasks", mock.Anything).Return(testTasks, nil)

	// When
	err := app.Initialize()

	// Then
	assert.NoError(t, err)
	mockTaskService.AssertExpectations(t)
}

func TestApp_SwitchView_ShouldChangeActiveWidget(t *testing.T) {
	// Given
	mockTaskService := &MockTaskService{}
	stateManager := service.NewStateManager()
	theme := NewTheme()
	app := NewApp(mockTaskService, stateManager, theme)

	mockTaskService.On("GetAllTasks", mock.Anything).Return([]*model.Task{}, nil)
	app.Initialize()

	// When & Then - リストビューに切り替え
	app.SwitchToListView()
	assert.Equal(t, ViewModeList, app.GetCurrentView())

	// When & Then - フォームビューに切り替え
	app.SwitchToFormView()
	assert.Equal(t, ViewModeForm, app.GetCurrentView())
}

func TestApp_RefreshTasks_ShouldUpdateTaskList(t *testing.T) {
	// Given
	mockTaskService := &MockTaskService{}
	stateManager := service.NewStateManager()
	theme := NewTheme()
	app := NewApp(mockTaskService, stateManager, theme)

	initialTasks := []*model.Task{
		{ID: "1", Title: "Initial Task", Priority: model.PriorityMedium},
	}
	updatedTasks := []*model.Task{
		{ID: "1", Title: "Initial Task", Priority: model.PriorityMedium},
		{ID: "2", Title: "New Task", Priority: model.PriorityHigh},
	}

	mockTaskService.On("GetAllTasks", mock.Anything).Return(initialTasks, nil).Once()
	mockTaskService.On("GetAllTasks", mock.Anything).Return(updatedTasks, nil).Once()

	app.Initialize()

	// When
	err := app.RefreshTasks()

	// Then
	assert.NoError(t, err)
	mockTaskService.AssertExpectations(t)
}

func TestApp_HandleCreateTask_ShouldCreateNewTask(t *testing.T) {
	// Given
	mockTaskService := &MockTaskService{}
	stateManager := service.NewStateManager()
	theme := NewTheme()
	app := NewApp(mockTaskService, stateManager, theme)

	mockTaskService.On("GetAllTasks", mock.Anything).Return([]*model.Task{}, nil)
	app.Initialize()

	newTask := &model.Task{
		ID:          "new-id",
		Title:       "New Task",
		Description: "New Description",
		Priority:    model.PriorityHigh,
	}

	createRequest := service.CreateTaskRequest{
		Title:       "New Task",
		Description: "New Description",
		Priority:    model.PriorityHigh,
		Tags:        []string{"test"},
	}

	mockTaskService.On("CreateTask", mock.Anything, createRequest).Return(newTask, nil)
	mockTaskService.On("GetAllTasks", mock.Anything).Return([]*model.Task{newTask}, nil)

	formData := FormData{
		Title:       "New Task",
		Description: "New Description",
		Priority:    model.PriorityHigh,
		Tags:        []string{"test"},
	}

	// When
	err := app.HandleCreateTask(formData)

	// Then
	assert.NoError(t, err)
	mockTaskService.AssertExpectations(t)
}

func TestApp_HandleUpdateTask_ShouldUpdateExistingTask(t *testing.T) {
	// Given
	mockTaskService := &MockTaskService{}
	stateManager := service.NewStateManager()
	theme := NewTheme()
	app := NewApp(mockTaskService, stateManager, theme)

	existingTask := &model.Task{
		ID:          "existing-id",
		Title:       "Existing Task",
		Description: "Existing Description",
		Priority:    model.PriorityMedium,
	}

	updatedTask := &model.Task{
		ID:          "existing-id",
		Title:       "Updated Task",
		Description: "Updated Description",
		Priority:    model.PriorityHigh,
	}

	mockTaskService.On("GetAllTasks", mock.Anything).Return([]*model.Task{existingTask}, nil)
	app.Initialize()

	updateRequest := service.UpdateTaskRequest{
		ID:          "existing-id",
		Title:       "Updated Task",
		Description: "Updated Description",
		Priority:    model.PriorityHigh,
		Status:      model.StatusInProgress,
		Tags:        []string{"updated"},
	}

	mockTaskService.On("UpdateTask", mock.Anything, updateRequest).Return(updatedTask, nil)
	mockTaskService.On("GetAllTasks", mock.Anything).Return([]*model.Task{updatedTask}, nil)

	formData := FormData{
		Title:       "Updated Task",
		Description: "Updated Description",
		Priority:    model.PriorityHigh,
		Status:      model.StatusInProgress,
		Tags:        []string{"updated"},
	}

	// When
	err := app.HandleUpdateTask("existing-id", formData)

	// Then
	assert.NoError(t, err)
	mockTaskService.AssertExpectations(t)
}

func TestApp_HandleDeleteTask_ShouldRemoveTask(t *testing.T) {
	// Given
	mockTaskService := &MockTaskService{}
	stateManager := service.NewStateManager()
	theme := NewTheme()
	app := NewApp(mockTaskService, stateManager, theme)

	existingTask := &model.Task{
		ID:    "task-to-delete",
		Title: "Task to Delete",
	}

	mockTaskService.On("GetAllTasks", mock.Anything).Return([]*model.Task{existingTask}, nil)
	app.Initialize()

	mockTaskService.On("DeleteTask", mock.Anything, "task-to-delete").Return(nil)
	mockTaskService.On("GetAllTasks", mock.Anything).Return([]*model.Task{}, nil)

	// When
	err := app.HandleDeleteTask("task-to-delete")

	// Then
	assert.NoError(t, err)
	mockTaskService.AssertExpectations(t)
}

func TestApp_HandleToggleTask_ShouldToggleTaskStatus(t *testing.T) {
	// Given
	mockTaskService := &MockTaskService{}
	stateManager := service.NewStateManager()
	theme := NewTheme()
	app := NewApp(mockTaskService, stateManager, theme)

	originalTask := &model.Task{
		ID:     "task-to-toggle",
		Title:  "Task to Toggle",
		Status: model.StatusTodo,
	}

	toggledTask := &model.Task{
		ID:          "task-to-toggle",
		Title:       "Task to Toggle",
		Status:      model.StatusCompleted,
		CompletedAt: &[]time.Time{time.Now()}[0],
	}

	mockTaskService.On("GetAllTasks", mock.Anything).Return([]*model.Task{originalTask}, nil)
	app.Initialize()

	mockTaskService.On("ToggleTaskStatus", mock.Anything, "task-to-toggle").Return(toggledTask, nil)
	mockTaskService.On("GetAllTasks", mock.Anything).Return([]*model.Task{toggledTask}, nil)

	// When
	err := app.HandleToggleTask("task-to-toggle")

	// Then
	assert.NoError(t, err)
	mockTaskService.AssertExpectations(t)
}

func TestApp_ApplyFilter_ShouldUpdateFilteredView(t *testing.T) {
	// Given
	mockTaskService := &MockTaskService{}
	stateManager := service.NewStateManager()
	theme := NewTheme()
	app := NewApp(mockTaskService, stateManager, theme)

	tasks := []*model.Task{
		{ID: "1", Title: "Todo Task", Status: model.StatusTodo, Priority: model.PriorityHigh},
		{ID: "2", Title: "Completed Task", Status: model.StatusCompleted, Priority: model.PriorityMedium},
	}

	mockTaskService.On("GetAllTasks", mock.Anything).Return(tasks, nil)
	app.Initialize()

	filter := service.TaskFilter{
		Status: &[]model.Status{model.StatusTodo}[0],
	}

	// When
	app.ApplyFilter(filter)

	// Then
	currentFilter := app.GetCurrentFilter()
	assert.Equal(t, model.StatusTodo, *currentFilter.Status)
}

func TestApp_GetCurrentView_ShouldReturnActiveView(t *testing.T) {
	// Given
	mockTaskService := &MockTaskService{}
	stateManager := service.NewStateManager()
	theme := NewTheme()
	app := NewApp(mockTaskService, stateManager, theme)

	mockTaskService.On("GetAllTasks", mock.Anything).Return([]*model.Task{}, nil)
	app.Initialize()

	// When & Then - 初期状態はリストビュー
	assert.Equal(t, ViewModeList, app.GetCurrentView())

	// フォームビューに切り替えてテスト
	app.SwitchToFormView()
	assert.Equal(t, ViewModeForm, app.GetCurrentView())
}