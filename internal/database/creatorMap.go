package database

import (
	"gorm.io/gorm"
	"supreme-flamego/internal/database/driver"
)

type Creator interface {
	Create(ip string, port string, userName string, password string, dbName string) (*gorm.DB, error)
}

type DbModel interface {
	DbKey() string
}

func init() {
	typeMap["mysql"] = driver.MySQLCreator{}
}

var typeMap = make(map[string]Creator)

func getCreatorByType(dbType string) Creator {
	return typeMap[dbType]
}
