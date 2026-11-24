// Package pprofservice
// @Title  性能性分析服务
// @Description  性能性分析服务
package pprofservice

import (
	"errors"
	"net/http"
	_ "net/http/pprof"

	"github.com/gin-gonic/gin"
	"github.com/velkrel/unitengine/engine/pkg/utils/httpx/router_center"

	systemConfig "github.com/velkrel/unitengine/engine/pkg/config"
	"github.com/velkrel/unitengine/engine/pkg/core"
	inf "github.com/velkrel/unitengine/engine/pkg/interfaces"
	"github.com/velkrel/unitengine/engine/pkg/services"
	"github.com/velkrel/unitengine/engine/pkg/sysModule/httpmodule"
	"github.com/velkrel/unitengine/engine/pkg/sysService/pprofservice/config"
)

func RegisterPprofService() {
	services.SetService("PprofService", func() inf.IService { return &PprofService{} })
	systemConfig.RegisterServiceConf(&systemConfig.ServiceConfig{
		ServiceName:   "PprofService",
		ConfName:      "pprof",
		ConfType:      "yaml",
		ConfPath:      "",
		CfgCreator:    func() interface{} { return &config.PprofConf{} },
		DefaultSetFun: config.SetPprofConfDefault,
		OnChangeFun: func() {
			// TODO 配置变更回调
		},
	})
}

type PprofService struct {
	core.Service

	httpModule *httpmodule.HttpModule
}

func (ps *PprofService) OnInit() error {
	conf, ok := ps.GetServiceCfg().(*config.PprofConf)
	if !ok || conf == nil || conf.PprofConf == nil || conf.PprofConf.Addr == "" {
		return errors.New("PprofService requires PprofConf.Addr")
	}
	status := systemConfig.Release
	if ctx := ps.GetNodeContext(); ctx != nil {
		if cfg := ctx.GetConfig(); cfg != nil {
			status = cfg.GetStatus()
		}
	}
	// 构造 Gin 引擎前设置模式，避免 release 模式的 worker 向 stdout 写调试提示。
	gin.SetMode(status)
	router := router_center.NewGroupHandlerPool()
	router.RegisterGroupHandler("", ps.routerHandler)
	ps.httpModule = httpmodule.NewHttpModule(conf.PprofConf, status).SetRouter(router)
	_, err := ps.AddModule(ps.httpModule)
	return err
}

func (ps *PprofService) routerHandler(r *gin.RouterGroup) {
	r.GET("/debug/pprof", gin.WrapH(http.DefaultServeMux))
	r.GET("/debug/pprof/*pprof", gin.WrapH(http.DefaultServeMux))
}

func (ps *PprofService) OnStart() error {
	return ps.httpModule.StartBackground()
}

func (ps *PprofService) OnRelease() {
	ps.ReleaseAllChildModule()
}
