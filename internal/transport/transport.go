package transport

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"connectrpc.com/connect"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	v1 "github.com/coreyvan/kid-dictionary/gen/kiddictionary/v1"
	"github.com/coreyvan/kid-dictionary/gen/kiddictionary/v1/kiddictionaryv1connect"
	intconnect "github.com/coreyvan/kid-dictionary/internal/connect"
	"github.com/coreyvan/kid-dictionary/internal/domain"
	conversationsvc "github.com/coreyvan/kid-dictionary/internal/service/conversation"
	messagesvc "github.com/coreyvan/kid-dictionary/internal/service/message"
)

type Server interface {
	Listen(addr string) error
	Ready() chan struct{}
	Shutdown(ctx context.Context) error
	Addr() (string, error)
	Handler() http.Handler
}

type server struct {
	addr            string
	readyChan       chan struct{}
	logger          slog.Logger
	httpServer      *http.Server
	handler         http.Handler
	conversationSvc *conversationsvc.Service
	messageSvc      *messagesvc.Service
	messageRepo     domain.MessageRepository
}

var _ Server = (*server)(nil)
var _ kiddictionaryv1connect.AuthServiceHandler = (*server)(nil)
var _ kiddictionaryv1connect.ConversationServiceHandler = (*server)(nil)
var _ kiddictionaryv1connect.MessageServiceHandler = (*server)(nil)

func NewServer(bindAddr, port string, logger slog.Logger, convSvc *conversationsvc.Service, msgSvc *messagesvc.Service, msgRepo domain.MessageRepository) Server {
	s := &server{
		addr:            fmt.Sprintf("%s:%s", bindAddr, port),
		readyChan:       make(chan struct{}),
		logger:          logger,
		conversationSvc: convSvc,
		messageSvc:      msgSvc,
		messageRepo:     msgRepo,
	}
	s.handler = s.buildHandler()
	return s
}

func (s *server) buildHandler() http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(s.requestLogger)
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
		_, err := w.Write([]byte("ok"))
		if err != nil {
			s.logger.Error("error writing to health check response", "error", err)
		}
	})

	// Mount Connect handlers
	authPath, authHandler := kiddictionaryv1connect.NewAuthServiceHandler(s)
	r.Mount(authPath, authHandler)

	convPath, convHandler := kiddictionaryv1connect.NewConversationServiceHandler(s)
	r.Mount(convPath, convHandler)

	msgPath, msgHandler := kiddictionaryv1connect.NewMessageServiceHandler(s)
	r.Mount(msgPath, msgHandler)

	return r
}

func (s *server) Listen(addr string) error {
	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: s.handler,
	}

	s.logger.Info("☕️ Chi server listening on " + addr)
	close(s.readyChan)
	return s.httpServer.ListenAndServe()
}

func (s *server) Handler() http.Handler {
	return s.handler
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

// requestLogger is a middleware that logs HTTP requests using slog at Debug level.
// This keeps request logs visible in development but silent during tests.
func (s *server) requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		defer func() {
			s.logger.Debug("http request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", ww.Status(),
				"bytes", ww.BytesWritten(),
				"duration", time.Since(start),
				"remote", r.RemoteAddr,
			)
		}()

		next.ServeHTTP(ww, r)
	})
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
	// Map proto age bracket to domain age bracket
	ageBracket := domain.AgeBracket(req.Msg.AgeBracket)

	conv, err := s.conversationSvc.CreateConversation(ctx, req.Msg.Title, ageBracket)
	if err != nil {
		return nil, intconnect.MapError(err)
	}

	return connect.NewResponse(&v1.CreateConversationResponse{
		Conversation: &v1.Conversation{
			Id:         conv.ID.String(),
			Title:      conv.Title,
			AgeBracket: v1.AgeBracket(conv.AgeBracket),
			CreatedAt:  timestamppb.New(conv.CreatedAt),
			UpdatedAt:  timestamppb.New(conv.UpdatedAt),
		},
	}), nil
}

func (s *server) GetConversation(ctx context.Context, req *connect.Request[v1.GetConversationRequest]) (*connect.Response[v1.GetConversationResponse], error) {
	convID, err := uuid.Parse(req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid id"))
	}

	conv, err := s.conversationSvc.GetConversation(ctx, convID)
	if err != nil {
		return nil, intconnect.MapError(err)
	}

	// Get messages for this conversation
	msgs, err := s.messageRepo.GetByConversationID(ctx, convID)
	if err != nil {
		return nil, intconnect.MapError(err)
	}

	// Convert messages to proto
	protoMsgs := make([]*v1.Message, 0, len(msgs))
	for _, msg := range msgs {
		protoMsgs = append(protoMsgs, &v1.Message{
			Id:             msg.ID.String(),
			ConversationId: msg.ConversationID.String(),
			Role:           v1.MessageRole(msg.Role),
			Content:        msg.Content,
			ContentTier:    v1.ContentTier(msg.ContentTier),
			CreatedAt:      timestamppb.New(msg.CreatedAt),
		})
	}

	return connect.NewResponse(&v1.GetConversationResponse{
		Conversation: &v1.Conversation{
			Id:         conv.ID.String(),
			Title:      conv.Title,
			AgeBracket: v1.AgeBracket(conv.AgeBracket),
			CreatedAt:  timestamppb.New(conv.CreatedAt),
			UpdatedAt:  timestamppb.New(conv.UpdatedAt),
		},
		Messages: protoMsgs,
	}), nil
}

func (s *server) ListConversations(ctx context.Context, req *connect.Request[v1.ListConversationsRequest]) (*connect.Response[v1.ListConversationsResponse], error) {
	pageSize := int(req.Msg.PageSize)
	if pageSize <= 0 {
		pageSize = 20
	}

	// Parse page token as offset (simple pagination)
	offset := 0
	if req.Msg.PageToken != "" {
		// Page token is the offset encoded
		_, err := fmt.Sscanf(req.Msg.PageToken, "%d", &offset)
		if err != nil {
			offset = 0
		}
	}

	convs, err := s.conversationSvc.ListConversations(ctx, nil, pageSize, offset)
	if err != nil {
		return nil, intconnect.MapError(err)
	}

	// Convert to proto
	protoConvs := make([]*v1.Conversation, 0, len(convs))
	for _, conv := range convs {
		protoConvs = append(protoConvs, &v1.Conversation{
			Id:         conv.ID.String(),
			Title:      conv.Title,
			AgeBracket: v1.AgeBracket(conv.AgeBracket),
			CreatedAt:  timestamppb.New(conv.CreatedAt),
			UpdatedAt:  timestamppb.New(conv.UpdatedAt),
		})
	}

	// Generate next page token if we got a full page
	var nextPageToken string
	if len(convs) == pageSize {
		nextPageToken = fmt.Sprintf("%d", offset+pageSize)
	}

	return connect.NewResponse(&v1.ListConversationsResponse{
		Conversations: protoConvs,
		NextPageToken: nextPageToken,
	}), nil
}

func (s *server) UpdateConversation(ctx context.Context, req *connect.Request[v1.UpdateConversationRequest]) (*connect.Response[v1.UpdateConversationResponse], error) {
	convID, err := uuid.Parse(req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid id"))
	}

	var title *string
	if req.Msg.Title != nil {
		title = req.Msg.Title
	}

	var ageBracket *domain.AgeBracket
	if req.Msg.AgeBracket != nil {
		ab := domain.AgeBracket(*req.Msg.AgeBracket)
		ageBracket = &ab
	}

	conv, err := s.conversationSvc.UpdateConversation(ctx, convID, title, ageBracket)
	if err != nil {
		return nil, intconnect.MapError(err)
	}

	return connect.NewResponse(&v1.UpdateConversationResponse{
		Conversation: &v1.Conversation{
			Id:         conv.ID.String(),
			Title:      conv.Title,
			AgeBracket: v1.AgeBracket(conv.AgeBracket),
			CreatedAt:  timestamppb.New(conv.CreatedAt),
			UpdatedAt:  timestamppb.New(conv.UpdatedAt),
		},
	}), nil
}

func (s *server) DeleteConversation(ctx context.Context, req *connect.Request[v1.DeleteConversationRequest]) (*connect.Response[v1.DeleteConversationResponse], error) {
	convID, err := uuid.Parse(req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid id"))
	}

	if err := s.conversationSvc.DeleteConversation(ctx, convID); err != nil {
		return nil, intconnect.MapError(err)
	}

	return connect.NewResponse(&v1.DeleteConversationResponse{}), nil
}

func (s *server) SendMessage(ctx context.Context, req *connect.Request[v1.SendMessageRequest]) (*connect.Response[v1.SendMessageResponse], error) {
	var conversationID uuid.UUID
	var err error

	// Parse conversation_id if provided; uuid.Nil signals auto-creation
	if req.Msg.ConversationId != "" {
		conversationID, err = uuid.Parse(req.Msg.ConversationId)
		if err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid conversation_id"))
		}
	}

	// Extract age_bracket for auto-creation (only used when conversation_id is empty)
	var ageBracket *domain.AgeBracket
	if req.Msg.AgeBracket != v1.AgeBracket_AGE_BRACKET_UNSPECIFIED {
		ab := domain.AgeBracket(req.Msg.AgeBracket)
		ageBracket = &ab
	}

	result, err := s.messageSvc.SendMessage(ctx, conversationID, req.Msg.Content, ageBracket)
	if err != nil {
		return nil, intconnect.MapError(err)
	}

	resp := &v1.SendMessageResponse{
		UserMessage: &v1.Message{
			Id:             result.UserMessage.ID.String(),
			ConversationId: result.UserMessage.ConversationID.String(),
			Role:           v1.MessageRole(result.UserMessage.Role),
			Content:        result.UserMessage.Content,
			CreatedAt:      timestamppb.New(result.UserMessage.CreatedAt),
		},
		AssistantMessage: &v1.Message{
			Id:             result.AssistantMessage.ID.String(),
			ConversationId: result.AssistantMessage.ConversationID.String(),
			Role:           v1.MessageRole(result.AssistantMessage.Role),
			Content:        result.AssistantMessage.Content,
			ContentTier:    v1.ContentTier(result.AssistantMessage.ContentTier),
			CreatedAt:      timestamppb.New(result.AssistantMessage.CreatedAt),
		},
	}

	// Include conversation in response if it was auto-created
	if result.Conversation != nil {
		resp.Conversation = &v1.Conversation{
			Id:         result.Conversation.ID.String(),
			Title:      result.Conversation.Title,
			AgeBracket: v1.AgeBracket(result.Conversation.AgeBracket),
			CreatedAt:  timestamppb.New(result.Conversation.CreatedAt),
			UpdatedAt:  timestamppb.New(result.Conversation.UpdatedAt),
		}
	}

	return connect.NewResponse(resp), nil
}
