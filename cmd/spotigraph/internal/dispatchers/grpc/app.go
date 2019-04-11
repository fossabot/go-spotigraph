package grpc

import (
	"context"
	"fmt"
	"net"

	"go.uber.org/zap"
	"go.zenithar.org/pkg/log"
	"google.golang.org/grpc"

	"go.zenithar.org/spotigraph/cmd/spotigraph/internal/config"
)

type application struct {
	cfg    *config.Configuration
	server *grpc.Server
}

// -----------------------------------------------------------------------------

// WaitForShutdown starts the server and wait for shutdown signal
func WaitForShutdown(ctx context.Context, cfg *config.Configuration) {
	// Initialize application
	app := &application{}

	// Apply configuration
	if err := app.ApplyConfiguration(cfg); err != nil {
		log.For(ctx).Fatal("Unable to initialize server settings", zap.Error(err))
	}

	// Fork the cancellation handler
	go func() {
		<-ctx.Done()
		err := app.Shutdown(context.Background())
		if err != nil {
			log.For(ctx).Error("Unable to shutdown spotigraph service", zap.Error(err))
		}
	}()

	// Start the server
	if err := app.Serve(ctx); err != nil {
		log.For(ctx).Fatal("Unable to start server", zap.Error(err))
	}
}

// -----------------------------------------------------------------------------

// ApplyConfiguration apply the configuration after checking it
func (s *application) ApplyConfiguration(cfg interface{}) error {
	// Check configuration validity
	if err := s.checkConfiguration(cfg); err != nil {
		return err
	}

	// Apply to current component (type assertion done if check)
	s.cfg, _ = cfg.(*config.Configuration)

	// No error
	return nil
}

// Serve starts the component
func (s *application) Serve(ctx context.Context) error {
	log.For(ctx).Info("Starting spotigraph service ...", zap.String("backend", s.cfg.Core.Mode))

	// Setup the gRPC server
	if err := s.setup(ctx); err != nil {
		log.For(ctx).Error("Unable to initialize spotigraph service", zap.Error(err))
		return err
	}

	// Initialize a listener
	lis, err := net.Listen(s.cfg.Server.GRPC.Network, s.cfg.Server.GRPC.Listen)
	if err != nil {
		return err
	}

	// Return no error
	log.For(ctx).Info("Spotigraph service listening ...", zap.String("listen", s.cfg.Server.GRPC.Listen))
	return s.server.Serve(lis)
}

// Shutdown the component
func (s *application) Shutdown(ctx context.Context) error {
	log.For(ctx).Info("Try to gracefully shutdown spotigraph ...")
	if s.server != nil {
		s.server.GracefulStop()
	}
	return nil
}

// -----------------------------------------------------------------------------

func (s *application) checkConfiguration(cfg interface{}) error {
	// Check via type assertion
	config, ok := cfg.(*config.Configuration)
	if !ok {
		return fmt.Errorf("server: invalid configuration")
	}

	switch config.Core.Mode {
	case "local":
		switch config.Core.Local.Type {
		case "rethinkdb":
			if config.Core.Local.Hosts == "" {
				return fmt.Errorf("server: database hosts list is mandatory")
			}
		case "mongodb":
			if config.Core.Local.Hosts == "" {
				return fmt.Errorf("server: database hosts list is mandatory")
			}
		case "postgresql":
			if config.Core.Local.Hosts == "" {
				return fmt.Errorf("server: database hosts list is mandatory")
			}
		default:
			return fmt.Errorf("server: invalid type (mongodb/rethinkdb/postgresql)")
		}
	default:
		return fmt.Errorf("server: invalid core mode, grpc only support 'local'")
	}

	// No error
	return nil
}

func (s *application) setup(ctx context.Context) error {
	var err error

	switch s.cfg.Core.Mode {
	case "local":
		switch s.cfg.Core.Local.Type {
		case "mongodb":
			s.server, err = setupLocalMongoDB(ctx, s.cfg)
		case "rethinkdb":
			s.server, err = setupLocalRethinkDB(ctx, s.cfg)
		case "postgresql":
			s.server, err = setupLocalPostgreSQL(ctx, s.cfg)
		default:

		}
	case "remote":
		log.For(ctx).Fatal("Invalid core mode, use 'local' only for gRPC service.")
	default:
		log.For(ctx).Fatal("Invalid core mode, use 'remote' or 'local'.")
	}

	// Return wired context
	return err
}