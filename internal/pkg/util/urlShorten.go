package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"snsDownload/internal/pkg/server/config"
	"strings"
)

func IsAlreadyShorten(targetUri string) bool {
	RedirectAttemptedError := errors.New("redirect attempted")

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return RedirectAttemptedError
		},
	}
	req, _ := http.NewRequest(http.MethodGet, targetUri, nil)
	_, err := client.Do(req)
	return err != nil
}

func ShortenUrl(targetUrl string) (string, error) {
	// targetUrl := `https://www.youtube.com/watch?v=qORaYudQ7Zc&ab_channel=STUDIOCHOOM%5B%EC%8A%A4%ED%8A%9C%EB%94%94%EC%98%A4%EC%B6%A4%5D`
	client := &http.Client{}
	var data = strings.NewReader(fmt.Sprintf("url=%v", targetUrl))
	req, err := http.NewRequest(http.MethodPost, "https://openapi.naver.com/v1/util/shorturl", data)
	if err != nil {
		return "", fmt.Errorf("%v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	req.Header.Set("X-Naver-Client-Id", config.Config.NaverId)
	req.Header.Set("X-Naver-Client-Secret", config.Config.NaverSecret)

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("%v", err)
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("%v", err)
	}
	dataMap := make(map[string]interface{})
	json.Unmarshal(bodyText, &dataMap)

	return dataMap["result"].(map[string]interface{})["url"].(string), nil
}
