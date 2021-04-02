package server

import (
	"context"
	"fmt"
	"github.com/arkdbp/jobservice/api"
	"google.golang.org/grpc/credentials"
	"log"
	"strconv"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// App - Poll API main app
type App struct {
	config        *Config
	serviceServer api.JobServiceServer
	logger        *logrus.Logger
}

// NewServer constructor for App
func NewServer(config *Config, service api.JobServiceServer, logger *logrus.Logger) *App {
	return &App{config: config, serviceServer: service, logger: logger}
}

// Run method to start the secure server
func (app *App) Run(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	tlsConfig := app.config.TLSConfig()

	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(app.ensureValidToken),
		grpc.Creds(credentials.NewTLS(tlsConfig)),
	}
	// Create the gRPC server
	gRPCServer := grpc.NewServer(opts...)

	// Register the handler object
	api.RegisterJobServiceServer(gRPCServer, app.serviceServer)

	// Create a channel to listen on
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", app.config.Port))
	if err != nil {
		app.logger.Error("failed to listen on port:", app.config.Port, "Error: ", err)
		panic(err)
	}

	fmt.Println("server will be running on ", app.config.ResolveListenAddress())
	if err := gRPCServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// ResolveListenAddress resolves the hostname:port combination for configs
func (c *Config) ResolveListenAddress() string {
	if host := c.ServerName; host != "" {
		return host + ":" + strconv.Itoa(int(c.Port)) // configs
	}
	return ":" + strconv.Itoa(int(defaultPort)) // default
}

func (app *App) ensureValidToken(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	return handler(ctx, req)
}

func (app *App) parseToken(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	token := md["authorization"]
	app.logger.Tracef("received authorization header from context: %+v", token)
	if len(token) < 1 {
		return ""
	}
	return token[0]
}

func (app *App) shutdownServer(ctx context.Context, srv *http.Server) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C
			_ = srv.Shutdown(ctx)
		}

		_, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		_ = srv.Shutdown(ctx)
	}()
}
