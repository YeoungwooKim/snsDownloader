package metadata

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"snsDownload/internal/pkg/log"
	"strings"
)

func executeMediaOptions(url string) (map[string]interface{}, error) {
	var outBuffer, errBuffer bytes.Buffer
	cmdLine := fmt.Sprintf(`yt-dlp --dump-json "%v" | jq '.formats'`, url)
	cmd := exec.Command("sh", "-c", cmdLine)
	defer func() {
		cmd = nil
	}()
	// bash command는 stderr에 모든 것을 작성한다.
	cmd.Stdout = &outBuffer
	cmd.Stderr = &errBuffer
	if err := cmd.Run(); err != nil {
		log.Info("command err %v", err)
		return nil, fmt.Errorf("cmd.Run error %v", err)
	}
	dataMapList := []map[string]interface{}{}
	if unmarshalErr := json.Unmarshal(outBuffer.Bytes(), &dataMapList); unmarshalErr != nil {
		panic(fmt.Sprintf("unmarshal err :%v", unmarshalErr))
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
