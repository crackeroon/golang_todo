package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	core_config "github.com/crackeroon/golang_todo/internal/config"
	core_logger "github.com/crackeroon/golang_todo/internal/core/logger"
	"github.com/crackeroon/golang_todo/internal/core/repository/postgres/pool/pqx"
	core_http_middleware "github.com/crackeroon/golang_todo/internal/core/transport/http/middleware"
	core_http_server "github.com/crackeroon/golang_todo/internal/core/transport/http/server"
	tasks_postgres_repository "github.com/crackeroon/golang_todo/internal/features/tasks/repository/postgres"
	tasks_service "github.com/crackeroon/golang_todo/internal/features/tasks/service"
	tasks_transport_http "github.com/crackeroon/golang_todo/internal/features/tasks/transport/http"
	users_postgres_repository "github.com/crackeroon/golang_todo/internal/features/users/repository/postgres"
	users_service "github.com/crackeroon/golang_todo/internal/features/users/service"
	users_transport_http "github.com/crackeroon/golang_todo/internal/features/users/transport/http"
	"go.uber.org/zap"
)

func main() {
	cfg := core_config.NewConfigMust()
	time.Local = cfg.TimeZone

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	logger.Debug("application TIME ZONE", zap.Any("timeZone", cfg.TimeZone))
	if err != nil {
		fmt.Println("failed to init app logger", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("initializing postgres connection pool")

	pool, err := core_pgx_pool.NewPool(ctx, core_pgx_pool.NewConfigMust())

	if err != nil {
		logger.Fatal("failed to init postgres connection pool", zap.Error(err))
	}
	defer pool.Close()

	logger.Debug("initializing feature", zap.String("feature", "users"))

	userRepository := users_postgres_repository.NewUsersRepository(pool)
	userService := users_service.NewUsersService(userRepository)
	usersTransportHTTP := users_transport_http.NewUsersHTTPHandler(userService)

	logger.Debug("initializing feature", zap.String("feature", "tasks"))

	taskRepository := tasks_postgres_repository.NewTaskRepository(pool)
	taskService := tasks_service.NewTasksService(taskRepository)
	tasksTransportHTTP := tasks_transport_http.NewTasksHTTPHandler(taskService)

	logger.Debug("initializing HTTP server")

	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)
	apiVersionRouterV1 := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouterV1.RegisterRoutes(usersTransportHTTP.Routes()...)
	apiVersionRouterV1.RegisterRoutes(tasksTransportHTTP.Routes()...)

	/* Example of usage apiVersionRouterV2 with separate Middleware
	apiVersionRouterV2 := core_http_server.NewAPIVersionRouter(
		core_http_server.ApiVersion2,
		core_http_middleware.Dummy("ApiVersion2 Dummy Middleware"),
	)
	apiVersionRouterV2.RegisterRoutes(usersTransportHTTP.Routes()...)
	*/

	httpServer.RegisterAPIRoutes(apiVersionRouterV1 /*, apiVersionRouterV2*/)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run error", zap.Error(err))
	}
}
