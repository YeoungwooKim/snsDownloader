package media

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

var mergerExist bool

//	message log output parsing feature
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
		// [info] abc: Downloading 1 format(s): 244+140
		// [info] 123: Downloading 1 format(s): http-2176
		expression = `\[info\]\ *[A-Za-z0-9]+\:\ *Downloading\ *[0-9]+\ *format\(s\)\:\ *([0-9]+\+[0-9]+|[A-Za-z0-9]+\-*[A-Za-z0-9]*)`
		OtherRegex := regexp.MustCompile(expression)
		if matched := OtherRegex.MatchString(message); matched {
			result := OtherRegex.FindStringSubmatch(message)
			for i, e := range result {
				fmt.Printf("\t\t result :[%v] %#v\n", i, e)
			}
			if strings.Contains(result[1], "+") {
				mergerExist = true
			}
		}

		if mergerExist == false {
			// expression = `\[ffmpeg\]*\ *[A-Za-z]*\ *[A-Za-z]*\ *[A-Za-z]*\ *((?:[^/]*/)*)(.*)`
			expression = `\[download\]*\ *[A-Za-z]*\:\ *((?:[^/]*/)*)(.*)`
			OtherRegex = regexp.MustCompile(expression)
			if matched := OtherRegex.MatchString(message); matched {
				message = strings.ReplaceAll(message, `"`, ``)
				result := OtherRegex.FindStringSubmatch(message)
				dataMap["location"] = result[1]
				dataMap["file_name"] = result[2]
				fmt.Printf("\t[NO-MERGER]%#v\n", dataMap)
				/*
					filenames := strings.Split(filename, ".")
					fmt.Printf("%v\n", filenames)
					filename = fmt.Sprintf("%v.%v", filenames[0], filenames[2])
					fmt.Printf("%v\n", filename)
				*/
			}
		} else {
			//[Merger] Merging formats into "/Users/kyw/Documents/git/mine/go/snsDownloader/data/1674379478974-244+140.mkv"
			expression = `\[Merger\]\ *Merging\ *formats\ *into\ *((?:[^/]*/)*)(.*)`
			OtherRegex = regexp.MustCompile(expression)
			if matched := OtherRegex.MatchString(message); matched {
				message = strings.ReplaceAll(message, `"`, ``)
				result := OtherRegex.FindStringSubmatch(message)
				dataMap["file_name"] = getFileName(result[2])
				fmt.Printf("\t[MergerExist]%#v\n", dataMap)
			}
		}

	}

	if len(dataMap) >= 1 {
		saveContent(uuid, dataMap)
		saveHistory(uuid, dataMap)
	}
	return dataMap
}

func getFileName(previousName string) string {
	filenames := strings.Split(previousName, ".")
	if len(filenames) >= 3 {
		return fmt.Sprintf("%v.%v", filenames[0], filenames[2])
	}
	return previousName
}

type YtDlp struct {
	cmd          *exec.Cmd
	err          error
	completeFlag bool
	cancelFlag   bool
}

func New() *YtDlp {
	return &YtDlp{
		cmd:          nil,
		err:          nil,
		completeFlag: false,
		cancelFlag:   false,
	}
}

func (y *YtDlp) Stop() {
	fmt.Printf("called stop...err %v\n", y.GetError())
	y.cmd.Process.Kill()
}

func (y *YtDlp) GetError() error {
	return y.err
}

func (y *YtDlp) SetError(errStr string) {
	y.err = fmt.Errorf("%v", errStr)
}

// execute command youtube-dl
func (y *YtDlp) ExecuteMedia(url string, dataMap map[string]interface{}) (<-chan string, error) {
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

	cmdLine := fmt.Sprintf(`yt-dlp -f "%v" -o "%v" "%v"`, quality, fileName, url)
	fmt.Printf("%v\n\n", cmdLine)
	y.cmd = exec.Command("sh", "-c", cmdLine)
	fmt.Printf("y.cmd %v\n", y.cmd)
	stdOut, err := y.cmd.StdoutPipe()
	if err != nil {
		panic(fmt.Sprintf("cmd.StdoutPipe err %v", err))
		return nil, err
	}

	// Start process (비동기)
	if err := y.cmd.Start(); err != nil {
		panic(fmt.Sprintf("fail. executing start. error=%v", err))
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
			if y.cancelFlag || y.GetError() != nil {
				fmt.Printf("i am in go routine.\n")
				y.cmd.Process.Kill()
				time.Sleep(time.Second * 2)
				close(msg)
				break
			}
			if <-progressFlag {
				break
			}
			time.Sleep(time.Second)
			fmt.Printf("waited 1sec in cmd.wait")
		}
		if waitError = y.cmd.Wait(); waitError != nil {
			fmt.Printf("wait start.. err %v \n", waitError)
		}
	}()

	return msg, waitError
}

// execute command youtube-dl
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

	cmdLine := fmt.Sprintf(`yt-dlp -f "%v" -o "%v" "%v"`, quality, fileName, url)
	fmt.Printf("%v\n\n", cmdLine)
	cmd := exec.Command("sh", "-c", cmdLine)
	stdOut, err := cmd.StdoutPipe()
	if err != nil {
		panic(fmt.Sprintf("cmd.StdoutPipe err %v", err))
		return nil, err
	}

	// Start process (비동기)
	if err := cmd.Start(); err != nil {
		panic(fmt.Sprintf("fail. executing start. error=%v", err))
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
			fmt.Printf("waited 1sec in cmd.wait")
		}
		if waitError = cmd.Wait(); waitError != nil {
			fmt.Printf("wait start.. err %v \n", waitError)
		}
	}()

	return msg, waitError
}

// sniff youtube-dl stdOut
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
		fmt.Printf("%v\n", scanner.Text())

		// Call Wait after reaching EOF.
		if err := scanner.Err(); err != nil {
			panic(fmt.Sprintf("scanner err : %v", err))
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
