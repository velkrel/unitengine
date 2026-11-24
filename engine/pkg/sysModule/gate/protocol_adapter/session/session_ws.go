// Package sess
// 模块名: 模块名
// 功能描述: 描述
package session

import (
	inf "github.com/velkrel/unitengine/engine/pkg/interfaces"
	"github.com/velkrel/unitengine/engine/pkg/log"
	"github.com/velkrel/unitengine/engine/pkg/utils/mpsc"
)

type WSSession struct {
	BaseSession
}

func NewWSSession(id uint64, conn inf.IConn, uid int64, logger log.ILoggerX) *WSSession {
	return &WSSession{
		BaseSession: BaseSession{
			id:           id,
			conn:         conn,
			logger:       logger,
			uid:          uid,
			msgCh:        mpsc.New[[]byte](),
			sendStrategy: 0,
			batchSize:    0,
		},
	}
}
