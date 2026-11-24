// Package msgenvelope
// 模块名: 模块名
// 功能描述: 描述
package msgenvelope

import (
	"sync"

	"github.com/velkrel/unitengine/engine/pkg/actor"
	"github.com/velkrel/unitengine/engine/pkg/utils/pool"
)

// TODO 需要做成配置
var msgPool pool.IPool[*actor.Message]
var msgInitOnce sync.Once

func getMsgPool() pool.IPool[*actor.Message] {
	msgInitOnce.Do(func() {
		// This process-wide pool outlives changes to GOMAXPROCS.
		msgPool = pool.NewSyncPoolWrapper(
			func() *actor.Message {
				return &actor.Message{}
			},
			func() pool.IStatsRecorder {
				if isDebug() {
					return pool.NewStatsRecorder("rpcMsgPool-syncPool")
				} else {
					return pool.NewNoStatsRecorder()
				}
			}(),
			pool.WithReset(func(msg *actor.Message) {
				msg.Reset()
			}),
		)
	})
	return msgPool
}

func NewMessage() *actor.Message {
	return getMsgPool().Get()
}

func ReleaseMessage(msg *actor.Message) {
	getMsgPool().Put(msg)
}
