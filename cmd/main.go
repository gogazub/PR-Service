package main

import (
	pullrequesthttp "PRService/internal/adapters/http/pullrequest"
	teamhttp "PRService/internal/adapters/http/team"
	userhttp "PRService/internal/adapters/http/user"
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

	_ "github.com/lib/pq"
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
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.PGUser,
		cfg.PGPassword,
		cfg.PGHost,
		cfg.PGPort,
		cfg.PGDatabase,
		cfg.SSLMode,
	)

	sugar.Info("try postgres conntection ", "address ", dsn)
	db, err := sql.Open("postgres", dsn)
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

	// Handlers
	userHandler := userhttp.NewHandler(svc, sugar)
	teamHandler := teamhttp.NewHandler(svc, sugar)
	pullrequestHandler := pullrequesthttp.NewHandler(svc, sugar)

	// Mux 
	mux := http.NewServeMux()

	mux.HandleFunc("/users/setIsActive", userHandler.SetIsActive)
	mux.HandleFunc("/users/getReview", userHandler.GetReview)

	mux.HandleFunc("/team/add", teamHandler.AddTeam)
	mux.HandleFunc("/team/get", teamHandler.GetTeam)

	mux.HandleFunc("/pullRequest/create", pullrequestHandler.CreatePullRequest)
	mux.HandleFunc("/pullRequest/merge", pullrequestHandler.MergePullRequest)
	mux.HandleFunc("/pullRequest/reassign", pullrequestHandler.ReassignReviewer)
	
	mux.HandleFunc("/health", HealtHandler)

	// Server
	
	port := ":"+cfg.HTTPPort
	sugar.Info("HTTP server started on ", port)
	err = http.ListenAndServe(port, mux)
	if err != nil {
		fmt.Println(err.Error())
	}

}

// TODO: refactor
func HealtHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Println("get a new request")
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-type", "text/plain")
	str := "Server is working!"
	_, _ = w.Write([]byte(str))
}