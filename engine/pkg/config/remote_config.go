package config

import (
	"bytes"
	"fmt"
	"log/slog"
	"path"
	"time"

	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
	"github.com/velkrel/unitengine/engine/pkg/config/remote"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type remoteConfig struct {
	client   *clientv3.Client
	key      string
	revision int64
	done     chan struct{}
}

func (c *Config) parseStartService() error {
	if c.ServiceConf == nil || !c.ServiceConf.OpenRemote {
		return nil
	}
	if c.ClusterConf == nil || c.ClusterConf.ETCDConf == nil {
		return fmt.Errorf("config validation: ETCDConf is required for remote configuration")
	}
	settings := c.ClusterConf.ETCDConf
	client, err := (remote.Config{
		Endpoints: settings.Endpoints,
		Username:  settings.UserName,
		Password:  settings.Password,
	}).NewClient()
	if err != nil {
		return fmt.Errorf("create remote config client: %w", err)
	}
	rc := &remoteConfig{client: client, key: path.Join(c.ServiceConf.RemoteConfPath, startServiceConfName)}
	c.remoteConfig = rc
	data, revision, err := remote.Read(client.Ctx(), client, rc.key)
	if err != nil {
		return fmt.Errorf("read remote start service config: %w", err)
	}
	rc.revision = revision
	initial, err := decodeRemoteServiceConf(data, c.ServiceConf)
	if err != nil {
		return err
	}
	c.ServiceConf = initial
	return nil
}

func decodeRemoteServiceConf(data []byte, current *ServiceConf) (*ServiceConf, error) {
	parser := viper.New()
	parser.SetConfigType("yaml")
	if err := parser.ReadConfig(bytes.NewReader(data)); err != nil {
		return nil, fmt.Errorf("parse remote start service config: %w", err)
	}
	next := *current
	// 未传入的字段沿用上一版，传入的容器重新分配，避免修改已发布的配置。
	if err := parser.Unmarshal(&next, func(dc *mapstructure.DecoderConfig) { dc.ZeroFields = true }); err != nil {
		return nil, fmt.Errorf("unmarshal start service config: %w", err)
	}
	return &next, nil
}

// CurrentServiceConf 返回当前服务配置快照；调用方只读，不修改返回值。
func (c *Config) CurrentServiceConf() *ServiceConf {
	if snapshot := c.serviceConfSnapshot.Load(); snapshot != nil {
		return snapshot
	}
	return c.ServiceConf
}

// Close 停止远程监听并释放连接，保留最后一份有效配置。
func (c *Config) Close() {
	if rc := c.remoteConfig; rc != nil {
		_ = rc.client.Close()
		if rc.done != nil {
			<-rc.done
		}
		c.remoteConfig = nil
	}
}

func (c *Config) watchRemoteConfig(rc *remoteConfig) {
	defer close(rc.done)
	ctx := rc.client.Ctx()
	revision := rc.revision + 1
	for ctx.Err() == nil {
		if revision == 0 {
			// 监听被服务端终止（如历史版本已压缩）时，重新读取后再接续。
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Second):
			}
			data, rev, err := remote.Read(ctx, rc.client, rc.key)
			if err != nil {
				continue
			}
			c.publishRemoteServiceConf(data)
			revision = rev + 1
		}
		for resp := range rc.client.Watch(ctx, rc.key, clientv3.WithRev(revision), clientv3.WithFilterDelete()) {
			if ctx.Err() != nil {
				return
			}
			if err := resp.Err(); err != nil {
				slog.Warn("remote service config watch stopped", "error", err)
				break
			}
			for _, event := range resp.Events {
				c.publishRemoteServiceConf(event.Kv.Value)
			}
		}
		revision = 0
	}
}

func (c *Config) publishRemoteServiceConf(data []byte) {
	next, err := decodeRemoteServiceConf(data, c.CurrentServiceConf())
	if err != nil {
		slog.Warn("remote service config update ignored", "error", err)
		return
	}
	c.serviceConfSnapshot.Store(next)
}
