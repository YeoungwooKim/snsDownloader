package main

import (
	"fmt"
	"headless/internal/pkg/colorLog"
	"headless/internal/pkg/crawler"
	"time"
)

func main() {
	colorLog.SetLogLevel(colorLog.INFO)
	// twitter
	startTime := time.Now().UnixMilli()
	url := `https://twitter.com/i/status/1273993946406907904`
	dataMapList, err := crawler.RunCrowler(url)
	if err != nil {
		colorLog.Fatal("twitter %v", err)
	}
	endTime := time.Now().UnixMilli()
	colorLog.Info("len %v \n%v\n", len(dataMapList), dataMapList)
	colorLog.Info("total Time : %vs", calculateTime(startTime, endTime))
	/*
		// youtube
		startTime = time.Now().UnixMilli()
		url = `https://www.youtube.com/watch?v=C0DPdy98e4c&ab_channel=SimonYapp`
		//`https://www.youtube.com/watch?v=XQPtaXguC3Q&ab_channel=ffeeco`
		dataMap, err := youtube.GetMediaOptions(url)
		if err != nil {
			colorLog.Fatal("%v", err)
		}
		endTime = time.Now().UnixMilli()
		colorLog.Info("total Time : %vs", calculateTime(startTime, endTime))

		videoId, audioId := getDummyIDs(dataMap)
		startTime = time.Now().UnixMilli()
		if err := youtube.GetMedia(url, videoId, audioId); err != nil {
			colorLog.Fatal("%v", err)
		}
		endTime = time.Now().UnixMilli()
		colorLog.Info("total Time : %vs", calculateTime(startTime, endTime))
	*/
	/*
		// instagram
		startTime = time.Now().UnixMilli()
		url = `https://www.instagram.com/reel/ClX6ZG8pM6A/?utm_source=ig_web_copy_link`

		if dataMap, err := crawler.RunCrowler(url); err != nil {
			colorLog.Error("%v", err)
		} else {
			fmt.Printf("%v\n", dataMap)
		}
		endTime = time.Now().UnixMilli()
		colorLog.Info("total Time : %vs", calculateTime(startTime, endTime))
	*/
}

func calculateTime(start, end int64) string {
	return fmt.Sprintf("%v", (float64)(end-start)/1000)
}

func getDummyIDs(dataMap map[string]interface{}) (string, string) {
	var videoId, audioId string
	for key, value := range dataMap["video"].(map[string]interface{}) {
		videoId = key
		fmt.Printf("video info %v\n", value)
		break
	}
	for key, value := range dataMap["audio"].(map[string]interface{}) {
		audioId = key
		fmt.Printf("audio info %v\n", value)
		break
	}
	return videoId, audioId
}
