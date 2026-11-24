// Package rx
// @Title  title
// @Description  desc
package rx

import (
	"github.com/velkrel/unitengine/engine/pkg/config"
	inf "github.com/velkrel/unitengine/engine/pkg/interfaces"
	"github.com/velkrel/unitengine/engine/pkg/log"
	"github.com/velkrel/unitengine/engine/pkg/rpc/remote/handler"
	"github.com/smallnest/rpcx/server"
)

type rpcxServer struct {
	svr      *server.Server
	listener *RpcxListener
	logger   log.ILoggerX
	handler  *handler.Handler
}

func NewRpcxServer() inf.IRemoteServer {
	return &rpcxServer{}
}

func (rs *rpcxServer) SetLogger(logger log.ILoggerX) {
	if logger != nil {
		rs.logger = logger
		if rs.listener != nil {
			rs.listener.logger = logger
		}
	}
}

func (rs *rpcxServer) Init(sf inf.IRpcSenderFactory) {
	rs.listener = &RpcxListener{
		cliFactory: sf,
		logger:     rs.logger,
		handler:    rs.handler,
	}
	rs.svr = server.NewServer()
}

func (rs *rpcxServer) SetHandler(h *handler.Handler) {
	rs.handler = h
	if rs.listener != nil {
		rs.listener.handler = h
	}
}

func (rs *rpcxServer) Serve(conf *config.RPCServer, nodeUid string) error {
	// 注册rpc监听服务
	if err := rs.svr.RegisterName("RpcxListener", rs.listener, ""); err != nil {
		return err
	}
	rs.logger.Infof("rpcx server listening at: %s", conf.Addr)
	return rs.svr.Serve(conf.Protoc, conf.Addr)
}

func (rs *rpcxServer) Close() {
	if rs.svr == nil {
		return
	}
	_ = rs.svr.Close()
	rs.svr = nil
}
