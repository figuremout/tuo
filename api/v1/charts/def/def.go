package def

import (
	nodes_def "github.com/githubzjm/tuo/api/v1/nodes/def"
)

const (
	ROUTE_CPU        = nodes_def.ROUTE_NODES + "/%s/cpu"
	SUBROUTE_PERCENT = "/percent"
)
