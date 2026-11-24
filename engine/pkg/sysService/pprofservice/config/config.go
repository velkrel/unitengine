// Package config
// @Title  请填写文件名称（需要改）
// @Description  请填写文件描述（需要改）
package config

import (
	"time"

	"github.com/spf13/viper"
	"github.com/velkrel/unitengine/engine/pkg/utils/httpx"
)

type PprofConf struct {
	PprofConf *httpx.Conf `binding:"required"`
}

func SetPprofConfDefault(parser *viper.Viper) {
	parser.SetDefault("PprofConf", &httpx.Conf{
		Addr:              "0.0.0.0:99",
		ReadHeaderTimeout: time.Second * 10,
		IdleTimeout:       time.Second * 30,
		CachePath:         "",
		ResourceRootPath:  "",
		HttpDir:           "",
		StaticDir:         "",
		Account:           nil,
	})
}
