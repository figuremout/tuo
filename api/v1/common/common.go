package common

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/githubzjm/tuo/api/v1/def"
)

// handler for bind request body failed
// func BindFailedHandler(c *gin.Context, err error) {
// 	c.JSON(http.StatusBadRequest, def.BaseResp{
// 		Error: "cannot parse data: " + err.Error(),
// 	})
// }

func BindHandler(obj interface{}, c *gin.Context) bool {
	if err := c.ShouldBind(obj); err != nil {
		c.JSON(http.StatusBadRequest, def.BaseResp{
			Error: "cannot parse data: " + err.Error(),
		})
		return false
	}
	return true
}

// handler for check ID in URL
// check if the param ID is number and
// also check if the user have the privilege to operate the model of param ID
func CheckParamIDHandler(visitor uint, idName string, checker func(uint, uint) (bool, error), c *gin.Context) int {
	target, err := strconv.Atoi(c.Param(idName))
	if target < 0 {
		return -1
	}
	if err != nil {
		// param ID is not a number
		c.JSON(http.StatusBadRequest, def.BaseResp{
			Error: fmt.Sprintf("bad request: %s URL %s for ROUTE %s", c.Request.Method, c.Request.RequestURI, c.FullPath()),
		})
		return -1
	}
	if isPass := CheckPrivilegeHandler(visitor, uint(target), checker, c); !isPass {
		return -1
	}
	return target
}

// check privilege, if fail, construct error resp
func CheckPrivilegeHandler(visitor, target uint, checker func(uint, uint) (bool, error), c *gin.Context) bool {
	isPass, err := checker(visitor, target)
	if err != nil {
		c.JSON(http.StatusInternalServerError, def.BaseResp{
			Error: err.Error(),
		})
		return false
	}
	if !isPass {
		c.JSON(http.StatusForbidden, def.BaseResp{
			Error: "you do not have the privilege",
		})
		return false
	}
	return true
}
