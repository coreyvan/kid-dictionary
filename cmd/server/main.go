package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/coreyvan/kid-dictionary/internal/config"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}

func run() error {
	cfg := config.GetConfig()
	l := config.NewLogger(cfg)
	l.Info("⚙️ Starting server", "config", cfg)

	wiring := config.NewWiring(cfg, *l)
	server := wiring.MustProvideServer()

	errChan := make(chan error)
	go func() {
		if err := server.Listen(fmt.Sprintf("%s:%s", cfg.BindAddr, cfg.Port)); err != nil {
			errChan <- err
		}
	}()

	<-server.Ready()

	// listen to os signals so we can gracefully shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errChan:
		return fmt.Errorf("server error: %w", err)
	case sig := <-sigChan:
		l.Info("received signal, shutting down", "signal", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			return fmt.Errorf("error during shutdown: %w", err)
		}
		return nil
	}
}
