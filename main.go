package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/FallenTaters/streepjes/backend/application/auth"
	"github.com/FallenTaters/streepjes/backend/application/order"
	"github.com/FallenTaters/streepjes/backend/global/settings"
	"github.com/FallenTaters/streepjes/backend/infrastructure/repo"
	"github.com/FallenTaters/streepjes/backend/infrastructure/repo/postgres"
	"github.com/FallenTaters/streepjes/backend/infrastructure/router"
	"github.com/FallenTaters/streepjes/domain"
	"github.com/FallenTaters/streepjes/domain/authdomain"
	"github.com/FallenTaters/streepjes/static"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	if err := newRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	var configFile string
	defaults := settings.DefaultConfig()

	cmd := &cobra.Command{
		Use:          "streepjes",
		Short:        "Streepjes POS server",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfg, err := loadConfig(cmd, configFile)
			if err != nil {
				return err
			}

			return run(cfg)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&configFile, "config", "", "path to TOML config file")
	flags.Bool("disable-secure", defaults.DisableSecure, "disable TLS (serve plain HTTP)")
	flags.Int("port", defaults.Port, "server port")
	flags.String("log-level", defaults.LogLevel, "log level (debug, info, warn, error)")
	flags.String("db-connection-string", defaults.DBConnectionString, "PostgreSQL connection string")
	flags.String("tls-cert-path", defaults.TLSCertPath, "path to TLS certificate")
	flags.String("tls-key-path", defaults.TLSKeyPath, "path to TLS key")

	return cmd
}

func loadConfig(cmd *cobra.Command, configFile string) (settings.Config, error) {
	v := viper.New()

	defaults := settings.DefaultConfig()
	v.SetDefault("disable_secure", defaults.DisableSecure)
	v.SetDefault("port", defaults.Port)
	v.SetDefault("log_level", defaults.LogLevel)
	v.SetDefault("db_connection_string", defaults.DBConnectionString)
	v.SetDefault("tls_cert_path", defaults.TLSCertPath)
	v.SetDefault("tls_key_path", defaults.TLSKeyPath)

	if configFile != "" {
		v.SetConfigFile(configFile)
	} else {
		v.SetConfigName("config")
		v.SetConfigType("toml")
		v.AddConfigPath(".")
	}

	if err := v.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if configFile != "" || !errors.As(err, &notFound) {
			return settings.Config{}, fmt.Errorf("reading config file: %w", err)
		}
	}

	v.SetEnvPrefix("STREEPJES")
	v.AutomaticEnv()

	_ = v.BindPFlag("disable_secure", cmd.Flags().Lookup("disable-secure"))
	_ = v.BindPFlag("port", cmd.Flags().Lookup("port"))
	_ = v.BindPFlag("log_level", cmd.Flags().Lookup("log-level"))
	_ = v.BindPFlag("db_connection_string", cmd.Flags().Lookup("db-connection-string"))
	_ = v.BindPFlag("tls_cert_path", cmd.Flags().Lookup("tls-cert-path"))
	_ = v.BindPFlag("tls_key_path", cmd.Flags().Lookup("tls-key-path"))

	var cfg settings.Config
	if err := v.Unmarshal(&cfg); err != nil {
		return settings.Config{}, fmt.Errorf("parsing config: %w", err)
	}

	return cfg, nil
}

func newLogger(levelStr string) (*zap.Logger, error) {
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(levelStr)); err != nil {
		return nil, fmt.Errorf("invalid log level %q: %w", levelStr, err)
	}

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderCfg),
		zapcore.AddSync(os.Stderr),
		level,
	)

	return zap.New(core), nil
}

func run(cfg settings.Config) error {
	logger, err := newLogger(cfg.LogLevel)
	if err != nil {
		return err
	}
	defer logger.Sync() //nolint:errcheck

	var lis net.Listener
	if !cfg.DisableSecure {
		cer, err := tls.LoadX509KeyPair(cfg.TLSCertPath, cfg.TLSKeyPath)
		if err != nil {
			return fmt.Errorf("loading TLS keypair: %w", err)
		}

		lis, err = tls.Listen("tcp", ":443", &tls.Config{Certificates: []tls.Certificate{cer}})
		if err != nil {
			return fmt.Errorf("TLS listen: %w", err)
		}

		go redirectHTTPS(logger)
	} else {
		lis, err = net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
		if err != nil {
			return fmt.Errorf("listen: %w", err)
		}
	}

	db, err := postgres.OpenDB(cfg.DBConnectionString)
	if err != nil {
		return fmt.Errorf("opening database: %w", err)
	}
	defer db.Close()

	if err := postgres.Migrate(db, logger); err != nil {
		return fmt.Errorf("running migrations: %w", err)
	}

	userRepo := postgres.NewUserRepo(db, logger)
	memberRepo := postgres.NewMemberRepo(db, logger)
	orderRepo := postgres.NewOrderRepo(db, logger)
	catalogRepo := postgres.NewCatalogRepo(db, logger)

	authService := auth.New(userRepo, orderRepo)
	checkNoUsers(userRepo, authService, logger)

	orderService := order.New(memberRepo, orderRepo, catalogRepo)

	handler := router.New(static.Get, authService, orderService, !cfg.DisableSecure, logger)
	srv := &http.Server{Handler: handler}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	logger.Info("starting server", zap.Int("port", cfg.Port))
	go func() {
		if err := srv.Serve(lis); err != nil && err != http.ErrServerClosed {
			logger.Fatal("server exited", zap.Error(err))
		}
	}()

	<-sigChan
	logger.Info("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	return nil
}

func redirectHTTPS(logger *zap.Logger) {
	err := http.ListenAndServe(":80", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://"+r.Host+r.URL.String(), http.StatusMovedPermanently)
	}))
	if err != nil {
		logger.Fatal("HTTPS redirect server failed", zap.Error(err))
	}
}

func checkNoUsers(userRepo repo.User, authService auth.Service, logger *zap.Logger) {
	users, err := userRepo.GetAll()
	if err != nil {
		logger.Error("failed to check for existing users", zap.Error(err))
		return
	}
	if len(users) == 0 {
		_ = authService.Register(authdomain.User{ //nolint:exhaustivestruct
			Username: `adminGladiators`,
			Club:     domain.ClubGladiators,
			Name:     `Gladiators Admin`,
			Role:     authdomain.RoleAdmin,
		}, `playlacrossebecauseitsfun`)
		_ = authService.Register(authdomain.User{ //nolint:exhaustivestruct
			Username: `adminParabool`,
			Club:     domain.ClubParabool,
			Name:     `Parabool Admin`,
			Role:     authdomain.RoleAdmin,
		}, `groningerstudentenkorfbalcommissie`)
		_ = authService.Register(authdomain.User{ //nolint:exhaustivestruct
			Username: `adminCalamari`,
			Club:     domain.ClubCalamari,
			Name:     `Calamari Admin`,
			Role:     authdomain.RoleAdmin,
		}, `blub`)
	}
}
