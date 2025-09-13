package validator

import (
	"errors"
	"fmt"

	"task-cli/internal/model"
)

// Validator はデータのバリデーションを行う
type Validator struct {
	maxTitleLength       int
	maxDescriptionLength int
}

// New は新しいValidatorを作成する
func New() *Validator {
	return &Validator{
		maxTitleLength:       100,
		maxDescriptionLength: 500,
	}
}

// ValidateTask はTaskの値を検証する
func (v *Validator) ValidateTask(task *model.Task) error {
	if task == nil {
		return errors.New("task cannot be nil")
	}

	// タイトルの検証
	if task.Title == "" {
		return errors.New("title is required")
	}
	if len(task.Title) > v.maxTitleLength {
		return fmt.Errorf("title must be %d characters or less", v.maxTitleLength)
	}

	// 説明の検証
	if len(task.Description) > v.maxDescriptionLength {
		return fmt.Errorf("description must be %d characters or less", v.maxDescriptionLength)
	}

	// ステータスの検証
	if !task.Status.IsValid() {
		return errors.New("invalid status")
	}

	// 優先度の検証
	if !task.Priority.IsValid() {
		return errors.New("invalid priority")
	}

	// IDの検証
	if task.ID == "" {
		return errors.New("id is required")
	}

	// 作成日時の検証
	if task.CreatedAt.IsZero() {
		return errors.New("created_at is required")
	}

	// 更新日時の検証
	if task.UpdatedAt.IsZero() {
		return errors.New("updated_at is required")
	}

	// 完了日時の検証（完了ステータスの場合）
	if task.Status == model.StatusCompleted && task.CompletedAt == nil {
		return errors.New("completed_at is required when status is completed")
	}

	return nil
}

// ValidateAppData はAppDataの値を検証する
func (v *Validator) ValidateAppData(appData *model.AppData) error {
	if appData == nil {
		return errors.New("appdata cannot be nil")
	}

	// IDの検証
	if appData.ID == "" {
		return errors.New("appdata id is required")
	}

	// 作成日時の検証
	if appData.CreatedAt.IsZero() {
		return errors.New("appdata created_at is required")
	}

	// 更新日時の検証
	if appData.UpdatedAt.IsZero() {
		return errors.New("appdata updated_at is required")
	}

	// 各タスクの検証
	for i, task := range appData.Tasks {
		if err := v.ValidateTask(task); err != nil {
			return fmt.Errorf("task at index %d is invalid: %w", i, err)
		}
	}

	return nil
}