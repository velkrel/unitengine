// Package nt
// @Title  title
// @Description  desc
package nt

import (
	"github.com/nats-io/nats.go"
	"github.com/velkrel/unitengine/engine/pkg/def"
	inf "github.com/velkrel/unitengine/engine/pkg/interfaces"
	"github.com/velkrel/unitengine/engine/pkg/log"
	"github.com/velkrel/unitengine/engine/pkg/rpc/message/msgenvelope"
	"github.com/velkrel/unitengine/engine/pkg/rpc/remote/handler"
	"github.com/velkrel/unitengine/engine/pkg/utils/codec"
)

type NatsListener struct {
	cliFactory inf.IRpcSenderFactory
	logger     log.ILoggerX
	handler    *handler.Handler
}

func (n *NatsListener) Handle(msg *nats.Msg) {
	req := msgenvelope.NewMessage()
	defer msgenvelope.ReleaseMessage(req)
	err := codec.Decode(def.ProtoBuf, msg.Data, req)
	if err != nil {
		if n.logger != nil {
			n.logger.Errorf("unmarshal nats message error: %v", err)
		}
		return
	}

	if n.handler == nil {
		return
	}
	if err = n.handler.RpcMessageHandler(n.cliFactory, req); err != nil {
		if n.logger != nil {
			n.logger.Errorf("handle nats message error: %v  req:%+v", err, req)
		}
	}
}
