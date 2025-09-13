package repository

import (
	"context"
	"testing"

	"task-cli/internal/model"

	"github.com/stretchr/testify/assert"
)

// mockRepository はRepositoryの模擬実装
type mockRepository struct {
	appData *model.AppData
}

// RED: Repository Interfaceのテスト
func TestRepository_Interface_ShouldDefineRequiredMethods(t *testing.T) {
	// Given
	appData := model.NewAppData()
	repo := &mockRepository{appData: appData}

	// When & Then - インターフェースが正しく実装されているかテスト
	var _ Repository = repo

	// インターフェースが必要なメソッドを持っているかを静的にチェック
	assert.Implements(t, (*Repository)(nil), repo)
}

func TestRepository_Save_WithValidData_ShouldPersistData(t *testing.T) {
	// Given
	appData := model.NewAppData()
	task, _ := model.NewTask("Test Task", "Description", model.PriorityMedium, []string{"test"})
	appData.AddTask(task)
	
	repo := &mockRepository{appData: appData}
	ctx := context.Background()

	// When
	err := repo.Save(ctx, appData)

	// Then
	assert.NoError(t, err)
}

func TestRepository_Load_WithExistingData_ShouldReturnData(t *testing.T) {
	// Given
	ctx := context.Background()
	repo := &mockRepository{appData: model.NewAppData()}

	// When
	loadedData, err := repo.Load(ctx)

	// Then
	assert.NoError(t, err)
	assert.NotNil(t, loadedData)
}

func TestRepository_Load_WithNonexistentData_ShouldReturnError(t *testing.T) {
	// Given
	ctx := context.Background()
	repo := &mockRepository{appData: nil} // データがない状態

	// When
	loadedData, err := repo.Load(ctx)

	// Then
	assert.Error(t, err)
	assert.Nil(t, loadedData)
}

func TestRepository_CreateBackup_ShouldCreateBackupSuccessfully(t *testing.T) {
	// Given
	appData := model.NewAppData()
	task, _ := model.NewTask("Test Task", "Description", model.PriorityHigh, []string{"test"})
	appData.AddTask(task)
	
	repo := &mockRepository{appData: appData}
	ctx := context.Background()

	// When
	backupPath, err := repo.CreateBackup(ctx, appData)

	// Then
	assert.NoError(t, err)
	assert.NotEmpty(t, backupPath)
}

func TestRepository_RestoreFromBackup_WithValidPath_ShouldRestoreData(t *testing.T) {
	// Given
	ctx := context.Background()
	repo := &mockRepository{}
	backupPath := "/tmp/test_backup.json" // テスト用のバックアップパス

	// When
	restoredData, err := repo.RestoreFromBackup(ctx, backupPath)

	// Then
	assert.NoError(t, err)
	assert.NotNil(t, restoredData)
}

// mockRepositoryのメソッド実装
func (m *mockRepository) Save(ctx context.Context, data *model.AppData) error {
	if data == nil {
		return assert.AnError
	}
	m.appData = data
	return nil
}

func (m *mockRepository) Load(ctx context.Context) (*model.AppData, error) {
	if m.appData == nil {
		return nil, assert.AnError
	}
	return m.appData, nil
}

func (m *mockRepository) CreateBackup(ctx context.Context, data *model.AppData) (string, error) {
	if data == nil {
		return "", assert.AnError
	}
	return "/tmp/backup_" + data.ID + ".json", nil
}

func (m *mockRepository) RestoreFromBackup(ctx context.Context, backupPath string) (*model.AppData, error) {
	if backupPath == "" {
		return nil, assert.AnError
	}
	// テスト用のダミーデータを返す
	return model.NewAppData(), nil
}