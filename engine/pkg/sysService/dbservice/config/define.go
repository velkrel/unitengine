package config

import (
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	mongodbmodule "github.com/velkrel/unitengine/engine/pkg/sysModule/mongomodule"
	"github.com/velkrel/unitengine/engine/pkg/sysModule/mysqlmodule"
)

type DBService struct {
	MysqlConf *mysqlmodule.Conf          `binding:""` // mysql配置
	RedisConf *redis.Options             `binding:""` // redis配置
	MongoConf *mongodbmodule.MongoConfig `binding:""` // mongodb配置
}

func DefaultDBService(parser *viper.Viper) {
	parser.SetDefault("MysqlConf", &mysqlmodule.Conf{
		UserName:           "",
		Passwd:             "",
		Net:                "tcp",
		Addr:               "0.0.0.0:3306",
		TimeZone:           "Local",
		Timeout:            time.Second * 10,
		ReadTimeout:        time.Second * 10,
		WriteTimeout:       time.Second * 10,
		SetConnMaxIdleTime: 0,
		SetConnMaxLifetime: 0,
		SetMaxIdleConns:    0,
		SetMaxOpenConns:    0,
	})
}
