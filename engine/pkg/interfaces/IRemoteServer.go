// Package interfaces
// @Title  title
// @Description  desc
package interfaces

import (
	"github.com/velkrel/unitengine/engine/pkg/config"
)

//type IRemoteServer interface {
//	Init(conf *config.RPCServer, cliFactory IRpcSenderFactory) IRemoteServer
//	Serve() error
//	Close()
//	GetAddress() string
//}

type IRemoteServer interface {
	Init(IRpcSenderFactory)
	Serve(*config.RPCServer, string) error
	Close()
}
