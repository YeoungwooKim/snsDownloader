package metadata

import (
	"encoding/json"
	"fmt"
	"net/http"
	"snsDownloader/internal/pkg/validation"

	"github.com/gofiber/fiber/v2"
)

/*
	메타데이터 받기전 단계로 검증 역할
*/
func ValidationRouter(c *fiber.Ctx) error {
	valid := *validation.Validate
	dataMap := make(map[string]interface{})
	// request body를 파싱.
	if err := json.Unmarshal(c.Body(), &dataMap); err != nil {
		fmt.Printf("err %v\n", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"code": http.StatusInternalServerError,
			"err":  fmt.Sprintf("%v", err),
		})
	}
	// request body를 검증/ 데이터 형식 & 데이터 길이 체크
	if err := valid.ValidateMap(dataMap, validation.VideoRules); len(err) > 0 {
		fmt.Printf("validation err %v\n", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"code": http.StatusInternalServerError,
			"err":  fmt.Sprintf("%v", err),
		})
	}
	fmt.Printf("%#v\n", dataMap)
	// if !strings.Contains(dataMap["uri"].(string), strings.ToLower(dataMap["platform"].(string))) {
	// 	return c.Status(http.StatusBadRequest).JSON(fiber.Map{
	// 		"code": http.StatusBadRequest,
	// 		"err":  "match platform and uri",
	// 	})
	// }
	return c.Next()
}

func GetMetaData(c *fiber.Ctx) error {
	dataMap := make(map[string]interface{})
	json.Unmarshal(c.Body(), &dataMap)
	mediaOptionMap, err := executeMediaOptions(dataMap["uri"].(string))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"code": http.StatusInternalServerError,
			"err":  fmt.Sprintf("%v", err),
		})
	}
	mediaOptionMap["platform"] = "youtube"
	return c.Status(http.StatusOK).JSON(mediaOptionMap)
}
