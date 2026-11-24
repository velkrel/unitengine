package config

import (
	"context"
	"fmt"
	"net"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.etcd.io/etcd/api/v3/etcdserverpb"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"google.golang.org/grpc"
)

const remoteTestKey = "/config/services.yaml"

type remoteTestEtcd struct {
	etcdserverpb.UnimplementedKVServer
	etcdserverpb.UnimplementedWatchServer
	watched chan *etcdserverpb.WatchCreateRequest
	updates chan string
	stopped chan struct{}
}

func (s *remoteTestEtcd) Range(_ context.Context, req *etcdserverpb.RangeRequest) (*etcdserverpb.RangeResponse, error) {
	var value string
	switch string(req.Key) {
	case remoteTestKey:
		value = "StartServices: [{ClassName: initial, Type: normal, Partition: 1}]"
	case "/config/demo.yaml":
		value = "Message: from-etcd"
	}
	resp := &etcdserverpb.RangeResponse{Header: &etcdserverpb.ResponseHeader{Revision: 40}}
	if value != "" {
		resp.Kvs = []*mvccpb.KeyValue{{Key: req.Key, Value: []byte(value), ModRevision: 40}}
	}
	return resp, nil
}

func (s *remoteTestEtcd) Watch(stream etcdserverpb.Watch_WatchServer) error {
	defer close(s.stopped)
	req, err := stream.Recv()
	if err != nil {
		return err
	}
	ctx := stream.Context()
	select {
	case s.watched <- req.GetCreateRequest():
	case <-ctx.Done():
		return ctx.Err()
	}
	if err := stream.Send(&etcdserverpb.WatchResponse{
		Header: &etcdserverpb.ResponseHeader{Revision: 40}, WatchId: 1, Created: true,
	}); err != nil {
		return err
	}
	for revision := int64(41); ; revision++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case data := <-s.updates:
			if err := stream.Send(&etcdserverpb.WatchResponse{
				Header:  &etcdserverpb.ResponseHeader{Revision: revision},
				WatchId: 1,
				Events: []*mvccpb.Event{{
					Type: mvccpb.PUT,
					Kv:   &mvccpb.KeyValue{Key: []byte(remoteTestKey), Value: []byte(data), ModRevision: revision},
				}},
			}); err != nil {
				return err
			}
		}
	}
}

func TestConfigLoadRemoteWatch(t *testing.T) {
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	etcd := &remoteTestEtcd{
		watched: make(chan *etcdserverpb.WatchCreateRequest, 1),
		updates: make(chan string, 3),
		stopped: make(chan struct{}),
	}
	server := grpc.NewServer()
	etcdserverpb.RegisterKVServer(server, etcd)
	etcdserverpb.RegisterWatchServer(server, etcd)
	go func() { _ = server.Serve(lis) }()
	t.Cleanup(server.Stop)

	t.Setenv("EMBER_CONF_PATH", "")
	tmp := filepath.ToSlash(t.TempDir())
	nodeYAML := fmt.Sprintf(`
NodeConf:
  NodeId: remote
  NodeType: test
  SystemStatus: debug
  PVCPath: %q
  PVPath: %q
  AntsPoolSize: 10
SystemLogger:
  Dir: %q
  Level: info
  OutputFormat: json
ClusterConf:
  ETCDConf:
    Endpoints: [%q]
ServiceConf:
  OpenRemote: true
  RemoteConfPath: /config
`, tmp+"/data", tmp+"/run", tmp+"/logs", lis.Addr().String())

	type demoConfig struct{ Message string }
	cfg := NewConfig()
	t.Cleanup(cfg.Close)
	cfg.RegisterServiceConf(&ServiceConfig{
		ServiceName: "demo",
		ConfName:    "demo",
		ConfType:    "yaml",
		CfgCreator:  func() interface{} { return new(demoConfig) },
	})
	invalidDir := writeYAMLTemp(t, strings.Replace(nodeYAML, "OutputFormat: json", "OutputFormat: invalid", 1))
	require.ErrorContains(t, cfg.Load(invalidDir), "OutputFormat")
	require.Nil(t, cfg.remoteConfig)
	require.Empty(t, etcd.watched)
	require.NoError(t, cfg.Load(writeYAMLTemp(t, nodeYAML)))
	initial := cfg.CurrentServiceConf()
	require.Len(t, initial.StartServices, 1)
	require.Equal(t, "initial", initial.StartServices[0].ClassName)
	require.Equal(t, "from-etcd", initial.ServicesConfMap["demo"].Cfg.(*demoConfig).Message)

	select {
	case watch := <-etcd.watched:
		require.Equal(t, remoteTestKey, string(watch.GetKey()))
		require.EqualValues(t, 41, watch.GetStartRevision())
	case <-time.After(3 * time.Second):
		t.Fatal("etcd watch did not start")
	}
	etcd.updates <- "StartServices: [{ClassName: updated, Type: normal, Partition: 2}]"
	etcd.updates <- "StartServices: [{ClassName: invalid, Partition: not-a-number}]"
	etcd.updates <- "RemoteConfPath: /updated"
	require.Eventually(t, func() bool {
		return cfg.CurrentServiceConf().RemoteConfPath == "/updated"
	}, 3*time.Second, 10*time.Millisecond)
	latest := cfg.CurrentServiceConf()
	require.Equal(t, "updated", latest.StartServices[0].ClassName)
	require.EqualValues(t, 2, latest.StartServices[0].Partition)
	require.Same(t, initial.ServicesConfMap["demo"], latest.ServicesConfMap["demo"])
	require.Equal(t, "initial", initial.StartServices[0].ClassName)
	require.Same(t, initial, cfg.ServiceConf)

	cfg.Close()
	select {
	case <-etcd.stopped:
	case <-time.After(3 * time.Second):
		t.Fatal("etcd watch did not stop")
	}
	require.Same(t, latest, cfg.CurrentServiceConf())
}
