package service

import (
	"sync"
	"testing"
	"time"

	"task-cli/internal/model"

	"github.com/stretchr/testify/assert"
)

// RED: StateManagerのテスト
func TestStateManager_New_ShouldCreateStateManager(t *testing.T) {
	// When
	stateManager := NewStateManager()

	// Then
	assert.NotNil(t, stateManager)
}

func TestStateManager_SetTasks_ShouldNotifySubscribers(t *testing.T) {
	// Given
	stateManager := NewStateManager()
	var receivedTasks []*model.Task
	var wg sync.WaitGroup
	wg.Add(1)

	// 通知を受け取るサブスクライバーを登録
	unsubscribe := stateManager.Subscribe(func(tasks []*model.Task, filter TaskFilter) {
		receivedTasks = tasks
		wg.Done()
	})
	defer unsubscribe()

	// When
	testTasks := []*model.Task{
		{ID: "1", Title: "Task 1"},
		{ID: "2", Title: "Task 2"},
	}
	stateManager.SetTasks(testTasks)

	// Then
	wg.Wait()
	assert.Len(t, receivedTasks, 2)
	assert.Equal(t, "Task 1", receivedTasks[0].Title)
}

func TestStateManager_Subscribe_ShouldReceiveNotifications(t *testing.T) {
	// Given
	stateManager := NewStateManager()
	notificationCount := 0
	var wg sync.WaitGroup

	// When
	wg.Add(2) // 2回の通知を期待
	unsubscribe := stateManager.Subscribe(func(tasks []*model.Task, filter TaskFilter) {
		notificationCount++
		wg.Done()
	})
	defer unsubscribe()

	// 2回の変更を通知
	stateManager.SetTasks([]*model.Task{{ID: "1", Title: "Task 1"}})
	stateManager.SetTasks([]*model.Task{{ID: "2", Title: "Task 2"}})

	// Then
	wg.Wait()
	assert.Equal(t, 2, notificationCount)
}

func TestStateManager_SetFilter_ShouldUpdateFilteredTasks(t *testing.T) {
	// Given
	stateManager := NewStateManager()
	var receivedFilter TaskFilter
	var wg sync.WaitGroup
	wg.Add(1)

	// サブスクライバーを登録
	unsubscribe := stateManager.Subscribe(func(tasks []*model.Task, filter TaskFilter) {
		receivedFilter = filter
		wg.Done()
	})
	defer unsubscribe()

	// When
	filter := TaskFilter{
		Status:   &[]model.Status{model.StatusTodo}[0],
		Priority: &[]model.Priority{model.PriorityHigh}[0],
		Query:    "test",
	}
	stateManager.SetFilter(filter)

	// Then
	wg.Wait()
	assert.Equal(t, model.StatusTodo, *receivedFilter.Status)
	assert.Equal(t, model.PriorityHigh, *receivedFilter.Priority)
	assert.Equal(t, "test", receivedFilter.Query)
}

func TestStateManager_GetCurrentTasks_ShouldReturnCurrentTasks(t *testing.T) {
	// Given
	stateManager := NewStateManager()
	testTasks := []*model.Task{
		{ID: "1", Title: "Task 1"},
		{ID: "2", Title: "Task 2"},
	}
	stateManager.SetTasks(testTasks)

	// When
	currentTasks := stateManager.GetCurrentTasks()

	// Then
	assert.Len(t, currentTasks, 2)
	assert.Equal(t, "Task 1", currentTasks[0].Title)
}

func TestStateManager_GetCurrentFilter_ShouldReturnCurrentFilter(t *testing.T) {
	// Given
	stateManager := NewStateManager()
	filter := TaskFilter{
		Status: &[]model.Status{model.StatusInProgress}[0],
		Query:  "important",
	}
	stateManager.SetFilter(filter)

	// When
	currentFilter := stateManager.GetCurrentFilter()

	// Then
	assert.Equal(t, model.StatusInProgress, *currentFilter.Status)
	assert.Equal(t, "important", currentFilter.Query)
}

func TestStateManager_ConcurrentAccess_ShouldBeSafe(t *testing.T) {
	// Given
	stateManager := NewStateManager()
	var wg sync.WaitGroup
	const goroutines = 10
	const operations = 100

	// When - 複数のgoroutineで同時にアクセス
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < operations; j++ {
				// タスクの設定と取得を繰り返す
				tasks := []*model.Task{{ID: string(rune(id)), Title: "Task"}}
				stateManager.SetTasks(tasks)
				stateManager.GetCurrentTasks()
				
				// フィルターの設定と取得を繰り返す
				filter := TaskFilter{Query: string(rune(id))}
				stateManager.SetFilter(filter)
				stateManager.GetCurrentFilter()
			}
		}(i)
	}

	// Then
	wg.Wait() // デッドロックがないことを確認
	assert.True(t, true) // ここに到達すればOK
}

func TestStateManager_Unsubscribe_ShouldStopReceivingNotifications(t *testing.T) {
	// Given
	stateManager := NewStateManager()
	notificationCount := 0
	
	// When
	unsubscribe := stateManager.Subscribe(func(tasks []*model.Task, filter TaskFilter) {
		notificationCount++
	})
	
	// 最初の通知
	stateManager.SetTasks([]*model.Task{{ID: "1", Title: "Task 1"}})
	time.Sleep(10 * time.Millisecond) // 通知の処理を待つ
	assert.Equal(t, 1, notificationCount)
	
	// サブスクライブ解除
	unsubscribe()
	
	// 2回目の通知（受け取られないはず）
	stateManager.SetTasks([]*model.Task{{ID: "2", Title: "Task 2"}})
	time.Sleep(10 * time.Millisecond) // 通知の処理を待つ
	
	// Then
	assert.Equal(t, 1, notificationCount) // まだ1のまま
}

func TestStateManager_ApplyFilter_ShouldReturnFilteredTasks(t *testing.T) {
	// Given
	stateManager := NewStateManager()
	
	task1, _ := model.NewTask("Buy milk", "Need to buy milk", model.PriorityHigh, []string{"shopping"})
	task2, _ := model.NewTask("Fix bug", "Fix login issue", model.PriorityMedium, []string{"dev"})
	task3, _ := model.NewTask("Meeting", "Team meeting", model.PriorityLow, []string{"work"})
	task2.Status = model.StatusInProgress
	
	allTasks := []*model.Task{task1, task2, task3}
	stateManager.SetTasks(allTasks)

	tests := []struct {
		name     string
		filter   TaskFilter
		expected int
	}{
		{
			name:     "Filter by status Todo",
			filter:   TaskFilter{Status: &[]model.Status{model.StatusTodo}[0]},
			expected: 2, // task1, task3
		},
		{
			name:     "Filter by priority High",
			filter:   TaskFilter{Priority: &[]model.Priority{model.PriorityHigh}[0]},
			expected: 1, // task1
		},
		{
			name:     "Filter by query",
			filter:   TaskFilter{Query: "bug"},
			expected: 1, // task2
		},
		{
			name:     "No filter",
			filter:   TaskFilter{},
			expected: 3, // all tasks
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// When
			filtered := stateManager.ApplyFilter(allTasks, tt.filter)

			// Then
			assert.Len(t, filtered, tt.expected)
		})
	}
}