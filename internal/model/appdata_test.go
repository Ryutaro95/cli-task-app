package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// RED: AppDataモデルのテスト
func TestAppData_NewAppData_ShouldReturnValidInstance(t *testing.T) {
	// When
	appData := NewAppData()

	// Then
	assert.NotNil(t, appData)
	assert.Empty(t, appData.Tasks)
	assert.NotEmpty(t, appData.ID)
	assert.False(t, appData.CreatedAt.IsZero())
	assert.False(t, appData.UpdatedAt.IsZero())
}

func TestAppData_AddTask_ShouldIncreaseTaskCount(t *testing.T) {
	// Given
	appData := NewAppData()
	task, err := NewTask("Test Task", "Description", PriorityMedium, []string{"test"})
	assert.NoError(t, err)

	// When
	err = appData.AddTask(task)

	// Then
	assert.NoError(t, err)
	assert.Len(t, appData.Tasks, 1)
	assert.Equal(t, task.ID, appData.Tasks[0].ID)
}

func TestAppData_GetTaskByID_WithExistentID_ShouldReturnTask(t *testing.T) {
	// Given
	appData := NewAppData()
	task, err := NewTask("Test Task", "Description", PriorityMedium, []string{"test"})
	assert.NoError(t, err)
	appData.AddTask(task)

	// When
	retrievedTask, err := appData.GetTaskByID(task.ID)

	// Then
	assert.NoError(t, err)
	assert.NotNil(t, retrievedTask)
	assert.Equal(t, task.ID, retrievedTask.ID)
}

func TestAppData_GetTaskByID_WithNonexistentID_ShouldReturnError(t *testing.T) {
	// Given
	appData := NewAppData()
	nonExistentID := "non-existent-id"

	// When
	retrievedTask, err := appData.GetTaskByID(nonExistentID)

	// Then
	assert.Error(t, err)
	assert.Nil(t, retrievedTask)
	assert.Contains(t, err.Error(), "not found")
}

func TestAppData_UpdateTask_WithValidTask_ShouldUpdateTask(t *testing.T) {
	// Given
	appData := NewAppData()
	task, err := NewTask("Test Task", "Description", PriorityMedium, []string{"test"})
	assert.NoError(t, err)
	appData.AddTask(task)

	// 変更するタスク
	updatedTask := *task
	updatedTask.Title = "Updated Task"
	updatedTask.Priority = PriorityHigh

	// When
	err = appData.UpdateTask(&updatedTask)

	// Then
	assert.NoError(t, err)
	retrievedTask, err := appData.GetTaskByID(task.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Task", retrievedTask.Title)
	assert.Equal(t, PriorityHigh, retrievedTask.Priority)
}

func TestAppData_DeleteTask_WithValidID_ShouldRemoveTask(t *testing.T) {
	// Given
	appData := NewAppData()
	task, err := NewTask("Test Task", "Description", PriorityMedium, []string{"test"})
	assert.NoError(t, err)
	appData.AddTask(task)
	assert.Len(t, appData.Tasks, 1)

	// When
	err = appData.DeleteTask(task.ID)

	// Then
	assert.NoError(t, err)
	assert.Len(t, appData.Tasks, 0)
}

func TestAppData_DeleteTask_WithNonexistentID_ShouldReturnError(t *testing.T) {
	// Given
	appData := NewAppData()
	nonExistentID := "non-existent-id"

	// When
	err := appData.DeleteTask(nonExistentID)

	// Then
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestAppData_GetTasksByStatus_ShouldReturnFilteredTasks(t *testing.T) {
	// Given
	appData := NewAppData()
	
	// 異なるステータスのタスクを作成
	task1, _ := NewTask("Task 1", "Desc 1", PriorityLow, []string{})
	task2, _ := NewTask("Task 2", "Desc 2", PriorityMedium, []string{})
	task3, _ := NewTask("Task 3", "Desc 3", PriorityHigh, []string{})
	
	// task2のステータスを変更
	task2.Status = StatusInProgress
	
	appData.AddTask(task1)
	appData.AddTask(task2)
	appData.AddTask(task3)

	// When
	todoTasks := appData.GetTasksByStatus(StatusTodo)
	inProgressTasks := appData.GetTasksByStatus(StatusInProgress)

	// Then
	assert.Len(t, todoTasks, 2) // task1とtask3
	assert.Len(t, inProgressTasks, 1) // task2
}

func TestAppData_GetTasksByPriority_ShouldReturnFilteredTasks(t *testing.T) {
	// Given
	appData := NewAppData()
	
	task1, _ := NewTask("Task 1", "Desc 1", PriorityLow, []string{})
	task2, _ := NewTask("Task 2", "Desc 2", PriorityHigh, []string{})
	task3, _ := NewTask("Task 3", "Desc 3", PriorityHigh, []string{})
	
	appData.AddTask(task1)
	appData.AddTask(task2)
	appData.AddTask(task3)

	// When
	highPriorityTasks := appData.GetTasksByPriority(PriorityHigh)
	lowPriorityTasks := appData.GetTasksByPriority(PriorityLow)

	// Then
	assert.Len(t, highPriorityTasks, 2) // task2とtask3
	assert.Len(t, lowPriorityTasks, 1) // task1
}

func TestAppData_GetTaskCount_ShouldReturnCorrectCount(t *testing.T) {
	// Given
	appData := NewAppData()
	task1, _ := NewTask("Task 1", "Desc 1", PriorityLow, []string{})
	task2, _ := NewTask("Task 2", "Desc 2", PriorityMedium, []string{})

	// When & Then - 初期状態
	assert.Equal(t, 0, appData.GetTaskCount())

	// タスク追加後
	appData.AddTask(task1)
	assert.Equal(t, 1, appData.GetTaskCount())

	appData.AddTask(task2)
	assert.Equal(t, 2, appData.GetTaskCount())
}

func TestAppData_GetCompletedTaskCount_ShouldReturnCorrectCount(t *testing.T) {
	// Given
	appData := NewAppData()
	task1, _ := NewTask("Task 1", "Desc 1", PriorityLow, []string{})
	task2, _ := NewTask("Task 2", "Desc 2", PriorityMedium, []string{})
	
	// task2を完了にする
	task2.Status = StatusCompleted
	
	appData.AddTask(task1)
	appData.AddTask(task2)

	// When & Then
	assert.Equal(t, 1, appData.GetCompletedTaskCount())
}

func TestAppData_GetActiveTaskCount_ShouldReturnCorrectCount(t *testing.T) {
	// Given
	appData := NewAppData()
	task1, _ := NewTask("Task 1", "Desc 1", PriorityLow, []string{})
	task2, _ := NewTask("Task 2", "Desc 2", PriorityMedium, []string{})
	task3, _ := NewTask("Task 3", "Desc 3", PriorityHigh, []string{})
	
	// task3を完了にする
	task3.Status = StatusCompleted
	
	appData.AddTask(task1)
	appData.AddTask(task2)
	appData.AddTask(task3)

	// When & Then
	assert.Equal(t, 3, appData.GetTaskCount())
	assert.Equal(t, 1, appData.GetCompletedTaskCount())
	assert.Equal(t, 2, appData.GetActiveTaskCount())
}

func TestAppData_SearchTasksByTitle_ShouldReturnMatchingTasks(t *testing.T) {
	// Given
	appData := NewAppData()
	task1, _ := NewTask("Buy groceries", "Need to buy milk and bread", PriorityMedium, []string{})
	task2, _ := NewTask("Fix Bug", "Fix login issue in the application", PriorityHigh, []string{})
	task3, _ := NewTask("Meeting", "Team meeting about project", PriorityLow, []string{})
	
	appData.AddTask(task1)
	appData.AddTask(task2)
	appData.AddTask(task3)

	tests := []struct {
		query    string
		expected int
	}{
		{"Buy", 1},      // task1にマッチ
		{"bug", 1},      // task2にマッチ（大文字小文字無視）
		{"team", 1},     // task3の説明にマッチ
		{"project", 1},  // task3の説明にマッチ
		{"xyz", 0},      // マッチしない
		{"", 0},         // 空文字列
	}

	for _, tt := range tests {
		t.Run(tt.query, func(t *testing.T) {
			// When
			results := appData.SearchTasksByTitle(tt.query)

			// Then
			assert.Len(t, results, tt.expected)
		})
	}
}