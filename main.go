//go:build !dev

package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
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
	"github.com/charmbracelet/log"
)

func main() {
	os.Exit(run())
}

func run() int {
	readSettings()

	logLevel := log.ErrorLevel
	if settings.Debug {
		logLevel = log.DebugLevel
	}
	log.Default().SetLevel(logLevel)

	var lis net.Listener
	var err error
	if !settings.DisableSecure {
		cer, err := tls.LoadX509KeyPair(settings.TLSCertPath, settings.TLSKeyPath)
		if err != nil {
			panic(err)
		}

		lis, err = tls.Listen("tcp", ":443", &tls.Config{Certificates: []tls.Certificate{cer}})

		go redirectHTTPS()
	} else {
		lis, err = net.Listen("tcp", fmt.Sprintf(":%d", settings.Port))
	}
	if err != nil {
		panic(err)
	}

	db, err := postgres.OpenDB(settings.DBConnectionString)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	postgres.Migrate(db)

	userRepo := postgres.NewUserRepo(db)
	memberRepo := postgres.NewMemberRepo(db)
	orderRepo := postgres.NewOrderRepo(db)
	catalogRepo := postgres.NewCatalogRepo(db)

	authService := auth.New(userRepo, orderRepo)
	checkNoUsers(userRepo, authService)

	orderService := order.New(memberRepo, orderRepo, catalogRepo)

	handler := router.New(static.Get, authService, orderService)
	srv := &http.Server{Handler: handler}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	log.Info("Starting server", "port", settings.Port)
	go func() {
		if err := srv.Serve(lis); err != nil && err != http.ErrServerClosed {
			log.Fatal("server exited", "error", err)
		}
	}()

	<-sigChan
	log.Info("Shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("server forced to shutdown", "error", err)
		return 1
	}

	return 0
}

func redirectHTTPS() {
	err := http.ListenAndServe(":80", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://"+r.Host+r.URL.String(), http.StatusMovedPermanently)
	}))
	if err != nil {
		panic(err)
	}
}

// check if there are no users in the database, if so, insert some
func checkNoUsers(userRepo repo.User, authService auth.Service) {
	if len(userRepo.GetAll()) == 0 {
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

// read settings from environment
func readSettings() {
	disableSecure, ok := os.LookupEnv(`STREEPJES_DISABLE_SECURE`)
	settings.DisableSecure = ok && disableSecure == `true`

	port, ok := os.LookupEnv(`STREEPJES_PORT`)
	if ok {
		portN, err := strconv.Atoi(port)
		if err != nil {
			panic(err)
		}

		settings.Port = portN
	}

	settings.Debug = os.Getenv("STREEPJES_DEBUG") == "true"
	dbConnectionString := os.Getenv("STREEPJES_DB_CONNECTION_STRING")
	if dbConnectionString != "" {
		settings.DBConnectionString = dbConnectionString
	}

	tlsCertPath := os.Getenv("STREEPJES_TLS_CERT_PATH")
	if tlsCertPath != "" {
		settings.TLSCertPath = tlsCertPath
	}

	tlsKeyPath := os.Getenv("STREEPJES_TLS_KEY_PATH")
	if tlsKeyPath != "" {
		settings.TLSKeyPath = tlsKeyPath
	}
}
