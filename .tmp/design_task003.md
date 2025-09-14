# 設計書 - ゴミ箱機能の実装

## システム設計概要
既存のタスク管理システムにゴミ箱機能を追加し、削除したタスクを一時保持・復元可能な仕組みを実装する。データ永続化、UI拡張、自動削除機能を包括的に設計する。

## アーキテクチャ設計

### システム構成の拡張
```
┌─────────────────┐  ┌──────────────────┐  ┌─────────────────┐
│   Task Model    │  │   Trash Service  │  │  Storage Layer  │
│                 │  │                  │  │                 │
│ - Status拡張    │→ │ - Move to Trash  │→ │ - Task DB       │
│ - Metadata追加  │  │ - Restore        │  │ - Trash DB      │
│                 │  │ - Auto Cleanup   │  │                 │
└─────────────────┘  └──────────────────┘  └─────────────────┘
         ↑                      ↑                     ↑
┌─────────────────┐  ┌──────────────────┐  ┌─────────────────┐
│   UI Layer      │  │  Controller      │  │   Background    │
│                 │  │                  │  │   Workers       │
│ - Main View     │  │ - Task CRUD      │  │                 │
│ - Trash View    │  │ - Trash CRUD     │  │ - Auto Cleanup  │
│ - Navigation    │  │ - View Switching │  │   Scheduler     │
└─────────────────┘  └──────────────────┘  └─────────────────┘
```

## データモデル設計

### Task構造体の拡張
```go
type Task struct {
    ID          string    `json:"id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    Status      TaskStatus `json:"status"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`

    // ゴミ箱機能用フィールド
    TrashData   *TrashMetadata `json:"trash_data,omitempty"`
}

type TaskStatus int

const (
    StatusActive TaskStatus = iota
    StatusCompleted
    StatusInTrash
)

type TrashMetadata struct {
    DeletedAt      time.Time  `json:"deleted_at"`
    AutoDeleteAt   time.Time  `json:"auto_delete_at"`
    OriginalStatus TaskStatus `json:"original_status"`
    DeletedBy      string     `json:"deleted_by,omitempty"` // 将来拡張用
}
```

### データベース設計
```sql
-- 既存のtasksテーブルの拡張
ALTER TABLE tasks ADD COLUMN status INTEGER DEFAULT 0;
ALTER TABLE tasks ADD COLUMN deleted_at TIMESTAMP NULL;
ALTER TABLE tasks ADD COLUMN auto_delete_at TIMESTAMP NULL;
ALTER TABLE tasks ADD COLUMN original_status INTEGER NULL;

-- インデックス追加
CREATE INDEX idx_tasks_status ON tasks(status);
CREATE INDEX idx_tasks_auto_delete ON tasks(auto_delete_at) WHERE status = 2;
```

## サービス層設計

### TrashService実装
```go
type TrashService interface {
    MoveToTrash(taskID string) error
    RestoreFromTrash(taskID string) error
    PermanentDelete(taskID string) error
    EmptyTrash() error
    ListTrashTasks() ([]*Task, error)
    AutoCleanup() error
}

type trashService struct {
    taskRepo TaskRepository
    config   TrashConfig
}

type TrashConfig struct {
    AutoDeleteDays int `default:"30"`
    MaxTrashItems  int `default:"1000"`
}
```

### 主要メソッドの実装設計

#### MoveToTrash
```go
func (s *trashService) MoveToTrash(taskID string) error {
    task, err := s.taskRepo.GetByID(taskID)
    if err != nil {
        return err
    }

    task.TrashData = &TrashMetadata{
        DeletedAt:      time.Now(),
        AutoDeleteAt:   time.Now().AddDate(0, 0, s.config.AutoDeleteDays),
        OriginalStatus: task.Status,
    }
    task.Status = StatusInTrash

    return s.taskRepo.Update(task)
}
```

#### RestoreFromTrash
```go
func (s *trashService) RestoreFromTrash(taskID string) error {
    task, err := s.taskRepo.GetByID(taskID)
    if err != nil {
        return err
    }

    if task.Status != StatusInTrash {
        return ErrTaskNotInTrash
    }

    task.Status = task.TrashData.OriginalStatus
    task.TrashData = nil

    return s.taskRepo.Update(task)
}
```

## UI層設計

### 画面遷移設計
```
Main List View ←→ Trash View
     ↓              ↓
Detail View    Restore/Delete
     ↓              ↓
Edit View      Confirmation
```

### Trash View Component
```go
type TrashView struct {
    tasks        []*Task
    selectedIdx  int
    confirmDialog *ConfirmDialog

    // Bubble Tea関連
    viewport.Model
    help.Model
}

func (m TrashView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "r":
            return m.restoreSelected()
        case "D":
            return m.showDeleteConfirm()
        case "C":
            return m.showEmptyTrashConfirm()
        case "q", "esc":
            return m.returnToMain()
        }
    }
    return m, nil
}
```

### 確認ダイアログの共通化
```go
type ConfirmDialog struct {
    title       string
    message     string
    confirmText string
    cancelText  string
    onConfirm   func() tea.Cmd
    onCancel    func() tea.Cmd
    focused     bool // true: confirm, false: cancel
}
```

## 自動削除機能設計

### バックグラウンドワーカー
```go
type AutoCleanupWorker struct {
    trashService TrashService
    interval     time.Duration
    logger       Logger
}

func (w *AutoCleanupWorker) Start() {
    ticker := time.NewTicker(w.interval)
    go func() {
        for range ticker.C {
            if err := w.trashService.AutoCleanup(); err != nil {
                w.logger.Error("Auto cleanup failed", err)
            }
        }
    }()
}

// アプリ起動時のクリーンアップ
func (s *trashService) AutoCleanup() error {
    expiredTasks, err := s.taskRepo.GetExpiredTrashTasks()
    if err != nil {
        return err
    }

    for _, task := range expiredTasks {
        if err := s.taskRepo.Delete(task.ID); err != nil {
            return err
        }
    }

    return nil
}
```

## 設定管理設計

### 設定ファイル拡張
```yaml
# config.yaml
trash:
  auto_delete_days: 30
  max_items: 1000
  cleanup_on_startup: true

ui:
  trash_keybind: "gt"
  show_trash_count: true
```

### 設定の動的変更（将来拡張）
- 設定画面での変更
- コマンドラインオプション
- 環境変数での上書き

## パフォーマンス最適化

### データベースクエリ最適化
```sql
-- 効率的なゴミ箱一覧取得
SELECT * FROM tasks
WHERE status = 2
ORDER BY deleted_at DESC
LIMIT 100;

-- 自動削除対象の効率的な取得
SELECT id FROM tasks
WHERE status = 2 AND auto_delete_at < NOW()
LIMIT 1000;
```

### メモリ使用量最適化
- 大量データの遅延読み込み
- ページネーション実装
- 不要なメタデータの遅延初期化

## エラーハンドリング設計

### カスタムエラー定義
```go
var (
    ErrTaskNotInTrash     = errors.New("task is not in trash")
    ErrTaskAlreadyInTrash = errors.New("task is already in trash")
    ErrTrashLimitExceeded = errors.New("trash limit exceeded")
    ErrAutoDeleteFailed   = errors.New("auto delete failed")
)
```

### 復旧戦略
- 削除操作の安全な失敗
- 部分的な復元処理
- データ整合性の保証

## テスト戦略

### 単体テスト
- TrashService の各メソッド
- UI コンポーネントの動作
- データベース操作

### 統合テスト
- 削除→復元のワークフロー
- 自動削除機能
- 画面遷移

### パフォーマンステスト
- 大量データでの動作
- メモリ使用量測定
- レスポンス時間測定