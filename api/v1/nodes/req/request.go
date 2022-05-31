package req

import (
	"fmt"
	"net/http"

	v1_def "github.com/githubzjm/tuo/api/v1/def"
	"github.com/githubzjm/tuo/api/v1/nodes/def"
	httpClient "github.com/githubzjm/tuo/internal/pkg/http/client"
)

// TODO the request.go should be like this func, pass in a whole struct not values
func ReqCreate(clusterID string, reqBody *def.CreateReq, token string) (int, *def.CreateResp, error) {
	var respBody def.CreateResp

	statusCode, err := httpClient.Request(http.MethodPost, fmt.Sprintf(def.ROUTE_NODES, clusterID),
		reqBody, map[string]string{v1_def.HEADER_TOKEN_KEY: token}, &respBody)
	return statusCode, &respBody, err
}

func ReqQueryAll(clusterID, token string) (int, *def.QueryAllResp, error) {
	var respBody def.QueryAllResp

	statusCode, err := httpClient.Request(http.MethodGet, fmt.Sprintf(def.ROUTE_NODES, clusterID),
		nil, map[string]string{v1_def.HEADER_TOKEN_KEY: token}, &respBody)
	return statusCode, &respBody, err
}
