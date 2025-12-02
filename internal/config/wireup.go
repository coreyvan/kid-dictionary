package config

import "github.com/coreyvan/kid-dictionary/internal/transport"

type Wireup struct {
	cfg Config

	Deps WireupDeps
}

type WireupDeps struct {
	Transport transport.Server
}

type Wiring interface {
	MustProvideServer() transport.Server
}

func NewWiring(cfg Config) Wiring {
	return &Wireup{
		cfg: cfg,
	}
}

func (w *Wireup) MustProvideServer() transport.Server {
	if w.Deps.Transport != nil {
		return w.Deps.Transport
	}

	t := transport.NewServer(w.cfg.BindAddr, w.cfg.Port)
	w.Deps.Transport = t
	return t
}
