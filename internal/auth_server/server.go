package auth_server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/salesforceanton/pocket-tg-bot/internal/logger"
	"github.com/salesforceanton/pocket-tg-bot/internal/repository"
	"github.com/zhashkevych/go-pocket-sdk"
)

const (
	AUTH_SERVER_PORT  = 80
	CHAT_ID_URL_PARAM = "chat_id"
	LOCATION_HEADER   = "Location"
)

type Server struct {
	httpServer   *http.Server
	pocketClient *pocket.Client
	repo         repository.Repository

	redirectUrl string
}

func NewServer(pocketClient *pocket.Client, repo repository.Repository, redirectUrl string) *Server {
	return &Server{
		pocketClient: pocketClient,
		repo:         repo,
		redirectUrl:  redirectUrl,
	}
}

func (s *Server) Run() error {
	s.httpServer = &http.Server{
		Addr:           fmt.Sprintf(":%d", AUTH_SERVER_PORT),
		Handler:        s,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	return s.httpServer.ListenAndServe()
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Validate request method
	if r.Method != http.MethodGet {
		logger.LogIssueWithPoint("auth server", errors.New("incorrect request method"))
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Retain chat id from request url
	chatId, err := s.retainChatIdValue(r)
	if err != nil {
		logger.LogIssueWithPoint("auth server", errors.New("incorrect param in URL"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Retain request token from repo
	requestToken, err := s.repo.GetRequestToken(chatId)
	if err != nil {
		logger.LogIssueWithPoint("auth server", errors.New("request token not found"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Make authorize callout for getting Access token
	authResp, err := s.pocketClient.Authorize(context.TODO(), requestToken)
	if err != nil {
		logger.LogIssueWithPoint("auth server - authorize to pocket", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	accessToken := authResp.AccessToken

	// Store Access token in repo
	s.repo.SaveAccessToken(chatId, accessToken)
	if err != nil {
		logger.LogIssueWithPoint("auth server - saving access token", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Make redirect to tg-bot chat forward
	w.Header().Set(LOCATION_HEADER, s.redirectUrl)
	w.WriteHeader(http.StatusMovedPermanently)
}

// Graceful shutdown
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) retainChatIdValue(r *http.Request) (int64, error) {
	id := r.URL.Query().Get(CHAT_ID_URL_PARAM)
	if id == "" {
		return 0, errors.New("Empty param: [chat_id]")
	}
	return strconv.ParseInt(id, 10, 64)
}
