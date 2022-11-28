package crawler

import (
	"context"
	"encoding/json"
	"fmt"
	"headless/internal/pkg/colorLog"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
)

var options = []chromedp.ExecAllocatorOption{
	chromedp.ExecPath(`/Applications/Google Chrome.app/Contents/MacOS/Google Chrome`),
	chromedp.UserAgent(`Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36`),
	chromedp.NoFirstRun,
	chromedp.NoDefaultBrowserCheck,
	chromedp.Headless,
	chromedp.DisableGPU,
}

func RunCrowler(url string) ([]map[string]interface{}, error) {
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

	if isTwitter(url) {
		listenNetworkEvent(ctx, channel, &responseBody, "TweetDetail")
	}

	var elem []string
	if err = chromedp.Run(ctx, getPostMetaData(url, channel, &elem)); err != nil {
		colorLog.Info("runError : %v", err)
		time.Sleep(time.Hour)
		return nil, err
	}
	// fmt.Printf(">>>>>>>>>>>\t elem is \n%v\n", elem)

	if isTwitter(url) {
		if dataMapList, err := validateData(responseBody); err != nil {
			return nil, err
		} else {
			return dataMapList, nil
		}
	}
	return []map[string]interface{}{
		convertMap(elem),
	}, nil
}

func convertMap(elem []string) map[string]interface{} {
	dataMap := map[string]interface{}{}
	for i := 0; i < len(elem); i += 2 {
		dataMap[elem[i]] = elem[i+1]
	}
	return dataMap
}

func validateData(responseBody string) ([]map[string]interface{}, error) {
	var err error
	if !strings.Contains(responseBody, "variants") {
		return nil, fmt.Errorf("there's no variants - twitter err")
	}

	if responseBody, err = filterJson(responseBody); err != nil {
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

func isTwitter(url string) bool {
	if strings.Contains(url, "https://twitter.com") {
		return true
	}
	return false
}

func filterJson(responseBody string) (string, error) {
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

func getPostMetaData(url string, channel chan bool, elem *[]string) chromedp.Tasks {
	if strings.Contains(url, "twitter") {
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
	var loc string

	return chromedp.Tasks{
		network.Enable(),
		chromedp.Navigate(url),
		chromedp.Location(&loc),
		chromedp.Sleep(time.Second * 1),
		chromedp.ActionFunc(func(ctx context.Context) error {
			// if strings.Contains(loc, "account") {
			chromedp.WaitReady(`#loginForm`, chromedp.ByID).Do(ctx)
			colorLog.Info("load complete....")

			chromedp.Click(`input[name="username"]`).Do(ctx)
			chromedp.SendKeys(`input[name="username"]`, INSTAGRAM_USERNAME+kb.Tab).Do(ctx)
			// chromedp.Sleep(time.Second * 3).Do(ctx)
			chromedp.SendKeys(`input[name="password"]`, INSTAGRAM_PASSWORD+kb.Enter).Do(ctx)
			// }
			return nil
		}),
		chromedp.Location(&loc),

		chromedp.ActionFunc(func(ctx context.Context) error {
			colorLog.Info("url : %v", loc)
			if strings.Contains(loc, "account") {
				colorLog.Info("its account page....")
				chromedp.Sleep(time.Second * 3).Do(ctx)
				chromedp.EvaluateAsDevTools(`document.querySelectorAll('button')[1].click()`, nil).Do(ctx)
			}
			return nil
		}),
		chromedp.Location(&loc),
		chromedp.ActionFunc(func(ctx context.Context) error {
			colorLog.Info("url : %v\n", loc)
			var projects []*cdp.Node
			chromedp.Nodes(`video[src]`, &projects).Do(ctx)
			// attributes = *&projects[0].Attributes
			*elem = *&projects[0].Attributes
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
					// fmt.Printf("\trequestID>%#v\n", reqId) // fmt.Printf("\trequest>%#v\n", ev.Request.URL)
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
