// Package router
// 模块名: 模块名
// 功能描述: 用于玩家路由查找
package router

import (
	"github.com/velkrel/unitengine/engine/pkg/core"
	"github.com/patrickmn/go-cache"
)

type Router struct {
	core.Module

	routerMap *cache.Cache // 玩家路由
}

func (r *Router) OnInit() error {
	r.routerMap = cache.New(cache.NoExpiration, cache.NoExpiration)
	return nil
}

func (r *Router) OnRelease() error {
	// 移除所有路由
	r.routerMap.Flush()
	return nil
}
