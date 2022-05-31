package completer

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"

	"github.com/c-bata/go-prompt"
	"github.com/githubzjm/tuo/internal/client/cache"
	"github.com/githubzjm/tuo/internal/client/completer/suggestions"
	jsonUtil "github.com/githubzjm/tuo/internal/pkg/utils/json"
)

const (
	regexGroupName = `command`
)

var commandExpression = regexp.MustCompile(`(?P<` + regexGroupName + `>` +
	`users info|users login|users` +
	`|clusters` +
	`|cache get --key|cache get|cache` +
	`)\s{1}`)

func getRegexGroups(text string) map[string]string {
	if !commandExpression.Match([]byte(text)) {
		return nil
	}

	match := commandExpression.FindStringSubmatch(text)
	result := make(map[string]string)
	for i, name := range commandExpression.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}
	return result
}

func userCompleter() []prompt.Suggest {
	bytes := cache.Get(cache.KeyUser)
	var user cache.User
	jsonUtil.Loads(bytes, &user)

	return []prompt.Suggest{
		{Text: strconv.Itoa(int(user.UserID)), Description: fmt.Sprintf("UserName: %s", user.Username)},
	}
}

func cacheKeyCompleter() []prompt.Suggest {
	var suggestions []prompt.Suggest
	// range over map is disorder, sort the map's keys to suggest in order
	keys := make([]string, 0)
	for k := range cache.KeyDescription {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		suggestions = append(suggestions, prompt.Suggest{
			Text:        key,
			Description: cache.KeyDescription[key],
		})
	}
	return suggestions
}

func Completer(d prompt.Document) []prompt.Suggest {
	word := d.GetWordBeforeCursor() // only the word without whitespace in the middle

	// all the term that needs suggestion should register in the regexp expression first
	group := getRegexGroups(d.Text)
	if group != nil {
		command := group[regexGroupName] // the whole matched string

		// dynamic suggestion should handle here
		if command == "users info" {
			return userCompleter()
		}
		if command == "cache get --key" {
			return cacheKeyCompleter()
		}

		// static sub suggestion
		if val, ok := suggestions.IsSubCommand(command); ok {
			return prompt.FilterHasPrefix(val, word, true)
		}
	}

	// static suggestion
	return prompt.FilterHasPrefix(suggestions.Suggestions, word, true)
}
