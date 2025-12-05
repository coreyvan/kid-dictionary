package config

import (
	"context"
	"log/slog"

	"github.com/coreyvan/kid-dictionary/internal/conversation"
	"github.com/coreyvan/kid-dictionary/internal/llm"
	"github.com/coreyvan/kid-dictionary/internal/message"
	"github.com/coreyvan/kid-dictionary/internal/transport"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Wireup struct {
	cfg    Config
	logger slog.Logger

	Deps WireupDeps
}

type WireupDeps struct {
	Transport            transport.Server
	LLMProvider          llm.Provider
	DBPool               *pgxpool.Pool
	ConversationRepo     conversation.Repository
	MessageRepo          message.Repository
	ConversationService  *conversation.Service
	MessageService       *message.Service
}

type Wiring interface {
	MustProvideServer() transport.Server
	MustProvideLLMProvider() llm.Provider
	MustProvideDBPool() *pgxpool.Pool
	MustProvideConversationRepo() conversation.Repository
	MustProvideMessageRepo() message.Repository
	MustProvideConversationService() *conversation.Service
	MustProvideMessageService() *message.Service
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

	convSvc := w.MustProvideConversationService()
	msgSvc := w.MustProvideMessageService()
	msgRepo := w.MustProvideMessageRepo()

	t := transport.NewServer(w.cfg.BindAddr, w.cfg.Port, w.logger, convSvc, msgSvc, msgRepo)
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

func (w *Wireup) MustProvideDBPool() *pgxpool.Pool {
	if w.Deps.DBPool != nil {
		return w.Deps.DBPool
	}

	if w.cfg.DatabaseURL == "" {
		panic("DATABASE_URL is required")
	}

	pool, err := pgxpool.New(context.Background(), w.cfg.DatabaseURL)
	if err != nil {
		panic("failed to connect to database: " + err.Error())
	}

	w.Deps.DBPool = pool
	return pool
}

func (w *Wireup) MustProvideConversationRepo() conversation.Repository {
	if w.Deps.ConversationRepo != nil {
		return w.Deps.ConversationRepo
	}

	pool := w.MustProvideDBPool()
	repo := conversation.NewPostgresRepository(pool)
	w.Deps.ConversationRepo = repo
	return repo
}

func (w *Wireup) MustProvideMessageRepo() message.Repository {
	if w.Deps.MessageRepo != nil {
		return w.Deps.MessageRepo
	}

	pool := w.MustProvideDBPool()
	repo := message.NewPostgresRepository(pool)
	w.Deps.MessageRepo = repo
	return repo
}

func (w *Wireup) MustProvideConversationService() *conversation.Service {
	if w.Deps.ConversationService != nil {
		return w.Deps.ConversationService
	}

	repo := w.MustProvideConversationRepo()
	svc := conversation.NewService(repo)
	w.Deps.ConversationService = svc
	return svc
}

func (w *Wireup) MustProvideMessageService() *message.Service {
	if w.Deps.MessageService != nil {
		return w.Deps.MessageService
	}

	msgRepo := w.MustProvideMessageRepo()
	convRepo := w.MustProvideConversationRepo()
	llmProvider := w.MustProvideLLMProvider()
	svc := message.NewService(msgRepo, convRepo, llmProvider)
	w.Deps.MessageService = svc
	return svc
}
