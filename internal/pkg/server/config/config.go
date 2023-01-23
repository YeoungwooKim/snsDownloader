package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"snsDownloader/internal/pkg/server/dbconn"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
)

type EnvConfig struct {
	HttpPort    int    `env:"HTTP_PORT" default:"8080"`
	DbUri       string `env:"DB_URI" default:""`
	QueueUri    string `env:"QUEUE_URI" default:"http://localhost:8090/api/v1/queue"`
	FiberConfig fiber.Config
}

func NewConfig() EnvConfig {
	pwd, _ := os.Getwd()
	templateEngine := html.New(fmt.Sprintf("%s/internal/pkg/views", pwd), ".html")
	dbConfig := getDbAuth()
	return EnvConfig{
		HttpPort: 8080,
		DbUri:    fmt.Sprintf("mongodb://%v:%v@%v:%v/?authMechanism=SCRAM-SHA-256&ssl=false", dbConfig["username"], dbConfig["password"], dbConfig["hostname"], dbConfig["port"]),
		QueueUri: "http://localhost:8090/api/v1/queue",
		FiberConfig: fiber.Config{
			CaseSensitive:     false,
			ColorScheme:       fiber.DefaultColors,
			EnablePrintRoutes: true,
			ErrorHandler:      fiber.DefaultErrorHandler,
			WriteTimeout:      time.Second * 10,
			Views:             templateEngine,
		},
	}

}

var Config = NewConfig()

func Load() {
	fmt.Printf("\tHttpPort\t:%v\n", Config.HttpPort)
	fmt.Printf("\tDbUri\t:%v\n", Config.DbUri)

	dbconn.Create(Config.DbUri)
	// cron.InitCron(3, Config.QueueUri)
	// go kafka.ConsumeMessage()
}

func getDbAuth() map[string]string {
	pwd, _ := os.Getwd()
	data, err := os.Open(fmt.Sprintf("%v/mongodb_auth_config.json", pwd))
	if err != nil {
		fmt.Printf("%v\n", err)
		return nil
	}
	var dbConfig map[string]string
	byteValue, _ := ioutil.ReadAll(data)
	json.Unmarshal(byteValue, &dbConfig)
	return dbConfig
}
