package database

import (
	"gorm.io/gorm"
	"supreme-flamego/config"
	"supreme-flamego/pkg/logger"
	"sync"
)

var (
	dbs         = make(map[string]*gorm.DB)
	mux         sync.RWMutex
	migrateList = make(map[string][]interface{})
)

func InitDB() {
	sources := config.GetConfig().Databases
	for _, source := range sources {
		setDbByKey(source.Key, mustCreateGorm(source))
		if source.Key == "" {
			source.Key = "*"
		}
		logger.NameSpace("database").Info("create datasource %s => %s:%s", source.Key, source.IP, source.PORT)
	}
	for key, models := range migrateList {
		db := GetDB(key)
		if db == nil {
			logger.NameSpace("database").Fatal("fail to find db for key:%s", key)
			return
		}
		err := db.AutoMigrate(models...)
		if err != nil {
			logger.NameSpace("database").Fatal(err)
			return
		}
		logger.NameSpace("database").Info("migrate datasource %s success", key)
	}
}

func GetDB(key string) *gorm.DB {
	mux.Lock()
	defer mux.Unlock()
	return dbs[key]
}

func setDbByKey(key string, db *gorm.DB) {
	if key == "" {
		key = "*"
	}
	if GetDB(key) != nil {
		logger.NameSpace("database").Error("duplicate db key: " + key)
	}
	mux.Lock()
	defer mux.Unlock()
	dbs[key] = db
}

func mustCreateGorm(database config.Datasource) *gorm.DB {
	var creator = getCreatorByType(database.Type)
	if creator == nil {
		logger.NameSpace("database").Fatalf("fail to find creator for types:%s", database.Type)
		return nil
	}
	db, err := creator.Create(database.IP, database.PORT, database.USER, database.PASSWORD, database.DATABASE)
	if err != nil {
		logger.NameSpace("database").Fatal(err)
		return nil
	}

	return db
}

// AutoMigrate 暂时注册一下数据库模型 将在InitDB的时候自动使用
func AutoMigrate(dbKey string, dst ...interface{}) error {
	migrateList[dbKey] = append(migrateList[dbKey], dst...)
	return nil
}
