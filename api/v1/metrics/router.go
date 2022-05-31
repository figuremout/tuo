package metrics

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/githubzjm/tuo/api/v1/common"
	v1_def "github.com/githubzjm/tuo/api/v1/def"
	"github.com/githubzjm/tuo/api/v1/metrics/def"
	"github.com/githubzjm/tuo/internal/pkg/jwt"
	cm "github.com/githubzjm/tuo/internal/server/clusters"
	mm "github.com/githubzjm/tuo/internal/server/metrics"
	nm "github.com/githubzjm/tuo/internal/server/nodes"
	"github.com/githubzjm/tuo/routers"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var (
	clusterIDParam = "clusterID"
	nodeIDParam    = "nodeID"
)

func cpuPercentHandler(c *gin.Context) {
	claims := c.MustGet("claims").(*jwt.CustomClaims)
	wsConn := c.MustGet("wsConn").(*websocket.Conn)
	defer wsConn.Close()
	defer wsConn.WriteControl(websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
		time.Now().Add(time.Second)) // TODO should not all return normal close, like ParamID check failed should return abnormal close

	targetClusterID := common.CheckParamIDHandler(claims.UserID, clusterIDParam, cm.CheckPrivilege, c)
	if targetClusterID < 0 {
		return
	}

	targetNodeID := common.CheckParamIDHandler(claims.UserID, nodeIDParam, nm.CheckPrivilege, c)
	if targetNodeID < 0 {
		return
	}

	// TODO should use targetNodeID to confirm which node to query metrics
	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()
	for range ticker.C {
		m, err := mm.CPUPercent()
		if err != nil {
			wsConn.WriteJSON(&def.CPUPercentResp{
				BaseResp: v1_def.BaseResp{
					Error: err.Error(),
				},
			})
			return
		}
		if m == nil {
			return
		}
		if err := wsConn.WriteJSON(&def.CPUPercentResp{ // TODO consider if client close conn
			CPU:     m["cpu"].(string),
			Percent: m["_value"].(float64),
			Time:    m["_time"].(time.Time),
		}); err != nil {
			log.Info(err)
			return
		}
	}

}

func Routers(e *gin.Engine) {
	cpu := e.Group(fmt.Sprintf(def.ROUTE_CPU, ":"+clusterIDParam, ":"+nodeIDParam))
	cpu.Use(jwt.WSJWTAuth())
	{
		percent := cpu.Group(def.SUBROUTE_PERCENT)
		percent.Any("", cpuPercentHandler)
	}
}

func init() {
	routers.Include(Routers)
}
