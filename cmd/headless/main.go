package main

import (
	"fmt"
	"headless/internal/pkg/colorLog"
	"headless/internal/pkg/youtube"
	"time"
)

func main() {
	colorLog.SetLogLevel(colorLog.INFO)
	// startTime := time.Now().UnixMilli()
	// url := `https://twitter.com/i/status/1273993946406907904`
	// dataMapList, err := twitter.RunCrowler(url)
	// if err != nil {
	// 	colorLog.Fatal("twitter %v", err)
	// }
	// endTime := time.Now().UnixMilli()
	// colorLog.Info("len %v \n%v\n", len(dataMapList), dataMapList)
	// colorLog.Info("total Time : %vs", calculateTime(startTime, endTime))

	startTime := time.Now().UnixMilli()
	url := `https://www.youtube.com/watch?v=XQPtaXguC3Q&ab_channel=ffeeco`
	dataMap, err := youtube.GetMediaOptions(url)
	if err != nil {
		colorLog.Fatal("%v", err)
	}
	endTime := time.Now().UnixMilli()
	colorLog.Info("total Time : %vs", calculateTime(startTime, endTime))

	videoId, audioId := getDummyIDs(dataMap)
	startTime = time.Now().UnixMilli()
	if err := youtube.GetMedia(url, videoId, audioId); err != nil {
		colorLog.Fatal("%v", err)
	}
	endTime = time.Now().UnixMilli()
	colorLog.Info("total Time : %vs", calculateTime(startTime, endTime))
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

func calculateTime(start, end int64) string {
	return fmt.Sprintf("%v", (float64)(end-start)/1000)
}
