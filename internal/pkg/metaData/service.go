package metadata

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"headless/internal/pkg/colorLog"
	"io"
	"os/exec"
	"strings"
	"time"
)

func executeMediaOptions(url string) (map[string]interface{}, error) {
	cmdLine := fmt.Sprintf(`youtube-dl --dump-json "%v" | jq '.formats'`, url)
	cmd := exec.Command("sh", "-c", cmdLine)
	// bash command는 stderr에 모든 것을 작성한다.
	stderrIn, stdErr := cmd.StdoutPipe()
	if stdErr != nil {
		colorLog.Fatal("cmd.StdoutPipe err %v", stdErr)
		return nil, stdErr
	}
	// Start process (wait 하지 않고 수행 - console 확인차)
	if startErr := cmd.Start(); startErr != nil {
		colorLog.Fatal("fail. executing start. error=%v", startErr)
		return nil, startErr
	}

	progressFlag := make(chan bool)
	jsonString := make(chan string)
	// log check
	go func() {
		cmdProgress(stderrIn, progressFlag, jsonString)
	}()
	go func() {
		for {
			if <-progressFlag {
				break
			}
			time.Sleep(time.Second)
			colorLog.Info("waited 1sec in cmd.wait")
		}
		if waitError := cmd.Wait(); waitError != nil {
			colorLog.Info("wait start.. err %v \n", waitError)
		} else {
			// colorLog.Info("command properly ended..")
		}
	}()
	dataMapList := []map[string]interface{}{}
	if unmarshalErr := json.Unmarshal([]byte(<-jsonString), &dataMapList); unmarshalErr != nil {
		colorLog.Fatal("unmarshal err :%v", unmarshalErr)
	}

	return exportMap(dataMapList), nil
}

func exportMap(dataMapList []map[string]interface{}) map[string]interface{} {
	dataMap := make(map[string]interface{})
	dataMap["audio"] = make(map[string]interface{})
	dataMap["video"] = make(map[string]interface{})
	for _, elem := range dataMapList {
		// fmt.Printf("%v/ %v/ %v\n", elem["format_id"], elem["format"], elem["filesize"])
		key := elem["format_id"].(string)
		if strings.Contains(elem["format"].(string), "audio") {
			dataMap["audio"].(map[string]interface{})[key] = map[string]interface{}{
				"format":   elem["format"],
				"filesize": elem["filesize"],
			}
		} else {
			dataMap["video"].(map[string]interface{})[key] = map[string]interface{}{
				"format":   elem["format"],
				"filesize": elem["filesize"],
			}
		}
	}
	// fmt.Printf("mapLen %v audioLen %v videoLen %v\n", len(dataMap), len(dataMap["audio"].(map[string]interface{})), len(dataMap["video"].(map[string]interface{})))
	return dataMap
}
func cmdProgress(stream io.ReadCloser, progressFlag chan bool, jsonString chan string) {
	var output string
	defer func() {
		stream.Close()
		progressFlag <- true
		jsonString <- output
	}()
	scanner := bufio.NewScanner(stream)
	scanner.Split(customSplit)

	buf := make([]byte, 2)
	scanner.Buffer(buf, bufio.MaxScanTokenSize)
	for scanner.Scan() {
		output += scanner.Text()
		// Call Wait after reaching EOF.
		if err := scanner.Err(); err != nil {
			colorLog.Fatal("scanner err : %v", err)
		}
	}
}

func customSplit(data []byte, eof bool) (advance int, token []byte, spliterror error) {
	if eof && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		// We have a full newline-terminated line.
		return i + 1, data[0:i], nil
	}
	if i := bytes.IndexByte(data, '\r'); i >= 0 {
		// We have a cr terminated line
		return i + 1, data[0:i], nil
	}
	if eof {
		return len(data), data, nil
	}
	return 0, nil, nil
}
