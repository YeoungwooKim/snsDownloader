package config

import (
	"encoding/json"
	"fmt"
	"headless/internal/pkg/colorLog"
	"headless/internal/pkg/server/dbconn"
	"io/ioutil"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
)

type EnvConfig struct {
	HttpPort    int    `env:"HTTP_PORT" default:"8080"`
	DbUri       string `env:"DB_URI" default:""`
	FiberConfig fiber.Config
}

func NewConfig() EnvConfig {
	pwd, _ := os.Getwd()
	templateEngine := html.New(fmt.Sprintf("%s/internal/pkg/views", pwd), ".html")
	dbConfig := getDbAuth()
	return EnvConfig{
		HttpPort: 8080,
		DbUri:    fmt.Sprintf("mongodb://%v:%v@%v:%v/?authMechanism=SCRAM-SHA-256&ssl=false", dbConfig["username"], dbConfig["password"], dbConfig["hostname"], dbConfig["port"]),
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
	colorLog.Info("\tHttpPort\t:%v", Config.HttpPort)
	colorLog.Info("\tDbUri\t:%v", Config.DbUri)

	dbconn.Create(Config.DbUri)
}

func getDbAuth() map[string]string {
	pwd, _ := os.Getwd()
	data, err := os.Open(fmt.Sprintf("%v/mongodb_auth_config.json", pwd))
	if err != nil {
		colorLog.Info("%v", err)
		return nil
	}
	var dbConfig map[string]string
	byteValue, _ := ioutil.ReadAll(data)
	json.Unmarshal(byteValue, &dbConfig)
	return dbConfig
}
