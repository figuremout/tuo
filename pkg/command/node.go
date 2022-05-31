package command

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/google/shlex"
)

type Node struct {
	Name        string
	FullName    string
	Description string
	root        bool
	parser      *Parser
	FS          *flag.FlagSet
	Childs      []*Node
	Handler     Handler
}

type Handler func(context.Context, *flag.FlagSet) error

// handler will be called only on those flags that have been set.
func (n *Node) Group(name, description string, handler Handler) *Node {
	fullName := strings.TrimSpace(strings.Join([]string{n.FullName, name}, " "))
	fs := flag.NewFlagSet(fullName, flag.ContinueOnError)
	fs.SetOutput(n.Output()) // do not allow flag package to print any error or usage message
	fs.Usage = func() {}     // overwirte the defaultUsage, do nothing
	child := &Node{
		Name:        strings.TrimSpace(name),
		FullName:    fullName,
		Description: description,
		root:        false,
		parser:      n.parser,
		FS:          fs,
		Childs:      nil,
		Handler:     handler,
	}
	n.Childs = append(n.Childs, child)
	return child
}

// search for the node that most match the input.
// base on BFS
func (n *Node) Search(names []string) *Node {
	if len(names) == 0 {
		return nil
	}
	queue := make(chan *Node, 100)
	defer close(queue)
	index := 0 // index for names
	length := len(names)
	var lastMatchedNode *Node = nil
	queue <- n
	for len(queue) > 0 {
		node := <-queue
		if node.root { // enqueue all root's childs
			for _, child := range node.Childs {
				queue <- child
			}
		} else {
			if node.Name == names[index] { // find the matched node
				lastMatchedNode = node
				// flush out the brother nodes
				for len(queue) > 0 {
					<-queue
				}
				index++              // to match later string
				if index == length { // reach the end of names
					return lastMatchedNode
				}
				for _, child := range node.Childs { // enqueue the matched node's childs
					queue <- child
				}
			}
		}

	}
	return lastMatchedNode
}

// base on DFS
func (n *Node) Print() {
	DFS(n, 0, func(n *Node, depth int) {
		indent := ""
		for i := 0; i < depth; i++ {
			indent += "   "
		}
		if depth > 0 {
			fmt.Printf("%s|- %s", indent, n.Name)
		} else {
			fmt.Printf("%s%s", indent, n.Name)
		}
		// print flags
		// not print dirctly, should use the settings output io.Writer
		if n.FS != nil {
			i := 0
			flags := ""
			n.FS.VisitAll(func(f *flag.Flag) {
				if i == 0 {
					flags += " {"
				} else {
					flags += ", "
				}
				flags += "-%s"
				flags = fmt.Sprintf(flags, f.Name)
				i++
			})
			if i > 0 {
				flags += "}"
			}
			fmt.Print(flags)
		}
		fmt.Print("\n")
	})
}

func DFS(n *Node, depth int, visit func(*Node, int)) {
	visit(n, depth)
	for _, child := range n.Childs {
		DFS(child, depth+1, visit)
	}
}

func (n *Node) Split(input string) ([]string, error) {
	return shlex.Split(input)
}

func (n *Node) Parse(ctx context.Context, input string) error {
	names, err := n.Split(input)
	if err != nil {
		return err
	}
	// parse command
	targetNode := n.Search(names)
	if targetNode == nil {
		return errors.New("unknown")
	}

	// parse flag
	flagStr := names[len(strings.Split(targetNode.FullName, " ")):]
	if err := targetNode.FS.Parse(flagStr); err != nil {
		return err
	}

	// exec handler
	if targetNode.Handler != nil {
		err = targetNode.Handler(ctx, targetNode.FS)
	} else {
		err = errors.New("no such command")
	}

	// reset flag value
	targetNode.FS.Visit(func(f *flag.Flag) {
		f.Value.Set(f.DefValue)
	})
	return err
}

func (n *Node) Output() io.Writer {
	return n.parser.output
}

func (n *Node) Usage() {
	if n.FS == nil {
		return
	}
	if n.Name == "" {
		fmt.Fprintf(n.Output(), "Usage:\n")
	} else {
		fmt.Fprintf(n.Output(), "Usage of %s:\n", n.FullName)
	}
	// TODO !!!!!! like root node parser, it's FS is init to nil, will cause error
	// NOTICE same error!!!!
	n.FS.PrintDefaults()
}

func (n *Node) GenerateSuggestions() []prompt.Suggest {
	suggestions := []prompt.Suggest{}
	for _, child := range n.Childs {
		suggestions = append(suggestions, prompt.Suggest{Text: child.Name, Description: child.Description})
	}
	// range map in particular order
	if n.FS != nil {
		n.FS.VisitAll(func(f *flag.Flag) {
			suggestions = append(suggestions, prompt.Suggest{Text: "-" + f.Name, Description: f.Usage})
		})
	}

	// keys := make([]string, 0, len(n.Flags))
	// for k := range n.Flags {
	// 	keys = append(keys, k)
	// }
	// sort.Strings(keys)
	// for _, k := range keys {
	// 	suggestions = append(suggestions, prompt.Suggest{Text: "-" + n.Flags[k].Name, Description: n.Flags[k].Description})
	// }
	return suggestions
}

// const (
// 	cmdNameRegexpGroup  = `name`
// 	cmdFlagsRegexpGroup = `flags`
// )
// use regexp to verify if the input consist with the CMD obj
// func (n *Node) Match(input string) (map[string]string, error) {
// 	var cmdExpression = regexp.MustCompile(
// 		`(?P<` + cmdNameRegexpGroup + `>^` + cmd.FirstName + `\s+` + cmd.SecondName + `)` +
// 			`(?P<` + cmdFlagsRegexpGroup + `>.*)`)

// 	matchedGroups := cmdExpression.FindStringSubmatch(input)
// 	if len(matchedGroups) == 0 {
// 		// the input not match this cmd's full name
// 		return nil, errors.New("wrong cmd parser")
// 	}
// 	groups := make(map[string]string)
// 	for i, group := range cmdExpression.SubexpNames() {
// 		if i != 0 && group != "" {
// 			groups[group] = matchedGroups[i]
// 		}
// 	}
// 	return groups, nil
// }

// func (cmd *CMDGroup) Parse(input string) error {
// 	groups, err := cmd.Match(input)
// 	if err != nil {
// 		return err
// 	}
// 	flags, err := shlex.Split(groups[cmdFlagsRegexpGroup])
// 	if err != nil {
// 		return err
// 	}
// 	cmd.FS.Parse(flags)
// 	return nil
// }

// type Completer struct {
// 	Suggestions    []prompt.Suggest
// 	SubSuggestions []map[string][]prompt.Suggest
// }

// func (c *Completer) Init(cmds []CMDGroup) {
// 	for _, cmd := range cmds {
// 		c.Suggestions = append(c.Suggestions, prompt.Suggest{Text: cmd.FirstName})
// 	}
// }
