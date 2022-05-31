package def

import (
	"time"

	cluster_def "github.com/githubzjm/tuo/api/v1/clusters/def"
	"github.com/githubzjm/tuo/api/v1/def"
)

const (
	ROUTE_NODES     = cluster_def.ROUTE_CLUSTERS + "/%s" + SUBROUTE_NODES // "/clusters/:clusterID/nodes"
	SUBROUTE_NODES  = "/nodes"
	SUBRUOTE_DEPLOY = "/deploy"
)

type CreateReq struct {
	// ClusterID uint   `json:"clusterID" binding:"required"`
	NodeName string `json:"nodeName" binding:"required"`
	Host     string `json:"host" binding:"required"`
	Port     uint   `json:"port" binding:"required"`
	User     string `json:"user" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type CreateResp struct {
	def.BaseResp
	NodeID  uint `json:"nodeID"`
	IsAlive bool `json:"isAlive"`
}

type QueryOneResp struct {
	def.BaseResp
	NodeID    uint      `json:"nodeID"`
	CreatedAt time.Time `json:"createAt"`
	UpdatedAt time.Time `json:"updateAt"`
	ClusterID uint      `json:"clusterID"`
	NodeName  string    `json:"nodeName"`
	Host      string    `json:"host"`
	Port      uint      `json:"port"`
	User      string    `json:"user"`
	IsAlive   bool      `json:"isAlive"`
}

type QueryAllResp struct {
	def.BaseResp
	Nodes []QueryOneResp `json:"nodes"`
}

type UpdateReq struct {
	NodeName string `json:"nodeName"`
	Host     string `json:"host"`
	Port     uint   `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
}
type UpdateResp struct {
	def.BaseResp
}

type DeployResp struct {
	def.BaseResp
	Step int `json:"step"`
}
