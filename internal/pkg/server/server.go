package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"snsDownloader/internal/pkg/media"
	metadata "snsDownloader/internal/pkg/metaData"
	"snsDownloader/internal/pkg/server/config"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func initialize() {
	fmt.Printf("=============================================\n")
	fmt.Printf("Initialize Server.\n")
	fmt.Printf("=============================================\n")

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
		dataMap := make(map[string]interface{})
		err := json.Unmarshal(c.Body(), &dataMap)
		if err == nil {
			shellInjections := []string{"&", "&&", "|", "||", ";", "0x0a", "\n", "$", "`", "'", `"`}
			for _, value := range dataMap {
				switch value.(type) {
				case string:
					// fmt.Printf("will be checked %v\n", value)
				default:
					// fmt.Printf("passed %v\n", value)
					continue
				}
				for _, injection := range shellInjections {
					if strings.Contains(value.(string), injection) {
						return c.Status(http.StatusConflict).JSON(fiber.Map{
							"code": http.StatusConflict,
							"err":  fmt.Sprintf("chracter(=%v) can be a command injection(=%v)", injection, value),
						})
					}
				}
			}
		}
		return c.Next()
	})

	videos := v1.Group("/media")
	videos.Post("", media.GetMedia)
	videos.Post("/:filename", media.DownloadMedia)
	// videos.Get("", media.GetAll)
	videos.Get("/:uuid<guid>", media.GetMediaStatus)
	videos.Delete("/:uuid<guid>", media.StopTask)

	videosInfo := v1.Group("/meta")
	videosInfo.Use(metadata.ValidationRouter)
	videosInfo.Post("", metadata.GetMetaData)

	app.Listen(":8080")
	return app
}
