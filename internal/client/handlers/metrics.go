package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/githubzjm/tuo/api/v1/charts"
	v1_def "github.com/githubzjm/tuo/api/v1/def"
	"github.com/githubzjm/tuo/api/v1/metrics/def"
	"github.com/githubzjm/tuo/internal/client/cache"
	"github.com/githubzjm/tuo/routers"
	"github.com/gorilla/websocket"
)

var (
	clusterIDParam = "clusterID"
	nodeIDParam    = "nodeID"
)

func CPUPercent(ctx context.Context, port string) {
	var clusterID, nodeID string
	fmt.Printf("cluster ID: ")
	fmt.Scanln(&clusterID)
	fmt.Printf("node ID: ")
	fmt.Scanln(&nodeID)

	token, err := cache.Get(cache.KeyToken)
	if err != nil {
		fmt.Println(err)
		return
	}
	if string(token) == "" {
		fmt.Printf("no cached token found\n")
		return
	}

	gin.SetMode(gin.ReleaseMode)

	router := routers.Init()
	router.GET(fmt.Sprintf(def.ROUTE_CPU+def.SUBROUTE_PERCENT, ":"+clusterIDParam, ":"+nodeIDParam), func(c *gin.Context) {
		c.HTML(http.StatusOK, "percent.html", gin.H{
			"title": "CPU Percent",
			"url":   "ws://localhost:8080/api/v1/clusters/" + c.Param(clusterIDParam) + "/nodes/" + c.Param(nodeIDParam) + "/cpu/percent", // TODO should not hard code
			"token": strings.Trim(string(token), "\""),
		})
	})
	fmt.Printf("charts displaying at http://localhost:"+port+def.ROUTE_CPU+def.SUBROUTE_PERCENT+"\n", clusterID, nodeID)

	// gracefully exit solution: https://pkg.go.dev/github.com/gin-gonic/gin@v1.7.7#readme-manually
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	<-ctx.Done()
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("Server forced to shutdown: %v\n", err)
	}

	fmt.Printf("Server exiting\n")
}

func CPUPercentTest(ctx context.Context, port string) {
	// clusterID := prompt.Input("cluster ID: ", func(d prompt.Document) []prompt.Suggest {
	// 	return []prompt.Suggest{}
	// })
	// nodeID := prompt.Input("node ID: ", func(d prompt.Document) []prompt.Suggest {
	// 	return []prompt.Suggest{}
	// })

	var clusterID, nodeID string
	fmt.Printf("cluster ID: ")
	fmt.Scanln(&clusterID)
	fmt.Printf("node ID: ")
	fmt.Scanln(&nodeID)
	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: fmt.Sprintf(def.ROUTE_CPU+def.SUBROUTE_PERCENT, clusterID, nodeID)}
	// fmt.Printf("connecting to %s\n", u.String())

	token, err := cache.Get(cache.KeyToken)
	if err != nil {
		fmt.Println(err)
		return
	}
	if string(token) == "" {
		fmt.Printf("no cached token found\n")
		return
	}

	header := http.Header{
		v1_def.HEADER_SecWebSocketProtocol_KEY: []string{string(token)},
	}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), header)
	if err != nil {
		fmt.Println("dial:", err)
		return
	}
	defer c.Close()

	// read message continuously
	for {
		select {
		case <-ctx.Done():
			c.Close() // inform server to close websocket
			return
		default:
			// _, message, err := c.ReadMessage()
			// if err != nil {
			// 	// if closeErr, ok := err.(*websocket.CloseError); ok {
			// 	// 	if closeErr.Code == websocket.ClosePolicyViolation { // server close connection due to auth fail
			// 	// 		return
			// 	// 	}
			// 	// }
			// 	fmt.Printf("Error: %s", err.Error())
			// 	return
			// }
			var message def.CPUPercentResp
			if err := c.ReadJSON(&message); err != nil {
				if closeErr, ok := err.(*websocket.CloseError); ok {
					if closeErr.Code == websocket.CloseNormalClosure { // server close normal
						return
					}
				}
				fmt.Printf("Error: %s\n", err.Error())
				return
			}

			fmt.Printf("recv: %#v\n", message)
		}
	}
}
