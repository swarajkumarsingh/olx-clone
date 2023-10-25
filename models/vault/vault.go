package vault

import (
	"context"
	"olx-clone/infra/db"
	"time"

	"github.com/jmoiron/sqlx"
)

var database = db.Mgr.DBConn

func GetDBTransaction() (*sqlx.Tx, error) {
	return database.Beginx()
}

func Get(context context.Context, key string) (Struct, error) {
	var keyValueData Struct
	query := `SELECT key, value, created_at, updated_at FROM vault where key = $1`
	err := database.GetContext(context, &keyValueData, query, key)
	return keyValueData, err
}

func Create(context context.Context, key string, value string, tx *sqlx.Tx) (*Struct, error) {
	keyValueData := Struct{
		Key:       key,
		Value:     value,
		UpdatedAt: time.Now(),
	}
	query := `INSERT INTO vault (key, value, updated_at) VALUES (:key, :value, :updated_at) ON CONFLICT (key) DO UPDATE SET value = :value, updated_at = :updated_at;`
	if tx != nil {
		_, err := tx.NamedExecContext(context, query, keyValueData)
		return &keyValueData, err
	}
	_, err := database.NamedExecContext(context, query, keyValueData)
	return &keyValueData, err
}

func DoesKeyExists(context context.Context, key string) (bool, error) {
	var doesExists bool
	err := database.GetContext(context, &doesExists, "SELECT COUNT(*)>0 FROM vault where key = $1", key)
	return doesExists, err
}
