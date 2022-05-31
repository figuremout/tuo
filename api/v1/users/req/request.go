package req

import (
	"net/http"

	"github.com/githubzjm/tuo/api/v1/users/def"
	httpClient "github.com/githubzjm/tuo/internal/pkg/http/client"
)

func ReqRegister(username, email, password string) (int, *def.RegisterResp, error) {
	var respBody def.RegisterResp
	reqBody := def.RegisterReq{
		Username: username,
		Email:    email,
		Password: password,
	}
	statusCode, err := httpClient.Request(http.MethodPost, def.ROUTE_REGISTER, reqBody, nil, &respBody)
	return statusCode, &respBody, err
}

func ReqLogin(username, password string) (int, *def.LoginResp, error) {
	var respBody def.LoginResp
	reqBody := def.LoginReq{
		Username: username,
		Password: password,
	}
	statusCode, err := httpClient.Request(http.MethodPost, def.ROUTE_LOGIN, reqBody, nil, &respBody)
	return statusCode, &respBody, err
}

func ReqQueryrUser(token string) {

}
