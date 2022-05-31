package alerts

import (
	"net/http"

	"github.com/gin-gonic/gin"
	v1_def "github.com/githubzjm/tuo/api/v1/def"
	"github.com/githubzjm/tuo/internal/pkg/jwt"
)

// pseudo code
type AlertManager struct{}

func (am *AlertManager) CreateCheck(name, level, typ, value, clusterID string) (string, error)
func (am *AlertManager) IsEndpointExist(userID string) (string, bool)
func (am *AlertManager) CreateEndpoint(name, method, url, token string) (string, error)
func (am *AlertManager) CreateRule(endpointID, ruleName, level string) (string, error)

func addAlertHandler(c *gin.Context) {
	am := &AlertManager{}
	var checkName, level, typ, value, clusterID, userID, endpointName, url, token, ruleName string

	// Create check
	// level optional value: "UNKNOWN", "OK", "INFO", "CRIT", "WARN"
	// typ optional value: "range", "lesser", "greater"
	_, err := am.CreateCheck(checkName, level, typ, value, clusterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, v1_def.BaseResp{
			Error: err.Error(),
		})
		return
	}

	// Check if notification endpoint exists
	var endpointID string
	endpointID, isExist := am.IsEndpointExist(userID)
	if !isExist {
		// create endpoint
		var err error
		endpointID, err = am.CreateEndpoint(endpointName, "GET", url, token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, v1_def.BaseResp{
				Error: err.Error(),
			})
			return
		}
	}

	// Create notification rule
	if _, err := am.CreateRule(ruleName, endpointID, level); err != nil {
		c.JSON(http.StatusInternalServerError, v1_def.BaseResp{
			Error: err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, v1_def.BaseResp{
		Error: "",
	})

}

func Routers(e *gin.Engine) {
	alerts := e.Group("")
	alerts.Use(jwt.JWTAuth())
	{
		alerts.POST("", addAlertHandler)
	}
}
