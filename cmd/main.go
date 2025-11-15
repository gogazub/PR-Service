package main

import (
	pullrequestrepo "PRService/internal/adapters/repo/pullrequest_repo"
	teamrepo "PRService/internal/adapters/repo/team_repo"
	userrepo "PRService/internal/adapters/repo/user_repo"
	"PRService/internal/app"
	"PRService/internal/config"
	pullrequest_usecase "PRService/internal/usecase/pullrequest"
	team_usecase "PRService/internal/usecase/team"
	user_usecase "PRService/internal/usecase/user"
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"go.uber.org/zap"
)

func main() {

	// Logger 
	logger, err  := zap.NewDevelopmentConfig().Build()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer logger.Sync()

	sugar := logger.Sugar()

	// Config
	cfg, err := config.LoadConfig()
	if err != nil {
		sugar.Fatalw("load config", "error", err)
	}

	// Postgres
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		cfg.PGUser,
		cfg.PGPassword,
		cfg.PGHost,
		cfg.PGPort,
		cfg.PGDatabase,
	)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		sugar.Fatalw("failed to open postgres", "error", err)
	}

	if err := db.Ping(); err != nil {
		sugar.Fatal("failed to ping postgres", "error", err)
	}

	sugar.Info("connected to postgres")

	// Repositories

	userRepo := userrepo.New(db)
	teamRepo := teamrepo.New(db)
	pullrequestRepo := pullrequestrepo.New(db)

	// Services
	userSvc := user_usecase.New(userRepo)
	teamSvc := team_usecase.New(teamRepo)
	pullrequestSvc := pullrequest_usecase.New(pullrequestRepo)
	svc := app.NewServices(userSvc, teamSvc, pullrequestSvc)

	// Server
	heatlzHandler := func(w http.ResponseWriter, _ *http.Request) {
		fmt.Println("get a new request")
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-type", "text/plain")
		str := "Server is working!"
		_, _ = w.Write([]byte(str))
	}

	http.HandleFunc("/healtz", heatlzHandler)

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err.Error())
	}

}
