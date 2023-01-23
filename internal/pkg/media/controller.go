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

	transcoder := New()
	if progress, err := transcoder.ExecuteMedia(dataMap["uri"].(string), dataMap); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"code": http.StatusInternalServerError,
			"err":  fmt.Sprintf("%v", err),
		})
	} else {
		go func() {
			// afterTenSecond := time.Now().Add(time.Second * 10)
			for msg := range progress {
				if dataMap := ProcessMessage(uuid, msg); dataMap["location"] != nil {
					// fileLocation = fmt.Sprintf("%v%v", dataMap["location"], dataMap["file_name"])
				}
				// if time.Now().After(afterTenSecond) {
				// 	transcoder.Stop()
				// 	// transcoder.cancelFlag = true
				// 	transcoder.SetError("hello 10 sec after.\n")
				// 	fmt.Printf("ended..\n")
				// 	break
				// }
			}
			commitProgress(uuid)
		}()
	}
	delete(dataMap, "location")
	// fmt.Printf("%v", fileLocation)

	dataMap["uuid"] = uuid
	return c.Status(http.StatusCreated).JSON(dataMap)
}

func StopTask(c *fiber.Ctx) error {

	return nil
}

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

func DownloadMedia(c *fiber.Ctx) error {
	dataMap := make(map[string]interface{})
	json.Unmarshal(c.Body(), &dataMap)
	filePath := c.Params("filename")

	fmt.Printf("parameter : %v\n", filePath)
	fmt.Printf("body : %v\n", dataMap)

	return c.Download(fmt.Sprintf("/Users/kyw/Documents/git/mine/go/snsDownloader/data/%v", filePath))
}
