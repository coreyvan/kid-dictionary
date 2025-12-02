package transport

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"connectrpc.com/connect"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	v1 "github.com/coreyvan/kid-dictionary/gen/kiddictionary/v1"
	"github.com/coreyvan/kid-dictionary/gen/kiddictionary/v1/kiddictionaryv1connect"
)

type Server interface {
	Listen(addr string) error
	Ready() chan struct{}
	Shutdown(ctx context.Context) error
	Addr() (string, error)
}

type server struct {
	addr       string
	readyChan  chan struct{}
	logger     slog.Logger
	httpServer *http.Server
}

var _ Server = (*server)(nil)
var _ kiddictionaryv1connect.AuthServiceHandler = (*server)(nil)
var _ kiddictionaryv1connect.ConversationServiceHandler = (*server)(nil)
var _ kiddictionaryv1connect.MessageServiceHandler = (*server)(nil)

func NewServer(bindAddr, port string, logger slog.Logger) Server {
	return &server{
		addr:      fmt.Sprintf("%s:%s", bindAddr, port),
		readyChan: make(chan struct{}),
		logger:    logger,
	}
}

func (s *server) Listen(addr string) error {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "Connect-Protocol-Version"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	// Mount Connect handlers
	authPath, authHandler := kiddictionaryv1connect.NewAuthServiceHandler(s)
	r.Mount(authPath, authHandler)

	convPath, convHandler := kiddictionaryv1connect.NewConversationServiceHandler(s)
	r.Mount(convPath, convHandler)

	msgPath, msgHandler := kiddictionaryv1connect.NewMessageServiceHandler(s)
	r.Mount(msgPath, msgHandler)

	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: r,
	}

	s.logger.Info("☕️ Chi server listening on " + addr)
	close(s.readyChan)
	return s.httpServer.ListenAndServe()
}

func (s *server) Ready() chan struct{} {
	return s.readyChan
}

func (s *server) Shutdown(ctx context.Context) error {
	if s.httpServer != nil {
		return s.httpServer.Shutdown(ctx)
	}
	return nil
}

func (s *server) Addr() (string, error) {
	return s.addr, nil
}

func (s *server) Register(ctx context.Context, req *connect.Request[v1.RegisterRequest]) (*connect.Response[v1.RegisterResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("not implemented"))
}

func (s *server) Login(ctx context.Context, req *connect.Request[v1.LoginRequest]) (*connect.Response[v1.LoginResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("not implemented"))
}

func (s *server) RefreshToken(ctx context.Context, req *connect.Request[v1.RefreshTokenRequest]) (*connect.Response[v1.RefreshTokenResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("not implemented"))
}

func (s *server) CreateConversation(ctx context.Context, req *connect.Request[v1.CreateConversationRequest]) (*connect.Response[v1.CreateConversationResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("not implemented"))
}

func (s *server) GetConversation(ctx context.Context, req *connect.Request[v1.GetConversationRequest]) (*connect.Response[v1.GetConversationResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("not implemented"))
}

func (s *server) ListConversations(ctx context.Context, req *connect.Request[v1.ListConversationsRequest]) (*connect.Response[v1.ListConversationsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("not implemented"))
}

func (s *server) DeleteConversation(ctx context.Context, req *connect.Request[v1.DeleteConversationRequest]) (*connect.Response[v1.DeleteConversationResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("not implemented"))
}

func (s *server) SendMessage(ctx context.Context, req *connect.Request[v1.SendMessageRequest]) (*connect.Response[v1.SendMessageResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("not implemented"))
}
