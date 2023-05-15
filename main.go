//go:build !dev

package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/FallenTaters/streepjes/backend/application/auth"
	"github.com/FallenTaters/streepjes/backend/application/order"
	"github.com/FallenTaters/streepjes/backend/global/settings"
	"github.com/FallenTaters/streepjes/backend/infrastructure/repo"
	"github.com/FallenTaters/streepjes/backend/infrastructure/repo/postgres"
	"github.com/FallenTaters/streepjes/backend/infrastructure/repo/sqlite"
	"github.com/FallenTaters/streepjes/backend/infrastructure/router"
	"github.com/FallenTaters/streepjes/domain"
	"github.com/FallenTaters/streepjes/domain/authdomain"
	"github.com/FallenTaters/streepjes/static"
	"github.com/charmbracelet/log"
)

func main() {
	// os.Exit(run())

	readSettings()

	oldDB, err := sqlite.OpenDB("streepjes.db")
	if err != nil {
		panic(err)
	}

	newDB, err := postgres.OpenDB(settings.DBConnectionString)
	if err != nil {
		panic(err)
	}

	postgres.Migrate(newDB)

	tx, err := newDB.Begin()
	if err != nil {
		panic(err)
	}
	defer tx.Rollback()

	users := sqlite.NewUserRepo(oldDB).GetAll()
	cat := sqlite.NewCatalogRepo(oldDB)
	categories := cat.GetCategories()
	items := cat.GetItems()

	members := sqlite.NewMemberRepo(oldDB).GetAll()
	orders := sqlite.NewOrderRepo(oldDB).Filter(repo.OrderFilter{
		Start: time.Date(2022, time.September, 1, 0, 0, 0, 0, time.Local),
		Limit: 100000000000,
	})

	q := "INSERT INTO users (id, username, password, club, name, role) VALUES "
	for i, user := range users {
		if i != 0 {
			q += ","
		}
		q += fmt.Sprintf("(%d,'%s','%s','%s','%s','%s')", user.ID, user.Username, user.PasswordHash, user.Club, user.Name, user.Role)
	}
	_, err = tx.Exec(q + ";")
	if err != nil {
		panic(err)
	}

	q = "INSERT INTO categories (id, name) VALUES "
	for i, cat := range categories {
		if i != 0 {
			q += ","
		}
		q += fmt.Sprintf("(%d,'%s')", cat.ID, cat.Name)
	}
	_, err = tx.Exec(q + ";")
	if err != nil {
		panic(err)
	}

	q = "INSERT INTO items (id, category_id, name, price_gladiators, price_parabool) VALUES "
	for i, item := range items {
		if item.CategoryID == 0 {
			fmt.Printf("SKIPPING ITEM %d (%s)\n", item.ID, item.Name)
			continue
		}
		if i != 0 {
			q += ","
		}
		q += fmt.Sprintf("(%d,%d,'%s',%d,%d)", item.ID, item.CategoryID, item.Name, item.PriceGladiators, item.PriceParabool)
	}
	_, err = tx.Exec(q + ";")
	if err != nil {
		panic(err)
	}

	format := "2006-01-02 15:04:05.000"

	memberIDs := map[int]bool{}
	q = "INSERT INTO members (id, club, name, last_order) VALUES "
	for i, member := range members {
		memberIDs[member.ID] = true
		if i != 0 {
			q += ","
		}
		q += fmt.Sprintf("(%d,'%s','%s','%s')", member.ID, member.Club, member.Name, member.LastOrder.Format(format))
	}
	_, err = tx.Exec(q + ";")
	if err != nil {
		panic(err)
	}

	q = "INSERT INTO orders (id, club, bartender_id, member_id, contents, price, order_time, status, status_time) VALUES "
	for i, o := range orders {
		if o.MemberID == 0 { // anonymous orders
			continue
		}

		if i != 0 {
			q += ","
		}
		q += fmt.Sprintf("(%d,'%s',%d,%d,'%s',%d,'%s','%s','%s')", o.ID, o.Club, o.BartenderID, o.MemberID, o.Contents, o.Price, o.OrderTime.Format(format), o.Status, o.StatusTime.Format(format))
	}
	_, err = tx.Exec(q + ";")
	if err != nil {
		panic(err)
	}

	err = tx.Commit()
	if err != nil {
		panic(err)
	}
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

	sigChan := make(chan os.Signal)
	shutdown := make(chan int)
	signal.Notify(sigChan, os.Interrupt, os.Kill)
	go func() {
		<-sigChan
		shutdown <- 0
	}()

	postgres.Migrate(db)

	userRepo := postgres.NewUserRepo(db)
	memberRepo := postgres.NewMemberRepo(db)
	orderRepo := postgres.NewOrderRepo(db)
	catalogRepo := postgres.NewCatalogRepo(db)

	authService := auth.New(userRepo, orderRepo)
	checkNoUsers(userRepo, authService)

	orderService := order.New(memberRepo, orderRepo, catalogRepo)

	handler := router.New(static.Get, authService, orderService)

	log.Info("Starting server", "port", settings.Port)
	go func() {
		err := http.Serve(lis, handler)
		log.Fatal("server exited", "error", err)
		shutdown <- 1
	}()

	return <-shutdown
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
