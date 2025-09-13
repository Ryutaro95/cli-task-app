package service

import (
	"strings"
	"sync"

	"task-cli/internal/model"
)

// TaskFilter はタスクのフィルタリング条件を定義
type TaskFilter struct {
	Status   *model.Status
	Priority *model.Priority
	Query    string
}

// SubscriberFunc は状態変更の通知を受け取る関数の型
type SubscriberFunc func(tasks []*model.Task, filter TaskFilter)

// StateManager はアプリケーションの状態を管理し、変更を通知する
type StateManager struct {
	mu          sync.RWMutex
	tasks       []*model.Task
	filter      TaskFilter
	subscribers map[int]SubscriberFunc
	nextID      int
}

// NewStateManager は新しいStateManagerを作成する
func NewStateManager() *StateManager {
	return &StateManager{
		tasks:       make([]*model.Task, 0),
		filter:      TaskFilter{},
		subscribers: make(map[int]SubscriberFunc),
		nextID:      1,
	}
}

// SetTasks は現在のタスクリストを設定し、サブスクライバーに通知する
func (sm *StateManager) SetTasks(tasks []*model.Task) {
	sm.mu.Lock()
	sm.tasks = make([]*model.Task, len(tasks))
	copy(sm.tasks, tasks)
	currentFilter := sm.filter
	subscribers := make(map[int]SubscriberFunc, len(sm.subscribers))
	for id, sub := range sm.subscribers {
		subscribers[id] = sub
	}
	sm.mu.Unlock()

	// フィルターを適用
	filteredTasks := sm.ApplyFilter(tasks, currentFilter)

	// サブスクライバーに通知（ロックの外で実行）
	for _, subscriber := range subscribers {
		go subscriber(filteredTasks, currentFilter)
	}
}

// SetFilter はフィルターを設定し、サブスクライバーに通知する
func (sm *StateManager) SetFilter(filter TaskFilter) {
	sm.mu.Lock()
	sm.filter = filter
	currentTasks := make([]*model.Task, len(sm.tasks))
	copy(currentTasks, sm.tasks)
	subscribers := make(map[int]SubscriberFunc, len(sm.subscribers))
	for id, sub := range sm.subscribers {
		subscribers[id] = sub
	}
	sm.mu.Unlock()

	// フィルターを適用
	filteredTasks := sm.ApplyFilter(currentTasks, filter)

	// サブスクライバーに通知（ロックの外で実行）
	for _, subscriber := range subscribers {
		go subscriber(filteredTasks, filter)
	}
}

// Subscribe は状態変更の通知を受け取るサブスクライバーを登録する
// 戻り値の関数を呼び出すことで購読を解除できる
func (sm *StateManager) Subscribe(subscriber SubscriberFunc) func() {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	id := sm.nextID
	sm.nextID++
	sm.subscribers[id] = subscriber

	// サブスクライブ解除用の関数を返す
	return func() {
		sm.mu.Lock()
		defer sm.mu.Unlock()
		delete(sm.subscribers, id)
	}
}

// GetCurrentTasks は現在のタスクリストを取得する
func (sm *StateManager) GetCurrentTasks() []*model.Task {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	tasks := make([]*model.Task, len(sm.tasks))
	copy(tasks, sm.tasks)
	return tasks
}

// GetCurrentFilter は現在のフィルターを取得する
func (sm *StateManager) GetCurrentFilter() TaskFilter {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	return sm.filter
}

// ApplyFilter はタスクリストにフィルターを適用する
func (sm *StateManager) ApplyFilter(tasks []*model.Task, filter TaskFilter) []*model.Task {
	var filtered []*model.Task

	for _, task := range tasks {
		// ステータスフィルター
		if filter.Status != nil && task.Status != *filter.Status {
			continue
		}

		// 優先度フィルター
		if filter.Priority != nil && task.Priority != *filter.Priority {
			continue
		}

		// クエリフィルター（タイトルと説明で部分一致検索）
		if filter.Query != "" {
			query := strings.ToLower(filter.Query)
			title := strings.ToLower(task.Title)
			description := strings.ToLower(task.Description)
			
			if !strings.Contains(title, query) && !strings.Contains(description, query) {
				continue
			}
		}

		filtered = append(filtered, task)
	}

	return filtered
}

// GetFilteredTasks は現在のフィルターを適用したタスクリストを取得する
func (sm *StateManager) GetFilteredTasks() []*model.Task {
	sm.mu.RLock()
	tasks := make([]*model.Task, len(sm.tasks))
	copy(tasks, sm.tasks)
	filter := sm.filter
	sm.mu.RUnlock()

	return sm.ApplyFilter(tasks, filter)
}

// ClearFilter はフィルターをクリアする
func (sm *StateManager) ClearFilter() {
	sm.SetFilter(TaskFilter{})
}

// GetTaskCount は現在のタスク数を取得する
func (sm *StateManager) GetTaskCount() int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return len(sm.tasks)
}

// GetFilteredTaskCount はフィルター適用後のタスク数を取得する
func (sm *StateManager) GetFilteredTaskCount() int {
	filtered := sm.GetFilteredTasks()
	return len(filtered)
}