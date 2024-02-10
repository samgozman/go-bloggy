package main

import (
	"github.com/labstack/echo/v4"
	"github.com/samgozman/go-bloggy/internal/handler"
	"github.com/samgozman/go-bloggy/pkg/client"
)

func main() {
	e := echo.New()
	h := handler.NewHandler()

	client.RegisterHandlers(e, h)

	e.Logger.Fatal(e.Start(":80"))
}
