// Package core
// 模块名: 模块名
// 功能描述: 描述
package core

import inf "github.com/velkrel/unitengine/engine/pkg/interfaces"

func (s *Service) AddMailboxMiddlewares(middlewares ...inf.IMailboxMiddleware) {
	s.mailboxMiddlewares = append(s.mailboxMiddlewares, middlewares...)
}
