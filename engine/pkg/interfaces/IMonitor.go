// Package interfaces
// @Title 调用监控
// @Description  用于监控一些具有超时属性的调用
package interfaces

type IMonitor interface {
	Init() IMonitor
	Start()
	Stop()
	Add(call IEnvelope)
	Remove(seq uint64) IEnvelope
	Get(seq uint64) IEnvelope
	GenSeq() uint64
}
