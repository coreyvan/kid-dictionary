package config

import (
	"log/slog"

	"github.com/coreyvan/kid-dictionary/internal/transport"
)

type Wireup struct {
	cfg    Config
	logger slog.Logger

	Deps WireupDeps
}

type WireupDeps struct {
	Transport transport.Server
}

type Wiring interface {
	MustProvideServer() transport.Server
}

func NewWiring(cfg Config, logger slog.Logger) Wiring {
	return &Wireup{
		cfg:    cfg,
		logger: logger,
	}
}

func (w *Wireup) MustProvideServer() transport.Server {
	if w.Deps.Transport != nil {
		return w.Deps.Transport
	}

	t := transport.NewServer(w.cfg.BindAddr, w.cfg.Port, w.logger)
	w.Deps.Transport = t
	return t
}
