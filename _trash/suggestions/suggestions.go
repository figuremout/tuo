package suggestions

import "github.com/c-bata/go-prompt"

var (
	// these are the first command term, should be less if possible
	Suggestions = []prompt.Suggest{
		{Text: "users", Description: "Manage users"},
		{Text: "clusters", Description: "Manage clusters"},
		{Text: "nodes", Description: "Manage nodes"},

		// consisder later
		{Text: "logs", Description: "Fetch the logs of a container"},

		{Text: "cache", Description: "Manage local cache"},
		{Text: "version", Description: "Show the system version information"},
		{Text: "exit", Description: "Exit command prompt"},
	}

	SubSuggestions = map[string][]prompt.Suggest{
		"users": {
			prompt.Suggest{Text: "register", Description: "Register user"},
			prompt.Suggest{Text: "login", Description: "Log in"},
			prompt.Suggest{Text: "logout", Description: "Log out"},
			prompt.Suggest{Text: "ls", Description: "List all users"},
			prompt.Suggest{Text: "info", Description: "Display user information"},
			prompt.Suggest{Text: "update", Description: "Update user information"},
			prompt.Suggest{Text: "rm", Description: "Remove user"},
		},

		"clusters": {
			prompt.Suggest{Text: "create", Description: "Create cluster"},
			prompt.Suggest{Text: "ls", Description: "List all clusters"},
			prompt.Suggest{Text: "info", Description: "Display cluster information"},
			prompt.Suggest{Text: "update", Description: "Update cluster information"},
			prompt.Suggest{Text: "delete", Description: "Delete cluster"},
			prompt.Suggest{Text: "plugin", Description: "Manage plugins"},
		},

		"nodes": {
			prompt.Suggest{Text: "create", Description: "Create node"},
			prompt.Suggest{Text: "ls", Description: "List all nodes"},
			prompt.Suggest{Text: "search", Description: "Search for nodes"}, // match nodes which name contains specified string
			prompt.Suggest{Text: "info", Description: "Display node information"},
			prompt.Suggest{Text: "update", Description: "Update node information"},
			prompt.Suggest{Text: "delete", Description: "Delete node"},
			prompt.Suggest{Text: "logs", Description: "Fetch the logs of a container"},
			prompt.Suggest{Text: "exec", Description: "Run a command in a running node"},
			prompt.Suggest{Text: "stats", Description: "Display a live stream of node resource usage statistics"}, // TODO on terminal like docker stats or web like npm start
		},
		"nodes search": {
			prompt.Suggest{Text: "--automated", Description: ""},
			prompt.Suggest{Text: "--filter", Description: "Filter output based on conditions provided"},
			prompt.Suggest{Text: "--format", Description: "Pretty-print search using a Go template"},
			prompt.Suggest{Text: "--limit", Description: "Max number of search results"},
		},
		"nodes logs": {
			prompt.Suggest{Text: "--details", Description: "Show extra details provided to logs"},
			prompt.Suggest{Text: "--follow", Description: "Follow log output"},
			prompt.Suggest{Text: "--since", Description: "Show logs since timestamp (e.g. 2013-01-02T13:23:37) or relative (e.g. 42m for 42 minutes)"},
			prompt.Suggest{Text: "--tail", Description: "Number of lines to show from the end of the logs"},
			prompt.Suggest{Text: "--timestamps", Description: "Show timestamps"},
			prompt.Suggest{Text: "--until", Description: ""},
		},

		"tasks": {},

		"cache": {
			prompt.Suggest{Text: "info", Description: "Display local cache information"},
			prompt.Suggest{Text: "reset", Description: "Empty all cache shards"},
			prompt.Suggest{Text: "get", Description: "Get cached value"},
		},
		"cache get": {
			prompt.Suggest{Text: "--key", Description: "Key of cache value"},
		},
	}
)

func IsCommand(kw string) bool {
	for _, cmd := range Suggestions {
		if cmd.Text == kw {
			return true
		}
	}
	return false
}

func IsSubCommand(kw string) ([]prompt.Suggest, bool) {
	val, ok := SubSuggestions[kw]
	return val, ok
}
