# command
A simple but powerful command parser for self-defined command string (NOT os.Args).\
Code design inpired by gin.RouterGroup.

中英文
知乎 CSDN Stackoverflow打广告
https://stackoverflow.com/questions/34118732/parse-a-command-line-string-into-flags-and-arguments-in-golang：回复：
实际上利用flag可以做到解析字符串，见https://pkg.go.dev/flag#Func 的example。我也利用这个功能实现了自己的命令解释器，见链接

## Features
- Support multi-level command, like `"users"`, `"users login"`, `"users login foo"`
- Support flags (reply on the built-in package "flag"), like `"users login -username foo -password bar"`
- Support [go-prompt](https://github.com/c-bata/go-prompt) suggestions generation
- Support shell-style quoting (rely on [shlex](https://github.com/google/shlex)), like `"users login -username \"foo bar\""`

## Shortcomings
- Take all flag values as string, if you need the value to be other types, you should convert it by yourself in handler.
- Do not support pass in argument directly without flag, like `docker exec 23dfsaf`

## Quick Start
```golang
	e := command.New()
	users := e.Group("users", "Manage users")
	{
		login := users.Group("login", "Login")
        login.AddFlag("username", "user's name", func(v string) error {
		    fmt.Printf("Login username: %s\n", v)
		    return nil
	    })
		login.AddFlag("pwd", "user's password", func(v string) error {
			fmt.Printf("Login pwd: %s\n", v)
			return nil
		})
	}
    {
        register := users.Group("register", "Register")
        login.AddFlag("username", "user's name", func(v string) error {
		    fmt.Printf("Register username: %s\n", v)
		    return nil
	    })
    }

    e.Parse("users login -username foo -pwd bar")
    // result:
    // Login username: foo
    // Login pwd: bar

    user.Parse("users register -username foo")
    // result:
    // Register username: foo

    e.Print()
    // result:
    // <root>
    // |- users
    //    |- login {-username, -pwd}
    //    |- register {-username}
```

## Examples
### Work with [go-prompt](https://github.com/c-bata/go-prompt)
```golang

```
动图、视频参考：https://github.com/Trendyol/docker-shell/tree/62605e9b1a8eb12d3a6f50242bdde942d9920831
https://github.com/golobby/repl/tree/da6c0b527cea52bb525fc544358de234382e3ee1