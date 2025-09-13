package ui

import (
	"sort"
	"strings"

	"task-cli/internal/model"
	"task-cli/internal/service"

	"github.com/rivo/tview"
)

// TaskListWidget はタスクリストを表示するウィジェット
type TaskListWidget struct {
	table             *tview.Table
	theme             *Theme
	allTasks          []*model.Task
	filteredTasks     []*model.Task
	selectedIndex     int
	selectionCallback func(*model.Task)
}

// NewTaskListWidget は新しいTaskListWidgetを作成する
func NewTaskListWidget(theme *Theme) *TaskListWidget {
	table := tview.NewTable().
		SetBorders(true).
		SetSelectable(true, false)

	widget := &TaskListWidget{
		table:         table,
		theme:         theme,
		allTasks:      make([]*model.Task, 0),
		filteredTasks: make([]*model.Task, 0),
		selectedIndex: 0,
	}

	// テーブルのスタイルを設定
	table.SetBackgroundColor(theme.GetBackgroundColor())

	// 選択変更時のコールバック
	table.SetSelectionChangedFunc(func(row, column int) {
		if row > 0 && row-1 < len(widget.filteredTasks) { // ヘッダー行を除く
			widget.selectedIndex = row - 1
			if widget.selectionCallback != nil {
				widget.selectionCallback(widget.filteredTasks[widget.selectedIndex])
			}
		}
	})

	// 初期ヘッダーを設定
	widget.setupHeader()

	return widget
}

// GetPrimitive はtview.Primitiveインターフェースを実装
func (w *TaskListWidget) GetPrimitive() tview.Primitive {
	return w.table
}

// SetTasks はタスクリストを設定する
func (w *TaskListWidget) SetTasks(tasks []*model.Task) {
	w.allTasks = make([]*model.Task, len(tasks))
	copy(w.allTasks, tasks)
	w.filteredTasks = make([]*model.Task, len(tasks))
	copy(w.filteredTasks, tasks)
	
	w.updateTable()
}

// GetTaskCount は表示中のタスク数を返す
func (w *TaskListWidget) GetTaskCount() int {
	return len(w.filteredTasks)
}

// GetSelectedTask は選択中のタスクを返す
func (w *TaskListWidget) GetSelectedTask() *model.Task {
	if w.selectedIndex >= 0 && w.selectedIndex < len(w.filteredTasks) {
		return w.filteredTasks[w.selectedIndex]
	}
	return nil
}

// GetSelectedIndex は選択中のインデックスを返す
func (w *TaskListWidget) GetSelectedIndex() int {
	return w.selectedIndex
}

// SelectTask は指定されたインデックスのタスクを選択する
func (w *TaskListWidget) SelectTask(index int) {
	if index >= 0 && index < len(w.filteredTasks) {
		w.selectedIndex = index
		w.table.Select(index+1, 0) // ヘッダー行を考慮して+1
		if w.selectionCallback != nil {
			w.selectionCallback(w.filteredTasks[index])
		}
	}
}

// SelectNext は次のタスクを選択する
func (w *TaskListWidget) SelectNext() {
	if w.selectedIndex < len(w.filteredTasks)-1 {
		w.SelectTask(w.selectedIndex + 1)
	}
}

// SelectPrevious は前のタスクを選択する
func (w *TaskListWidget) SelectPrevious() {
	if w.selectedIndex > 0 {
		w.SelectTask(w.selectedIndex - 1)
	}
}

// ApplyFilter はフィルターを適用する
func (w *TaskListWidget) ApplyFilter(filter service.TaskFilter) {
	w.filteredTasks = w.applyFilterToTasks(w.allTasks, filter)
	
	// 選択インデックスを調整
	if w.selectedIndex >= len(w.filteredTasks) {
		w.selectedIndex = len(w.filteredTasks) - 1
	}
	if w.selectedIndex < 0 && len(w.filteredTasks) > 0 {
		w.selectedIndex = 0
	}
	
	w.updateTable()
}

// ClearFilter はフィルターをクリアする
func (w *TaskListWidget) ClearFilter() {
	w.filteredTasks = make([]*model.Task, len(w.allTasks))
	copy(w.filteredTasks, w.allTasks)
	w.updateTable()
}

// SortByPriority は優先度でソートする
func (w *TaskListWidget) SortByPriority() {
	sort.Slice(w.filteredTasks, func(i, j int) bool {
		return w.getPriorityOrder(w.filteredTasks[i].Priority) < w.getPriorityOrder(w.filteredTasks[j].Priority)
	})
	w.updateTable()
}

// SortByStatus はステータスでソートする
func (w *TaskListWidget) SortByStatus() {
	sort.Slice(w.filteredTasks, func(i, j int) bool {
		return w.getStatusOrder(w.filteredTasks[i].Status) < w.getStatusOrder(w.filteredTasks[j].Status)
	})
	w.updateTable()
}

// SetSelectionChangedCallback は選択変更時のコールバックを設定する
func (w *TaskListWidget) SetSelectionChangedCallback(callback func(*model.Task)) {
	w.selectionCallback = callback
}

// setupHeader はテーブルヘッダーを設定する
func (w *TaskListWidget) setupHeader() {
	headerStyle := tview.NewTableCell("Status").
		SetTextColor(w.theme.GetHighlightColor()).
		SetBackgroundColor(w.theme.GetBackgroundColor()).
		SetSelectable(false).
		SetAlign(tview.AlignCenter)

	w.table.SetCell(0, 0, headerStyle)
	w.table.SetCell(0, 1, tview.NewTableCell("Priority").
		SetTextColor(w.theme.GetHighlightColor()).
		SetBackgroundColor(w.theme.GetBackgroundColor()).
		SetSelectable(false).
		SetAlign(tview.AlignCenter))
	w.table.SetCell(0, 2, tview.NewTableCell("Title").
		SetTextColor(w.theme.GetHighlightColor()).
		SetBackgroundColor(w.theme.GetBackgroundColor()).
		SetSelectable(false).
		SetAlign(tview.AlignLeft))
	w.table.SetCell(0, 3, tview.NewTableCell("Description").
		SetTextColor(w.theme.GetHighlightColor()).
		SetBackgroundColor(w.theme.GetBackgroundColor()).
		SetSelectable(false).
		SetAlign(tview.AlignLeft))
}

// updateTable はテーブルの内容を更新する
func (w *TaskListWidget) updateTable() {
	// 既存の行をクリア（ヘッダー以外）
	rowCount := w.table.GetRowCount()
	for i := rowCount - 1; i > 0; i-- {
		w.table.RemoveRow(i)
	}

	// タスク行を追加
	for i, task := range w.filteredTasks {
		row := i + 1 // ヘッダー行を考慮

		// ステータス列
		statusCell := tview.NewTableCell(w.getStatusSymbol(task.Status)).
			SetTextColor(w.theme.GetStatusColor(task.Status)).
			SetBackgroundColor(w.theme.GetBackgroundColor()).
			SetAlign(tview.AlignCenter)

		// 優先度列
		priorityCell := tview.NewTableCell(w.getPrioritySymbol(task.Priority)).
			SetTextColor(w.theme.GetPriorityColor(task.Priority)).
			SetBackgroundColor(w.theme.GetBackgroundColor()).
			SetAlign(tview.AlignCenter)

		// タイトル列
		titleCell := tview.NewTableCell(task.Title).
			SetTextColor(w.theme.GetForegroundColor()).
			SetBackgroundColor(w.theme.GetBackgroundColor()).
			SetAlign(tview.AlignLeft)

		// 説明列
		description := task.Description
		if len(description) > 50 {
			description = description[:47] + "..."
		}
		descCell := tview.NewTableCell(description).
			SetTextColor(w.theme.GetForegroundColor()).
			SetBackgroundColor(w.theme.GetBackgroundColor()).
			SetAlign(tview.AlignLeft)

		w.table.SetCell(row, 0, statusCell)
		w.table.SetCell(row, 1, priorityCell)
		w.table.SetCell(row, 2, titleCell)
		w.table.SetCell(row, 3, descCell)
	}

	// 選択状態を復元
	if len(w.filteredTasks) > 0 && w.selectedIndex >= 0 && w.selectedIndex < len(w.filteredTasks) {
		w.table.Select(w.selectedIndex+1, 0)
	}
}

// applyFilterToTasks はタスクリストにフィルターを適用する
func (w *TaskListWidget) applyFilterToTasks(tasks []*model.Task, filter service.TaskFilter) []*model.Task {
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

		// クエリフィルター
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

// getStatusSymbol はステータスに対応するシンボルを取得する
func (w *TaskListWidget) getStatusSymbol(status model.Status) string {
	switch status {
	case model.StatusTodo:
		return "◯"
	case model.StatusInProgress:
		return "◐"
	case model.StatusCompleted:
		return "●"
	default:
		return "?"
	}
}

// getPrioritySymbol は優先度に対応するシンボルを取得する
func (w *TaskListWidget) getPrioritySymbol(priority model.Priority) string {
	switch priority {
	case model.PriorityHigh:
		return "!!!"
	case model.PriorityMedium:
		return " !! "
	case model.PriorityLow:
		return " !  "
	default:
		return "   "
	}
}

// getPriorityOrder は優先度の並び順を取得する
func (w *TaskListWidget) getPriorityOrder(priority model.Priority) int {
	switch priority {
	case model.PriorityHigh:
		return 0
	case model.PriorityMedium:
		return 1
	case model.PriorityLow:
		return 2
	default:
		return 3
	}
}

// getStatusOrder はステータスの並び順を取得する
func (w *TaskListWidget) getStatusOrder(status model.Status) int {
	switch status {
	case model.StatusTodo:
		return 0
	case model.StatusInProgress:
		return 1
	case model.StatusCompleted:
		return 2
	default:
		return 3
	}
}