package ui

import (
	"testing"

	"task-cli/internal/model"
	"task-cli/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/rivo/tview"
)

// RED: TaskListWidgetのテスト
func TestTaskListWidget_New_ShouldCreateWidget(t *testing.T) {
	// Given
	theme := NewTheme()

	// When
	widget := NewTaskListWidget(theme)

	// Then
	assert.NotNil(t, widget)
	assert.Implements(t, (*tview.Primitive)(nil), widget.GetPrimitive())
}

func TestTaskListWidget_SetTasks_ShouldDisplayTasks(t *testing.T) {
	// Given
	theme := NewTheme()
	widget := NewTaskListWidget(theme)
	
	task1, _ := model.NewTask("Task 1", "Description 1", model.PriorityHigh, []string{"test"})
	task2, _ := model.NewTask("Task 2", "Description 2", model.PriorityMedium, []string{"test"})
	tasks := []*model.Task{task1, task2}

	// When
	widget.SetTasks(tasks)

	// Then
	assert.Equal(t, 2, widget.GetTaskCount())
}

func TestTaskListWidget_GetSelectedTask_ShouldReturnSelectedTask(t *testing.T) {
	// Given
	theme := NewTheme()
	widget := NewTaskListWidget(theme)
	
	task1, _ := model.NewTask("Task 1", "Description 1", model.PriorityHigh, []string{"test"})
	task2, _ := model.NewTask("Task 2", "Description 2", model.PriorityMedium, []string{"test"})
	tasks := []*model.Task{task1, task2}
	widget.SetTasks(tasks)

	// When
	widget.SelectTask(0) // 最初のタスクを選択
	selectedTask := widget.GetSelectedTask()

	// Then
	assert.NotNil(t, selectedTask)
	assert.Equal(t, "Task 1", selectedTask.Title)
}

func TestTaskListWidget_SelectNext_ShouldMoveSelection(t *testing.T) {
	// Given
	theme := NewTheme()
	widget := NewTaskListWidget(theme)
	
	task1, _ := model.NewTask("Task 1", "Description 1", model.PriorityHigh, []string{"test"})
	task2, _ := model.NewTask("Task 2", "Description 2", model.PriorityMedium, []string{"test"})
	tasks := []*model.Task{task1, task2}
	widget.SetTasks(tasks)

	// When
	widget.SelectTask(0) // 最初のタスクを選択
	initialTask := widget.GetSelectedTask()
	
	widget.SelectNext() // 次のタスクに移動
	nextTask := widget.GetSelectedTask()

	// Then
	assert.NotEqual(t, initialTask.ID, nextTask.ID)
	assert.Equal(t, "Task 2", nextTask.Title)
}

func TestTaskListWidget_SelectPrevious_ShouldMoveSelection(t *testing.T) {
	// Given
	theme := NewTheme()
	widget := NewTaskListWidget(theme)
	
	task1, _ := model.NewTask("Task 1", "Description 1", model.PriorityHigh, []string{"test"})
	task2, _ := model.NewTask("Task 2", "Description 2", model.PriorityMedium, []string{"test"})
	tasks := []*model.Task{task1, task2}
	widget.SetTasks(tasks)

	// When
	widget.SelectTask(1) // 2番目のタスクを選択
	initialTask := widget.GetSelectedTask()
	
	widget.SelectPrevious() // 前のタスクに移動
	previousTask := widget.GetSelectedTask()

	// Then
	assert.NotEqual(t, initialTask.ID, previousTask.ID)
	assert.Equal(t, "Task 1", previousTask.Title)
}

func TestTaskListWidget_ApplyFilter_ShouldShowOnlyMatchingTasks(t *testing.T) {
	// Given
	theme := NewTheme()
	widget := NewTaskListWidget(theme)
	
	task1, _ := model.NewTask("Buy groceries", "Need to buy milk", model.PriorityHigh, []string{"shopping"})
	task2, _ := model.NewTask("Fix bug", "Fix login issue", model.PriorityMedium, []string{"dev"})
	task3, _ := model.NewTask("Meeting", "Team meeting", model.PriorityLow, []string{"work"})
	task2.Status = model.StatusInProgress
	tasks := []*model.Task{task1, task2, task3}
	widget.SetTasks(tasks)

	// When - ステータスフィルターを適用
	filter := service.TaskFilter{Status: &[]model.Status{model.StatusTodo}[0]}
	widget.ApplyFilter(filter)

	// Then
	filteredCount := widget.GetTaskCount()
	assert.Equal(t, 2, filteredCount) // task1とtask3のみ表示される
}

func TestTaskListWidget_ClearFilter_ShouldShowAllTasks(t *testing.T) {
	// Given
	theme := NewTheme()
	widget := NewTaskListWidget(theme)
	
	task1, _ := model.NewTask("Task 1", "Description 1", model.PriorityHigh, []string{"test"})
	task2, _ := model.NewTask("Task 2", "Description 2", model.PriorityMedium, []string{"test"})
	task2.Status = model.StatusInProgress
	tasks := []*model.Task{task1, task2}
	widget.SetTasks(tasks)

	// フィルターを適用
	filter := service.TaskFilter{Status: &[]model.Status{model.StatusTodo}[0]}
	widget.ApplyFilter(filter)
	assert.Equal(t, 1, widget.GetTaskCount())

	// When - フィルターをクリア
	widget.ClearFilter()

	// Then
	assert.Equal(t, 2, widget.GetTaskCount()) // 全タスクが表示される
}

func TestTaskListWidget_GetSelectedIndex_ShouldReturnCorrectIndex(t *testing.T) {
	// Given
	theme := NewTheme()
	widget := NewTaskListWidget(theme)
	
	task1, _ := model.NewTask("Task 1", "Description 1", model.PriorityHigh, []string{"test"})
	task2, _ := model.NewTask("Task 2", "Description 2", model.PriorityMedium, []string{"test"})
	tasks := []*model.Task{task1, task2}
	widget.SetTasks(tasks)

	// When
	widget.SelectTask(1)
	selectedIndex := widget.GetSelectedIndex()

	// Then
	assert.Equal(t, 1, selectedIndex)
}

func TestTaskListWidget_SortByPriority_ShouldOrderTasks(t *testing.T) {
	// Given
	theme := NewTheme()
	widget := NewTaskListWidget(theme)
	
	task1, _ := model.NewTask("Low Priority Task", "Description", model.PriorityLow, []string{"test"})
	task2, _ := model.NewTask("High Priority Task", "Description", model.PriorityHigh, []string{"test"})
	task3, _ := model.NewTask("Medium Priority Task", "Description", model.PriorityMedium, []string{"test"})
	tasks := []*model.Task{task1, task2, task3}
	widget.SetTasks(tasks)

	// When
	widget.SortByPriority()

	// Then
	// 優先度順（High -> Medium -> Low）にソートされる
	widget.SelectTask(0)
	firstTask := widget.GetSelectedTask()
	assert.Equal(t, model.PriorityHigh, firstTask.Priority)
}

func TestTaskListWidget_SortByStatus_ShouldOrderTasks(t *testing.T) {
	// Given
	theme := NewTheme()
	widget := NewTaskListWidget(theme)
	
	task1, _ := model.NewTask("Completed Task", "Description", model.PriorityMedium, []string{"test"})
	task2, _ := model.NewTask("Todo Task", "Description", model.PriorityMedium, []string{"test"})
	task3, _ := model.NewTask("In Progress Task", "Description", model.PriorityMedium, []string{"test"})
	
	task1.Status = model.StatusCompleted
	task3.Status = model.StatusInProgress
	
	tasks := []*model.Task{task1, task2, task3}
	widget.SetTasks(tasks)

	// When
	widget.SortByStatus()

	// Then
	// ステータス順（Todo -> InProgress -> Completed）にソートされる
	widget.SelectTask(0)
	firstTask := widget.GetSelectedTask()
	assert.Equal(t, model.StatusTodo, firstTask.Status)
}

func TestTaskListWidget_SetSelectionChangedCallback_ShouldInvokeCallback(t *testing.T) {
	// Given
	theme := NewTheme()
	widget := NewTaskListWidget(theme)
	
	task1, _ := model.NewTask("Task 1", "Description 1", model.PriorityHigh, []string{"test"})
	tasks := []*model.Task{task1}
	widget.SetTasks(tasks)

	callbackInvoked := false
	var selectedTask *model.Task
	
	// When
	widget.SetSelectionChangedCallback(func(task *model.Task) {
		callbackInvoked = true
		selectedTask = task
	})
	
	widget.SelectTask(0)

	// Then
	assert.True(t, callbackInvoked)
	assert.NotNil(t, selectedTask)
	assert.Equal(t, "Task 1", selectedTask.Title)
}