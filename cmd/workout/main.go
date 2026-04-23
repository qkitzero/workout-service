package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	authv1 "github.com/qkitzero/auth-service/gen/go/auth/v1"
	exercisev1 "github.com/qkitzero/workout-service/gen/go/exercise/v1"
	setv1 "github.com/qkitzero/workout-service/gen/go/set/v1"
	appexercise "github.com/qkitzero/workout-service/internal/application/exercise"
	appset "github.com/qkitzero/workout-service/internal/application/set"
	apiauth "github.com/qkitzero/workout-service/internal/infrastructure/api/auth"
	"github.com/qkitzero/workout-service/internal/infrastructure/db"
	infraexercise "github.com/qkitzero/workout-service/internal/infrastructure/exercise"
	infraset "github.com/qkitzero/workout-service/internal/infrastructure/set"
	grpcexercise "github.com/qkitzero/workout-service/internal/interface/grpc/exercise"
	grpcset "github.com/qkitzero/workout-service/internal/interface/grpc/set"
	"github.com/qkitzero/workout-service/util"
)

func main() {
	db, err := db.Init(
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
	setRepository := infraset.NewSetRepository(db)
	exerciseRepository := infraexercise.NewExerciseRepository(db)

	authService := apiauth.NewAuthService(authServiceClient)
	setUsecase := appset.NewSetUsecase(authService, setRepository, exerciseRepository)
	exerciseUsecase := appexercise.NewExerciseUsecase(exerciseRepository)

	healthServer := health.NewServer()
	setHandler := grpcset.NewSetHandler(setUsecase)
	exerciseHandler := grpcexercise.NewExerciseHandler(exerciseUsecase)

	grpc_health_v1.RegisterHealthServer(server, healthServer)
	setv1.RegisterSetServiceServer(server, setHandler)
	exercisev1.RegisterExerciseServiceServer(server, exerciseHandler)

	healthServer.SetServingStatus("set", grpc_health_v1.HealthCheckResponse_SERVING)
	healthServer.SetServingStatus("exercise", grpc_health_v1.HealthCheckResponse_SERVING)

	if util.GetEnv("ENV", "development") == "development" {
		reflection.Register(server)
	}

	if err = server.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
