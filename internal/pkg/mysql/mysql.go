package mysql

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/VividCortex/mysqlerr"
	"github.com/githubzjm/tuo/internal/pkg/mysql/models"

	// mysql driver, its init func will register itself into database/sql, then can be reached by sql.Open()
	mysqlDriver "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var SqlDB *sql.DB   // can exec sql, query, ping and so on
var GormDB *gorm.DB // can map db table into struct, can also do CRUD

func InitDB(userName, password, ip, port, dbName string) error {
	// dsn form: "user:password@tcp(IP:port)/database"
	// To handle time.Time correctly, you need to include parseTime as a parameter. (more parameters)
	// To fully support UTF-8 encoding, you need to change charset=utf8 to charset=utf8mb4. See this article for a detailed explanation
	dsn := strings.Join([]string{userName, ":", password, "@tcp(", ip, ":", port, ")/", dbName,
		"?charset=utf8mb4&parseTime=True&loc=Local"}, "")
	var err error
	SqlDB, err = sql.Open("mysql", dsn) // use mysql driver
	if err != nil {
		return err
	}
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	SqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	SqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	SqlDB.SetConnMaxLifetime(time.Hour)

	// init orm rely on the SqlDB connection
	// TODO need to config logger to be consist with logrus
	GormDB, err = gorm.Open(mysql.New(mysql.Config{
		Conn: SqlDB,
	}), &gorm.Config{})
	if err != nil {
		return err
	}

	// Migrate the schema
	// If table not exist will create automatically
	err = GormDB.AutoMigrate(
		&models.User{},
		&models.Cluster{},
		&models.Monitor{},
		&models.Node{},
	)
	if err != nil {
		return err
	}
	return nil
}

// value is a struct pointer or a slice of pointers, which will be used to create records and also be filled in auto fileds like ID
func Create(value models.IModel) (isConflict bool, err error) {
	// (*gorm.DB).Create() func will use a point of struct to create,
	// but also fill in some fields like ID, CreatedAt field to the struct automatically,
	// so no need to query database to get these fields
	res := GormDB.Create(value)
	if res.Error != nil {
		if driverErr, ok := res.Error.(*mysqlDriver.MySQLError); ok {
			if driverErr.Number == mysqlerr.ER_NO_REFERENCED_ROW_2 {
				// create fails due to association, like attempting to create a record while its owner not exists
				return false, errors.New("attempt to create with non-existent owner")
			}
			if driverErr.Number == mysqlerr.ER_DUP_ENTRY {
				// unique fields already exist
				return true, nil
			}
		}
		return false, res.Error
	}
	// create success
	return false, nil
}

// will not query associations
func Query(query, dest models.IModel) (isExist bool, err error) {
	res := GormDB.Where(query).First(dest)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			// not exist
			// distinguish this error because it should be coped with diffrent way from other internal errors
			return false, nil
		} else {
			// other error
			return false, res.Error
		}
	}
	// dest is assigned values that queried
	return true, nil
}

func QueryAll(query models.IModel, dest interface{}) (isExist bool, err error) {
	res := GormDB.Where(query).Find(dest)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			// not exist
			return false, nil
		} else {
			// other error
			return false, res.Error
		}
	}
	return true, nil
}

func Update(query, values models.IModel) error {
	//  updating with struct it will only update non-zero fields by default
	res := GormDB.Where(query).Updates(values)
	if res.Error != nil {
		// TODO what to do with recordnotfound err
		return res.Error
	}
	return nil
}
func Delete(query, value models.IModel) error {
	// value can be a empty instance, only used to indicate the operated table
	// .Unscoped(): delete permanently
	// cannot use .Clauses(clause.Returning{}) to return deleted data, MySQL do not support returning syntax
	res := GormDB.Unscoped().Where(query).Delete(value)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

// Associations will only be deleted if the deleting records’s primary key is not zero
// https://gorm.io/docs/associations.html#Delete-with-Select
func DeleteWithAssociatons(value models.IModel) error {
	if value.GetID() <= 0 { // autoIncrement start from 1
		return errors.New("record ID is illegal")
	}
	res := GormDB.Unscoped().Select(clause.Associations).Delete(value)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
