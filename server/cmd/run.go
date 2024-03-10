package main

import (
	"context"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/eonias189/grpc-test/gen/go/proto"
	"github.com/eonias189/grpc-test/server/internal/app"
	"github.com/eonias189/grpc-test/server/internal/config"
	"google.golang.org/grpc"
)

func SetUpLogger() *slog.Logger {
	return slog.Default()
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	logger := SetUpLogger()
	reserver := app.NewReserver(logger)

	cfg, err := config.Get()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(0)
	}

	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(0)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	proto.RegisterReverserServer(grpcServer, reserver)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGSTOP,
		syscall.SIGKILL,
	)
	go func() {
		<-sigChan
		logger.Info("closing")
		cancel()
		reserver.Close()
		os.Exit(0)
	}()

	reserver.Run(ctx)
	err = grpcServer.Serve(listener)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(0)
	}

}
