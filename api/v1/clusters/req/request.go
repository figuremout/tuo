package req

import (
	"net/http"

	"github.com/githubzjm/tuo/api/v1/clusters/def"
	v1_def "github.com/githubzjm/tuo/api/v1/def"
	httpClient "github.com/githubzjm/tuo/internal/pkg/http/client"
)

func ReqCreate(clusterName, token string) (int, *def.CreateResp, error) {
	var respBody def.CreateResp
	reqBody := def.CreateReq{
		ClusterName: clusterName,
	}

	statusCode, err := httpClient.Request(http.MethodPost, def.ROUTE_CLUSTERS,
		reqBody, map[string]string{v1_def.HEADER_TOKEN_KEY: token}, &respBody)
	return statusCode, &respBody, err
}

func ReqQueryAll(token string) (int, *def.QueryAllResp, error) {
	var respBody def.QueryAllResp

	statusCode, err := httpClient.Request(http.MethodGet, def.ROUTE_CLUSTERS,
		nil, map[string]string{v1_def.HEADER_TOKEN_KEY: token}, &respBody)
	return statusCode, &respBody, err
}
