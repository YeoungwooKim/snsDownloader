package media

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

/*
여기서 바디데이터 체크(미디어 영상/음성 옵션)하는데
암것도 없으면 베스트로
*/
func GetMedia(c *fiber.Ctx) error {
	dataMap := make(map[string]interface{})
	json.Unmarshal(c.Body(), &dataMap)
	uuid := utils.UUIDv4()
	// fileLocation := ""

	if progress, err := ExecuteMedia(dataMap["uri"].(string), dataMap); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"code": http.StatusInternalServerError,
			"err":  fmt.Sprintf("%v", err),
		})
	} else {
		go func() {
			for msg := range progress {
				if dataMap := ProcessMessage(uuid, msg); dataMap["location"] != nil {
					// fileLocation = fmt.Sprintf("%v%v", dataMap["location"], dataMap["file_name"])
				}
				// colorLog.Info("%v", msg)
			}
			commitProgress(uuid)
		}()
	}
	delete(dataMap, "location")
	// colorLog.Info("%v", fileLocation)

	dataMap["uuid"] = uuid
	return c.Status(http.StatusCreated).JSON(dataMap)
}

// func GetAll(c *fiber.Ctx) error {
// 	var dataMapList []map[string]interface{}
// 	var err error

// 	if dataMapList, err = getNotCompleteHistory(); err != nil {
// 		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
// 			"code": http.StatusInternalServerError,
// 			"err":  err.Error(),
// 		})
// 	}

// 	return c.Status(http.StatusOK).JSON(fiber.Map{
// 		"code": http.StatusOK,
// 		"list": dataMapList,
// 	})
// }

func GetMediaStatus(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	if uuid == "" {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"code": http.StatusInternalServerError,
			"err":  "requested uuid is not exist or cant mapped.. wait for second..",
		})
	}
	if dataMap, err := getLatestHistory(uuid); err != nil {
		return c.Status(http.StatusConflict).JSON(fiber.Map{
			"code": http.StatusConflict,
			"err":  fmt.Sprintf("%v", err),
		})
	} else {
		return c.Status(http.StatusOK).JSON(dataMap)
	}
}

// func GetMedia(c *fiber.Ctx) error {
// 	dataMap := make(map[string]interface{})
// 	json.Unmarshal(c.Body(), &dataMap)
// 	uuid := utils.UUIDv4()
// 	fileLocation := ""

// 	if progress, err := ExecuteMedia(dataMap["uri"].(string), dataMap); err != nil {
// 		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
// 			"code": http.StatusInternalServerError,
// 			"err":  fmt.Sprintf("%v", err),
// 		})
// 	} else {
// 		for msg := range progress {
// 			if dataMap := ProcessMessage(uuid, msg); dataMap["location"] != nil {
// 				fileLocation = fmt.Sprintf("%v%v", dataMap["location"], dataMap["file_name"])
// 			}
// 			// colorLog.Info("%v", msg)
// 		}
// 	}
// 	// colorLog.Info("%v", fileLocation)

// 	return c.Download(fileLocation)
// }
