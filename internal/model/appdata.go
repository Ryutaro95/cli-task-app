package model

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

// AppData はアプリケーションのデータ全体を管理する構造体
type AppData struct {
	ID        string     `json:"id"`
	Tasks     []*Task    `json:"tasks"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// NewAppData は新しいAppDataインスタンスを作成する
func NewAppData() *AppData {
	now := time.Now()
	return &AppData{
		ID:        uuid.New().String(),
		Tasks:     make([]*Task, 0),
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// AddTask はタスクを追加する
func (a *AppData) AddTask(task *Task) error {
	if task == nil {
		return errors.New("task cannot be nil")
	}
	
	a.Tasks = append(a.Tasks, task)
	a.UpdatedAt = time.Now()
	return nil
}

// GetTaskByID はIDでタスクを取得する
func (a *AppData) GetTaskByID(id string) (*Task, error) {
	for _, task := range a.Tasks {
		if task.ID == id {
			return task, nil
		}
	}
	return nil, errors.New("task not found")
}

// UpdateTask はタスクを更新する
func (a *AppData) UpdateTask(updatedTask *Task) error {
	if updatedTask == nil {
		return errors.New("updated task cannot be nil")
	}

	for i, task := range a.Tasks {
		if task.ID == updatedTask.ID {
			a.Tasks[i] = updatedTask
			a.UpdatedAt = time.Now()
			return nil
		}
	}
	return errors.New("task not found")
}

// DeleteTask はタスクを削除する
func (a *AppData) DeleteTask(id string) error {
	for i, task := range a.Tasks {
		if task.ID == id {
			// スライスから要素を削除
			a.Tasks = append(a.Tasks[:i], a.Tasks[i+1:]...)
			a.UpdatedAt = time.Now()
			return nil
		}
	}
	return errors.New("task not found")
}

// GetTasksByStatus は指定されたステータスのタスクを返す
func (a *AppData) GetTasksByStatus(status Status) []*Task {
	var filteredTasks []*Task
	for _, task := range a.Tasks {
		if task.Status == status {
			filteredTasks = append(filteredTasks, task)
		}
	}
	return filteredTasks
}

// GetTasksByPriority は指定された優先度のタスクを返す
func (a *AppData) GetTasksByPriority(priority Priority) []*Task {
	var filteredTasks []*Task
	for _, task := range a.Tasks {
		if task.Priority == priority {
			filteredTasks = append(filteredTasks, task)
		}
	}
	return filteredTasks
}

// GetTaskCount は総タスク数を返す
func (a *AppData) GetTaskCount() int {
	return len(a.Tasks)
}

// GetCompletedTaskCount は完了済みタスク数を返す
func (a *AppData) GetCompletedTaskCount() int {
	count := 0
	for _, task := range a.Tasks {
		if task.IsCompleted() {
			count++
		}
	}
	return count
}

// GetActiveTaskCount はアクティブ（未完了）タスク数を返す
func (a *AppData) GetActiveTaskCount() int {
	return a.GetTaskCount() - a.GetCompletedTaskCount()
}

// SearchTasksByTitle はタイトルまたは説明でタスクを検索する
func (a *AppData) SearchTasksByTitle(query string) []*Task {
	if query == "" {
		return []*Task{}
	}

	var matchedTasks []*Task
	lowerQuery := strings.ToLower(query)
	
	for _, task := range a.Tasks {
		// 大文字小文字を無視してタイトルと説明を検索
		if strings.Contains(strings.ToLower(task.Title), lowerQuery) ||
		   strings.Contains(strings.ToLower(task.Description), lowerQuery) {
			matchedTasks = append(matchedTasks, task)
		}
	}
	return matchedTasks
}