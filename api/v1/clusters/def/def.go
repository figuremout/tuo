package def

import (
	"time"

	"github.com/githubzjm/tuo/api/v1/def"
)

const (
	ROUTE_CLUSTERS = def.ROUTE_V1 + "/clusters"
)

type CreateReq struct {
	ClusterName string `json:"clusterName" binding:"required"`
}
type CreateResp struct {
	def.BaseResp
	ClusterID uint `json:"clusterID"`
}

type QueryOneResp struct {
	def.BaseResp
	CreatedAt   time.Time `json:"createAt"`
	UpdatedAt   time.Time `json:"updateAt"`
	ClusterID   uint      `json:"clusterID"`
	UserID      uint      `json:"userID"`
	ClusterName string    `json:"clusterName"`
}
type QueryAllResp struct {
	def.BaseResp
	Clusters []QueryOneResp `json:"clusters"`
}

type UpdateReq struct {
	ClusterName string `json:"clusterName"`
}
type UpdateResp struct {
	def.BaseResp
}

type DelResp struct {
	def.BaseResp
}
