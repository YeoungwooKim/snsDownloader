package cron

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"snsDownload/internal/pkg/media"
	"time"
)

func InitCron(everySecond int64, targetUri string) {
	go func() {
		targetTime := time.Now().Add(time.Second * time.Duration(everySecond))
		for {
			if time.Now().After(targetTime) {
				findMediaWorkLoad(targetUri)
				targetTime = time.Now().Add(time.Second * time.Duration(everySecond))
			}
			time.Sleep(time.Second * time.Duration(everySecond-5))
		}
	}()

	// go func() {
	// 	targetTime := time.Now().Add(time.Second * time.Duration(everySecond))
	// 	for {
	// 		if time.Now().After(targetTime) {
	// 			// if dataMapList, err := media.GetNotCompleteHistory(); err == nil && len(dataMapList) < 2 {
	// 			fmt.Println("hello - right before [consume message]")
	// 			kafka.ConsumeMessage()
	// 			// }
	// 			targetTime = time.Now().Add(time.Second * time.Duration(everySecond))
	// 		}
	// 		time.Sleep(time.Second * time.Duration(everySecond-1))
	// 	}
	// }()

}

func findMediaWorkLoad(targetUri string) {
	var dataMapList []map[string]interface{}
	var err error
	if dataMapList, err = media.GetNotCompleteHistory(); err != nil {
		panic(err)
	}

	if len(dataMapList) < 2 {
		fmt.Printf("we have little jobs.. (%v)\n", len(dataMapList))
		return
	}

	// fmt.Printf("job left %v\n", len(dataMapList))
	// fmt.Printf("will posted to %v\n", targetUri)
	// dataMap := map[string]interface{}{
	// 	"jobCount":    len(dataMapList),
	// 	"postedFrom":  "downloader",
	// 	"dataMapList": dataMapList,
	// }
	// // 요거하고난 뒤에는 큐에 넣었으니 플래그 처리도 해야할듯?
	// if err := enqeueWorkload(targetUri, dataMap); err != nil {
	// 	return
	// }

}

func enqeueWorkload(targetUri string, dataMap map[string]interface{}) error {
	jsonData, _ := json.Marshal(dataMap)
	response, err := http.Post(targetUri, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	fmt.Printf("response... %v\n", string(responseBody))
	return nil
}
