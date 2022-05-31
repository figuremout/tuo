package models

import (
	"time"

	"gorm.io/gorm"
)

// one User has many Clusters
// has many association: https://gorm.io/zh_CN/docs/has_many.html

var _ IModel = &User{}
var _ IModel = &Cluster{}
var _ IModel = &Monitor{}
var _ IModel = &Node{}

type IModel interface {
	GetID() uint
}

type User struct {
	gorm.Model
	ID       uint      `gorm:"primarykey;autoIncrement;autoIncrementIncrement:1"`
	Name     string    `gorm:"column:name;not null;unique"`
	Email    string    `gorm:"column:email;not null"`
	Password string    `gorm:"column:password;not null"`
	Clusters []Cluster `gorm:"foreignKey:UserID;references:ID"` // has many association
}

func (u *User) GetID() uint {
	return u.ID
}

// has one association: https://gorm.io/docs/has_one.html
type Cluster struct {
	gorm.Model
	ID uint `gorm:"primarykey;autoIncrement;autoIncrementIncrement:1"`
	// uniqueIndexes of same name will form a composite unique index, https://gorm.io/docs/indexes.html#uniqueIndex
	UserID uint `gorm:"column:userID;not null;uniqueIndex:ui"` // user who create the cluster
	// go type string will be mapped into longtext in mysql, which will cause MySQL Error 1170 when uniqueIndex is set, trans to varchar type will solve
	// set size:255 will trans the type from longtext to varchar(256), issue like: https://github.com/go-gorm/gorm/issues/3744
	Name    string  `gorm:"column:name;not null;uniqueIndex:ui;size:255"`
	Monitor Monitor `gorm:"foreignKey:ClusterID;references:ID"` // has one association
	Nodes   []Node  `gorm:"foreignKey:ClusterID;references:ID"` // has many association
}

func (c *Cluster) GetID() uint {
	return c.ID
}

type Monitor struct {
	gorm.Model
	ID        uint          `gorm:"primarykey;autoIncrement;autoIncrementIncrement:1"`
	ClusterID uint          `gorm:"column:clusterID;not null"`
	Interval  time.Duration `gorm:"column:interval;not null"`
}

func (m *Monitor) GetID() uint {
	return m.ID
}

type Node struct {
	gorm.Model
	ID        uint   `gorm:"primarykey;autoIncrement;autoIncrementIncrement:1"`
	ClusterID uint   `gorm:"column:clusterID;not null;uniqueIndex:ui"`
	Name      string `gorm:"column:name;not null;uniqueIndex:ui;size:255"`

	Host     string `gorm:"column:host;not null"`     // ssh host
	Port     uint   `gorm:"column:port;not null"`     // ssh port
	User     string `gorm:"column:user;not null"`     // ssh user
	Password string `gorm:"column:password;not null"` // ssh pwd

	// IsAlive bool `gorm:"column:isAlive;not null"`
	// use influxdb's checks deadman function to realize heartbeat
}

func (n *Node) GetID() uint {
	return n.ID
}

// one monitor has one alert or just put alert in monitor
// use influxdb checks to realize alert
// type Alert struct {
// 	gorm.Model
// 	ID        uint   `gorm:"primarykey;autoIncrement"`
// 	ClusterID uint   `gorm:"column:clusterID;not null;uniqueIndex:ui"`
// 	Name      string `gorm:"column:name;not null;uniqueIndex:ui;size:255"`
// 	Metric    string `gorm:"column:metric;not null"`

// 	Condition string `gorm:"column:condition;not null"`
// }

// func CreateTable() error {
// 	_, err := mysql.SqlDB.Exec(
// 		"CREATE TABLE IF NOT EXISTS user (" +
// 			"id INT UNSIGNED AUTO_INCREMENT, " +
// 			"username VARCHAR(20) NOT NULL, " +
// 			"email VARCHAR(50) NOT NULL, " +
// 			"password VARCHART(30) NOT NULL" +
// 			"PRIMARY KEY (id));")
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func Register(username, email string) (int64, error) {
// 	res, err := mysql.SqlDB.Exec("INSERT INTO user (username, email, password) VALUES (?, ?, ?)", username, email)
// 	if err != nil {
// 		return -1, err
// 	}
// 	id, err := res.LastInsertId()
// 	if err != nil {
// 		return -1, err
// 	}
// 	return id, nil
// }
