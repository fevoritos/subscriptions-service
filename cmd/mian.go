package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"subs-service/config"
	"subs-service/internal/lib/logger/slogpretty"
	"subs-service/internal/middleware"
	"subs-service/internal/transport/http/handlers"

	subusecese "subs-service/internal/usecase/subscription"
	"syscall"
	"time"

	infrastructurepostgres "subs-service/internal/infrastucture/postgres"
	subrepo "subs-service/internal/reposotory/postgres"
	transporthttp "subs-service/internal/transport/http"
)

func main() {
	log := slogpretty.SetupPrettySlog()

	cfg, err := config.Load()
	if err != nil {
		log.Error("load config", "error", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	defer stop()

	pool, err := infrastructurepostgres.Open(ctx, cfg.DatabaseDSN)
	if err != nil {
		log.Error("open postgres", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	subrepo := subrepo.New(log, pool)
	subservice := subusecese.NewSubscriptionService(log, subrepo)
	subhandler := handlers.NewSubscriptionHandler(log, subservice)

	stack := middleware.Chain(
		middleware.CORS,
		middleware.NewLogger(log),
	)
	router := transporthttp.NewRouter(subhandler)

	server := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           stack(router),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Error("shutdown http server", "error", err)
		}
	}()

	log.Info("http server started", "addr", cfg.HTTPAddr)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error("listen and serve", "error", err)
		os.Exit(1)
	}
}
