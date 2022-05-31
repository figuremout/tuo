package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/c-bata/go-prompt"
	"github.com/pkg/term/termios"
	"golang.org/x/sys/unix"
)

var fd int
var originalTermios *unix.Termios

func Excutor(input string) {
	// restore the original settings to allow ctrl-c to generate signal
	if err := termios.Tcsetattr(uintptr(fd), termios.TCSANOW, (*unix.Termios)(originalTermios)); err != nil {
		panic(err)
	}

	if input == "test" {
		ctx, cancel := context.WithCancel(context.Background())
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)

		go func() {
			select {
			case <-c:
				cancel()
			}
		}()
		go func() {
			defer cancel()
			for { // long task
			}
		}()
		select {
		case <-ctx.Done():
			return
		}
	}
}
func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func main() {
	// test
	var err error
	fd, err = syscall.Open("/dev/tty", syscall.O_RDONLY, 0)
	if err != nil {
		panic(err)
	}
	// get the original settings
	originalTermios, err = termios.Tcgetattr(uintptr(fd))
	if err != nil {
		panic(err)
	}

	p := prompt.New(Excutor, completer)

	p.Run()
}

// func Read() ([]byte, error) {
// 	buf := make([]byte, 1024)
// 	n, err := syscall.Read(fd, buf)
// 	if err != nil {
// 		return []byte{}, err
// 	}
// 	return buf[:n], nil
// }

// func main() {

// 	var err error
// 	fd, err = syscall.Open("/dev/tty", syscall.O_RDONLY, 0)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// Set NonBlocking mode because if syscall.Read block this goroutine, it cannot receive data from stopCh.
// 	// if err := syscall.SetNonblock(fd, true); err != nil {
// 	// 	panic(err)
// 	// }

// 	n, err := termios.Tcgetattr(uintptr(fd))
// 	if err != nil {
// 		panic(err)
// 	}

// 	n.Iflag &^= syscall.IGNBRK | syscall.BRKINT | syscall.PARMRK |
// 		syscall.ISTRIP | syscall.INLCR | syscall.IGNCR |
// 		syscall.ICRNL | syscall.IXON
// 	n.Lflag &^= syscall.ECHO | syscall.ICANON | syscall.IEXTEN | syscall.ECHONL //  | syscall.ISIG
// 	n.Cflag &^= syscall.CSIZE | syscall.PARENB
// 	n.Cflag |= syscall.CS8 // Set to 8-bit wide.  Typical value for displaying characters.
// 	n.Cc[syscall.VMIN] = 1
// 	n.Cc[syscall.VTIME] = 0
// 	n.Cc[syscall.VINTR] = uint8(0x03)

// 	if err := termios.Tcsetattr(uintptr(fd), termios.TCSANOW, (*unix.Termios)(n)); err != nil {
// 		panic(err)
// 	}

// 	for {
// 		bs, err := Read()
// 		if err != nil {
// 			fmt.Printf("Error: %v", err)
// 		}
// 		fmt.Printf("%v", bs)
// 	}
// }
