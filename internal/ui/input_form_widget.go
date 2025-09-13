package ui

import (
	"errors"
	"strings"

	"task-cli/internal/model"

	"github.com/rivo/tview"
)

// FormMode はフォームのモードを定義
type FormMode int

const (
	FormModeCreate FormMode = iota
	FormModeEdit
)

// FormData はフォームのデータを定義
type FormData struct {
	Title       string
	Description string
	Priority    model.Priority
	Status      model.Status
	Tags        []string
}

// InputFormWidget はタスク入力フォームのウィジェット
type InputFormWidget struct {
	form           *tview.Form
	theme          *Theme
	mode           FormMode
	enabled        bool
	errorMessage   string
	submitCallback func(FormData)
	cancelCallback func()
	
	// フィールド
	titleField       *tview.InputField
	descriptionField *tview.TextArea
	priorityField    *tview.DropDown
	statusField      *tview.DropDown
	tagsField        *tview.InputField
	errorLabel       *tview.TextView
}

// NewInputFormWidget は新しいInputFormWidgetを作成する
func NewInputFormWidget(theme *Theme) *InputFormWidget {
	form := tview.NewForm()
	
	widget := &InputFormWidget{
		form:    form,
		theme:   theme,
		mode:    FormModeCreate,
		enabled: true,
	}
	
	widget.initializeFields()
	widget.setupForm()
	
	return widget
}

// GetPrimitive はtview.Primitiveインターフェースを実装
func (w *InputFormWidget) GetPrimitive() tview.Primitive {
	return w.form
}

// initializeFields はフォームフィールドを初期化する
func (w *InputFormWidget) initializeFields() {
	// タイトルフィールド
	w.titleField = tview.NewInputField().
		SetLabel("Title: ").
		SetFieldWidth(50)

	// 説明フィールド
	w.descriptionField = tview.NewTextArea().
		SetLabel("Description: ").
		SetWrap(true).
		SetSize(3, 50)

	// 優先度フィールド
	w.priorityField = tview.NewDropDown().
		SetLabel("Priority: ").
		SetOptions([]string{"Low", "Medium", "High"}, nil).
		SetCurrentOption(1) // デフォルトはMedium

	// ステータスフィールド（編集モード時のみ表示）
	w.statusField = tview.NewDropDown().
		SetLabel("Status: ").
		SetOptions([]string{"Todo", "In Progress", "Completed"}, nil).
		SetCurrentOption(0) // デフォルトはTodo

	// タグフィールド
	w.tagsField = tview.NewInputField().
		SetLabel("Tags: ").
		SetFieldWidth(50).
		SetPlaceholder("Comma separated tags")

	// エラーラベル
	w.errorLabel = tview.NewTextView().
		SetTextColor(w.theme.GetPriorityColor(model.PriorityHigh)). // 赤色でエラー表示
		SetDynamicColors(true)
}

// setupForm はフォームを設定する
func (w *InputFormWidget) setupForm() {
	w.form.SetBackgroundColor(w.theme.GetBackgroundColor())
	
	// 基本フィールドを追加
	w.form.AddFormItem(w.titleField)
	w.form.AddFormItem(w.descriptionField)
	w.form.AddFormItem(w.priorityField)
	w.form.AddFormItem(w.tagsField)
	w.form.AddFormItem(w.errorLabel)

	// ボタンを追加
	w.form.AddButton("Submit", func() {
		w.Submit()
	}).AddButton("Cancel", func() {
		w.Cancel()
	})

	w.form.SetButtonsAlign(tview.AlignCenter)
}

// SetMode はフォームのモードを設定する
func (w *InputFormWidget) SetMode(mode FormMode) {
	w.mode = mode
	w.updateFormForMode()
}

// GetMode は現在のフォームモードを取得する
func (w *InputFormWidget) GetMode() FormMode {
	return w.mode
}

// updateFormForMode はモードに応じてフォームを更新する
func (w *InputFormWidget) updateFormForMode() {
	// フォームをクリアして再構築
	w.form.Clear(true)
	
	// 基本フィールドを追加
	w.form.AddFormItem(w.titleField)
	w.form.AddFormItem(w.descriptionField)
	w.form.AddFormItem(w.priorityField)
	
	// 編集モードの場合はステータスフィールドを追加
	if w.mode == FormModeEdit {
		w.form.AddFormItem(w.statusField)
	}
	
	w.form.AddFormItem(w.tagsField)
	w.form.AddFormItem(w.errorLabel)

	// ボタンのテキストをモードに応じて変更
	submitText := "Create"
	if w.mode == FormModeEdit {
		submitText = "Update"
	}
	
	w.form.AddButton(submitText, func() {
		w.Submit()
	}).AddButton("Cancel", func() {
		w.Cancel()
	})

	w.form.SetButtonsAlign(tview.AlignCenter)
}

// LoadTask は編集対象のタスクを読み込む
func (w *InputFormWidget) LoadTask(task *model.Task) {
	w.SetTitle(task.Title)
	w.SetDescription(task.Description)
	w.SetPriority(task.Priority)
	w.SetStatus(task.Status)
	w.SetTags(strings.Join(task.Tags, ","))
}

// SetTitle はタイトルを設定する
func (w *InputFormWidget) SetTitle(title string) {
	w.titleField.SetText(title)
}

// GetTitle はタイトルを取得する
func (w *InputFormWidget) GetTitle() string {
	return w.titleField.GetText()
}

// SetDescription は説明を設定する
func (w *InputFormWidget) SetDescription(description string) {
	w.descriptionField.SetText(description, true)
}

// GetDescription は説明を取得する
func (w *InputFormWidget) GetDescription() string {
	return w.descriptionField.GetText()
}

// SetPriority は優先度を設定する
func (w *InputFormWidget) SetPriority(priority model.Priority) {
	index := w.priorityToIndex(priority)
	w.priorityField.SetCurrentOption(index)
}

// GetPriority は優先度を取得する
func (w *InputFormWidget) GetPriority() model.Priority {
	_, option := w.priorityField.GetCurrentOption()
	return w.indexToPriority(option)
}

// SetStatus はステータスを設定する
func (w *InputFormWidget) SetStatus(status model.Status) {
	index := w.statusToIndex(status)
	w.statusField.SetCurrentOption(index)
}

// GetStatus はステータスを取得する
func (w *InputFormWidget) GetStatus() model.Status {
	_, option := w.statusField.GetCurrentOption()
	return w.indexToStatus(option)
}

// SetTags はタグを設定する
func (w *InputFormWidget) SetTags(tags string) {
	w.tagsField.SetText(tags)
}

// GetTags はタグを取得する
func (w *InputFormWidget) GetTags() string {
	return w.tagsField.GetText()
}

// Validate はフォームの入力を検証する
func (w *InputFormWidget) Validate() error {
	if strings.TrimSpace(w.GetTitle()) == "" {
		return errors.New("title is required")
	}
	
	if len(w.GetTitle()) > 100 {
		return errors.New("title must be 100 characters or less")
	}
	
	if len(w.GetDescription()) > 500 {
		return errors.New("description must be 500 characters or less")
	}
	
	return nil
}

// Clear はフォームをクリアする
func (w *InputFormWidget) Clear() {
	w.SetTitle("")
	w.SetDescription("")
	w.SetPriority(model.PriorityMedium) // デフォルト値
	w.SetStatus(model.StatusTodo) // デフォルト値
	w.SetTags("")
	w.ClearError()
}

// Focus はフォームにフォーカスを設定する
func (w *InputFormWidget) Focus() {
	w.form.SetFocus(0) // 最初のフィールドにフォーカス
}

// HasFocus はフォームがフォーカスを持っているかを返す
func (w *InputFormWidget) HasFocus() bool {
	return w.form.HasFocus()
}

// Submit はフォームを送信する
func (w *InputFormWidget) Submit() {
	if err := w.Validate(); err != nil {
		w.SetErrorMessage(err.Error())
		return
	}
	
	w.ClearError()
	
	if w.submitCallback != nil {
		data := w.GetFormData()
		w.submitCallback(data)
	}
}

// Cancel はフォームをキャンセルする
func (w *InputFormWidget) Cancel() {
	w.ClearError()
	if w.cancelCallback != nil {
		w.cancelCallback()
	}
}

// SetSubmitCallback は送信時のコールバックを設定する
func (w *InputFormWidget) SetSubmitCallback(callback func(FormData)) {
	w.submitCallback = callback
}

// SetCancelCallback はキャンセル時のコールバックを設定する
func (w *InputFormWidget) SetCancelCallback(callback func()) {
	w.cancelCallback = callback
}

// GetFormData はフォームデータを取得する
func (w *InputFormWidget) GetFormData() FormData {
	tags := []string{}
	if tagText := strings.TrimSpace(w.GetTags()); tagText != "" {
		for _, tag := range strings.Split(tagText, ",") {
			if trimmed := strings.TrimSpace(tag); trimmed != "" {
				tags = append(tags, trimmed)
			}
		}
	}
	
	return FormData{
		Title:       strings.TrimSpace(w.GetTitle()),
		Description: strings.TrimSpace(w.GetDescription()),
		Priority:    w.GetPriority(),
		Status:      w.GetStatus(),
		Tags:        tags,
	}
}

// SetEnabled はフォームの有効/無効を設定する
func (w *InputFormWidget) SetEnabled(enabled bool) {
	w.enabled = enabled
	// 実際の有効/無効の実装は、必要に応じてフィールド毎に設定
}

// IsEnabled はフォームが有効かを返す
func (w *InputFormWidget) IsEnabled() bool {
	return w.enabled
}

// SetErrorMessage はエラーメッセージを設定する
func (w *InputFormWidget) SetErrorMessage(message string) {
	w.errorMessage = message
	w.errorLabel.SetText("[red]" + message + "[white]")
}

// GetErrorMessage はエラーメッセージを取得する
func (w *InputFormWidget) GetErrorMessage() string {
	return w.errorMessage
}

// ClearError はエラーメッセージをクリアする
func (w *InputFormWidget) ClearError() {
	w.errorMessage = ""
	w.errorLabel.SetText("")
}

// priorityToIndex は優先度をインデックスに変換する
func (w *InputFormWidget) priorityToIndex(priority model.Priority) int {
	switch priority {
	case model.PriorityLow:
		return 0
	case model.PriorityMedium:
		return 1
	case model.PriorityHigh:
		return 2
	default:
		return 1 // デフォルト
	}
}

// indexToPriority はインデックスを優先度に変換する
func (w *InputFormWidget) indexToPriority(option string) model.Priority {
	switch option {
	case "Low":
		return model.PriorityLow
	case "Medium":
		return model.PriorityMedium
	case "High":
		return model.PriorityHigh
	default:
		return model.PriorityMedium
	}
}

// statusToIndex はステータスをインデックスに変換する
func (w *InputFormWidget) statusToIndex(status model.Status) int {
	switch status {
	case model.StatusTodo:
		return 0
	case model.StatusInProgress:
		return 1
	case model.StatusCompleted:
		return 2
	default:
		return 0 // デフォルト
	}
}

// indexToStatus はインデックスをステータスに変換する
func (w *InputFormWidget) indexToStatus(option string) model.Status {
	switch option {
	case "Todo":
		return model.StatusTodo
	case "In Progress":
		return model.StatusInProgress
	case "Completed":
		return model.StatusCompleted
	default:
		return model.StatusTodo
	}
}