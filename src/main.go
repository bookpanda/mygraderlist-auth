package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	ar "github.com/bookpanda/mygraderlist-auth/src/app/repository/auth"
	"github.com/bookpanda/mygraderlist-auth/src/app/repository/cache"
	as "github.com/bookpanda/mygraderlist-auth/src/app/service/auth"
	js "github.com/bookpanda/mygraderlist-auth/src/app/service/jwt"
	ts "github.com/bookpanda/mygraderlist-auth/src/app/service/token"
	"github.com/bookpanda/mygraderlist-auth/src/app/service/user"
	jsg "github.com/bookpanda/mygraderlist-auth/src/app/strategy"
	"github.com/bookpanda/mygraderlist-auth/src/config"
	"github.com/bookpanda/mygraderlist-auth/src/database"
	auth_proto "github.com/bookpanda/mygraderlist-proto/MyGraderList/auth"
	user_proto "github.com/bookpanda/mygraderlist-proto/MyGraderList/backend/user"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type operation func(ctx context.Context) error

func gracefulShutdown(ctx context.Context, timeout time.Duration, ops map[string]operation) <-chan struct{} {
	wait := make(chan struct{})
	go func() {
		s := make(chan os.Signal, 1)

		signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		sig := <-s

		log.Info().
			Str("service", "graceful shutdown").
			Msgf("got signal \"%v\" shutting down service", sig)

		timeoutFunc := time.AfterFunc(timeout, func() {
			log.Error().
				Str("service", "graceful shutdown").
				Msgf("timeout %v ms has been elapsed, force exit", timeout.Milliseconds())
			os.Exit(0)
		})

		defer timeoutFunc.Stop()

		var wg sync.WaitGroup

		for key, op := range ops {
			wg.Add(1)
			innerOp := op
			innerKey := key
			go func() {
				defer wg.Done()

				log.Info().
					Str("service", "graceful shutdown").
					Msgf("cleaning up: %v", innerKey)
				if err := innerOp(ctx); err != nil {
					log.Error().
						Str("service", "graceful shutdown").
						Err(err).
						Msgf("%v: clean up failed: %v", innerKey, err.Error())
					return
				}

				log.Info().
					Str("service", "graceful shutdown").
					Msgf("%v was shutdown gracefully", innerKey)
			}()
		}

		wg.Wait()
		close(wait)
	}()

	return wait
}

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal().
			Err(err).
			Str("service", "auth").
			Msg("Failed to start service")
	}

	db, err := database.InitDatabase(&conf.Database)
	if err != nil {
		log.Fatal().
			Err(err).
			Str("service", "auth").
			Msg("Failed to start service")
	}

	cacheDB, err := database.InitRedisConnect(&conf.Redis)
	if err != nil {
		log.Fatal().
			Err(err).
			Str("service", "auth").
			Msg("Failed to start service")
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", conf.App.Port))
	if err != nil {
		log.Fatal().
			Err(err).
			Str("service", "auth").
			Msg("Failed to start service")
	}

	backendConn, err := grpc.Dial(conf.Service.Backend, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal().
			Err(err).
			Str("service", "mgl-backend").
			Msg("Cannot connect to service")
	}

	grpcServer := grpc.NewServer()

	cacheRepo := cache.NewRepository(cacheDB)

	usrClient := user_proto.NewUserServiceClient(backendConn)
	usrSrv := user.NewUserService(usrClient)

	stg := jsg.NewJwtStrategy(conf.Jwt.Secret)
	jtSrv := js.NewJwtService(conf.Jwt, stg)

	tkSrv := ts.NewTokenService(jtSrv, cacheRepo)

	aRepo := ar.NewRepository(db)
	aSrv := as.NewService(aRepo, tkSrv, usrSrv, conf.App)

	grpc_health_v1.RegisterHealthServer(grpcServer, health.NewServer())
	auth_proto.RegisterAuthServiceServer(grpcServer, aSrv)

	reflection.Register(grpcServer)
	go func() {
		log.Info().
			Str("service", "auth").
			Msgf("MyGraderList auth starting at port %v", conf.App.Port)

		if err = grpcServer.Serve(lis); err != nil {
			log.Fatal().
				Err(err).
				Str("service", "auth").
				Msg("Failed to start service")
		}
	}()

	wait := gracefulShutdown(context.Background(), 2*time.Second, map[string]operation{
		"database": func(ctx context.Context) error {
			sqlDb, err := db.DB()
			if err != nil {
				return err
			}
			return sqlDb.Close()
		},
		"server": func(ctx context.Context) error {
			grpcServer.GracefulStop()
			return nil
		},
		"cache": func(ctx context.Context) error {
			return cacheDB.Close()
		},
	})

	<-wait

	grpcServer.GracefulStop()
	log.Info().
		Str("service", "auth").
		Msg("Closing the listener")
	lis.Close()
	log.Info().
		Str("service", "auth").
		Msg("End of Program")
}
