package server

import (
	"github.com/stretchr/testify/assert"
	"testing"

	jwtMock "github.com/samgozman/go-bloggy/mocks/jwt"
)

func TestProvideServer(t *testing.T) {
	t.Run("ProvideServer", func(t *testing.T) {
		// Arrange
		jwtService := jwtMock.NewMockServiceInterface(t)

		// Act
		got := ProvideServer(jwtService)

		// Assert
		assert.NotNil(t, got)
	})
}
