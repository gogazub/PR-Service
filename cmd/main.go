package main

import (
	"PRService/config"
	pullreqhandler "PRService/internal/adapters/http/pullrequest/handlers"
	teamhandlers "PRService/internal/adapters/http/team/handlers"
	userhandler "PRService/internal/adapters/http/user/handler"
	pullrequestrepo "PRService/internal/adapters/repo/pullrequest_repo"
	teamrepo "PRService/internal/adapters/repo/team_repo"
	userrepo "PRService/internal/adapters/repo/user_repo"
	"PRService/internal/app"
	pullrequest_usecase "PRService/internal/usecase/pullrequest"
	team_usecase "PRService/internal/usecase/team"
	user_usecase "PRService/internal/usecase/user"
	"PRService/pkg/logger"
	"database/sql"
	"fmt"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

func main() {

	logger, err := logger.New()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer logger.Sync()

	// Config
	cfg, err := config.NewConfig()

	if err != nil {
		logger.Fatal("load config", "error", err)
	}

	// Postgres
	dsn := cfg.PG.URL
	

	logger.Info("try postgres conntection ", "address ", dsn)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.Fatal("failed to open postgres", "error", err)
	}

	if err := db.Ping(); err != nil {
		logger.Fatal("failed to ping postgres", "error", err)
	}

	logger.Info("connected to postgres")

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
	userHandler := userhandler.NewHandler(svc, logger)
	teamHandler := teamhandlers.NewHandler(svc, logger)
	pullrequestHandler := pullreqhandler.NewHandler(svc, logger)

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
	
	port := cfg.HTTP.PORT
	port = ":"+port
	logger.Info("HTTP server started on ", port)
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