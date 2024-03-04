package handler

import (
	"github.com/oapi-codegen/testutil"
	"github.com/samgozman/go-bloggy/pkg/server"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func Test_GetHealth(t *testing.T) {
	e, _, _ := registerHandlers(nil, nil)

	res := testutil.NewRequest().Get("/health").GoWithHTTPHandler(t, e)

	var body server.HealthCheckResponse
	err := res.UnmarshalBodyToObject(&body)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, res.Code())
	assert.Equal(t, "OK", body.Status)
}