# 技術設計書 - TUI ToDo CLI アプリケーション (Go言語実装)

## 1. システム概要

### 1.1 アーキテクチャ概要

```
┌─────────────────────────────────────────────────────────────┐
│                         TUI Layer                           │
│  ┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐ │
│  │   TaskList      │ │  InputForm      │ │  StatusBar      │ │
│  │   Widget        │ │   Widget        │ │   Widget        │ │
│  └─────────────────┘ └─────────────────┘ └─────────────────┘ │
└─────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────┐
│                      Application Layer                      │
│  ┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐ │
│  │  TaskService    │ │  AppState       │ │  EventHandler   │ │
│  │  (Business)     │ │  (State)        │ │  (Controller)   │ │
│  └─────────────────┘ └─────────────────┘ └─────────────────┘ │
└─────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────┐
│                       Data Layer                            │
│  ┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐ │
│  │  Repository     │ │   Models        │ │   Validator     │ │
│  │  (Storage)      │ │   (Domain)      │ │   (Rules)       │ │
│  └─────────────────┘ └─────────────────┘ └─────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

### 1.2 技術選定

- **TUIライブラリ**: `tview` (rivo/tview)
  - 理由: Goで最も成熟したTUIライブラリ、豊富なウィジェット、カスタマイズ性
- **データストレージ**: JSON ファイル (`~/.task-cli/data.json`)
- **設定管理**: `viper` (spf13/viper)
- **CLI管理**: `cobra` (spf13/cobra) 
- **バリデーション**: カスタム実装 + `validator` パッケージ
- **テスト**: 標準 `testing` パッケージ + `testify`
- **ビルド**: 標準 `go build` + マルチプラットフォーム対応

## 2. データ設計

### 2.1 データモデル

```go
// Task はタスクの基本構造を定義
type Task struct {
    ID          string    `json:"id" validate:"required,uuid4"`
    Title       string    `json:"title" validate:"required,min=1,max=100"`
    Description string    `json:"description" validate:"max=500"`
    Status      Status    `json:"status" validate:"required"`
    Priority    Priority  `json:"priority" validate:"required"`
    Tags        []string  `json:"tags" validate:"dive,min=1,max=20"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
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

// AppData はアプリケーション全体のデータ構造
type AppData struct {
    Tasks    []Task      `json:"tasks"`
    Settings AppSettings `json:"settings"`
    Version  string      `json:"version"`
}

// AppSettings はアプリケーション設定
type AppSettings struct {
    Theme          string            `json:"theme"`
    KeyBindings    map[string]string `json:"key_bindings"`
    DisplayOptions DisplayOptions    `json:"display_options"`
}

// DisplayOptions は表示オプション
type DisplayOptions struct {
    ShowCompleted   bool `json:"show_completed"`
    ShowDescription bool `json:"show_description"`
    SortBy          string `json:"sort_by"`
    SortAscending   bool `json:"sort_ascending"`
}
```

### 2.2 データアクセス層

```go
// Repository はデータアクセスのインターフェース
type Repository interface {
    Load() (*AppData, error)
    Save(data *AppData) error
    Backup() error
    Restore() (*AppData, error)
}

// FileRepository はファイルベースのRepository実装
type FileRepository struct {
    dataPath   string
    backupPath string
    mu         sync.RWMutex
}

func NewFileRepository() *FileRepository {
    homeDir, _ := os.UserHomeDir()
    dataDir := filepath.Join(homeDir, ".task-cli")
    
    return &FileRepository{
        dataPath:   filepath.Join(dataDir, "data.json"),
        backupPath: filepath.Join(dataDir, "data.backup.json"),
    }
}

func (r *FileRepository) Load() (*AppData, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    // 実装詳細...
}
```

## 3. UI/UX 設計

### 3.1 画面構成（tview実装）

```go
// App はメインアプリケーション構造
type App struct {
    *tview.Application
    
    // UI Components
    headerBar  *tview.TextView
    filterBar  *tview.Form
    taskList   *tview.Table
    statusBar  *tview.TextView
    
    // Services
    taskService *TaskService
    state       *AppState
    repo        Repository
    
    // Layout
    layout     *tview.Flex
    pages      *tview.Pages
}

// TaskListWidget はタスクリスト表示ウィジェット
type TaskListWidget struct {
    *tview.Table
    tasks       []Task
    selectedIdx int
    filter      TaskFilter
    app         *App
}

// InputFormWidget はタスク入力フォーム
type InputFormWidget struct {
    *tview.Form
    mode     FormMode // Create or Edit
    taskID   string   // Edit時のタスクID
    callback func(Task) error
}

type FormMode int

const (
    FormModeCreate FormMode = iota
    FormModeEdit
)
```

### 3.2 カラーテーマ

```go
// Theme はカラーテーマを定義
type Theme struct {
    Primary     tcell.Color
    Success     tcell.Color
    Warning     tcell.Color
    Danger      tcell.Color
    Muted       tcell.Color
    Background  tcell.Color
    Foreground  tcell.Color
}

var DefaultTheme = Theme{
    Primary:     tcell.ColorBlue,
    Success:     tcell.ColorGreen,
    Warning:     tcell.ColorYellow,
    Danger:      tcell.ColorRed,
    Muted:       tcell.ColorGray,
    Background:  tcell.ColorBlack,
    Foreground:  tcell.ColorWhite,
}

// GetPriorityColor は優先度に応じた色を返す
func (t *Theme) GetPriorityColor(priority Priority) tcell.Color {
    switch priority {
    case PriorityHigh:
        return t.Danger
    case PriorityMedium:
        return t.Warning
    case PriorityLow:
        return t.Success
    default:
        return t.Muted
    }
}

// GetStatusColor はステータスに応じた色を返す
func (t *Theme) GetStatusColor(status Status) tcell.Color {
    switch status {
    case StatusCompleted:
        return t.Success
    case StatusInProgress:
        return t.Primary
    case StatusTodo:
        return t.Muted
    default:
        return t.Muted
    }
}
```

### 3.3 キーバインド

```go
// KeyBinding はキーバインドを定義
type KeyBinding struct {
    Key         tcell.Key
    Char        rune
    Description string
    Handler     func(*App) error
}

// DefaultKeyBindings はデフォルトキーバインド
var DefaultKeyBindings = []KeyBinding{
    // Navigation
    {tcell.KeyUp, 0, "Move up", (*App).moveUp},
    {tcell.KeyDown, 0, "Move down", (*App).moveDown},
    {0, 'k', "Move up (vim)", (*App).moveUp},
    {0, 'j', "Move down (vim)", (*App).moveDown},
    
    // Actions
    {tcell.KeyEnter, 0, "Edit task", (*App).editTask},
    {0, ' ', "Toggle status", (*App).toggleStatus},
    {tcell.KeyDelete, 0, "Delete task", (*App).deleteTask},
    {0, 'n', "New task", (*App).newTask},
    {0, 'd', "Delete task", (*App).deleteTask},
    
    // Filters
    {0, 'a', "Show all", (*App).showAll},
    {0, 't', "Show todo", (*App).showTodo},
    {0, 'p', "Show in progress", (*App).showInProgress},
    {0, 'c', "Show completed", (*App).showCompleted},
    
    // Utils
    {0, '/', "Search", (*App).search},
    {0, '?', "Help", (*App).showHelp},
    {tcell.KeyEscape, 0, "Cancel/Back", (*App).cancel},
    {0, 'q', "Quit", (*App).quit},
    {tcell.KeyCtrlS, 0, "Save", (*App).save},
}
```

## 4. ビジネスロジック設計

### 4.1 TaskService

```go
// TaskService はタスク管理のビジネスロジック
type TaskService struct {
    repo      Repository
    validator *Validator
    mu        sync.RWMutex
}

func NewTaskService(repo Repository) *TaskService {
    return &TaskService{
        repo:      repo,
        validator: NewValidator(),
    }
}

func (s *TaskService) CreateTask(req CreateTaskRequest) (*Task, error) {
    if err := s.validator.ValidateCreateRequest(req); err != nil {
        return nil, fmt.Errorf("validation failed: %w", err)
    }
    
    task := &Task{
        ID:          uuid.New().String(),
        Title:       req.Title,
        Description: req.Description,
        Status:      StatusTodo,
        Priority:    req.Priority,
        Tags:        req.Tags,
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }
    
    data, err := s.repo.Load()
    if err != nil {
        return nil, fmt.Errorf("failed to load data: %w", err)
    }
    
    data.Tasks = append(data.Tasks, *task)
    
    if err := s.repo.Save(data); err != nil {
        return nil, fmt.Errorf("failed to save task: %w", err)
    }
    
    return task, nil
}

func (s *TaskService) UpdateTask(id string, req UpdateTaskRequest) (*Task, error) {
    // 実装詳細...
}

func (s *TaskService) DeleteTask(id string) error {
    // 実装詳細...
}

func (s *TaskService) ToggleTask(id string) (*Task, error) {
    // 実装詳細...
}

func (s *TaskService) SearchTasks(query string, filter TaskFilter) ([]Task, error) {
    // 実装詳細...
}
```

### 4.2 状態管理

```go
// AppState はアプリケーション状態を管理
type AppState struct {
    tasks         []Task
    selectedTask  *Task
    filter        TaskFilter
    searchQuery   string
    currentView   ViewType
    isLoading     bool
    error         error
    mu            sync.RWMutex
    subscribers   []StateChangeHandler
}

type StateChangeHandler func(*AppState)

type ViewType int

const (
    ViewTypeList ViewType = iota
    ViewTypeEdit
    ViewTypeHelp
    ViewTypeSearch
)

// TaskFilter はタスクフィルタ条件
type TaskFilter struct {
    Status   []Status
    Priority []Priority
    Tags     []string
}

func NewAppState() *AppState {
    return &AppState{
        tasks:       make([]Task, 0),
        filter:      TaskFilter{},
        currentView: ViewTypeList,
        subscribers: make([]StateChangeHandler, 0),
    }
}

func (s *AppState) Subscribe(handler StateChangeHandler) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.subscribers = append(s.subscribers, handler)
}

func (s *AppState) notifySubscribers() {
    for _, handler := range s.subscribers {
        handler(s)
    }
}

func (s *AppState) SetTasks(tasks []Task) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.tasks = tasks
    s.notifySubscribers()
}
```

## 5. パッケージ構成

```
cmd/
├── root.go              # Cobraルートコマンド
└── run.go               # アプリ実行コマンド

internal/
├── app/
│   ├── app.go           # メインアプリケーション
│   └── config.go        # 設定管理
├── ui/
│   ├── widgets/
│   │   ├── tasklist.go  # タスクリストウィジェット
│   │   ├── form.go      # フォームウィジェット
│   │   └── statusbar.go # ステータスバーウィジェット
│   ├── theme.go         # テーマ管理
│   └── keybindings.go   # キーバインド管理
├── service/
│   ├── task.go          # タスクサービス
│   └── state.go         # 状態管理サービス
├── repository/
│   ├── interface.go     # Repository interface
│   └── file.go          # ファイルRepository実装
├── model/
│   ├── task.go          # タスクモデル
│   ├── app.go           # アプリデータモデル
│   └── request.go       # リクエスト/レスポンスモデル
└── validator/
    └── validator.go     # バリデーション

pkg/
└── errors/
    └── errors.go        # カスタムエラー型

go.mod                   # Go modules
go.sum                   # 依存関係ロック
main.go                  # エントリーポイント
```

## 6. エラー処理

### 6.1 カスタムエラー型

```go
// AppError はアプリケーション固有のエラー
type AppError struct {
    Type       ErrorType
    Message    string
    Cause      error
    Recoverable bool
}

type ErrorType string

const (
    ErrorTypeValidation    ErrorType = "validation_error"
    ErrorTypeFileSystem    ErrorType = "file_system_error"
    ErrorTypeDataCorruption ErrorType = "data_corruption_error"
    ErrorTypeUnknown       ErrorType = "unknown_error"
)

func (e *AppError) Error() string {
    if e.Cause != nil {
        return fmt.Sprintf("%s: %s (caused by: %v)", e.Type, e.Message, e.Cause)
    }
    return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

func NewValidationError(message string, cause error) *AppError {
    return &AppError{
        Type:        ErrorTypeValidation,
        Message:     message,
        Cause:       cause,
        Recoverable: true,
    }
}
```

### 6.2 エラー回復戦略

```go
// RecoveryManager はエラー回復を管理
type RecoveryManager struct {
    repo Repository
}

func (rm *RecoveryManager) HandleDataCorruption() (*AppData, error) {
    log.Warn("Data corruption detected, attempting recovery from backup")
    
    data, err := rm.repo.Restore()
    if err != nil {
        log.Error("Failed to restore from backup", "error", err)
        return rm.createEmptyData(), nil
    }
    
    return data, nil
}

func (rm *RecoveryManager) createEmptyData() *AppData {
    return &AppData{
        Tasks:    make([]Task, 0),
        Settings: getDefaultSettings(),
        Version:  version.Current,
    }
}
```

## 7. TDD実装戦略

### 7.1 テストファースト開発フロー

```
1. RED Phase (失敗するテストを書く)
   ├── 期待する動作を定義
   ├── テストが失敗することを確認
   └── 失敗理由が正しいことを確認

2. GREEN Phase (テストを通す最小限のコードを書く) 
   ├── テストを通すためだけのコード
   ├── 他のテストを壊さない
   └── 過度な実装をしない

3. REFACTOR Phase (コードを改善する)
   ├── テストを維持しながらコード改善
   ├── 重複排除
   └── 設計の改善
```

### 7.2 TDDサイクル適用例

```go
// RED: 失敗するテストを先に書く
func TestTaskService_CreateTask_ShouldReturnValidTask(t *testing.T) {
    // Given
    repo := &mockRepository{}
    service := NewTaskService(repo)
    
    req := CreateTaskRequest{
        Title:       "Test Task",
        Description: "Test Description", 
        Priority:    PriorityHigh,
        Tags:        []string{"test"},
    }
    
    // When
    task, err := service.CreateTask(req)
    
    // Then
    assert.NoError(t, err)
    assert.Equal(t, req.Title, task.Title)
    assert.Equal(t, StatusTodo, task.Status)
    assert.NotEmpty(t, task.ID)
    assert.WithinDuration(t, time.Now(), task.CreatedAt, time.Second)
}

// GREEN: テストを通す最小限の実装
func (s *TaskService) CreateTask(req CreateTaskRequest) (*Task, error) {
    task := &Task{
        ID:        "some-id", // まず固定値で通す
        Title:     req.Title,
        Status:    StatusTodo,
        CreatedAt: time.Now(),
    }
    return task, nil
}

// REFACTOR: より良い実装に改善
func (s *TaskService) CreateTask(req CreateTaskRequest) (*Task, error) {
    if err := s.validator.ValidateCreateRequest(req); err != nil {
        return nil, fmt.Errorf("validation failed: %w", err)
    }
    
    task := &Task{
        ID:          uuid.New().String(), // 実際のUUID生成
        Title:       req.Title,
        Description: req.Description,
        Status:      StatusTodo,
        Priority:    req.Priority,
        Tags:        req.Tags,
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }
    
    return task, nil
}
```

### 7.3 テストピラミッド構成

```
        /\
       /  \
      / E2E \     <- 10%: ユーザーシナリオテスト
     /______\
    /        \
   /Integration\ <- 20%: コンポーネント間連携テスト  
  /____________\
 /              \
/  Unit Tests    \ <- 70%: 関数・メソッド単体テスト
/________________\
```

### 7.4 TDD品質指標

- **ユニットテストカバレッジ**: >90%
- **統合テストカバレッジ**: >80%
- **テスト実行時間**: <30秒（全テスト）
- **テストの独立性**: 各テストが他に依存しない
- **テストの可読性**: Given-When-Then構造

### 7.5 モック戦略

```go
// Repository のモック
type mockRepository struct {
    data     *AppData
    saveErr  error
    loadErr  error
}

func (m *mockRepository) Save(data *AppData) error {
    if m.saveErr != nil {
        return m.saveErr
    }
    m.data = data
    return nil
}

func (m *mockRepository) Load() (*AppData, error) {
    if m.loadErr != nil {
        return nil, m.loadErr
    }
    return m.data, nil
}
```

### 7.6 テスト分類とネーミング

```go
// ユニットテスト: TestTargetFunction_Scenario_ExpectedResult
func TestTaskService_CreateTask_WithValidData_ShouldReturnTask(t *testing.T)
func TestTaskService_CreateTask_WithInvalidTitle_ShouldReturnError(t *testing.T)
func TestTaskService_UpdateTask_WithNonExistentID_ShouldReturnError(t *testing.T)

// 統合テスト: TestIntegration_Component_Scenario
func TestIntegration_FileRepository_SaveAndLoad_ShouldPersistData(t *testing.T)
func TestIntegration_TaskService_WithRealRepository_ShouldHandleFullCycle(t *testing.T)

// E2Eテスト: TestE2E_UserScenario
func TestE2E_CreateEditDeleteTask_ShouldCompleteSuccessfully(t *testing.T)
```

## 8. ビルドとデプロイ

### 8.1 Makefile

```makefile
.PHONY: build test clean install release

BINARY_NAME=task-cli
VERSION=$(shell git describe --tags --always --dirty)
LDFLAGS=-ldflags "-X main.version=${VERSION}"

build:
	go build ${LDFLAGS} -o bin/${BINARY_NAME} ./cmd

test:
	go test -v ./...

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

clean:
	rm -rf bin/
	go clean

install: build
	cp bin/${BINARY_NAME} ${GOPATH}/bin/

release:
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o bin/${BINARY_NAME}-linux-amd64 ./cmd
	GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o bin/${BINARY_NAME}-darwin-amd64 ./cmd
	GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o bin/${BINARY_NAME}-windows-amd64.exe ./cmd
```

### 8.2 依存関係

```go
// go.mod
module github.com/user/task-cli

go 1.21

require (
    github.com/gdamore/tcell/v2 v2.6.0
    github.com/rivo/tview v0.0.0-20230916063804-5a2d8b05c73e
    github.com/spf13/cobra v1.7.0
    github.com/spf13/viper v1.16.0
    github.com/google/uuid v1.3.0
    github.com/stretchr/testify v1.8.4
)
```

この設計書に基づいて、Goで高品質なTUI ToDoアプリケーションを実装することができます。モジュラーで保守しやすい設計により、要件を満たしながら拡張性も確保されています。