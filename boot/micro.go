package boot

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gzjjyz/micro"
	"github.com/gzjjyz/micro/discovery"
	"github.com/gzjjyz/micro/factory"
	"github.com/gzjjyz/srvlib/handler/grpchandler"
	"github.com/gzjjyz/srvlib/logger"
	"github.com/gzjjyz/srvlib/pb3/health"
	"github.com/gzjjyz/srvlib/utils"
	"google.golang.org/grpc"
	"net"
	"net/http"
	_ "net/http/pprof"
)

func InitMicroSuiteWithGrpc(ctx context.Context, metaFilePath string) error {
	if err := micro.InitSuitWithGrpc(ctx, metaFilePath); err != nil {
		return err
	}
	return nil
}

func ServeGrpc(ctx context.Context, srvName string, ipVar string, port, pprofPort int, registerCustomServiceServer func(*grpc.Server)) error {
	go func() {
		err := http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", pprofPort), nil)
		if err != nil {
			logger.Errorf(err.Error())
		}
	}()

	ip, err := utils.EvalVarToParseIp(ipVar)
	if err != nil {
		logger.Errorf(err.Error())
		return err
	}

	nodeExtra, err := json.Marshal(&health.NodeHealthDesc{})
	if err != nil {
		logger.Errorf(err.Error())
		return err
	}
	node := discovery.NewNode(ip, port)
	node.Extra = string(nodeExtra)

	grpcServer := grpc.NewServer()
	registerCustomServiceServer(grpcServer)
	health.RegisterHealthServer(grpcServer, &grpchandler.Health{})

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		return err
	}

	discover, err := factory.GetOrMakeDiscovery()
	if err != nil {
		return err
	}

	defer func() {
		err = discover.Unregister(ctx, srvName, node, true)
		if err != nil {
			logger.Errorf(err.Error())
		}
	}()
	err = discover.Register(ctx, srvName, node)
	if err != nil {
		return err
	}

	err = grpcServer.Serve(listener)
	if err != nil {
		return err
	}

	return nil
}
