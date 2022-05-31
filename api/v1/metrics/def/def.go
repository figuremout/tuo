package def

import (
	"time"

	"github.com/githubzjm/tuo/api/v1/def"
	nodes_def "github.com/githubzjm/tuo/api/v1/nodes/def"
)

const (
	ROUTE_CPU        = nodes_def.ROUTE_NODES + "/%s" + "/cpu"
	SUBROUTE_PERCENT = "/percent"
)

type CPUPercentResp struct {
	def.BaseResp
	CPU     string    `json:"cpu"`
	Percent float64   `json:"percent"`
	Time    time.Time `json:"time"`
}
