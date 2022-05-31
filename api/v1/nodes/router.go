package nodes

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/githubzjm/tuo/api/v1/common"
	v1_def "github.com/githubzjm/tuo/api/v1/def"
	"github.com/githubzjm/tuo/api/v1/nodes/def"
	"github.com/githubzjm/tuo/internal/pkg/jwt"
	"github.com/githubzjm/tuo/internal/pkg/mysql/models"
	sshUtil "github.com/githubzjm/tuo/internal/pkg/ssh"
	cm "github.com/githubzjm/tuo/internal/server/clusters"
	nm "github.com/githubzjm/tuo/internal/server/nodes"
	"github.com/githubzjm/tuo/routers"
	"github.com/gorilla/websocket"
)

var (
	clusterIDParam = "clusterID"
	nodeIDParam    = "nodeID"
)

func createHandler(c *gin.Context) {
	claims := c.MustGet("claims").(*jwt.CustomClaims)
	var data def.CreateReq
	if !common.BindHandler(&data, c) {
		return
	}

	targetCluster := common.CheckParamIDHandler(claims.UserID, clusterIDParam, cm.CheckPrivilege, c)
	if targetCluster < 0 {
		return
	}

	newNode := models.Node{
		ClusterID: uint(targetCluster),
		Name:      data.NodeName,
		Host:      data.Host,
		Port:      data.Port,
		User:      data.User,
		Password:  data.Password,
	}
	isConflict, err := nm.CreateNode(&newNode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, def.CreateResp{
			BaseResp: v1_def.BaseResp{
				Error: err.Error(),
			},
		})
		return
	}
	if isConflict {
		c.JSON(http.StatusBadRequest, def.CreateResp{
			BaseResp: v1_def.BaseResp{
				Error: "node name already taken",
			},
		})
		return
	} else {
		c.JSON(http.StatusCreated, def.CreateResp{
			NodeID:  newNode.ID,
			IsAlive: true,
		})
		return
	}
}

func queryAllHandler(c *gin.Context) {
	claims := c.MustGet("claims").(*jwt.CustomClaims)

	targetCluster := common.CheckParamIDHandler(claims.UserID, clusterIDParam, cm.CheckPrivilege, c)
	if targetCluster < 0 {
		return
	}

	nodes, err := nm.QueryAllNodes(&models.Node{ClusterID: uint(targetCluster)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, def.QueryAllResp{
			BaseResp: v1_def.BaseResp{
				Error: err.Error(),
			},
		})
		return
	}
	if nodes == nil {
		// not exist
		c.JSON(http.StatusUnauthorized, def.QueryAllResp{
			BaseResp: v1_def.BaseResp{
				Error: "nodes not exist",
			},
		})
		return
	}

	var data []def.QueryOneResp
	for _, node := range nodes {
		data = append(data, def.QueryOneResp{
			CreatedAt: node.CreatedAt,
			UpdatedAt: node.UpdatedAt,
			NodeID:    node.ID,
			ClusterID: node.ID,
			NodeName:  node.Name,
			Host:      node.Host,
			Port:      node.Port,
			User:      node.User,
			// IsAlive:   node.IsAlive,
		})
	}
	c.JSON(http.StatusOK, def.QueryAllResp{
		Nodes: data,
	})
}

func deployHandler(c *gin.Context) {
	claims := c.MustGet("claims").(*jwt.CustomClaims)
	wsConn := c.MustGet("wsConn").(*websocket.Conn)
	defer wsConn.Close() // 1006 abnormal closure, last choice
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

	// TODO to optimize the query time, checkParam, checkPrivilege may cause some waste query times
	node, err := nm.QueryNode(&models.Node{ID: uint(targetNodeID)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, def.DeployResp{
			BaseResp: v1_def.BaseResp{
				Error: err.Error(),
			},
		})
		return
	}

	// deploy agent to node
	sshClient, err := sshUtil.InitSshClient(node.User, node.Password, node.Host, node.Port)
	if err != nil {
		wsConn.WriteJSON(&def.DeployResp{
			BaseResp: v1_def.BaseResp{
				Error: err.Error(),
			},
			Step: -1,
		})
		return
	}
	sftpClient, err := sshUtil.InitSftpClient(sshClient)
	if err != nil {
		wsConn.WriteJSON(&def.DeployResp{
			BaseResp: v1_def.BaseResp{
				Error: err.Error(),
			},
			Step: -1,
		})
		return
	}
	// kill agent
	if _, err := sshUtil.Exec(sshClient,
		"pid=$(ps -ef | grep /tuo/agent | grep -v grep | awk '{print $2}');[[ -n ${pid} ]] && kill -9 ${pid};exit 0"); err != nil {
		wsConn.WriteJSON(&def.DeployResp{
			BaseResp: v1_def.BaseResp{
				Error: err.Error(),
			},
			Step: -1,
		})
		return
	} else {
		wsConn.WriteJSON(&def.DeployResp{
			Step: 1,
		})
	}

	// upload agent bin
	if err := sshUtil.UploadFile("./bin/agent", "/tuo", sftpClient); err != nil {
		wsConn.WriteJSON(&def.DeployResp{
			BaseResp: v1_def.BaseResp{
				Error: err.Error(),
			},
			Step: -1,
		})
		return
	} else {
		wsConn.WriteJSON(&def.DeployResp{
			Step: 2,
		})
	}
	// start agent
	if _, err := sshUtil.Exec(sshClient, "nohup /tuo/agent > /tuo/agent.log 2>&1 &"); err != nil {
		wsConn.WriteJSON(&def.DeployResp{
			BaseResp: v1_def.BaseResp{
				Error: err.Error(),
			},
			Step: -1,
		})
		return
	} else {
		wsConn.WriteJSON(&def.DeployResp{
			Step: 3,
		})
	}

}

func Routers(e *gin.Engine) {
	nodes := e.Group(fmt.Sprintf(def.ROUTE_NODES, ":"+clusterIDParam))
	nodes.Use(jwt.JWTAuth())
	{
		nodes.GET("", queryAllHandler)

		nodes.POST("", createHandler)

		// id := nodes.Group("/:"+nodeIDParam)

		// nested group
		// id := clusters.Group("/:ID")
		// {
		// 	id.GET("", queryOneHandler) // query cluster
		// 	id.DELETE("", delHandler)   // delete cluster
		// 	id.PATCH("", updateHandler) // update cluster
		// }
	}

	// should not nested in the /nodes route, because it uses special auth middleware
	deploy := e.Group(fmt.Sprintf(def.ROUTE_NODES, ":"+clusterIDParam) + "/:" + nodeIDParam + def.SUBRUOTE_DEPLOY) // websocket
	deploy.Use(jwt.WSJWTAuth())
	{
		deploy.Any("", deployHandler)
	}
}

func init() {
	routers.Include(Routers)
}
