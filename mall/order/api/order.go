package main

import (
	"flag"
	"fmt"

	"go-zero-demo/mall/order/api/internal/config"
	"go-zero-demo/mall/order/api/internal/handler"
	"go-zero-demo/mall/order/api/internal/svc"

	"github.com/gogf/gf/frame/g"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/zero-contrib/zrpc/registry/nacos"
)

var configFile = flag.String("f", "etc/order.yaml", "the config file")
var gConfigFile = flag.String("gFile", "config.yaml", "the config file")
var gConfigDir = flag.String("gDir", "etc", "the config directory")

func main() {
	flag.Parse()
	g.Cfg().SetPath(*gConfigDir)
	g.Cfg().SetFileName(*gConfigFile)
	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)

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
	opts := nacos.NewNacosConfig(c.Name, fmt.Sprintf("%s:%d", c.Host, c.Port), sc, cc)
	_ = nacos.RegisterService(opts)

	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
