package service

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/micro/go-micro/v2/logger"
)

type MysqlService struct {
	dbs map[string]map[string]*gorm.DB //example -> master 2级
}

func NewMysqlService(cfgMap map[string]interface{}) (*MysqlService, error) {
	dbs := make(map[string]map[string]*gorm.DB)

	for key, val := range cfgMap { //example -> map
		v, ok := val.(map[string]interface{})
		if !ok {
			log.Fatalf("read mysql config failed: %s", key)
		}

		for key2, val2 := range v { // master -> map
			v2, ok := val2.(map[string]interface{})
			if !ok {
				log.Fatalf("read mysql config failed: %s", key)
			}

			if _, ok := dbs[key]; !ok {
				dbs[key] = make(map[string]*gorm.DB)
			}

			db, err := instanceDb(v2)
			if err != nil {
				log.Fatalf("instanceDb mysql failed: %s => %s", key, key2)
			}

			dbs[key][key2] = db
		}
	}

	return &MysqlService{
		dbs: dbs,
	}, nil
}

func (srv *MysqlService) ChooseDb(index string, ms string) *gorm.DB {
	if _, ok := srv.dbs[index]; !ok {
		log.Errorf("choose database failed: %s\n", index)
		return nil
	}

	if _, ok := srv.dbs[index][ms]; !ok {
		log.Errorf("choose database failed: %s, %s\n", index, ms)
		return nil
	}

	return srv.dbs[index][ms]
}

func instanceDb(cfg map[string]interface{}) (*gorm.DB, error) {
	user := cfg["user"].(string)
	pass := cfg["pass"].(string)
	host := cfg["host"].(string)
	port := cfg["port"].(string)
	database := cfg["database"].(string)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, port, database)
	db, err := gorm.Open("mysql", dsn)
	return db, err
}
