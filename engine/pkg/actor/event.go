// Package actor
// 提供 actor 协议消息的运行时辅助方法。
package actor

import (
	"github.com/velkrel/unitengine/engine/pkg/def"
)

func (e *Event) GetEventType() def.EventType {
	return def.EventType(e.Type)
}
