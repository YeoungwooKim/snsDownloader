package server

import (
	"fmt"
	"snsDownload/internal/pkg/media"
	metadata "snsDownload/internal/pkg/metaData"
	"snsDownload/internal/pkg/server/config"

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
		// shellInjections := []string{"&", "&&", "|", "||", ";", "0x0a", "\n", "$", "`", "'", `"`}
		fmt.Printf("%s\n", c.Request().URI())
		// dataMap := make(map[string]interface{})
		// json.Unmarshal(c.Body(), &dataMap)

		// if util.IsAlreadyShorten(dataMap["uri"].(string)) == true {
		// 	fmt.Printf("already short-uri %v\n", dataMap["uri"].(string))
		// 	return c.Next()
		// }
		// if shortenUrl, err := util.ShortenUrl(dataMap["uri"].(string)); err != nil {
		// 	return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
		// 		"code":        http.StatusInternalServerError,
		// 		"description": err,
		// 	})
		// } else {
		// 	fmt.Printf("origin %v\n", dataMap["uri"])
		// 	dataMap["uri"] = shortenUrl
		// }
		// fmt.Printf("shorten %v\n", dataMap["uri"])

		return c.Next()
	})

	videosInfo := v1.Group("/meta")
	videosInfo.Use(metadata.ValidationRouter)
	videosInfo.Post("", metadata.GetMetaData)

	videos := v1.Group("/media")
	videos.Post("", media.GetMedia)
	videos.Post("/:filename", media.DownloadMedia)
	// videos.Get("", media.GetAll)
	videos.Get("/:uuid<guid>", media.GetMediaStatus)
	videos.Delete("/:uuid<guid>", media.StopTask)

	// pwd, _ := os.Getwd()
	//8443
	// app.ListenTLS(":8443", fmt.Sprintf("%v/cert.pem", pwd), fmt.Sprintf("%v/cert.key", pwd))
	app.Listen(":8443")

	return app
}
