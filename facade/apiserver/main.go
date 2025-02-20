package main

import (
	"context"
	"flag"
	api "github.com/CloudNativeGame/kruise-game-api/facade/apiserver/grpc"
	grpcapi "github.com/CloudNativeGame/kruise-game-api/facade/apiserver/proto"
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
	"log/slog"
	"net"
	"os"
)

var opts = &slog.HandlerOptions{AddSource: false, Level: slog.LevelInfo}
var logger = slog.New(slog.NewJSONHandler(os.Stdout, opts))
var isPrintBody *bool

func init() {
	slog.SetDefault(logger)
}

func main() {
	// 是否打印请求和响应数据
	isPrintBody = flag.Bool("print-body", false, "Enable HTTP/gRPC body logging")
	flag.Parse()

	r := gin.New()

	if *isPrintBody {
		r.Use(sloggin.NewWithConfig(logger, sloggin.Config{
			WithRequestBody:  true,
			WithResponseBody: true,
			Filters:          []sloggin.Filter{sloggin.IgnorePath("/healthz")},
		}))
	} else {
		r.Use(sloggin.NewWithConfig(logger, sloggin.Config{
			WithRequestBody:  false,
			WithResponseBody: false,
			Filters:          []sloggin.Filter{sloggin.IgnorePath("/healthz")},
		}))
	}

	r.Use(gin.Recovery())
	registerRoutes(r)

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(serverLoggingInterceptor))
	grpcapi.RegisterGameServerServiceServer(grpcServer, &api.GameServerGrpcService{})
	grpcapi.RegisterGameServerSetServiceServer(grpcServer, &api.GameServerSetGrpcService{})

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		slog.Error("failed to listen", "err", err.Error())
		os.Exit(1)
	}

	mux := cmux.New(lis)
	grpcLis := mux.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
	httpLis := mux.Match(cmux.Any())

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

func serverLoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if *isPrintBody {
		slog.Info("Received gRPC request", "method", info.FullMethod, "request", req)
	} else {
		slog.Info("Received gRPC request", "method", info.FullMethod)
	}

	resp, err := handler(ctx, req)
	if err != nil {
		slog.Error("gRPC request failed", "method", info.FullMethod, "error", err)
	} else {
		if *isPrintBody {
			slog.Info("gRPC request succeeded", "method", info.FullMethod, "response", resp)
		} else {
			slog.Info("gRPC request succeeded", "method", info.FullMethod)
		}
	}
	return resp, err
}
