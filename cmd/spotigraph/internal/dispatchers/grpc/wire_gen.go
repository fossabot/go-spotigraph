// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package grpc

import (
	"context"
	"crypto/tls"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats/view"
	"go.uber.org/zap"
	"go.zenithar.org/pkg/db/adapter/mongodb"
	"go.zenithar.org/pkg/db/adapter/postgresql"
	"go.zenithar.org/pkg/db/adapter/rethinkdb"
	"go.zenithar.org/pkg/log"
	"go.zenithar.org/pkg/tlsconfig"
	"go.zenithar.org/spotigraph/cmd/spotigraph/internal/config"
	"go.zenithar.org/spotigraph/cmd/spotigraph/internal/core"
	mongodb2 "go.zenithar.org/spotigraph/internal/repositories/pkg/mongodb"
	postgresql2 "go.zenithar.org/spotigraph/internal/repositories/pkg/postgresql"
	rethinkdb2 "go.zenithar.org/spotigraph/internal/repositories/pkg/rethinkdb"
	"go.zenithar.org/spotigraph/internal/services"
	"go.zenithar.org/spotigraph/internal/services/pkg/chapter"
	"go.zenithar.org/spotigraph/internal/services/pkg/graph"
	"go.zenithar.org/spotigraph/internal/services/pkg/guild"
	"go.zenithar.org/spotigraph/internal/services/pkg/squad"
	"go.zenithar.org/spotigraph/internal/services/pkg/tribe"
	"go.zenithar.org/spotigraph/internal/services/pkg/user"
	"go.zenithar.org/spotigraph/pkg/grpc/v1/spotigraph/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

// Injectors from wire.go:

func setupLocalMongoDB(ctx context.Context, cfg *config.Configuration) (*grpc.Server, error) {
	configuration := core.MongoDBConfig(cfg)
	wrappedClient, err := mongodb.Connection(ctx, configuration)
	if err != nil {
		return nil, err
	}
	repositoriesUser := mongodb2.NewUserRepository(configuration, wrappedClient)
	servicesUser := user.New(repositoriesUser)
	repositoriesChapter := mongodb2.NewChapterRepository(configuration, wrappedClient)
	servicesChapter := chapter.New(repositoriesChapter)
	repositoriesGuild := mongodb2.NewGuildRepository(configuration, wrappedClient)
	servicesGuild := guild.New(repositoriesGuild)
	repositoriesSquad := mongodb2.NewSquadRepository(configuration, wrappedClient)
	servicesSquad := squad.New(repositoriesSquad)
	repositoriesTribe := mongodb2.NewTribeRepository(configuration, wrappedClient)
	servicesTribe := tribe.New(repositoriesTribe)
	servicesGraph := graph.New(repositoriesUser, repositoriesSquad, repositoriesChapter, repositoriesGuild, repositoriesTribe)
	server, err := grpcServer(ctx, cfg, servicesUser, servicesChapter, servicesGuild, servicesSquad, servicesTribe, servicesGraph)
	if err != nil {
		return nil, err
	}
	return server, nil
}

func setupLocalRethinkDB(ctx context.Context, cfg *config.Configuration) (*grpc.Server, error) {
	configuration := core.RethinkDBConfig(cfg)
	session, err := rethinkdb.Connection(ctx, configuration)
	if err != nil {
		return nil, err
	}
	repositoriesUser := rethinkdb2.NewUserRepository(configuration, session)
	servicesUser := user.New(repositoriesUser)
	repositoriesChapter := rethinkdb2.NewChapterRepository(configuration, session)
	servicesChapter := chapter.New(repositoriesChapter)
	repositoriesGuild := rethinkdb2.NewGuildRepository(configuration, session)
	servicesGuild := guild.New(repositoriesGuild)
	repositoriesSquad := rethinkdb2.NewSquadRepository(configuration, session)
	servicesSquad := squad.New(repositoriesSquad)
	repositoriesTribe := rethinkdb2.NewTribeRepository(configuration, session)
	servicesTribe := tribe.New(repositoriesTribe)
	servicesGraph := graph.New(repositoriesUser, repositoriesSquad, repositoriesChapter, repositoriesGuild, repositoriesTribe)
	server, err := grpcServer(ctx, cfg, servicesUser, servicesChapter, servicesGuild, servicesSquad, servicesTribe, servicesGraph)
	if err != nil {
		return nil, err
	}
	return server, nil
}

func setupLocalPostgreSQL(ctx context.Context, cfg *config.Configuration) (*grpc.Server, error) {
	configuration := core.PosgreSQLConfig(cfg)
	db, err := postgresql.Connection(ctx, configuration)
	if err != nil {
		return nil, err
	}
	repositoriesUser := postgresql2.NewUserRepository(configuration, db)
	servicesUser := user.New(repositoriesUser)
	repositoriesChapter := postgresql2.NewChapterRepository(configuration, db)
	servicesChapter := chapter.New(repositoriesChapter)
	repositoriesGuild := postgresql2.NewGuildRepository(configuration, db)
	servicesGuild := guild.New(repositoriesGuild)
	repositoriesSquad := postgresql2.NewSquadRepository(configuration, db)
	servicesSquad := squad.New(repositoriesSquad)
	repositoriesTribe := postgresql2.NewTribeRepository(configuration, db)
	servicesTribe := tribe.New(repositoriesTribe)
	servicesGraph := graph.New(repositoriesUser, repositoriesSquad, repositoriesChapter, repositoriesGuild, repositoriesTribe)
	server, err := grpcServer(ctx, cfg, servicesUser, servicesChapter, servicesGuild, servicesSquad, servicesTribe, servicesGraph)
	if err != nil {
		return nil, err
	}
	return server, nil
}

// wire.go:

func grpcServer(ctx context.Context, cfg *config.Configuration, users services.User, chapters services.Chapter, guilds services.Guild, squads services.Squad, tribes services.Tribe, graph2 services.Graph) (*grpc.Server, error) {
	sopts := []grpc.ServerOption{}
	grpc_zap.ReplaceGrpcLogger(zap.L())

	sopts = append(sopts, grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(grpc_zap.StreamServerInterceptor(zap.L()), grpc_recovery.StreamServerInterceptor())), grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(grpc_recovery.UnaryServerInterceptor(), grpc_zap.UnaryServerInterceptor(zap.L()))), grpc.StatsHandler(&ocgrpc.ServerHandler{}),
	)

	if cfg.Server.GRPC.UseTLS {

		clientAuth := tls.VerifyClientCertIfGiven
		if cfg.Server.GRPC.TLS.ClientAuthenticationRequired {
			clientAuth = tls.RequireAndVerifyClientCert
		}

		tlsConfig, err := tlsconfig.Server(tlsconfig.Options{
			KeyFile:    cfg.Server.GRPC.TLS.PrivateKeyPath,
			CertFile:   cfg.Server.GRPC.TLS.CertificatePath,
			CAFile:     cfg.Server.GRPC.TLS.CACertificatePath,
			ClientAuth: clientAuth,
		})
		if err != nil {
			log.For(ctx).Error("Unable to build TLS configuration from settings", zap.Error(err))
			return nil, err
		}

		sopts = append(sopts, grpc.Creds(credentials.NewTLS(tlsConfig)))
	} else {
		log.For(ctx).Info("No transport authentication enabled")
	}

	server := grpc.NewServer(sopts...)

	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(server, healthServer)
	pb.RegisterUserServer(server, users)
	pb.RegisterChapterServer(server, chapters)
	pb.RegisterGuildServer(server, guilds)
	pb.RegisterSquadServer(server, squads)
	pb.RegisterTribeServer(server, tribes)
	pb.RegisterGraphServer(server, graph2)
	reflection.Register(server)

	err := view.Register(ochttp.ServerRequestCountView, ochttp.ServerRequestBytesView, ochttp.ServerResponseBytesView, ochttp.ServerLatencyView, ochttp.ServerRequestCountByMethod, ochttp.ServerResponseCountByStatusCode)
	if err != nil {
		log.For(ctx).Fatal("Unable to register HTTP stat views", zap.Error(err))
	}

	err = view.Register(ocgrpc.DefaultServerViews...)
	if err != nil {
		log.For(ctx).Fatal("Unable to register gRPC stat views", zap.Error(err))
	}

	return server, nil
}
