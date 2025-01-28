package grpcapp

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/1abobik1/Cloud-Storage/file_upload_service/internal/middleware"
	"github.com/1abobik1/Cloud-Storage/file_upload_service/internal/server"
	"google.golang.org/grpc"
)

type App struct {
	log           *slog.Logger
	gRPCServer     *grpc.Server
	port          int
	publicKeyPath string
}

func New(log *slog.Logger, fileUploaderService server.FileUploaderServiceI, port int, publicKeyPath string) *App {

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.TokenInterceptor(publicKeyPath)),
	)

	server.RegisterFileUploaderServ(grpcServer, fileUploaderService)

	return &App{
		log:           log,
		gRPCServer:     grpcServer,
		port:          port,
		publicKeyPath: publicKeyPath,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "location: internal.app.grpc.app.Run()"

	log := a.log.With(
		slog.String("operation", op),
		slog.Int("port", a.port),
	)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%v, %s", err, op)
	}

	log.Info("gRPC server is running", slog.String("addr", lis.Addr().String()))

	if err := a.gRPCServer.Serve(lis); err != nil {
		return fmt.Errorf("%v, %s", err, op)
	}

	return nil
}

func (a *App) Stop() {
	const op = "location: internal.app.grpc.app.Stop()"

	
	a.log.With(slog.String("operation", op), slog.Int("port", a.port)).Info("stopping gRPC server")

	a.gRPCServer.GracefulStop()
}