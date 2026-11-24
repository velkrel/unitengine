// Package interfaces
// @Title  title
// @Description  desc
package interfaces

import (
	"github.com/velkrel/unitengine/engine/pkg/config"
)

type IDiscovery interface {
	Init(conf *config.ClusterConf, eventProcessor IEventProcessor, evtCh IEventChannel) error
	Start()
	Close()
}
