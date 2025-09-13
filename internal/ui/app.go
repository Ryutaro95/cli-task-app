package ui

import (
	"context"
	"fmt"

	"task-cli/internal/model"
	"task-cli/internal/service"

	"github.com/rivo/tview"
	"github.com/gdamore/tcell/v2"
)

// ViewMode はアプリケーションのビューモードを定義
type ViewMode int

const (
	ViewModeList ViewMode = iota
	ViewModeForm
)

// TaskServiceInterface はTaskServiceのインターフェース
type TaskServiceInterface interface {
	CreateTask(ctx context.Context, request service.CreateTaskRequest) (*model.Task, error)
	UpdateTask(ctx context.Context, request service.UpdateTaskRequest) (*model.Task, error)
	DeleteTask(ctx context.Context, taskID string) error
	ToggleTaskStatus(ctx context.Context, taskID string) (*model.Task, error)
	GetAllTasks(ctx context.Context) ([]*model.Task, error)
	SearchTasks(ctx context.Context, query string) ([]*model.Task, error)
	GetTasksByStatus(ctx context.Context, status model.Status) ([]*model.Task, error)
	GetTasksByPriority(ctx context.Context, priority model.Priority) ([]*model.Task, error)
}

// App はメインアプリケーション
type App struct {
	// Core components
	tviewApp     *tview.Application
	taskService  TaskServiceInterface
	stateManager *service.StateManager
	theme        *Theme
	
	// UI Components
	taskListWidget *TaskListWidget
	inputFormWidget *InputFormWidget
	pages          *tview.Pages
	
	// State
	currentView    ViewMode
	editingTaskID  string
	ctx           context.Context
}

// NewApp は新しいAppを作成する
func NewApp(taskService TaskServiceInterface, stateManager *service.StateManager, theme *Theme) *App {
	app := &App{
		tviewApp:     tview.NewApplication(),
		taskService:  taskService,
		stateManager: stateManager,
		theme:        theme,
		currentView:  ViewModeList,
		ctx:         context.Background(),
	}
	
	app.setupUI()
	app.setupEventHandlers()
	
	return app
}

// Initialize はアプリケーションを初期化する
func (a *App) Initialize() error {
	return a.RefreshTasks()
}

// setupUI はUIコンポーネントを設定する
func (a *App) setupUI() {
	// ウィジェットを作成
	a.taskListWidget = NewTaskListWidget(a.theme)
	a.inputFormWidget = NewInputFormWidget(a.theme)
	
	// ページコンテナを作成
	a.pages = tview.NewPages()
	
	// リストビューを作成
	listLayout := a.createListLayout()
	a.pages.AddPage("list", listLayout, true, true)
	
	// フォームビューを作成
	formLayout := a.createFormLayout()
	a.pages.AddPage("form", formLayout, true, false)
	
	// メインレイアウトを設定
	a.tviewApp.SetRoot(a.pages, true)
}

// createListLayout はリストビューのレイアウトを作成する
func (a *App) createListLayout() tview.Primitive {
	// ヘルプテキストを作成
	helpText := tview.NewTextView().
		SetText("Keys: n=New, e=Edit, d=Delete, t=Toggle, q=Quit, /=Search").
		SetTextColor(a.theme.GetHighlightColor()).
		SetBackgroundColor(a.theme.GetBackgroundColor())
	
	// ボーダーを作成
	border := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(a.taskListWidget.GetPrimitive(), 0, 1, true).
		AddItem(helpText, 1, 0, false)
	
	border.SetBackgroundColor(a.theme.GetBackgroundColor())
	
	return border
}

// createFormLayout はフォームビューのレイアウトを作成する
func (a *App) createFormLayout() tview.Primitive {
	// ヘルプテキストを作成
	helpText := tview.NewTextView().
		SetText("Keys: Ctrl+S=Submit, Escape=Cancel").
		SetTextColor(a.theme.GetHighlightColor()).
		SetBackgroundColor(a.theme.GetBackgroundColor())
	
	// ボーダーを作成
	border := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(a.inputFormWidget.GetPrimitive(), 0, 1, true).
		AddItem(helpText, 1, 0, false)
	
	border.SetBackgroundColor(a.theme.GetBackgroundColor())
	
	return border
}

// setupEventHandlers はイベントハンドラーを設定する
func (a *App) setupEventHandlers() {
	// タスクリストの選択変更イベント
	a.taskListWidget.SetSelectionChangedCallback(func(task *model.Task) {
		// 必要に応じて詳細表示などを実装
	})
	
	// フォームのイベント
	a.inputFormWidget.SetSubmitCallback(func(data FormData) {
		a.handleFormSubmit(data)
	})
	
	a.inputFormWidget.SetCancelCallback(func() {
		a.handleFormCancel()
	})
	
	// StateManagerのイベント
	a.stateManager.Subscribe(func(tasks []*model.Task, filter service.TaskFilter) {
		a.taskListWidget.SetTasks(tasks)
		a.taskListWidget.ApplyFilter(filter)
	})
	
	// キーボードイベント
	a.tviewApp.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		return a.handleKeyPress(event)
	})
}

// handleKeyPress はキーボード入力を処理する
func (a *App) handleKeyPress(event *tcell.EventKey) *tcell.EventKey {
	switch a.currentView {
	case ViewModeList:
		return a.handleListViewKeyPress(event)
	case ViewModeForm:
		return a.handleFormViewKeyPress(event)
	}
	return event
}

// handleListViewKeyPress はリストビューでのキーボード入力を処理する
func (a *App) handleListViewKeyPress(event *tcell.EventKey) *tcell.EventKey {
	switch event.Rune() {
	case 'q':
		a.tviewApp.Stop()
		return nil
	case 'n':
		a.StartCreateTask()
		return nil
	case 'e':
		a.StartEditTask()
		return nil
	case 'd':
		a.DeleteSelectedTask()
		return nil
	case 't':
		a.ToggleSelectedTask()
		return nil
	case '/':
		// 検索機能は将来的に実装
		return nil
	}
	
	switch event.Key() {
	case tcell.KeyEscape:
		a.tviewApp.Stop()
		return nil
	}
	
	return event
}

// handleFormViewKeyPress はフォームビューでのキーボード入力を処理する
func (a *App) handleFormViewKeyPress(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyEscape:
		a.handleFormCancel()
		return nil
	case tcell.KeyCtrlS:
		data := a.inputFormWidget.GetFormData()
		a.handleFormSubmit(data)
		return nil
	}
	
	return event
}

// SwitchToListView はリストビューに切り替える
func (a *App) SwitchToListView() {
	a.currentView = ViewModeList
	a.pages.SwitchToPage("list")
}

// SwitchToFormView はフォームビューに切り替える
func (a *App) SwitchToFormView() {
	a.currentView = ViewModeForm
	a.pages.SwitchToPage("form")
}

// GetCurrentView は現在のビューモードを取得する
func (a *App) GetCurrentView() ViewMode {
	return a.currentView
}

// RefreshTasks はタスクリストを更新する
func (a *App) RefreshTasks() error {
	tasks, err := a.taskService.GetAllTasks(a.ctx)
	if err != nil {
		return fmt.Errorf("failed to refresh tasks: %w", err)
	}
	
	a.stateManager.SetTasks(tasks)
	return nil
}

// StartCreateTask は新しいタスク作成を開始する
func (a *App) StartCreateTask() {
	a.editingTaskID = ""
	a.inputFormWidget.SetMode(FormModeCreate)
	a.inputFormWidget.Clear()
	a.SwitchToFormView()
}

// StartEditTask は選択されたタスクの編集を開始する
func (a *App) StartEditTask() {
	selectedTask := a.taskListWidget.GetSelectedTask()
	if selectedTask == nil {
		return
	}
	
	a.editingTaskID = selectedTask.ID
	a.inputFormWidget.SetMode(FormModeEdit)
	a.inputFormWidget.LoadTask(selectedTask)
	a.SwitchToFormView()
}

// DeleteSelectedTask は選択されたタスクを削除する
func (a *App) DeleteSelectedTask() {
	selectedTask := a.taskListWidget.GetSelectedTask()
	if selectedTask == nil {
		return
	}
	
	if err := a.HandleDeleteTask(selectedTask.ID); err != nil {
		// エラー処理（将来的にダイアログで表示）
		return
	}
}

// ToggleSelectedTask は選択されたタスクのステータスを切り替える
func (a *App) ToggleSelectedTask() {
	selectedTask := a.taskListWidget.GetSelectedTask()
	if selectedTask == nil {
		return
	}
	
	if err := a.HandleToggleTask(selectedTask.ID); err != nil {
		// エラー処理（将来的にダイアログで表示）
		return
	}
}

// handleFormSubmit はフォーム送信を処理する
func (a *App) handleFormSubmit(data FormData) {
	var err error
	
	if a.editingTaskID == "" {
		// 新規作成
		err = a.HandleCreateTask(data)
	} else {
		// 更新
		err = a.HandleUpdateTask(a.editingTaskID, data)
	}
	
	if err != nil {
		a.inputFormWidget.SetErrorMessage(err.Error())
		return
	}
	
	a.SwitchToListView()
}

// handleFormCancel はフォームキャンセルを処理する
func (a *App) handleFormCancel() {
	a.inputFormWidget.Clear()
	a.SwitchToListView()
}

// HandleCreateTask は新しいタスクを作成する
func (a *App) HandleCreateTask(data FormData) error {
	request := service.CreateTaskRequest{
		Title:       data.Title,
		Description: data.Description,
		Priority:    data.Priority,
		Tags:        data.Tags,
	}
	
	_, err := a.taskService.CreateTask(a.ctx, request)
	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}
	
	return a.RefreshTasks()
}

// HandleUpdateTask は既存のタスクを更新する
func (a *App) HandleUpdateTask(taskID string, data FormData) error {
	request := service.UpdateTaskRequest{
		ID:          taskID,
		Title:       data.Title,
		Description: data.Description,
		Priority:    data.Priority,
		Status:      data.Status,
		Tags:        data.Tags,
	}
	
	_, err := a.taskService.UpdateTask(a.ctx, request)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}
	
	return a.RefreshTasks()
}

// HandleDeleteTask はタスクを削除する
func (a *App) HandleDeleteTask(taskID string) error {
	err := a.taskService.DeleteTask(a.ctx, taskID)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}
	
	return a.RefreshTasks()
}

// HandleToggleTask はタスクのステータスを切り替える
func (a *App) HandleToggleTask(taskID string) error {
	_, err := a.taskService.ToggleTaskStatus(a.ctx, taskID)
	if err != nil {
		return fmt.Errorf("failed to toggle task: %w", err)
	}
	
	return a.RefreshTasks()
}

// ApplyFilter はフィルターを適用する
func (a *App) ApplyFilter(filter service.TaskFilter) {
	a.stateManager.SetFilter(filter)
}

// GetCurrentFilter は現在のフィルターを取得する
func (a *App) GetCurrentFilter() service.TaskFilter {
	return a.stateManager.GetCurrentFilter()
}

// Run はアプリケーションを実行する
func (a *App) Run() error {
	return a.tviewApp.Run()
}

// Stop はアプリケーションを停止する
func (a *App) Stop() {
	a.tviewApp.Stop()
}