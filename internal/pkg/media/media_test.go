package media_test

import (
	"fmt"
	"regexp"
	"snsDownload/internal/pkg/media"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2/utils"
)

func TestYoutubeSingle(t *testing.T) {
	// uri := `https://twitter.com/NASA/status/1606686673915584512`
	// uri := `https://twitter.com/NASA/status/1606686673915584512"; ls -ashl; echo "hello`
	//`https://www.youtube.com/watch?v=C0DPdy98e4c&ab_channel=SimonYapp`
	var wg sync.WaitGroup

	startTime := time.Now().UnixMilli()
	optionMap := map[string]interface{}{
		"platform": "Twitter",
		"videoId":  "http-2176",
		"audioId":  nil,
		"uri":      `https://www.youtube.com/watch?v=qORaYudQ7Zc&ab_channel=STUDIOCHOOM%5B%EC%8A%A4%ED%8A%9C%EB%94%94%EC%98%A4%EC%B6%A4%5D`,
	}
	shellInjections := []string{"&", "&&", "|", "||", ";", " ", "0x0a", "\n", "$", "`", "'", `"`}
	for key, value := range optionMap {
		switch value.(type) {
		case string:
			fmt.Printf("will be checked %v\n", value)
		default:
			continue
		}
		for _, injection := range shellInjections {
			if strings.Contains(value.(string), injection) {
				optionMap[key] = strings.ReplaceAll(optionMap[key].(string), injection, "")
				fmt.Printf("\t[changed] %v\n", optionMap[key])
			}
		}
	}

	wg.Add(1)
	go func(uri string) {
		defer wg.Done()
		startTime = time.Now().UnixMilli()
		uuid := utils.UUIDv4()
		fmt.Printf("uuid : %v\n", uuid)

		if msg, err := media.ExecuteMedia(uri, optionMap); err != nil {
			fmt.Printf("%v", err)
			return
		} else {
			fmt.Println("================================================================")
			// media := "video"
			for message := range msg {
				media.ProcessMessage(uuid, message)
			}
			fmt.Println("================================================================")
		}
	}(optionMap["uri"].(string))
	wg.Wait()
	endTime := time.Now().UnixMilli()
	fmt.Printf("total Time : %vs\n", calculateTime(startTime, endTime))
}

func TestRegex(t *testing.T) {
	msg := `[youtube] C0DPdy98e4c: Downloading webpage`

	expression := `(\[[A-Za-z]+\])+\ *[A-Za-z]+\:\ *((?:[^/]*/)*)(.*)`
	OtherRegex := regexp.MustCompile(expression)
	if matched := OtherRegex.MatchString(msg); matched {
		result := OtherRegex.FindStringSubmatch(msg)
		fmt.Printf("%v\n", result)
	} else {
		fmt.Println("not match")
	}

	fmt.Printf("\n\n")
	msg = `[download] Destination: /Users/kyw/Documents/git/mine/go/headless/internal/pkg/youtube/data/1670762468979-TEST VIDEO.f244.f244.webm`
	// expression = `(\[[A-Za-z]+\])+\ *[A-Za-z]+\:\ *((?:[^/]*/)*)(.*)`
	// OtherRegex = regexp.MustCompile(expression)
	if matched := OtherRegex.MatchString(msg); matched {
		result := OtherRegex.FindStringSubmatch(msg)
		for _, elem := range result {
			fmt.Printf("%#v\n", elem)
		}
	} else {
		fmt.Printf("not match \n")
	}
	fmt.Printf("\n\n")

	expression = `\[ffmpeg\]*\ *[A-Za-z]*\ *[A-Za-z]*\ *[A-Za-z]*\ *((?:[^/]*/)*)(.*)`
	OtherRegex = regexp.MustCompile(expression)
	msg = `[ffmpeg] Merging formats into "/Users/kyw/Documents/git/mine/go/headless/internal/pkg/youtube/data/1670762468979-TEST VIDEO.f244+140.mkv"`
	msg = strings.ReplaceAll(msg, `"`, ``)
	if matched := OtherRegex.MatchString(msg); matched {
		result := OtherRegex.FindStringSubmatch(msg)
		for _, elem := range result {
			fmt.Printf("%#v\n", elem)
		}
	} else {
		fmt.Printf("not match \n")
	}

	fmt.Printf("\n\n")
	msg = `Deleting original file /Users/kyw/Documents/git/mine/go/headless/internal/pkg/youtube/data/1670762468979-TEST VIDEO.f244.f244.webm (pass -k to keep)`
	if matched := OtherRegex.MatchString(msg); matched {
		result := OtherRegex.FindStringSubmatch(msg)
		for _, elem := range result {
			fmt.Printf("%#v\n", elem)
		}
	} else {
		fmt.Printf("not match \n")
	}
}

// func TestYoutubeMultiple(t *testing.T) {
// 	uri := `https://www.youtube.com/watch?v=C0DPdy98e4c&ab_channel=SimonYapp`
// 	var wg sync.WaitGroup
// 	numberOfGoRoutine := 6
// 	startTime := time.Now().UnixMilli()
// 	wg.Add(numberOfGoRoutine)
// 	for i := 0; i < numberOfGoRoutine; i++ {
// 		go func(uri string) {
// 			defer func() {
// 				// fmt.Printf("%vth routine end..", i)
// 				wg.Done()
// 			}()
// 			getRndOptions := func(uri string) map[string]interface{} {
// 				if dataMap, err := youtube.GetMediaOptions(uri); err != nil {
// 					fmt.Printf("failed to get media Option %v\n", err)
// 					t.Error()
// 					return nil
// 				} else {
// 					v, a := getDummyIDs(dataMap)
// 					return map[string]interface{}{
// 						"videoId": v,
// 						"audioId": a,
// 					}
// 				}
// 			}
// 			// videoId, audioId := getDummyIDs(getOptions(uri))
// 			startTime = time.Now().UnixMilli()
// 			if filePath, err := youtube.GetMedia(uri, getRndOptions(uri)); err != nil {
// 				fmt.Printf("%v", err)
// 				return
// 			} else {
// 				fmt.Printf("result is located at %v\n", filePath)
// 			}
// 		}(uri)
// 	}
// 	wg.Wait()
// 	endTime := time.Now().UnixMilli()
// 	fmt.Printf("total Time : %vs\n", calculateTime(startTime, endTime))
// }

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
