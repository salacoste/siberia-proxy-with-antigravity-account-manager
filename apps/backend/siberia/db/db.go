package db

import (
	"path/filepath"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	Conn *gorm.DB
}

func Init(configDir string, key string) (*Database, error) {
	SetMasterKey(key)

	dbPath := filepath.Join(configDir, "siberia.db")

	// Open SQLite connection
	// Using glebarez/sqlite for pure Go implementation
	conn, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	// AutoMigrate models
	if err := conn.AutoMigrate(&Account{}); err != nil {
		return nil, err
	}

	return &Database{Conn: conn}, nil
}
