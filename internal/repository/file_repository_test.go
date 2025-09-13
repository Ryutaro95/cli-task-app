package repository

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"task-cli/internal/model"

	"github.com/stretchr/testify/assert"
)

// RED: FileRepositoryのテスト
func TestFileRepository_New_ShouldCreateRepository(t *testing.T) {
	// Given
	tempDir := "/tmp/test_tasks"

	// When
	repo := NewFileRepository(tempDir)

	// Then
	assert.NotNil(t, repo)
	assert.Implements(t, (*Repository)(nil), repo)
}

func TestFileRepository_Save_WithValidData_ShouldCreateFile(t *testing.T) {
	// Given
	tempDir, err := ioutil.TempDir("", "task_test_")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	repo := NewFileRepository(tempDir)
	appData := model.NewAppData()
	task, _ := model.NewTask("Test Task", "Description", model.PriorityMedium, []string{"test"})
	appData.AddTask(task)
	ctx := context.Background()

	// When
	err = repo.Save(ctx, appData)

	// Then
	assert.NoError(t, err)
	
	// ファイルが作成されているかチェック
	dataFile := filepath.Join(tempDir, "tasks.json")
	_, err = os.Stat(dataFile)
	assert.NoError(t, err, "tasks.json file should exist")
}

func TestFileRepository_Load_WithExistingFile_ShouldReturnData(t *testing.T) {
	// Given
	tempDir, err := ioutil.TempDir("", "task_test_")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	repo := NewFileRepository(tempDir)
	ctx := context.Background()

	// 先にデータを保存
	originalData := model.NewAppData()
	task, _ := model.NewTask("Existing Task", "Description", model.PriorityHigh, []string{"existing"})
	originalData.AddTask(task)
	repo.Save(ctx, originalData)

	// When
	loadedData, err := repo.Load(ctx)

	// Then
	assert.NoError(t, err)
	assert.NotNil(t, loadedData)
	assert.Len(t, loadedData.Tasks, 1)
	assert.Equal(t, "Existing Task", loadedData.Tasks[0].Title)
}

func TestFileRepository_Load_WithNonexistentFile_ShouldReturnError(t *testing.T) {
	// Given
	tempDir, err := ioutil.TempDir("", "task_test_")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	repo := NewFileRepository(tempDir)
	ctx := context.Background()

	// When
	loadedData, err := repo.Load(ctx)

	// Then
	assert.Error(t, err)
	assert.Nil(t, loadedData)
}

func TestFileRepository_Load_WithCorruptedFile_ShouldReturnError(t *testing.T) {
	// Given
	tempDir, err := ioutil.TempDir("", "task_test_")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// 壊れたJSONファイルを作成
	dataFile := filepath.Join(tempDir, "tasks.json")
	err = ioutil.WriteFile(dataFile, []byte("invalid json content"), 0644)
	assert.NoError(t, err)

	repo := NewFileRepository(tempDir)
	ctx := context.Background()

	// When
	loadedData, err := repo.Load(ctx)

	// Then
	assert.Error(t, err)
	assert.Nil(t, loadedData)
}

func TestFileRepository_CreateBackup_ShouldCreateBackupFile(t *testing.T) {
	// Given
	tempDir, err := ioutil.TempDir("", "task_test_")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	repo := NewFileRepository(tempDir)
	appData := model.NewAppData()
	task, _ := model.NewTask("Backup Task", "Description", model.PriorityLow, []string{"backup"})
	appData.AddTask(task)
	ctx := context.Background()

	// When
	backupPath, err := repo.CreateBackup(ctx, appData)

	// Then
	assert.NoError(t, err)
	assert.NotEmpty(t, backupPath)
	
	// バックアップファイルが存在するかチェック
	_, err = os.Stat(backupPath)
	assert.NoError(t, err, "Backup file should exist")
}

func TestFileRepository_RestoreFromBackup_WithValidBackup_ShouldRestoreData(t *testing.T) {
	// Given
	tempDir, err := ioutil.TempDir("", "task_test_")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	repo := NewFileRepository(tempDir)
	ctx := context.Background()

	// バックアップを作成
	originalData := model.NewAppData()
	task, _ := model.NewTask("Backup Task", "Description", model.PriorityLow, []string{"backup"})
	originalData.AddTask(task)
	backupPath, err := repo.CreateBackup(ctx, originalData)
	assert.NoError(t, err)

	// When
	restoredData, err := repo.RestoreFromBackup(ctx, backupPath)

	// Then
	assert.NoError(t, err)
	assert.NotNil(t, restoredData)
	assert.Len(t, restoredData.Tasks, 1)
	assert.Equal(t, "Backup Task", restoredData.Tasks[0].Title)
}

func TestFileRepository_Save_WithNilData_ShouldReturnError(t *testing.T) {
	// Given
	tempDir, err := ioutil.TempDir("", "task_test_")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	repo := NewFileRepository(tempDir)
	ctx := context.Background()

	// When
	err = repo.Save(ctx, nil)

	// Then
	assert.Error(t, err)
}

func TestFileRepository_CreateBackup_WithNilData_ShouldReturnError(t *testing.T) {
	// Given
	tempDir, err := ioutil.TempDir("", "task_test_")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	repo := NewFileRepository(tempDir)
	ctx := context.Background()

	// When
	backupPath, err := repo.CreateBackup(ctx, nil)

	// Then
	assert.Error(t, err)
	assert.Empty(t, backupPath)
}