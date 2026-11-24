// Package gr
// @Title  grpc服务端监听器
// @Description  desc
package gr

import (
	"context"

	"github.com/velkrel/unitengine/engine/pkg/actor"
	inf "github.com/velkrel/unitengine/engine/pkg/interfaces"
	"github.com/velkrel/unitengine/engine/pkg/rpc/remote/handler"
)

type GrpcListener struct {
	actor.UnimplementedGrpcListenerServer
	cliFactory inf.IRpcSenderFactory
	handler    *handler.Handler
}

func (g *GrpcListener) RPCCall(_ context.Context, req *actor.Message) (*actor.RpcCallResponse, error) {
	if g.handler == nil {
		return &actor.RpcCallResponse{}, nil
	}
	return &actor.RpcCallResponse{}, g.handler.RpcMessageHandler(g.cliFactory, req)
}
