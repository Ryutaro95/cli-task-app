package repository

import (
	"context"

	"task-cli/internal/model"
)

// Repository はデータの永続化を担当するインターフェース
type Repository interface {
	// Save はAppDataを保存する
	Save(ctx context.Context, data *model.AppData) error

	// Load は保存されたAppDataを読み込む
	Load(ctx context.Context) (*model.AppData, error)

	// CreateBackup はデータのバックアップを作成する
	CreateBackup(ctx context.Context, data *model.AppData) (string, error)

	// RestoreFromBackup はバックアップからデータを復元する
	RestoreFromBackup(ctx context.Context, backupPath string) (*model.AppData, error)
}