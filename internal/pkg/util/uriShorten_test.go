package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"testing"
)

func TestUrlShorten(t *testing.T) {
	// orgUrl https://www.youtube.com/watch?v=qORaYudQ7Zc
	// url https://me2.do/FXrDLRDO
	// code 200
	// message ok
	targetUrl := `https://www.youtube.com/watch?v=qORaYudQ7Zc&ab_channel=STUDIOCHOOM%5B%EC%8A%A4%ED%8A%9C%EB%94%94%EC%98%A4%EC%B6%A4%5D\";ls -ashl;echo\"`
	//`https://dummyjson.com/products/1`

	client := &http.Client{}
	var data = strings.NewReader(fmt.Sprintf("url=%v", targetUrl))
	fmt.Printf("body %#v\n", data)
	req, err := http.NewRequest(http.MethodPost, "https://openapi.naver.com/v1/util/shorturl", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("X-Naver-Client-Id", "fTOvEJQMRgyI7oWl6e2n")
	req.Header.Set("X-Naver-Client-Secret", "4dhnvJivv5")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	dataMap := make(map[string]interface{})
	json.Unmarshal(bodyText, &dataMap)
	fmt.Printf("%s\n", dataMap)

	fmt.Printf("%v\n", dataMap["result"].(map[string]interface{})["url"].(string))
}
