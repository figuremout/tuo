package command

import (
	"io"
	"os"
)

// group mechanism inspired by gin router
// inspired from: https://pkg.go.dev/flag#Func
type Engine struct {
	Node
	output io.Writer
}

func New() *Engine {
	engine := &Engine{
		Node: Node{
			Name:  "<root>",
			root:  true,
			FS:    nil,
			Flags: nil,
		},
		output: os.Stdout,
	}
	engine.Node.engine = engine
	return engine
}

// func (e *Engine) SetOutput(output io.Writer) {
// 	e.output = output
// }

// func Test() {
// 	e := New()
// 	docker := e.Group("docker", "Manage users", func(flags Flags) error {
// 		if flags["v"].IsSet {
// 			fmt.Println("0.0.1")
// 		}
// 		if flags["D"].IsSet {
// 			fmt.Println("Debug mode on")
// 		}
// 		return nil
// 	})
// 	docker.AddFlag("v", "Print version information and quit", false)
// 	docker.AddFlag("D", "Enable debug mode", false)
// 	{
// 		build := docker.Group("build", "Build an image from a Dockerfile", func(flags Flags) error {
// 			if flags["f"].IsSet {
// 				if flags["f"].Value == "" {
// 					flags["f"].Owner.Usage()
// 					return errors.New("dockerfile path illegal")
// 				}
// 				fmt.Printf("Building image from %s\n", flags["f"].Value)
// 			}
// 			return nil
// 		})
// 		build.AddFlag("f", "Name of the Dockerfile (Default is 'PATH/Dockerfile')", true)
// 	}
// 	e.Print()
// 	// fmt.Printf("Error: %v\n", e.Parse("docker -v"))
// 	// fmt.Printf("Error: %v\n", e.Parse("docker -v 1"))
// 	// fmt.Printf("Error: %v\n", e.Parse("docker -f"))
// 	// fmt.Printf("Error: %v\n", e.Parse("docker build -f file"))
// 	if err := e.Parse("docker build -f \"\""); err != nil {
// 		fmt.Printf("Error: %v\n", err)
// 	}
// }
