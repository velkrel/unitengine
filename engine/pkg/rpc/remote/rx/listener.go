// Package rx
// @Title  rpcx的服务端监听器
// @Description  desc
package rx

import (
	"context"

	"github.com/velkrel/unitengine/engine/pkg/actor"
	"github.com/velkrel/unitengine/engine/pkg/dto"
	inf "github.com/velkrel/unitengine/engine/pkg/interfaces"
	"github.com/velkrel/unitengine/engine/pkg/log"
	"github.com/velkrel/unitengine/engine/pkg/rpc/remote/handler"
)

type RpcxListener struct {
	cliFactory inf.IRpcSenderFactory
	logger     log.ILoggerX
	handler    *handler.Handler
}

func (rm *RpcxListener) RPCCall(_ context.Context, req *actor.Message, _ *dto.RPCResponse) error {
	// rm.logger.Debugf("rpcx call: %+v", req)
	if rm.handler == nil {
		return nil
	}
	return rm.handler.RpcMessageHandler(rm.cliFactory, req)
}
