package main

import (
	"PRService/config"
	httperror "PRService/internal/adapters/http/error"
	pullreqhandler "PRService/internal/adapters/http/pullrequest/handlers"
	teamhandlers "PRService/internal/adapters/http/team/handlers"
	userhandler "PRService/internal/adapters/http/user/handler"
	pullrequestrepo "PRService/internal/adapters/repo/pullrequest"
	teamrepo "PRService/internal/adapters/repo/team"
	"PRService/internal/adapters/repo/transactor"
	userrepo "PRService/internal/adapters/repo/user"
	"PRService/internal/app"
	pullrequestusecase "PRService/internal/usecase/pullrequest"
	teamusecase "PRService/internal/usecase/team"
	userusecase "PRService/internal/usecase/user"
	"PRService/pkg/logger"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
)

func main() {

	// 0. Init Logger
	logger, err := logger.New()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	defer func() {
		_ = logger.Sync()
	}()

	// 1. Init Config
	cfg, err := config.NewConfig()

	if err != nil {
		logger.Fatal("new config", "error", err)
	}

	// 2. Connect to Postgres
	dsn := cfg.PG.URL

	logger.Info("try postgres connection", "address", dsn)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.Fatal("failed to open postgres", "error", err)
	}

	// Ping DB
	if err := db.Ping(); err != nil {
		logger.Fatal("failed to ping postgres", "error", err)
	}

	defer func() {
		logger.Info("closing database connection")
		if err := db.Close(); err != nil {
			logger.Error("failed to close db", "error", err)
		}
	}()

	logger.Info("connected to postgres")

	// 3. Init Layers

	// Repositories
	userRepo := userrepo.New(db)
	teamRepo := teamrepo.New(db)
	pullrequestRepo := pullrequestrepo.New(db)

	// Services
	userSvc := userusecase.New(userRepo)
	teamSvc := teamusecase.New(teamRepo)
	pullrequestSvc := pullrequestusecase.New(pullrequestRepo)

	svc := app.NewServices(userSvc, teamSvc, pullrequestSvc, transactor.NewTransactor(db))

	// Handlers
	userHandler := userhandler.NewHandler(svc, logger)
	teamHandler := teamhandlers.NewHandler(svc, logger)
	pullrequestHandler := pullreqhandler.NewHandler(svc, logger)

	// 4. Router
	mux := http.NewServeMux()

	mux.HandleFunc("/users/setIsActive", userHandler.SetIsActive)
	mux.HandleFunc("/users/getReview", userHandler.GetReview)

	mux.HandleFunc("/team/add", teamHandler.AddTeam)
	mux.HandleFunc("/team/get", teamHandler.GetTeam)

	mux.HandleFunc("/pullRequest/create", pullrequestHandler.CreatePullRequest)
	mux.HandleFunc("/pullRequest/merge", pullrequestHandler.MergePullRequest)
	mux.HandleFunc("/pullRequest/reassign", pullrequestHandler.ReassignReviewer)

	mux.HandleFunc("/health", HealthHandler)

	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mux.ServeHTTP(w, r)
	})

	// 5. Setup HTTP Server
	port := cfg.HTTP.PORT
	port = ":" + port

	srv := &http.Server{
		Addr:         port,
		Handler:      finalHandler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 6. Start server in goroutine
	go func() {
		logger.Info("HTTP server started", "port", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("listen and serve failed", "error", err)
		}
	}()

	// 7. Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("server forced to shutdown", "error", err)
	}

	logger.Info("server exited properly")

}

type healthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Time    string `json:"timestamp"`
}

func HealthHandler(w http.ResponseWriter, _ *http.Request) {
	response := healthResponse{
		Status:  "healthy",
		Message: "Server is working properly",
		Time:    time.Now().UTC().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		httperror.WriteErrorResponse(w, http.StatusInternalServerError, httperror.ErrorCodeInternal, "health check failed")
	}
}
