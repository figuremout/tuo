package common

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

const (
	DefaultCliPrefix = "@CLI> "
)

var (
	UserPrefix string
	Prefix     = DefaultCliPrefix
)

func LivePrefix() (prefix string, useLivePrefix bool) {
	return Prefix, true
}
func ChangePrefix(username string) {
	Prefix = username + DefaultCliPrefix
}

func GracefulExit(code int) {
	fmt.Println("Bye")
	beforeExit()
	os.Exit(code)
}

// Terminal does not display command after go-prompt exit: https://github.com/c-bata/go-prompt/issues/228
func beforeExit() {
	rawModeOff := exec.Command("/bin/stty", "-raw", "echo")
	rawModeOff.Stdin = os.Stdin
	_ = rawModeOff.Run()
	rawModeOff.Wait()
}

func ErrorColor(info string) string {
	return "\033[0;31m" + info + "\033[0m" // dark red text color
}
func NormalColor(info string) string {
	return "\033[0;36m" + info + "\033[0m" // dark cyan text color
}

func Progress(progressChan chan int, isShowPercent bool, loadingMsg, successMsg string) {
	var progress, charIndex int
	chars := []string{"-", "\\", "|", "/"}
	for {
		select {
		case i := <-progressChan: // a new progress is passed in
			if i == 100 {
				fmt.Printf(NormalColor("\rDONE ")+"%s\n", successMsg)
				return
			} else if i > 100 {
				fmt.Printf(ErrorColor("\rError ")+"%s\n", loadingMsg)
				return
			}
			progress = i
			charIndex++
			if isShowPercent {
				fmt.Printf(NormalColor("\r%s ")+"%d%% %s", chars[charIndex%len(chars)], progress, loadingMsg)
			} else {
				fmt.Printf(NormalColor("\r%s ")+"%s", chars[charIndex%len(chars)], loadingMsg)
			}

		default: // no new progress
			charIndex++
			if isShowPercent {
				fmt.Printf(NormalColor("\r%s ")+"%d%% %s", chars[charIndex%len(chars)], progress, loadingMsg)
			} else {
				fmt.Printf(NormalColor("\r%s ")+"%s", chars[charIndex%len(chars)], loadingMsg)
			}

		}
		time.Sleep(time.Second / 10)
	}
}
