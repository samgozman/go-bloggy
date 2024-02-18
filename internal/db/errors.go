package db

import "errors"

var (
	ErrFailedToConnectDatabase = errors.New("ERR_FAILED_TO_CONNECT_DATABASE")
	ErrFailedToMigrateDatabase = errors.New("ERR_FAILED_TO_MIGRATE_DATABASE")
)
