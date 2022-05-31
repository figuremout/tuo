package nodes

import (
	"github.com/githubzjm/tuo/internal/pkg/mysql"
	"github.com/githubzjm/tuo/internal/pkg/mysql/models"
	"github.com/githubzjm/tuo/internal/server/clusters"
)

// read heartbeat from influxdb
// TODO
func IsAlive() {

}

func CreateNode(value *models.Node) (isConflict bool, err error) {
	return mysql.Create(value)
}

func QueryNode(query *models.Node) (*models.Node, error) {
	var node models.Node
	isExist, err := mysql.Query(query, &node)
	if err != nil {
		return nil, err // internal error
	}
	if !isExist {
		return nil, nil // no match
	}
	return &node, nil // matched
}

func QueryAllNodes(query *models.Node) ([]models.Node, error) {
	var nodes []models.Node
	isExist, err := mysql.QueryAll(query, &nodes)
	if err != nil {
		return nil, err // internal error
	}
	if !isExist {
		return nil, nil
	}
	return nodes, nil
}

func CheckPrivilege(visitor uint, targetNode uint) (bool, error) {
	// get target node
	node, err := QueryNode(&models.Node{ID: targetNode})
	if err != nil {
		return false, err
	}
	if node == nil {
		// for now, user can only operate its own info
		// may adapt privilege system later
		return false, nil
	}

	isPass, err := clusters.CheckPrivilege(visitor, node.ClusterID)
	if err != nil {
		return false, err
	}
	return isPass, nil
}
