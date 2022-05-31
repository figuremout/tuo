package clusters

import (
	"github.com/githubzjm/tuo/internal/pkg/mysql"
	"github.com/githubzjm/tuo/internal/pkg/mysql/models"
)

func CreateCluster(value *models.Cluster) (isConflict bool, err error) {
	return mysql.Create(value)
}
func QueryCluster(query *models.Cluster) (*models.Cluster, error) {
	var cluster models.Cluster
	isExist, err := mysql.Query(query, &cluster)
	if err != nil {
		return nil, err // internal error
	}
	if !isExist {
		return nil, nil // no match
	}
	return &cluster, nil // matched
}

func QueryAllClusters(query *models.Cluster) ([]models.Cluster, error) {
	var clusters []models.Cluster
	isExist, err := mysql.QueryAll(query, &clusters)
	if err != nil {
		return nil, err // internal error
	}
	if !isExist {
		return nil, nil
	}
	return clusters, nil
}

func UpdateCluster(query, values *models.Cluster) error {
	return mysql.Update(query, values)
}
func DeleteCluster(query *models.Cluster) error {
	return mysql.Delete(query, &models.Cluster{})
}
func DeleteClusterWithAssociation(query *models.Cluster) error {
	return mysql.DeleteWithAssociatons(query)
}

// if internal error occurs, return err != nil
// if pass, return true, else return false
func CheckPrivilege(visitor uint, targetCluster uint) (bool, error) {
	cluster, err := QueryCluster(&models.Cluster{ID: targetCluster, UserID: visitor})
	if err != nil {
		return false, err
	}
	if cluster == nil {
		// for now, user can only operate its own info
		// may adapt privilege system later
		return false, nil
	} else {
		return true, nil
	}
}
