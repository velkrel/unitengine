// Package pool
// @Title  title
// @Description  desc
package pool

import (
	"github.com/velkrel/unitengine/engine/pkg/def"
	inf "github.com/velkrel/unitengine/engine/pkg/interfaces"
	"github.com/velkrel/unitengine/engine/pkg/rpc/remote/gr"
	"github.com/velkrel/unitengine/engine/pkg/rpc/remote/nt"
	"github.com/velkrel/unitengine/engine/pkg/rpc/remote/rx"
)

type ListenerCreator func(cliFactory inf.IRpcSenderFactory) interface{}

// remoteFactory 存储无状态工厂函数（非实例），每次调用创建新实例。
// 属于 init() 后只读的注册表，多 Node 安全共享。
var remoteFactory = map[string]func() inf.IRemoteServer{
	def.RpcTypeRpcx: func() inf.IRemoteServer { return rx.NewRpcxServer() },
	def.RpcTypeGrpc: func() inf.IRemoteServer { return gr.NewGrpcServer() },
	def.RpcTypeNats: func() inf.IRemoteServer { return nt.NewNatsServer() },
}

// Register 注册远程服务器工厂函数
func Register(tp string, factory func() inf.IRemoteServer) {
	remoteFactory[tp] = factory
}

// CreateServer 创建新的远程服务器实例（每次调用返回新实例）
func CreateServer(tp string) inf.IRemoteServer {
	if f, ok := remoteFactory[tp]; ok {
		return f()
	}
	return nil
}
