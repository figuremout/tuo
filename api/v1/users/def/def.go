package def

import (
	"time"

	"github.com/githubzjm/tuo/api/v1/def"
)

const (
	ROUTE_USERS    = def.ROUTE_V1 + "/users"
	ROUTE_REGISTER = ROUTE_USERS + SUBROUTE_REGISTER
	ROUTE_LOGIN    = ROUTE_USERS + SUBROUTE_LOGIN

	SUBROUTE_REGISTER       = "/register"
	SUBROUTE_LOGIN          = "/login"
	TOKEN_VALID_DURATION_NS = time.Hour * 24
)

// If a field is decorated with binding:"required" and has a empty value when binding, an error will be returned.
type RegisterReq struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type RegisterResp struct {
	def.BaseResp
}

type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type LoginResp struct {
	def.BaseResp
	UserID   uint   `json:"userid"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Token    string `json:"token"`
}

type QueryResp struct {
	def.BaseResp
	CreatedAt time.Time `json:"createAt"`
	UpdatedAt time.Time `json:"updateAt"`
	UserID    uint      `json:"userid"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
}

type UpdateReq struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type UpdateResp struct {
	def.BaseResp
}

type DelResp struct {
	def.BaseResp
}
