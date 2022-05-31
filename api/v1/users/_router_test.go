package users

import (
	"net/http"
	"strings"
	"testing"

	jsonUtil "github.com/githubzjm/tuo/internal/pkg/utils/json"
	"github.com/githubzjm/tuo/test"
)

const (
	username string = "testusername"
	email    string = "testemail"
	password string = "testpwd"
)

func TestRegister(t *testing.T) {
	data := registerReqData{
		Username: username,
		Email:    email,
		Password: password,
	}
	payload := strings.NewReader(jsonUtil.Stringfy(data, ""))
	resp, err := http.Post("http://localhost:8080/api/v1/users/register", "application/json", payload)
	if err != nil {
		t.Fatal(err)
	}
	test.CheckAPIResp(resp, t)
	t.Log("test register pass")
}

func TestLogin(t *testing.T) {
	Login(t)
	t.Log("test login pass")
}

func Login(t *testing.T) string {
	data := loginReqData{
		Username: username,
		Password: password,
	}
	payload := strings.NewReader(jsonUtil.Stringfy(data, ""))
	resp, err := http.Post("http://localhost:8080/api/v1/users/login", "application/json", payload)
	if err != nil {
		t.Fatal(err)
	}
	respStruct := test.CheckAPIResp(resp, t)

	respData, ok := respStruct.Data.(map[string]interface{})
	if !ok {
		t.Fatal("Response data parse failed")
	}
	token, ok := respData["token"].(string)
	if !ok {
		t.Fatal("Response data token parse failed")
	}
	return token
}

func TestUpdate(t *testing.T) {
	token := Login(t)

	data := updateReqData{
		Password: "testpwd",
		Email:    "",
	}
	payload := strings.NewReader(jsonUtil.Stringfy(data, ""))
	req, err := http.NewRequest("PATCH", "http://localhost:8080/api/v1/users", payload)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("token", token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	test.CheckAPIResp(resp, t)
	t.Log("test update pass")
}
