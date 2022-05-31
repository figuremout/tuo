package clusters

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/githubzjm/tuo/api/v1/clusters/def"
	"github.com/githubzjm/tuo/api/v1/common"
	v1_def "github.com/githubzjm/tuo/api/v1/def"
	"github.com/githubzjm/tuo/internal/pkg/jwt"
	"github.com/githubzjm/tuo/internal/pkg/mysql/models"
	cm "github.com/githubzjm/tuo/internal/server/clusters"
	"github.com/githubzjm/tuo/routers"
)

var clusterIDParam string = "clusterID"

func createHandler(c *gin.Context) {
	claims := c.MustGet("claims").(*jwt.CustomClaims)
	var data def.CreateReq
	if !common.BindHandler(&data, c) {
		return
	}

	newCluster := models.Cluster{Name: data.ClusterName, UserID: claims.UserID}
	isConflict, err := cm.CreateCluster(&newCluster)
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
				Error: "cluster name already taken",
			},
		})
		return
	} else {
		c.JSON(http.StatusCreated, def.CreateResp{
			ClusterID: newCluster.ID,
		})
		return
	}
}

func queryAllHandler(c *gin.Context) {
	claims := c.MustGet("claims").(*jwt.CustomClaims)

	clusters, err := cm.QueryAllClusters(&models.Cluster{UserID: claims.UserID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, def.QueryAllResp{
			BaseResp: v1_def.BaseResp{
				Error: err.Error(),
			},
		})
		return
	}
	if clusters == nil {
		// not exist
		c.JSON(http.StatusUnauthorized, def.QueryAllResp{
			BaseResp: v1_def.BaseResp{
				Error: "clusters not exist",
			},
		})
		return
	}

	var data []def.QueryOneResp
	for _, cluster := range clusters {
		data = append(data, def.QueryOneResp{
			CreatedAt:   cluster.CreatedAt,
			UpdatedAt:   cluster.UpdatedAt,
			ClusterID:   cluster.ID,
			UserID:      cluster.UserID,
			ClusterName: cluster.Name,
		})
	}
	c.JSON(http.StatusOK, def.QueryAllResp{
		Clusters: data,
	})
}

func queryOneHandler(c *gin.Context) {
	claims := c.MustGet("claims").(*jwt.CustomClaims)

	target := common.CheckParamIDHandler(claims.UserID, clusterIDParam, cm.CheckPrivilege, c)
	if target < 0 {
		return
	}

	cluster, err := cm.QueryCluster(&models.Cluster{ID: uint(target)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, v1_def.BaseResp{
			Error: err.Error(),
		})
		return
	}
	if cluster == nil {
		// not exist
		c.JSON(http.StatusUnauthorized, v1_def.BaseResp{
			Error: "cluster not exist",
		})
		return
	}

	c.JSON(http.StatusOK, def.QueryOneResp{
		CreatedAt:   cluster.CreatedAt,
		UpdatedAt:   cluster.UpdatedAt,
		ClusterID:   cluster.ID,
		UserID:      cluster.UserID,
		ClusterName: cluster.Name,
	})
}

func updateHandler(c *gin.Context) {
	var data def.UpdateReq
	if !common.BindHandler(&data, c) {
		return
	}

	claims := c.MustGet("claims").(*jwt.CustomClaims)

	target := common.CheckParamIDHandler(claims.UserID, clusterIDParam, cm.CheckPrivilege, c)
	if target < 0 {
		return
	}

	if err := cm.UpdateCluster(&models.Cluster{ID: uint(target)}, &models.Cluster{
		Name: data.ClusterName,
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

	target := common.CheckParamIDHandler(claims.UserID, clusterIDParam, cm.CheckPrivilege, c)
	if target < 0 {
		return
	}

	err := cm.DeleteClusterWithAssociation(&models.Cluster{ID: uint(target)})
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

func Routers(e *gin.Engine) {
	clusters := e.Group(def.ROUTE_CLUSTERS)
	clusters.Use(jwt.JWTAuth())
	{
		clusters.GET("", queryAllHandler) // get user's all clusters
		clusters.POST("", createHandler)  // create a cluster

		// nested group
		clusterID := clusters.Group("/:" + clusterIDParam)
		{
			clusterID.GET("", queryOneHandler) // query cluster
			clusterID.DELETE("", delHandler)   // delete cluster
			clusterID.PATCH("", updateHandler) // update cluster
		}
	}
}

func init() {
	routers.Include(Routers)
}
