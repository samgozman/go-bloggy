package db

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestInitDatabase(t *testing.T) {
	testDSN := "init_database_test.db"

	t.Cleanup(func() {
		// Clean up the database file after the test
		err := os.Remove(testDSN)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("should return error if failed to connect to database", func(t *testing.T) {
		_, err := InitDatabase("//error")
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrFailedToConnectDatabase)
	})

	t.Run("should return database connection", func(t *testing.T) {
		db, err := InitDatabase(testDSN)
		assert.NoError(t, err)
		assert.NotNil(t, db)
		assert.NotNil(t, db.conn)
		assert.NotNil(t, db.Models)
		assert.NotNil(t, db.Models.Users)
	})
}
