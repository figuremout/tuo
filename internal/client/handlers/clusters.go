package handlers

import (
	"fmt"

	"github.com/c-bata/go-prompt"
	"github.com/githubzjm/tuo/api/v1/clusters/req"
	"github.com/githubzjm/tuo/internal/client/cache"
)

func CreateCluster() {
	clusterName := prompt.Input("cluster name: ", func(d prompt.Document) []prompt.Suggest {
		return []prompt.Suggest{}
	})
	token, err := cache.Get(cache.KeyToken)
	if err != nil {
		fmt.Println(err)
		return
	}
	if string(token) == "" {
		fmt.Printf("no cached token found\n")
		return
	}

	status, resp, err := req.ReqCreate(clusterName, string(token))
	if isSuccess := handleRespErr(status, &resp.BaseResp, err); isSuccess {
		fmt.Printf("cluster %s created\n", clusterName)
	}
}

func QueryAllClusters() {
	token, err := cache.Get(cache.KeyToken)
	if err != nil {
		fmt.Println(err)
		return
	}
	if string(token) == "" {
		fmt.Printf("no cached token found\n")
		return
	}
	status, resp, err := req.ReqQueryAll(string(token))
	if isSuccess := handleRespErr(status, &resp.BaseResp, err); isSuccess {
		if len(resp.Clusters) == 0 {
			fmt.Printf("No clusters\n")
		} else {
			fmt.Printf("%-8s %-8s %-8s %-35s %-35s\n", "ID", "UserID", "Name", "CreateAt", "UpdateAt")
			for _, cluster := range resp.Clusters {
				fmt.Printf("%-8d %-8d %-8s %-35s %-35s\n", cluster.ClusterID, cluster.UserID, cluster.ClusterName, cluster.CreatedAt, cluster.UpdatedAt)
			}
		}
	}
}
