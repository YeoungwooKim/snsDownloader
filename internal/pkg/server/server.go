package server

import (
	"fmt"
	"headless/internal/pkg/colorLog"
	"headless/internal/pkg/media"
	metadata "headless/internal/pkg/metaData"
	"headless/internal/pkg/server/config"

	"github.com/gofiber/fiber/v2"
)

func initialize() {
	colorLog.Info("=============================================")
	colorLog.Info("Initialize Server.")
	colorLog.Info("=============================================")

	config.Load()
}

func Create() *fiber.App {
	initialize()

	app := fiber.New(config.Config.FiberConfig)
	app.Static("/", "internal/pkg/views")
	app.Get("/home", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"Title": "This is title from layout/main",
		}, "layouts/main")
	})

	api := app.Group("/api")

	v1 := api.Group("/v1")
	v1.Use(func(c *fiber.Ctx) error {
		// uri 확인용.
		fmt.Printf("%s\n", c.Request().URI())
		return c.Next()
	})

	videos := v1.Group("/media")
	videos.Post("", media.GetMedia)

	videosInfo := v1.Group("/meta")
	videosInfo.Use(metadata.ValidationRouter)
	videosInfo.Post("", metadata.GetMetaData)

	app.Listen(":8080")
	return app
}
