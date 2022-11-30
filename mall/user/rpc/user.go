package main

import (
	"flag"
	"fmt"

	"go-zero-demo/mall/user/rpc/internal/config"
	"go-zero-demo/mall/user/rpc/internal/server"
	"go-zero-demo/mall/user/rpc/internal/svc"
	"go-zero-demo/mall/user/rpc/types/user"

	"github.com/gogf/gf/frame/g"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/zeromicro/zero-contrib/zrpc/registry/nacos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/user.yaml", "the config file")
var gConfigFile = flag.String("gFile", "config.yaml", "the config file")
var gConfigDir = flag.String("gDir", "etc", "the config directory")

func main() {
	flag.Parse()
	g.Cfg().SetPath(*gConfigDir)
	g.Cfg().SetFileName(*gConfigFile)
	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)
	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		user.RegisterUserServer(grpcServer, server.NewUserServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})

	// 注册服务
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(g.Cfg().GetString("nacos_server.host", "127.0.0.1"), g.Cfg().GetUint64("nacos_server.port", 8848)),
	}

	cc := &constant.ClientConfig{
		NamespaceId:         "public",
		TimeoutMs:           50000,
		NotLoadCacheAtStart: true,
		LogDir:              fmt.Sprintf("/tmp/nacos/%s/log", c.Name),
		CacheDir:            fmt.Sprintf("/tmp/nacos/%s/cache", c.Name),
		LogLevel:            "debug",
		AppName:             c.Name,
	}
	opts := nacos.NewNacosConfig(c.Name, c.ListenOn, sc, cc)
	_ = nacos.RegisterService(opts)

	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
