package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/c-bata/go-prompt"
	v1_def "github.com/githubzjm/tuo/api/v1/def"
	"github.com/githubzjm/tuo/api/v1/nodes/def"
	"github.com/githubzjm/tuo/api/v1/nodes/req"
	"github.com/githubzjm/tuo/internal/client/cache"
	"github.com/githubzjm/tuo/internal/client/common"
	"github.com/gorilla/websocket"
)

func DeployNode() {
	clusterID := prompt.Input("cluster ID: ", func(d prompt.Document) []prompt.Suggest {
		return []prompt.Suggest{}
	})
	nodeID := prompt.Input("node ID: ", func(d prompt.Document) []prompt.Suggest {
		return []prompt.Suggest{}
	})
	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: fmt.Sprintf(def.ROUTE_NODES+"/%s"+def.SUBRUOTE_DEPLOY, clusterID, nodeID)}
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

	kill_progress := make(chan int)
	upload_progress := make(chan int)
	start_progress := make(chan int)
	defer close(kill_progress)
	defer close(upload_progress)
	defer close(start_progress)

	go common.Progress(kill_progress, false, "killing agent", "kill agent success")

	// read message continuously
	for {
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
		var message def.DeployResp
		if err := c.ReadJSON(&message); err != nil {
			if closeErr, ok := err.(*websocket.CloseError); ok {
				if closeErr.Code == websocket.CloseNormalClosure { // server close normal
					kill_progress <- 100
					upload_progress <- 100
					start_progress <- 100
					return
				}
			}
			kill_progress <- 101
			upload_progress <- 101
			start_progress <- 101
			fmt.Printf("Error: %s\n", err.Error())
			return
		}

		if message.Step <= 0 {
			kill_progress <- 101
			upload_progress <- 101
			start_progress <- 101
			fmt.Printf("Error: %s\n", err.Error())
			return
		} else if message.Step == 1 {
			kill_progress <- 100
			// fmt.Printf("\n")
			go common.Progress(upload_progress, false, "uploading agent", "upload agent success")
		} else if message.Step == 2 {
			upload_progress <- 100
			// fmt.Printf("\n")
			go common.Progress(start_progress, false, "starting agent", "start agent success")
		} else if message.Step == 3 {
			start_progress <- 100
			return
		}

		// if !message.IsUpload {
		// 	kill_progress <- 101
		// 	fmt.Printf("Upload Error: %s\n", message.Error)
		// 	return
		// } else {
		// 	if message.IsUp {
		// 		kill_progress <- 100
		// 		return
		// 	} else {
		// 		kill_progress <- 50
		// 	}
		// }
		// fmt.Printf("recv: %#v\n", message)
	}
}

func CreateNode() {
	clusterID := prompt.Input("cluster ID: ", func(d prompt.Document) []prompt.Suggest {
		return []prompt.Suggest{}
	})
	// clusterID_uint, err := strconv.ParseUint(clusterID, 10, 64)
	// if err != nil {
	// 	fmt.Printf("cluster ID should be a uint")
	// 	return
	// }
	nodeName := prompt.Input("node name: ", func(d prompt.Document) []prompt.Suggest {
		return []prompt.Suggest{}
	})
	host := prompt.Input("host: ", func(d prompt.Document) []prompt.Suggest {
		return []prompt.Suggest{}
	})
	port := prompt.Input("ssh port: ", func(d prompt.Document) []prompt.Suggest {
		return []prompt.Suggest{}
	})
	port_uint, err := strconv.ParseUint(port, 10, 64)
	if err != nil {
		fmt.Printf("port should be a uint")
		return
	}
	user := prompt.Input("ssh user: ", func(d prompt.Document) []prompt.Suggest {
		return []prompt.Suggest{}
	})
	password := prompt.Input("ssh password: ", func(d prompt.Document) []prompt.Suggest {
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

	status, resp, err := req.ReqCreate(clusterID, &def.CreateReq{
		NodeName: nodeName,
		Host:     host,
		Port:     uint(port_uint),
		User:     user,
		Password: password,
	}, string(token))
	if isSuccess := handleRespErr(status, &resp.BaseResp, err); isSuccess {
		fmt.Printf("node %s created\n", nodeName)
	}
}

func QueryAllNodes(clusterID string) {
	// clusterID := prompt.Input("cluster ID: ", func(d prompt.Document) []prompt.Suggest {
	// 	return []prompt.Suggest{}
	// })
	// clusterID_uint, err := strconv.ParseUint(clusterID, 10, 64)
	// if err != nil {
	// 	fmt.Printf("cluster ID should be a uint")
	// 	return
	// }

	token, err := cache.Get(cache.KeyToken)
	if err != nil {
		fmt.Println(err)
		return
	}
	if string(token) == "" {
		fmt.Printf("no cached token found\n")
		return
	}

	status, resp, err := req.ReqQueryAll(clusterID, string(token))

	if isSuccess := handleRespErr(status, &resp.BaseResp, err); isSuccess {
		if len(resp.Nodes) == 0 {
			fmt.Printf("No nodes\n")
		} else {
			// TODO only show part of info, if need whole set flag -l
			fmt.Printf("%-8s %-8s %-8s %-35s %-35s\n", "NodeID", "NodeName", "IsAlive", "CreateAt", "UpdateAt")
			for _, node := range resp.Nodes {
				fmt.Printf("%-8d %-8s %-8v %-35s %-35s\n", node.NodeID, node.NodeName, node.IsAlive, node.CreatedAt, node.UpdatedAt)
			}
		}
	}
}
