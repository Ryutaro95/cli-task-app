package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"task-cli/internal/model"
)

// FileRepository はファイルベースのRepository実装
type FileRepository struct {
	dataDir  string
	fileName string
}

// NewFileRepository は新しいFileRepositoryを作成する
func NewFileRepository(dataDir string) Repository {
	return &FileRepository{
		dataDir:  dataDir,
		fileName: "tasks.json",
	}
}

// Save はAppDataをファイルに保存する
func (f *FileRepository) Save(ctx context.Context, data *model.AppData) error {
	if data == nil {
		return errors.New("data cannot be nil")
	}

	// ディレクトリが存在しない場合は作成
	if err := f.ensureDataDir(); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	// JSONにエンコード
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	// ファイルに書き込み
	filePath := f.getDataFilePath()
	err = ioutil.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// Load はファイルからAppDataを読み込む
func (f *FileRepository) Load(ctx context.Context) (*model.AppData, error) {
	filePath := f.getDataFilePath()

	// ファイルが存在するかチェック
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, errors.New("data file does not exist")
	}

	// ファイルを読み込み
	jsonData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// JSONをデコード
	var appData model.AppData
	err = json.Unmarshal(jsonData, &appData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal data: %w", err)
	}

	return &appData, nil
}

// CreateBackup はデータのバックアップを作成する
func (f *FileRepository) CreateBackup(ctx context.Context, data *model.AppData) (string, error) {
	if data == nil {
		return "", errors.New("data cannot be nil")
	}

	// バックアップディレクトリを確保
	backupDir := filepath.Join(f.dataDir, "backups")
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create backup directory: %w", err)
	}

	// タイムスタンプを含むバックアップファイル名を生成
	timestamp := time.Now().Format("20060102_150405")
	backupFileName := fmt.Sprintf("tasks_backup_%s.json", timestamp)
	backupPath := filepath.Join(backupDir, backupFileName)

	// JSONにエンコード
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal backup data: %w", err)
	}

	// バックアップファイルに書き込み
	err = ioutil.WriteFile(backupPath, jsonData, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write backup file: %w", err)
	}

	return backupPath, nil
}

// RestoreFromBackup はバックアップからデータを復元する
func (f *FileRepository) RestoreFromBackup(ctx context.Context, backupPath string) (*model.AppData, error) {
	if backupPath == "" {
		return nil, errors.New("backup path cannot be empty")
	}

	// バックアップファイルが存在するかチェック
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("backup file does not exist: %s", backupPath)
	}

	// バックアップファイルを読み込み
	jsonData, err := ioutil.ReadFile(backupPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read backup file: %w", err)
	}

	// JSONをデコード
	var appData model.AppData
	err = json.Unmarshal(jsonData, &appData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal backup data: %w", err)
	}

	return &appData, nil
}

// ensureDataDir はデータディレクトリが存在することを確認し、必要に応じて作成する
func (f *FileRepository) ensureDataDir() error {
	return os.MkdirAll(f.dataDir, 0755)
}

// getDataFilePath はデータファイルのフルパスを返す
func (f *FileRepository) getDataFilePath() string {
	return filepath.Join(f.dataDir, f.fileName)
}