package main

import (
	api "github.com/CloudNativeGame/kruise-game-api/facade/apiserver/grpc"
	grpcapi "github.com/CloudNativeGame/kruise-game-api/facade/apiserver/proto"
	"github.com/gin-gonic/gin"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
	"log/slog"
	"net"
	"os"
)

var opts = &slog.HandlerOptions{AddSource: true, Level: slog.LevelInfo}
var logger = slog.New(slog.NewJSONHandler(os.Stdout, opts))

func init() {
	slog.SetDefault(logger)
}

func main() {
	r := gin.New()

	r.Use(gin.Recovery())
	registerRoutes(r)

	grpcServer := grpc.NewServer()
	grpcapi.RegisterGameServerServiceServer(grpcServer, &api.GameServerGrpcService{})
	grpcapi.RegisterGameServerSetServiceServer(grpcServer, &api.GameServerSetGrpcService{})

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		slog.Error("failed to listen", "err", err.Error())
		os.Exit(1)
	}

	mux := cmux.New(lis)
	grpcLis := mux.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
	httpLis := mux.Match(cmux.HTTP1Fast())

	go func() {
		slog.Info("gRPC server listen at " + grpcLis.Addr().String())
		if err := grpcServer.Serve(grpcLis); err != nil {
			slog.Error("failed to serve gRPC", "err", err.Error())
		}
	}()

	go func() {
		slog.Info("HTTP server listen at " + httpLis.Addr().String())
		if err := r.RunListener(httpLis); err != nil {
			slog.Error("failed to serve HTTP", "err", err.Error())
		}
	}()

	err = mux.Serve()
	if err != nil {
		slog.Error("failed to serve mux", "err", err.Error())
	}
}
