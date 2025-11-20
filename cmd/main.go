package main

import (
	"PRService/config"
	pullreqhandler "PRService/internal/adapters/http/pullrequest/handlers"
	teamhandlers "PRService/internal/adapters/http/team/handlers"
	userhandler "PRService/internal/adapters/http/user/handler"
	pullrequestrepo "PRService/internal/adapters/repo/pullrequest"
	teamrepo "PRService/internal/adapters/repo/team"
	"PRService/internal/adapters/repo/transactor"
	userrepo "PRService/internal/adapters/repo/user"
	"PRService/internal/app"
	pullrequest_usecase "PRService/internal/usecase/pullrequest"
	team_usecase "PRService/internal/usecase/team"
	user_usecase "PRService/internal/usecase/user"
	"PRService/pkg/logger"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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
		logger.Fatal("new config", "error", err)
	}

	// Postgres
	dsn := cfg.PG.URL

	logger.Info("try postgres connection", "address", dsn)
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


	svc := app.NewServices(userSvc, teamSvc, pullrequestSvc, transactor.NewTransactor(db))

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

	mux.HandleFunc("/health", HealthHandler)

	// Server
	port := cfg.HTTP.PORT
	port = ":" + port
	logger.Info("HTTP server started on", "port", port)
	err = http.ListenAndServe(port, mux)
	if err != nil {
		fmt.Println(err.Error())
	}

}

type healthResponse struct {
    Status  string `json:"status"`
    Message string `json:"message"`
    Time    string `json:"timestamp"`
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
    log.Printf("Health check from %s", r.RemoteAddr)
    
    response := healthResponse{
        Status:  "healthy",
        Message: "Server is working properly",
        Time:    time.Now().UTC().Format(time.RFC3339),
    }
    
    w.Header().Set("Content-Type", "application/json")
    
    if err := json.NewEncoder(w).Encode(response); err != nil {
        log.Printf("Failed to encode health response: %v", err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
    }
}
