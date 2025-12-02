package config

import (
	"log/slog"

	"github.com/coreyvan/kid-dictionary/internal/llm"
	"github.com/coreyvan/kid-dictionary/internal/transport"
)

type Wireup struct {
	cfg    Config
	logger slog.Logger

	Deps WireupDeps
}

type WireupDeps struct {
	Transport   transport.Server
	LLMProvider llm.Provider
}

type Wiring interface {
	MustProvideServer() transport.Server
	MustProvideLLMProvider() llm.Provider
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

func (w *Wireup) MustProvideLLMProvider() llm.Provider {
	if w.Deps.LLMProvider != nil {
		return w.Deps.LLMProvider
	}

	opts := []llm.OpenAIOption{}
	if w.cfg.OpenAIModel != "" {
		opts = append(opts, llm.WithModel(w.cfg.OpenAIModel))
	}

	p := llm.NewOpenAIProvider(w.cfg.OpenAIAPIKey, opts...)
	w.Deps.LLMProvider = p
	return p
}
