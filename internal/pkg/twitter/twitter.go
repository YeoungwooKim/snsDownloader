package twitter

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

var options = []chromedp.ExecAllocatorOption{
	chromedp.ExecPath(`/Applications/Google Chrome.app/Contents/MacOS/Google Chrome`),
	chromedp.UserAgent(`Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36`),
	chromedp.NoFirstRun,
	chromedp.NoDefaultBrowserCheck,
	chromedp.Headless,
	chromedp.DisableGPU,
}

func defineRequestNameByPlatforms(url string) string {
	if strings.Contains(url, "twitter.com") {
		return "TweetDetail"
	}
	return "XXXXXXX"
}

func RunCrowler(url string) ([]map[string]interface{}, error) {
	// url := `https://twitter.com/i/status/1273993946406907904`
	channel := make(chan bool)
	var responseBody string
	var err error

	ctx, cancel := chromedp.NewExecAllocator(context.Background(), options...)
	defer cancel()
	ctx, cancel = chromedp.NewContext(
		ctx,
	// chromedp.WithDebugf(log.Printf),
	)
	defer cancel()
	listenNetworkEvent(ctx, channel, &responseBody, defineRequestNameByPlatforms(url))

	if err = chromedp.Run(ctx, getPostMetaData(url, channel)); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("page retrieve complete!\n")

	if !strings.Contains(responseBody, "variants") {
		return nil, fmt.Errorf("there's no variants - twitter err")
	}

	if responseBody, err = parseURI(responseBody); err != nil {
		return nil, fmt.Errorf("parsing error - regex")
	}
	responseBody = fmt.Sprintf(`{%v}`, responseBody)
	var postInfoMap map[string]interface{}
	if err = json.Unmarshal([]byte(responseBody), &postInfoMap); err != nil {
		return nil, fmt.Errorf("unmarshal error %v", err)
	}
	var dataMapList []map[string]interface{}
	for _, dataMap := range postInfoMap["variants"].([]interface{}) {
		if dataMap.(map[string]interface{})["bitrate"] == nil {
			continue
		}
		dataMapList = append(dataMapList, dataMap.(map[string]interface{}))
	}

	sort.Slice(dataMapList, func(i, j int) bool {
		return dataMapList[i]["bitrate"].(float64) < dataMapList[j]["bitrate"].(float64)
	})

	return dataMapList, nil
}

func parseURI(responseBody string) (string, error) {
	responseBody = strings.ReplaceAll(responseBody, " ", "")
	//"variants":[{"bitrate":632000,"content_type":"video/mp4","url":"https://video.twimg.com/ext_tw_video/1593749370519707648/pu/vid/320x690/87DtVdvn-vAe2oVH.mp4?tag=12"},
	//{"bitrate":2176000,"content_type":"video/mp4","url":"https://video.twimg.com/ext_tw_video/1593749370519707648/pu/vid/592x1280/zG4CO3aA57JY0EBT.mp4?tag=12"},
	//{"content_type":"application/x-mpegURL","url":"https://video.twimg.com/ext_tw_video/1593749370519707648/pu/pl/740BJCkv_UysgKJZ.m3u8?tag=12&container=fmp4"},
	//{"bitrate":950000,"content_type":"video/mp4","url":"https://video.twimg.com/ext_tw_video/1593749370519707648/pu/vid/480x1036/EHELd7KDjGm_kjMg.mp4?tag=12"}]}
	variants := `\{("bitrate":[0-9]*\,)?"content_type":"(video\/mp4|application\/x\-mpegURL)"\,"url":"https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)"\}\,?`
	regex := regexp.MustCompile(fmt.Sprintf(`"variants":\[(%v)*\]`, variants))

	if matched := regex.MatchString(responseBody); matched {
		result := regex.FindStringSubmatch(responseBody)
		// fmt.Printf("%v %v \n", len(result), result[0])
		fmt.Println("parse complete")
		return result[0], nil
	}
	return "", fmt.Errorf("can't find from ResponseBody...")
}

func getPostMetaData(url string, channel chan bool) chromedp.Tasks {
	return chromedp.Tasks{
		network.Enable(),
		chromedp.Navigate(url),
		chromedp.ActionFunc(func(ctx context.Context) error {
			maxWaitTime := 10
			for {
				if <-channel || maxWaitTime > 10 {
					// fmt.Printf("\ntotal %v sec waited... break\n", 10-maxWaitTime)
					break
				}
				maxWaitTime -= 1
				time.Sleep(time.Second * 1)
			}
			return nil
		}),
	}
}

func listenNetworkEvent(ctx context.Context, channel chan bool, responseBody *string, requestName string) {
	var reqId string
	chromedp.ListenTarget(ctx, func(v interface{}) {
		switch ev := v.(type) {
		case *network.EventRequestWillBeSent:
			if ev.Type == "XHR" {
				if strings.Contains(ev.Request.URL, requestName) {
					reqId = ev.RequestID.String()
					// fmt.Printf("\trequestID>%#v\n", reqId)
					// fmt.Printf("\trequest>%#v\n", ev.Request.URL)
				}
			}
		case *network.EventResponseReceived:
			if ev.Type == "XHR" {
				go func() {
					c := chromedp.FromContext(ctx)
					respBody := network.GetResponseBody(ev.RequestID)
					body, _ := respBody.Do(cdp.WithExecutor(ctx, c.Target))
					if ev.RequestID.String() == reqId {
						// fmt.Printf("\tresponse>%s\n", body)
						*responseBody = fmt.Sprintf("%s", body)
						channel <- true
					}
				}()
			}
		}
	})
}
