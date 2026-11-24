// Package main
// @Title  title
// @Description  desc
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	inf "github.com/velkrel/unitengine/engine/pkg/interfaces"
	"github.com/velkrel/unitengine/engine/pkg/node"
	"github.com/velkrel/unitengine/engine/pkg/services"
	"github.com/velkrel/unitengine/example/comm"
)

func init() {
	services.SetService("Service3", func() inf.IService {
		return &comm.Service3{}
	})
}

func main() {
	n, err := node.New().Start(node.WithConfPath("./example/configs/node2"))
	if err != nil {
		panic(err)
	}
	exitCh := make(chan os.Signal, 1)
	signal.Notify(exitCh, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	<-exitCh
	fmt.Println("exit signal received")
	n.Stop()
}
