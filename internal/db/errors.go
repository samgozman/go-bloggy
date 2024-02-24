package db

import "errors"

var (
	ErrFailedToConnectDatabase             = errors.New("ERR_FAILED_TO_CONNECT_DATABASE")
	ErrFailedToMigrateDatabase             = errors.New("ERR_FAILED_TO_MIGRATE_DATABASE")
	ErrFailedToGetDatabaseConnection       = errors.New("ERR_FAILED_TO_GET_DATABASE_CONNECTION")
	ErrFailedToEnableForeignKeyConstraints = errors.New("ERR_FAILED_TO_ENABLE_FOREIGN_KEY_CONSTRAINTS")
)
