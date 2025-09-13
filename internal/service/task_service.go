package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"task-cli/internal/model"
	"task-cli/internal/repository"
	"task-cli/internal/validator"
)

// TaskService はタスク関連のビジネスロジックを提供する
type TaskService struct {
	repo      repository.Repository
	validator *validator.Validator
}

// CreateTaskRequest はタスク作成のリクエスト
type CreateTaskRequest struct {
	Title       string
	Description string
	Priority    model.Priority
	Tags        []string
	DueDate     *time.Time
}

// UpdateTaskRequest はタスク更新のリクエスト
type UpdateTaskRequest struct {
	ID          string
	Title       string
	Description string
	Priority    model.Priority
	Status      model.Status
	Tags        []string
	DueDate     *time.Time
}

// NewTaskService は新しいTaskServiceを作成する
func NewTaskService(repo repository.Repository, validator *validator.Validator) *TaskService {
	return &TaskService{
		repo:      repo,
		validator: validator,
	}
}

// CreateTask は新しいタスクを作成する
func (s *TaskService) CreateTask(ctx context.Context, request CreateTaskRequest) (*model.Task, error) {
	// リクエストの基本バリデーション
	if request.Title == "" {
		return nil, errors.New("title is required")
	}

	// データを読み込み
	appData, err := s.loadAppData(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load data: %w", err)
	}

	// 新しいタスクを作成
	task, err := model.NewTask(request.Title, request.Description, request.Priority, request.Tags)
	if err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	// 期限日が指定されている場合は設定
	if request.DueDate != nil {
		task.DueDate = request.DueDate
	}

	// バリデーション
	if err := s.validator.ValidateTask(task); err != nil {
		return nil, fmt.Errorf("task validation failed: %w", err)
	}

	// データに追加
	if err := appData.AddTask(task); err != nil {
		return nil, fmt.Errorf("failed to add task: %w", err)
	}

	// データを保存
	if err := s.repo.Save(ctx, appData); err != nil {
		return nil, fmt.Errorf("failed to save data: %w", err)
	}

	return task, nil
}

// UpdateTask は既存のタスクを更新する
func (s *TaskService) UpdateTask(ctx context.Context, request UpdateTaskRequest) (*model.Task, error) {
	// データを読み込み
	appData, err := s.loadAppData(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load data: %w", err)
	}

	// 既存のタスクを取得
	existingTask, err := appData.GetTaskByID(request.ID)
	if err != nil {
		return nil, fmt.Errorf("task not found: %w", err)
	}

	// タスクを更新
	existingTask.Title = request.Title
	existingTask.Description = request.Description
	existingTask.Priority = request.Priority
	existingTask.Status = request.Status
	existingTask.Tags = request.Tags
	existingTask.DueDate = request.DueDate
	existingTask.UpdatedAt = time.Now()

	// ステータスが完了に変更された場合、完了日時を設定
	if request.Status == model.StatusCompleted && existingTask.CompletedAt == nil {
		now := time.Now()
		existingTask.CompletedAt = &now
	}

	// バリデーション
	if err := s.validator.ValidateTask(existingTask); err != nil {
		return nil, fmt.Errorf("task validation failed: %w", err)
	}

	// データを更新
	if err := appData.UpdateTask(existingTask); err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	// データを保存
	if err := s.repo.Save(ctx, appData); err != nil {
		return nil, fmt.Errorf("failed to save data: %w", err)
	}

	return existingTask, nil
}

// DeleteTask はタスクを削除する
func (s *TaskService) DeleteTask(ctx context.Context, taskID string) error {
	// データを読み込み
	appData, err := s.loadAppData(ctx)
	if err != nil {
		return fmt.Errorf("failed to load data: %w", err)
	}

	// タスクが存在するか確認
	_, err = appData.GetTaskByID(taskID)
	if err != nil {
		return fmt.Errorf("task not found: %w", err)
	}

	// タスクを削除
	if err := appData.DeleteTask(taskID); err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	// データを保存
	if err := s.repo.Save(ctx, appData); err != nil {
		return fmt.Errorf("failed to save data: %w", err)
	}

	return nil
}

// ToggleTaskStatus はタスクのステータスを切り替える
func (s *TaskService) ToggleTaskStatus(ctx context.Context, taskID string) (*model.Task, error) {
	// データを読み込み
	appData, err := s.loadAppData(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load data: %w", err)
	}

	// 既存のタスクを取得
	task, err := appData.GetTaskByID(taskID)
	if err != nil {
		return nil, fmt.Errorf("task not found: %w", err)
	}

	// ステータスを切り替え
	switch task.Status {
	case model.StatusTodo:
		task.Status = model.StatusCompleted
		now := time.Now()
		task.CompletedAt = &now
	case model.StatusInProgress:
		task.Status = model.StatusCompleted
		now := time.Now()
		task.CompletedAt = &now
	case model.StatusCompleted:
		task.Status = model.StatusTodo
		task.CompletedAt = nil
	}

	task.UpdatedAt = time.Now()

	// データを更新
	if err := appData.UpdateTask(task); err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	// データを保存
	if err := s.repo.Save(ctx, appData); err != nil {
		return nil, fmt.Errorf("failed to save data: %w", err)
	}

	return task, nil
}

// GetAllTasks は全てのタスクを取得する
func (s *TaskService) GetAllTasks(ctx context.Context) ([]*model.Task, error) {
	appData, err := s.loadAppData(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load data: %w", err)
	}

	return appData.Tasks, nil
}

// GetTaskByID はIDでタスクを取得する
func (s *TaskService) GetTaskByID(ctx context.Context, taskID string) (*model.Task, error) {
	appData, err := s.loadAppData(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load data: %w", err)
	}

	task, err := appData.GetTaskByID(taskID)
	if err != nil {
		return nil, fmt.Errorf("task not found: %w", err)
	}

	return task, nil
}

// SearchTasks はクエリでタスクを検索する
func (s *TaskService) SearchTasks(ctx context.Context, query string) ([]*model.Task, error) {
	appData, err := s.loadAppData(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load data: %w", err)
	}

	return appData.SearchTasksByTitle(query), nil
}

// GetTasksByStatus はステータスでタスクを取得する
func (s *TaskService) GetTasksByStatus(ctx context.Context, status model.Status) ([]*model.Task, error) {
	appData, err := s.loadAppData(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load data: %w", err)
	}

	return appData.GetTasksByStatus(status), nil
}

// GetTasksByPriority は優先度でタスクを取得する
func (s *TaskService) GetTasksByPriority(ctx context.Context, priority model.Priority) ([]*model.Task, error) {
	appData, err := s.loadAppData(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load data: %w", err)
	}

	return appData.GetTasksByPriority(priority), nil
}

// loadAppData はAppDataを読み込み、存在しない場合は新しいインスタンスを作成する
func (s *TaskService) loadAppData(ctx context.Context) (*model.AppData, error) {
	appData, err := s.repo.Load(ctx)
	if err != nil {
		// データが存在しない場合は新しいインスタンスを作成
		return model.NewAppData(), nil
	}
	return appData, nil
}