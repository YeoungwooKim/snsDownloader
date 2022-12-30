package media

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

/**
여기서 바디데이터 체크(미디어 영상/음성 옵션)하는데
암것도 없으면 베스트로
*/
func GetMedia(c *fiber.Ctx) error {
	dataMap := make(map[string]interface{})
	json.Unmarshal(c.Body(), &dataMap)
	uuid := utils.UUIDv4()
	fileLocation := ""

	if progress, err := ExecuteMedia(dataMap["uri"].(string), dataMap); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"code": http.StatusInternalServerError,
			"err":  fmt.Sprintf("%v", err),
		})
	} else {
		for msg := range progress {
			if dataMap := ProcessMessage(uuid, msg); dataMap["location"] != nil {
				fileLocation = fmt.Sprintf("%v%v", dataMap["location"], dataMap["file_name"])
			}
			// colorLog.Info("%v", msg)
		}
	}
	// colorLog.Info("%v", fileLocation)

	return c.Download(fileLocation)
}
