// Package remote
// @Title  viper的远程配置
// @Description  用于读取远程配置
package remote

import (
	"context"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

type Config struct {
	Endpoints []string
	Username  string
	Password  string
}

func (c Config) NewClient() (*clientv3.Client, error) {
	return clientv3.New(clientv3.Config{
		Endpoints: c.Endpoints,
		Username:  c.Username,
		Password:  c.Password,
		Logger:    zap.NewNop(),
	})
}

func Read(ctx context.Context, client *clientv3.Client, key string) ([]byte, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	resp, err := client.Get(ctx, key)
	if err != nil {
		return nil, 0, err
	}
	if len(resp.Kvs) == 0 {
		return nil, resp.Header.Revision, nil
	}
	return resp.Kvs[0].Value, resp.Header.Revision, nil
}
