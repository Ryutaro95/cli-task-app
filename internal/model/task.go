package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Task はタスクの基本構造を定義
type Task struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      Status     `json:"status"`
	Priority    Priority   `json:"priority"`
	Tags        []string   `json:"tags"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"`
}

// Status はタスクのステータスを定義
type Status string

const (
	StatusTodo       Status = "todo"
	StatusInProgress Status = "in_progress"
	StatusCompleted  Status = "completed"
)

// Priority はタスクの優先度を定義
type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

// NewTask は新しいTaskを作成する
func NewTask(title, description string, priority Priority, tags []string) (*Task, error) {
	now := time.Now()
	task := &Task{
		ID:          uuid.New().String(),
		Title:       title,
		Description: description,
		Status:      StatusTodo,
		Priority:    priority,
		Tags:        tags,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := task.Validate(); err != nil {
		return nil, err
	}

	return task, nil
}

// Validate はTaskの値を検証する
func (t *Task) Validate() error {
	if t.Title == "" {
		return errors.New("title is required")
	}
	if len(t.Title) > 100 {
		return errors.New("title must be 100 characters or less")
	}
	if len(t.Description) > 500 {
		return errors.New("description must be 500 characters or less")
	}
	if !t.Status.IsValid() {
		return errors.New("invalid status")
	}
	if !t.Priority.IsValid() {
		return errors.New("invalid priority")
	}
	return nil
}

// String はTaskの文字列表現を返す
func (t *Task) String() string {
	return t.Title
}

// IsValid はStatusが有効かを検証する
func (s Status) IsValid() bool {
	switch s {
	case StatusTodo, StatusInProgress, StatusCompleted:
		return true
	default:
		return false
	}
}

// String はStatusの文字列表現を返す
func (s Status) String() string {
	return string(s)
}

// IsValid はPriorityが有効かを検証する
func (p Priority) IsValid() bool {
	switch p {
	case PriorityLow, PriorityMedium, PriorityHigh:
		return true
	default:
		return false
	}
}

// String はPriorityの文字列表現を返す
func (p Priority) String() string {
	return string(p)
}

// IsCompleted はタスクが完了しているかを返す
func (t *Task) IsCompleted() bool {
	return t.Status == StatusCompleted
}