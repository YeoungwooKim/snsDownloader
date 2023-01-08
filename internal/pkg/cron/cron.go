package cron

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"snsDownloader/internal/pkg/media"
	"time"
)

func InitCron(everySecond int64, targetUri string) {
	go func() {
		targetTime := time.Now().Add(time.Second * time.Duration(everySecond))
		for {
			if time.Now().After(targetTime) {
				findMediaWorkLoad(targetUri)
				/*
					작업량이 적을때 이제 일해야지?
				*/
				targetTime = time.Now().Add(time.Second * time.Duration(everySecond))
			}

			time.Sleep(time.Second * time.Duration(everySecond))
		}
	}()

}

func findMediaWorkLoad(targetUri string) {
	var dataMapList []map[string]interface{}
	var err error
	if dataMapList, err = media.GetNotCompleteHistory(); err != nil {
		panic(err)
	}

	if len(dataMapList) < 2 {
		fmt.Printf("no queue")
		return
	}

	fmt.Printf("job left %v\n", len(dataMapList))
	fmt.Printf("will posted to %v\n", targetUri)
	dataMap := map[string]interface{}{
		"jobCount":    len(dataMapList),
		"postedFrom":  "downloader",
		"dataMapList": dataMapList,
	}

	jsonData, _ := json.Marshal(dataMap)
	response, err := http.Post(targetUri, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}
	defer func() {
		response.Body.Close()
	}()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	fmt.Printf("response... %v\n", string(responseBody))

}
