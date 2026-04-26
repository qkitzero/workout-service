package main

import (
	"context"
	"log"
	"net"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	authv1 "github.com/qkitzero/auth-service/gen/go/auth/v1"
	exercisev1 "github.com/qkitzero/workout-service/gen/go/exercise/v1"
	musclev1 "github.com/qkitzero/workout-service/gen/go/muscle/v1"
	setv1 "github.com/qkitzero/workout-service/gen/go/set/v1"
	appexercise "github.com/qkitzero/workout-service/internal/application/exercise"
	appmuscle "github.com/qkitzero/workout-service/internal/application/muscle"
	appset "github.com/qkitzero/workout-service/internal/application/set"
	apiauth "github.com/qkitzero/workout-service/internal/infrastructure/api/auth"
	"github.com/qkitzero/workout-service/internal/infrastructure/db"
	infraexercise "github.com/qkitzero/workout-service/internal/infrastructure/exercise"
	inframuscle "github.com/qkitzero/workout-service/internal/infrastructure/muscle"
	infraset "github.com/qkitzero/workout-service/internal/infrastructure/set"
	grpcexercise "github.com/qkitzero/workout-service/internal/interface/grpc/exercise"
	grpcmuscle "github.com/qkitzero/workout-service/internal/interface/grpc/muscle"
	grpcset "github.com/qkitzero/workout-service/internal/interface/grpc/set"
	"github.com/qkitzero/workout-service/util"
)

const shutdownTimeout = 15 * time.Second

func main() {
	gormDB, err := db.Init(
		util.GetEnv("DB_HOST", ""),
		util.GetEnv("DB_USER", ""),
		util.GetEnv("DB_PASSWORD", ""),
		util.GetEnv("DB_NAME", ""),
		util.GetEnv("DB_PORT", ""),
		util.GetEnv("DB_SSL_MODE", ""),
	)
	if err != nil {
		log.Fatal(err)
	}
	sqlDB, err := gormDB.DB()
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	listener, err := net.Listen("tcp", ":"+util.GetEnv("PORT", ""))
	if err != nil {
		log.Fatal(err)
	}

	authTarget := util.GetEnv("AUTH_SERVICE_HOST", "") + ":" + util.GetEnv("AUTH_SERVICE_PORT", "")

	var opts grpc.DialOption
	switch util.GetEnv("ENV", "development") {
	case "production":
		opts = grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, ""))
	default:
		opts = grpc.WithTransportCredentials(insecure.NewCredentials())
	}

	conn, err := grpc.NewClient(authTarget, opts)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	server := grpc.NewServer()

	authServiceClient := authv1.NewAuthServiceClient(conn)
	setRepository := infraset.NewSetRepository(gormDB)
	exerciseRepository := infraexercise.NewExerciseRepository(gormDB)
	muscleRepository := inframuscle.NewMuscleRepository(gormDB)

	authService := apiauth.NewAuthService(authServiceClient)
	setUsecase := appset.NewSetUsecase(authService, setRepository, exerciseRepository)
	exerciseUsecase := appexercise.NewExerciseUsecase(exerciseRepository)
	muscleUsecase := appmuscle.NewMuscleUsecase(muscleRepository)

	healthServer := health.NewServer()
	setHandler := grpcset.NewSetHandler(setUsecase)
	exerciseHandler := grpcexercise.NewExerciseHandler(exerciseUsecase)
	muscleHandler := grpcmuscle.NewMuscleHandler(muscleUsecase)

	grpc_health_v1.RegisterHealthServer(server, healthServer)
	setv1.RegisterSetServiceServer(server, setHandler)
	exercisev1.RegisterExerciseServiceServer(server, exerciseHandler)
	musclev1.RegisterMuscleServiceServer(server, muscleHandler)

	healthServer.SetServingStatus("set", grpc_health_v1.HealthCheckResponse_SERVING)
	healthServer.SetServingStatus("exercise", grpc_health_v1.HealthCheckResponse_SERVING)
	healthServer.SetServingStatus("muscle", grpc_health_v1.HealthCheckResponse_SERVING)

	if util.GetEnv("ENV", "development") == "development" {
		reflection.Register(server)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	serveErr := make(chan error, 1)
	go func() {
		log.Printf("gRPC server listening on %s", listener.Addr().String())
		serveErr <- server.Serve(listener)
	}()

	select {
	case err := <-serveErr:
		if err != nil {
			log.Fatalf("gRPC server failed: %v", err)
		}
	case <-ctx.Done():
		log.Println("shutdown signal received, starting graceful stop")
		healthServer.Shutdown()

		stopped := make(chan struct{})
		go func() {
			server.GracefulStop()
			close(stopped)
		}()

		select {
		case <-stopped:
			log.Println("gRPC server stopped gracefully")
		case <-time.After(shutdownTimeout):
			log.Printf("graceful stop timed out after %s, forcing stop", shutdownTimeout)
			server.Stop()
			<-stopped
		}
	}
}
