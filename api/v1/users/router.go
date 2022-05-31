package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/githubzjm/tuo/api/v1/common"
	v1_def "github.com/githubzjm/tuo/api/v1/def"
	"github.com/githubzjm/tuo/api/v1/users/def"
	"github.com/githubzjm/tuo/internal/pkg/jwt"
	"github.com/githubzjm/tuo/internal/pkg/mysql/models"
	um "github.com/githubzjm/tuo/internal/server/users" // users manager
	"github.com/githubzjm/tuo/routers"
)

var userIDParam = "userID"

//curl -X GET "localhost:8080/ping"
// curl -X POST "localhost:8080/api/v1/users/register" -d '{"username": "user1", "email": "testemail", "password": "pwd"}' -i  -H 'Content-Type: application/json'
// curl -X GET "localhost:8080/api/v1/auth/time" -i  -H 'token:...'
func registerHandler(c *gin.Context) {
	var data def.RegisterReq
	if !common.BindHandler(&data, c) {
		return
	}

	isConflict, err := um.CreateUser(&models.User{Name: data.Username, Email: data.Email, Password: data.Password})
	if err != nil {
		c.JSON(http.StatusInternalServerError, def.RegisterResp{
			BaseResp: v1_def.BaseResp{
				Error: err.Error(),
			},
		})
		return
	}
	if isConflict {
		c.JSON(http.StatusBadRequest, def.RegisterResp{
			BaseResp: v1_def.BaseResp{
				Error: "username already registered",
			},
		})
		return
	} else {
		c.JSON(http.StatusCreated, def.RegisterResp{})
		return
	}
}

func loginHandler(c *gin.Context) {
	var data def.LoginReq
	if !common.BindHandler(&data, c) {
		return
	}

	user := getUser(c, &models.User{Name: data.Username})
	if user == nil {
		return
	}

	if user.Password == data.Password {
		// password correct
		token, err := um.GenerateToken(c, user)
		if err != nil {
			// generate token failed
			c.JSON(http.StatusInternalServerError, def.LoginResp{
				BaseResp: v1_def.BaseResp{
					Error: err.Error(),
				},
			})
			return
		}
		// login success
		c.JSON(http.StatusOK, def.LoginResp{
			UserID:   user.ID,
			Username: user.Name,
			Email:    user.Email,
			Token:    token,
		})
		return
	} else {
		// password wrong
		c.JSON(http.StatusUnauthorized, def.LoginResp{
			BaseResp: v1_def.BaseResp{
				Error: "password wrong",
			},
		})
		return
	}
}

func queryHandler(c *gin.Context) {
	claims := c.MustGet("claims").(*jwt.CustomClaims)

	target := common.CheckParamIDHandler(claims.UserID, userIDParam, um.CheckPrivilege, c)
	if target < 0 {
		return
	}

	user := getUser(c, &models.User{ID: uint(target)})
	if user == nil {
		return
	}
	c.JSON(http.StatusOK, def.QueryResp{
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		UserID:    user.ID,
		Username:  user.Name,
		Email:     user.Email,
	})
}

func updateHandler(c *gin.Context) {
	var data def.UpdateReq
	if !common.BindHandler(&data, c) {
		return
	}

	claims := c.MustGet("claims").(*jwt.CustomClaims)
	target := common.CheckParamIDHandler(claims.UserID, userIDParam, um.CheckPrivilege, c)
	if target < 0 {
		return
	}

	if err := um.UpdateUser(&models.User{ID: uint(target)}, &models.User{
		Name:     data.Username,
		Email:    data.Email,
		Password: data.Password,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, def.UpdateResp{
			BaseResp: v1_def.BaseResp{
				Error: err.Error(),
			},
		})
		return
	}
	c.JSON(http.StatusOK, def.UpdateResp{})
}

func delHandler(c *gin.Context) {
	claims := c.MustGet("claims").(*jwt.CustomClaims)
	target := common.CheckParamIDHandler(claims.UserID, userIDParam, um.CheckPrivilege, c)
	if target < 0 {
		return
	}

	err := um.DeleteUserWithAssociation(&models.User{ID: uint(target)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, def.DelResp{
			BaseResp: v1_def.BaseResp{
				Error: err.Error(),
			},
		})
		return
	}

	// gin do not allow http.StatusNoContent to return body
	c.JSON(http.StatusNoContent, def.DelResp{})
}

func getUser(c *gin.Context, query *models.User) *models.User {
	user, err := um.QueryUser(query)

	if err != nil {
		c.JSON(http.StatusInternalServerError, v1_def.BaseResp{
			Error: err.Error(),
		})
		return nil
	}
	if user == nil {
		// not exist
		c.JSON(http.StatusUnauthorized, v1_def.BaseResp{
			Error: "user not exist",
		})
		return nil
	}
	return user
}

func Routers(e *gin.Engine) {
	users := e.Group(def.ROUTE_USERS)
	{
		users.POST(def.SUBROUTE_REGISTER, registerHandler)
		users.POST(def.SUBROUTE_LOGIN, loginHandler)
		// nested group
		id := users.Group("/:" + userIDParam)
		// TODO api like PATCH /register will match to this group
		id.Use(jwt.JWTAuth())
		{
			id.PATCH("", updateHandler) // update user
			id.GET("", queryHandler)    // query user
			id.DELETE("", delHandler)   // delete user
		}
	}
}

func init() {
	routers.Include(Routers)
}
