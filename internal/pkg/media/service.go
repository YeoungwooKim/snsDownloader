package media

import (
	"bufio"
	"bytes"
	"fmt"
	"headless/internal/pkg/colorLog"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

func ProcessMessage(uuid, message string) map[string]interface{} {
	if len(strings.ReplaceAll(message, ` `, ``)) == 0 {
		return nil
	}

	expression := `\[download\]\ *([0-9]+\.[0-9]+%)\ *of\ *([0-9]+\.[0-9]+[A-Za-z]+)\ *at\ *([0-9]*\.[0-9]+[A-Za-z]+\/s)\ *ETA\ *([0-9]+\:[0-9]+)`
	ProgressRegex := regexp.MustCompile(expression)

	dataMap := make(map[string]interface{})
	if matched := ProgressRegex.MatchString(message); matched {
		// fmt.Printf("[progress catched] ")
		result := ProgressRegex.FindStringSubmatch(message)
		dataMap["download_percent"] = result[1]
		dataMap["file_size"] = result[2]
		dataMap["download_speed"] = result[3]
		dataMap["time_left"] = result[4]
	} else {
		// fmt.Printf("not matched -download ") //%v\n", message)
		// expression = `\[ffmpeg\]*\ *[A-Za-z]*\ *[A-Za-z]*\ *[A-Za-z]*\ *((?:[^/]*/)*)(.*)`
		expression = `\[download\]*\ *[A-Za-z]*\:\ *((?:[^/]*/)*)(.*)`
		OtherRegex := regexp.MustCompile(expression)
		if matched := OtherRegex.MatchString(message); matched {
			// fmt.Printf("[file_location catched] ")
			message = strings.ReplaceAll(message, `"`, ``)
			result := OtherRegex.FindStringSubmatch(message)
			dataMap["location"] = result[1]
			dataMap["file_name"] = result[2]
		}
	}

	if len(dataMap) >= 1 {
		saveHistory(uuid, dataMap)
	}
	return dataMap
}

func ExecuteMedia(url string, dataMap map[string]interface{}) (<-chan string, error) {
	var quality string
	fmt.Printf("%v %v\n", url, dataMap)
	pwd, _ := os.Getwd()
	timaStamp := time.Now().UnixMilli()
	fileName := fmt.Sprintf("%v/data/%v-", pwd, timaStamp) + `%(format_id)s.%(ext)s`

	if dataMap["videoId"] == nil || dataMap["audioId"] == nil {
		quality = "best"
	} else {
		quality = fmt.Sprintf("%v+%v", dataMap["videoId"], dataMap["audioId"])
	}
	//"worstvideo[ext=mp4]+worstaudio[ext=m4a]"

	cmdLine := fmt.Sprintf(`youtube-dl -f "%v" -o "%v" "%v"`, quality, fileName, url)
	fmt.Printf("%v\n\n", cmdLine)
	cmd := exec.Command("sh", "-c", cmdLine)
	stdOut, err := cmd.StdoutPipe()
	if err != nil {
		colorLog.Fatal("cmd.StdoutPipe err %v", err)
		return nil, err
	}

	// Start process (wait 하지 않고 수행 - console 확인차)
	if err := cmd.Start(); err != nil {
		colorLog.Fatal("fail. executing start. error=%v", err)
		return nil, err
	}

	progressFlag := make(chan bool)
	msg := make(chan string)

	go func() {
		defer func() {
			fmt.Println("[msg channel will be close....]")
			close(msg)
		}()
		cmdDownloadProgress(stdOut, msg, progressFlag)
	}()

	var waitError error
	go func() {
		for {
			if <-progressFlag {
				break
			}
			time.Sleep(time.Second)
			colorLog.Info("waited 1sec in cmd.wait")
		}
		if waitError = cmd.Wait(); waitError != nil {
			colorLog.Info("wait start.. err %v \n", waitError)
		}
	}()

	return msg, waitError
}

func cmdDownloadProgress(stream io.ReadCloser, msg chan<- string, progressFlag chan bool) {
	defer func() {
		stream.Close()
		progressFlag <- true
	}()
	scanner := bufio.NewScanner(stream)
	scanner.Split(customSplit)

	buf := make([]byte, 2)
	scanner.Buffer(buf, bufio.MaxScanTokenSize)
	for scanner.Scan() {
		msg <- scanner.Text()

		// Call Wait after reaching EOF.
		if err := scanner.Err(); err != nil {
			colorLog.Fatal("scanner err : %v", err)
		}
	}
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
