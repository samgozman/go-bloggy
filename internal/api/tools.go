//go:build tools

package api

import (
	_ "github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen"
)

// This file is used to declare the tools.go build tag. It is used to include the
// oapi-codegen tool in the go.mod file. This is necessary to ensure that the
// oapi-codegen tool is available to the go generate command.
