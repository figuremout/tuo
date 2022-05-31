package handlers

import (
	"fmt"

	"github.com/c-bata/go-prompt"
	"github.com/githubzjm/tuo/api/v1/users/req"
	"github.com/githubzjm/tuo/internal/client/cache"
	"github.com/githubzjm/tuo/internal/client/common"
)

func Register() {
	username := prompt.Input("username: ", func(d prompt.Document) []prompt.Suggest {
		return []prompt.Suggest{}
	})
	email := prompt.Input("email: ", func(d prompt.Document) []prompt.Suggest {
		return []prompt.Suggest{}
	})
	password := prompt.Input("password: ", func(d prompt.Document) []prompt.Suggest {
		return []prompt.Suggest{}
	})
	status, resp, err := req.ReqRegister(username, email, password)
	if isSuccess := handleRespErr(status, &resp.BaseResp, err); isSuccess {
		fmt.Printf("Register success\n")
	}

}

func Login() {
	username := prompt.Input("username: ", func(d prompt.Document) []prompt.Suggest {
		return []prompt.Suggest{}
	})
	password := prompt.Input("password: ", func(d prompt.Document) []prompt.Suggest {
		return []prompt.Suggest{}
	})
	status, resp, err := req.ReqLogin(username, password)
	if isSuccess := handleRespErr(status, &resp.BaseResp, err); isSuccess {
		common.ChangePrefix(username)
		cache.Set(cache.KeyToken, resp.Token)
		cache.Append(cache.KeyUser, &cache.User{
			Username: resp.Username,
			UserID:   resp.UserID,
		})
		fmt.Printf("Welcome %s!\n", username)
	}
}

func Logout() {
	cache.Del(cache.KeyToken)
	fmt.Print("Logout success\n")
	common.ChangePrefix("")
}
