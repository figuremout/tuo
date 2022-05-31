package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"sync"
	"syscall"

	"github.com/c-bata/go-prompt"
	"github.com/githubzjm/tuo/api/v1/users/def"
	"github.com/githubzjm/tuo/internal/client/cache"
	"github.com/githubzjm/tuo/internal/client/common"
	"github.com/pkg/term/termios"
	"golang.org/x/sys/unix"

	"github.com/githubzjm/tuo/internal/client/handlers"
	httpClient "github.com/githubzjm/tuo/internal/pkg/http/client"
	"github.com/githubzjm/tuo/pkg/command"
)

var parser *command.Parser
var fd int
var originalTermios *unix.Termios

func Executor(input string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Panic: %v\n%s", err, debug.Stack())
		}
	}()

	if input == "" {
		return
	}

	// restore the original settings to allow ctrl-c to generate signal
	// Notice: (*prompt.Prompt).Run() and prompt.Input() will call (*Prompt).setUp(), which will change termios,
	// so the cmd handler called by parser.Parse() better not use them
	if err := termios.Tcsetattr(uintptr(fd), termios.TCSANOW, (*unix.Termios)(originalTermios)); err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	var wg sync.WaitGroup // when cancel() called, the Executor func may end before handler func, need to make sure the end order

	go func() {
		<-c
		cancel() // if signal received, finish task
		fmt.Printf("\n")
	}()
	go func() {
		defer cancel()
		wg.Add(1)
		if err := parser.Parse(ctx, input); err != nil {
			fmt.Printf(common.ErrorColor("Unknown command: %s\n"), input)
		}
		wg.Done()
	}()

	<-ctx.Done() // block here until cancel() called
	wg.Wait()    // block here until handler is canceled
}

func Completer(d prompt.Document) []prompt.Suggest {
	names, err := parser.Split(d.Text)
	if err != nil {
		return []prompt.Suggest{}
	}
	target := parser.Search(names)
	var suggestions []prompt.Suggest
	if target == nil {
		suggestions = parser.GenerateSuggestions()
	} else {
		suggestions = target.GenerateSuggestions()
	}
	return prompt.FilterHasPrefix(suggestions, d.GetWordBeforeCursor(), true)
}

// serve page for data chart, and open browser automatically
func main() {
	var err error
	fd, err = syscall.Open("/dev/tty", syscall.O_RDONLY, 0)
	if err != nil {
		panic(err)
	}

	originalTermios, err = termios.Tcgetattr(uintptr(fd))
	if err != nil {
		panic(err)
	}

	// init cache
	err = cache.InitCache(def.TOKEN_VALID_DURATION_NS, -1)
	if err != nil {
		fmt.Print(err)
	}
	defer cache.Close()

	// init http client
	httpClient.InitHttpClient("http://localhost:8080")

	// init command parser
	parser = command.New()
	// register cmds
	users := parser.Group("users", "Manage users", nil)
	users.Group("register", "Register user", func(ctx context.Context, fs *flag.FlagSet) error {
		handlers.Register()
		return nil
	})
	users.Group("login", "Log in", func(ctx context.Context, fs *flag.FlagSet) error {
		handlers.Login()
		return nil
	})
	users.Group("logout", "Log out", func(ctx context.Context, fs *flag.FlagSet) error {
		handlers.Logout()
		return nil
	})
	users.Group("update", "Update user information", nil)
	users.Group("rm", "Remove user", nil)
	{
		ls := users.Group("ls", "List information about users (the current user by default)", nil)
		ls.FS.Bool("a", false, "List all users (only admin allowed)")
	}

	clusters := parser.Group("clusters", "Manage clusters", nil)
	{
		clusters.Group("ls", "List information about users", func(ctx context.Context, fs *flag.FlagSet) error {
			handlers.QueryAllClusters()
			return nil
		})
		clusters.Group("create", "Create cluster", func(ctx context.Context, fs *flag.FlagSet) error {
			handlers.CreateCluster()
			return nil
		})
		clusters.Group("update", "Update cluster information", nil)
		clusters.Group("delete", "Delete cluster", nil)
		clusters.Group("plugin", "Manage plugins", nil)
	}

	nodes := parser.Group("nodes", "Manage nodes", nil)
	{
		ls := nodes.Group("ls", "List information about nodes", func(ctx context.Context, fs *flag.FlagSet) error {
			c := fs.Lookup("c")
			handlers.QueryAllNodes(c.Value.String())
			return nil
		}) // TODO flag for search
		ls.FS.Int("c", -1, "Cluster ID to query (required)") // id is actually uint, but default value should be a invalid value

		nodes.Group("create", "Create node", func(ctx context.Context, fs *flag.FlagSet) error {
			handlers.CreateNode()
			return nil
		})
		nodes.Group("deploy", "Deploy agent to the node", func(ctx context.Context, fs *flag.FlagSet) error {
			handlers.DeployNode()
			return nil
		})

		stats := nodes.Group("stats", "Display a live stream of node resource usage statistics", nil)
		{
			cpu := stats.Group("cpu", "CPU Percent", func(ctx context.Context, fs *flag.FlagSet) error {
				port := fs.Lookup("p")
				handlers.CPUPercent(ctx, port.Value.String())
				return nil
			})
			cpu.FS.Uint("p", 80, "Serve web chart at localhost:<port> (80 by default)")
		}

		// search.AddFlag("filter", "Filter output based on conditions provided", true)
		// search.AddFlag("limit", "Max number of search results", true)

		nodes.Group("update", "Update node information", nil)
		nodes.Group("delete", "Delete node", nil)

		// nodes.Group("logs", "Fetch the logs of a container", nil)
		// logs.AddFlag("details", "Show extra details provided to logs", false)
		// logs.AddFlag("tail", "Number of lines to show from the end of the logs", true)
		// logs.AddFlag("since", "Show logs since timestamp (e.g. 2013-01-02T13:23:37) or relative (e.g. 42m for 42 minutes)", true)
		// logs.AddFlag("until", "Show logs until timestamp", true)

		nodes.Group("exec", "Run a command in a running node", nil)

	}

	cache := parser.Group("cache", "Manage local cache", nil)
	{
		cache.Group("ls", "List all cached key-value", func(ctx context.Context, fs *flag.FlagSet) error { // ls -a -l
			handlers.CacheList()
			return nil
		})
		cache.Group("info", "Display local cache information", func(ctx context.Context, fs *flag.FlagSet) error {
			handlers.CacheInfo()
			return nil
		})
		cache.Group("reset", "Empty all cache shards", func(ctx context.Context, fs *flag.FlagSet) error {
			handlers.CacheReset()
			return nil
		})
	}

	cmd := parser.Group("cmd", "Manage cmd parser", nil)
	{
		cmd.Group("tree", "Print cmd tree", func(ctx context.Context, fs *flag.FlagSet) error {
			parser.Print()
			return nil
		})
	}

	parser.Group("exit", "Exit command prompt", func(ctx context.Context, fs *flag.FlagSet) error {
		common.GracefulExit(0)
		return nil
	})

	// parser.Group("test", "test", func(ctx context.Context, fs *flag.FlagSet) error {
	// 	ws.Test()
	// 	return nil
	// })

	p := prompt.New(
		Executor, Completer,
		// prompt.OptionPrefix(prefix),
		prompt.OptionTitle("Client Prompt"),
		prompt.OptionLivePrefix(common.LivePrefix),
		prompt.OptionCompletionOnDown(),
		prompt.OptionMaxSuggestion(uint16(len(parser.Childs))), // allow all suggestions show
		prompt.OptionShowCompletionAtStart(),
		prompt.OptionPrefixTextColor(prompt.Cyan),
	)

	p.Run()
	fmt.Println("\nBye")

}
