# TDDタスク分割書 - TUI ToDo CLI アプリケーション

## TDD実装原則

各機能は以下のRed-Green-Refactorサイクルに従って実装します：

1. **Red**: 失敗するテストを先に書く
2. **Green**: テストを通す最小限のコードを書く
3. **Refactor**: テストを維持しながらコードを改善する

## フェーズ1: プロジェクト基盤構築

### 1.1 プロジェクト初期化
- **優先度**: 高 | **推定時間**: 45分
- **説明**: Go modulesとディレクトリ構造、テスト環境の構築

#### 1.1.1 RED: プロジェクト構造テスト
- プロジェクト構造の検証テストを作成
- go.modファイル存在確認テスト
- 必要ディレクトリ存在確認テスト

#### 1.1.2 GREEN: 基本構造実装
- `go mod init`でプロジェクト初期化
- ディレクトリ構造作成 (`cmd/`, `internal/`, `pkg/`)
- `main.go`作成

#### 1.1.3 REFACTOR: 構造最適化
- ディレクトリ構成の見直し
- `.gitignore`、`README.md`の追加
- `Makefile`でテスト・ビルドコマンド設定

### 1.2 依存関係管理
- **優先度**: 高 | **推定時間**: 30分
- **説明**: 必要なライブラリの追加とテスト環境構築

#### 1.2.1 RED: 依存関係テスト
- 必要なパッケージが利用可能かテスト
- import文のコンパイルテスト

#### 1.2.2 GREEN: ライブラリ追加
- tview, cobra, viper, testify, uuid追加
- go.modとgo.sumの更新

#### 1.2.3 REFACTOR: 依存関係最適化
- 不要な依存関係の削除
- バージョン固定の検討

## フェーズ2: データモデル層（TDD）

### 2.1 Taskモデル実装
- **優先度**: 高 | **推定時間**: 90分
- **説明**: 基本データ構造をTDDで実装

#### 2.1.1 RED: Taskモデルテスト
```go
func TestTask_NewTask_WithValidData_ShouldCreateTask(t *testing.T)
func TestTask_Validate_WithInvalidTitle_ShouldReturnError(t *testing.T) 
func TestTask_ToJSON_ShouldReturnValidJSON(t *testing.T)
```

#### 2.1.2 GREEN: 基本Task構造体
- Task構造体の最小実装
- 必要なフィールドのみ定義
- JSON tagの追加

#### 2.1.3 REFACTOR: モデル改善
- バリデーション機能強化
- メソッドの追加（String(), IsCompleted()など）
- ドキュメント追加

### 2.2 Status・Priority型実装
- **優先度**: 高 | **推定時間**: 60分

#### 2.2.1 RED: Enum型テスト
```go
func TestStatus_String_ShouldReturnCorrectValue(t *testing.T)
func TestPriority_IsValid_WithInvalidValue_ShouldReturnFalse(t *testing.T)
```

#### 2.2.2 GREEN: 基本Enum実装
- Status、Priority型の定義
- 基本的なString()メソッド

#### 2.2.3 REFACTOR: Enum機能強化
- バリデーション機能追加
- カラー情報の関連付け
- 並び替え優先度の定義

### 2.3 AppDataモデル実装
- **優先度**: 高 | **推定時間**: 75分

#### 2.3.1 RED: AppDataテスト
```go
func TestAppData_NewAppData_ShouldReturnValidInstance(t *testing.T)
func TestAppData_AddTask_ShouldIncreaseTaskCount(t *testing.T)
func TestAppData_GetTaskByID_WithNonexistentID_ShouldReturnError(t *testing.T)
```

#### 2.3.2 GREEN: AppData基本実装
- AppData構造体定義
- 基本的なタスク操作メソッド

#### 2.3.3 REFACTOR: AppData機能強化
- 検索機能の追加
- フィルタリング機能
- エラーハンドリング改善

## フェーズ3: Repository層（TDD）

### 3.1 Repository Interface実装
- **優先度**: 高 | **推定時間**: 45分

#### 3.1.1 RED: Repository Interfaceテスト
```go
func TestRepository_Interface_ShouldDefineRequiredMethods(t *testing.T)
```

#### 3.1.2 GREEN: Repository Interface定義
- Repository interfaceの定義
- 必要なメソッドシグネチャ

#### 3.1.3 REFACTOR: Interface設計改善
- エラーハンドリング統一
- コンテキスト対応検討

### 3.2 FileRepository実装
- **優先度**: 高 | **推定時間**: 120分

#### 3.2.1 RED: FileRepositoryテスト
```go
func TestFileRepository_Save_WithValidData_ShouldCreateFile(t *testing.T)
func TestFileRepository_Load_WithNonexistentFile_ShouldReturnError(t *testing.T)
func TestFileRepository_Load_WithCorruptedFile_ShouldReturnError(t *testing.T)
func TestFileRepository_Backup_ShouldCreateBackupFile(t *testing.T)
```

#### 3.2.2 GREEN: FileRepository基本実装
- ファイル読み書き基本機能
- エラーハンドリング
- ディレクトリ作成機能

#### 3.2.3 REFACTOR: FileRepository改善
- アトミックな書き込み実装
- バックアップ機能強化
- ファイルロック機能
- パフォーマンス最適化

## フェーズ4: バリデーション層（TDD）

### 4.1 Validator実装
- **優先度**: 高 | **推定時間**: 105分

#### 4.1.1 RED: Validatorテスト
```go
func TestValidator_ValidateTask_WithValidTask_ShouldReturnNil(t *testing.T)
func TestValidator_ValidateTask_WithEmptyTitle_ShouldReturnError(t *testing.T)
func TestValidator_ValidateTask_WithTooLongTitle_ShouldReturnError(t *testing.T)
func TestValidator_ValidateTask_WithInvalidStatus_ShouldReturnError(t *testing.T)
```

#### 4.1.2 GREEN: Validator基本実装
- 基本バリデーション機能
- エラーメッセージ生成

#### 4.1.3 REFACTOR: Validator改善
- カスタムバリデーションルール
- 国際化対応
- バリデーション結果詳細化

## フェーズ5: ビジネスロジック層（TDD）

### 5.1 TaskService実装
- **優先度**: 高 | **推定時間**: 180分

#### 5.1.1 RED: TaskServiceテスト
```go
func TestTaskService_CreateTask_WithValidRequest_ShouldReturnTask(t *testing.T)
func TestTaskService_CreateTask_WithInvalidRequest_ShouldReturnError(t *testing.T)  
func TestTaskService_UpdateTask_WithValidData_ShouldUpdateTask(t *testing.T)
func TestTaskService_UpdateTask_WithNonexistentID_ShouldReturnError(t *testing.T)
func TestTaskService_DeleteTask_WithValidID_ShouldRemoveTask(t *testing.T)
func TestTaskService_ToggleTask_ShouldChangeStatus(t *testing.T)
func TestTaskService_SearchTasks_WithQuery_ShouldReturnMatchingTasks(t *testing.T)
```

#### 5.1.2 GREEN: TaskService基本実装
- CRUD操作の基本実装
- Repository連携
- 基本的なビジネスルール

#### 5.1.3 REFACTOR: TaskService改善
- エラーハンドリング統一
- ログ機能追加
- パフォーマンス最適化
- トランザクション処理

### 5.2 StateManager実装
- **優先度**: 高 | **推定時間**: 135分

#### 5.2.1 RED: StateManagerテスト
```go
func TestStateManager_SetTasks_ShouldNotifySubscribers(t *testing.T)
func TestStateManager_Subscribe_ShouldReceiveNotifications(t *testing.T)
func TestStateManager_SetFilter_ShouldUpdateFilteredTasks(t *testing.T)
func TestStateManager_ConcurrentAccess_ShouldBeSafe(t *testing.T)
```

#### 5.2.2 GREEN: StateManager基本実装
- 状態管理基本機能
- 変更通知機能
- スレッドセーフ実装

#### 5.2.3 REFACTOR: StateManager改善
- イベント系統の細分化
- パフォーマンス最適化
- メモリリーク対策

## フェーズ6: UI層（TDD）

### 6.1 テーマシステム実装
- **優先度**: 中 | **推定時間**: 90分

#### 6.1.1 RED: テーマテスト
```go
func TestTheme_GetPriorityColor_ShouldReturnCorrectColor(t *testing.T)
func TestTheme_GetStatusColor_ShouldReturnCorrectColor(t *testing.T)
```

#### 6.1.2 GREEN: テーマ基本実装
- 基本カラーパレット
- 優先度・ステータス色マッピング

#### 6.1.3 REFACTOR: テーマ改善
- 複数テーマ対応
- カスタムテーマ機能

### 6.2 TaskListWidget実装
- **優先度**: 高 | **推定時間**: 165分

#### 6.2.1 RED: TaskListWidgetテスト
```go
func TestTaskListWidget_SetTasks_ShouldDisplayTasks(t *testing.T)
func TestTaskListWidget_SelectNext_ShouldMoveSelection(t *testing.T)
func TestTaskListWidget_ApplyFilter_ShouldShowOnlyMatchingTasks(t *testing.T)
```

#### 6.2.2 GREEN: TaskListWidget基本実装
- tviewテーブル基本機能
- タスク表示
- 選択機能

#### 6.2.3 REFACTOR: TaskListWidget改善
- 表示フォーマット改善
- ソート機能
- 仮想化対応

### 6.3 InputFormWidget実装
- **優先度**: 高 | **推定時間**: 135分

#### 6.3.1 RED: InputFormWidgetテスト  
```go
func TestInputFormWidget_CreateMode_ShouldAllowInput(t *testing.T)
func TestInputFormWidget_EditMode_ShouldPrePopulateFields(t *testing.T)
func TestInputFormWidget_Validate_WithInvalidInput_ShouldShowError(t *testing.T)
```

#### 6.3.2 GREEN: InputFormWidget基本実装
- フォーム基本機能
- 入力フィールド
- 送信処理

#### 6.3.3 REFACTOR: InputFormWidget改善
- バリデーション表示
- UX改善
- キーボードナビゲーション

### 6.4 メインApp実装
- **優先度**: 高 | **推定時間**: 210分

#### 6.4.1 RED: Appテスト
```go
func TestApp_Initialize_ShouldSetupComponents(t *testing.T)
func TestApp_HandleKeyPress_ShouldRouteToCorrectHandler(t *testing.T)
func TestApp_SwitchView_ShouldChangeActiveWidget(t *testing.T)
```

#### 6.4.2 GREEN: App基本実装
- tviewアプリケーション初期化
- 基本レイアウト
- イベントハンドリング

#### 6.4.3 REFACTOR: App改善
- エラーハンドリング
- グレースフル終了
- パフォーマンス最適化

## フェーズ7: CLI統合（TDD）

### 7.1 Cobraコマンド実装
- **優先度**: 高 | **推定時間**: 90分

#### 7.1.1 RED: CLIテスト
```go
func TestCLI_RunCommand_ShouldStartApp(t *testing.T)
func TestCLI_ParseFlags_ShouldSetCorrectOptions(t *testing.T)
```

#### 7.1.2 GREEN: CLI基本実装
- Cobraコマンド定義
- フラグ処理
- アプリケーション起動

#### 7.1.3 REFACTOR: CLI改善
- ヘルプ文改善
- エラーメッセージ改善
- 設定オプション追加

## フェーズ8: 統合テスト・E2Eテスト

### 8.1 統合テスト実装
- **優先度**: 高 | **推定時間**: 120分

#### 統合テストケース
```go
func TestIntegration_TaskLifecycle_ShouldWorkEndToEnd(t *testing.T)
func TestIntegration_FileStorage_ShouldPersistData(t *testing.T)
func TestIntegration_ErrorRecovery_ShouldHandleCorruption(t *testing.T)
```

### 8.2 E2Eテスト実装
- **優先度**: 中 | **推定時間**: 90分

#### E2Eテストケース
```go
func TestE2E_NewUserWorkflow_ShouldCompleteTaskManagement(t *testing.T)
func TestE2E_PowerUserWorkflow_ShouldHandleComplexScenarios(t *testing.T)
```

## フェーズ9: 最終調整

### 9.1 パフォーマンス最適化（TDD）
- **優先度**: 中 | **推定時間**: 75分

#### 9.1.1 RED: パフォーマンステスト
- レスポンス時間テスト
- メモリ使用量テスト
- 大量データ処理テスト

#### 9.1.2 GREEN: 基本最適化
- 明らかなボトルネック解消

#### 9.1.3 REFACTOR: 詳細最適化
- プロファイリング結果に基づく改善

### 9.2 ドキュメント・デプロイ準備
- **優先度**: 低 | **推定時間**: 90分
- README.md作成
- ビルドスクリプト整備
- リリース準備

## TDD実装スケジュール

### 推奨開発順序

1. **フェーズ1-2**: データ基盤（3-4時間）
2. **フェーズ3-4**: 永続化・バリデーション（4-5時間） 
3. **フェーズ5**: ビジネスロジック（5-6時間）
4. **フェーズ6**: UI実装（7-8時間）
5. **フェーズ7-8**: 統合・テスト（3-4時間）
6. **フェーズ9**: 最終調整（3時間）

**総推定時間: 25-30時間**

## TDD成功指標

- [ ] 全機能でRed-Green-Refactorサイクル完了
- [ ] ユニットテストカバレッジ90%以上
- [ ] 統合テストカバレッジ80%以上  
- [ ] 全テスト実行時間30秒以下
- [ ] コードレビューでTDD品質確認済み

各タスクでは必ずテストファーストを徹底し、小さなサイクルで確実に品質を積み上げていきます。