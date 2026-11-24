package httpx

import (
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"sync"
)

// StartBackground 在 Init 后校验证书并绑定端口，然后后台处理请求；由 Stop 关闭。
func (gs *GinServer) StartBackground() error {
	serve := gs.server.Serve
	if ca := gs.conf.CAFile; ca != nil {
		cert, err := tls.LoadX509KeyPair(ca.CertFile, ca.KeyFile)
		if err != nil {
			return err
		}
		gs.server.TLSConfig = &tls.Config{Certificates: []tls.Certificate{cert}, MinVersion: tls.VersionTLS12}
		serve = func(ln net.Listener) error {
			defer func() { _ = ln.Close() }() // ServeTLS 可能在接管监听器前失败。
			return gs.server.ServeTLS(ln, "", "")
		}
	}
	ln, err := net.Listen("tcp", gs.server.Addr)
	if err != nil {
		return err
	}
	var hooks sync.WaitGroup
	hooks.Add(len(gs.runHook))
	for _, hook := range gs.runHook {
		go func() {
			defer hooks.Done()
			hook()
		}()
	}
	hooks.Wait()

	gs.wg.Add(1)
	go func() {
		defer gs.wg.Done()
		if err := serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
			gs.logger.Error(err)
		}
	}()
	return nil
}
