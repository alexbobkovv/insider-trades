package httpapi

import (
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/alexbobkovv/insider-trades/trades-receiver-service/docs"
	"github.com/alexbobkovv/insider-trades/trades-receiver-service/internal/service"
	"github.com/alexbobkovv/insider-trades/trades-receiver-service/pkg/logger"
	"github.com/gorilla/mux"
)

const (
	receiverURL = "/insider-trades/receiver"
	tradesURL   = "/trades/api/v1"
	rootURL     = "/"
)

type handler struct {
	s service.InsiderTrade
	l *logger.Logger
}

func NewHandler(service service.InsiderTrade, logger *logger.Logger) *handler {
	return &handler{s: service, l: logger}
}

// Register handlers
// Comment for swaggo/swag
// @title       Insider-trades trades-receiver API
// @version     1.0
// @description Receives insider trades sec forms from external api and serves out structured trades information
// @host        localhost:8080
// @BasePath    /
// @accept json
// @produce json
// @schemes http https
func (h *handler) Register(router *mux.Router) http.Handler {

	router.Use(h.setHeadersMiddleware)
	router.HandleFunc(receiverURL, h.receiveTrades).Methods("POST", "OPTIONS")
	router.HandleFunc(tradesURL, h.getAllTransactions).Methods("GET", "OPTIONS")
	router.HandleFunc(rootURL, h.handleHomePage).Methods("GET")

	return router
}

func (h *handler) handleHomePage(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "Under Construction")
	if err != nil {
		h.Error(w, r, http.StatusInternalServerError, err)
		return
	}
}

func (h *handler) setHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Content-Type", "application/json")
		if r.Method == "OPTIONS" {
			return
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func (h *handler) Respond(w http.ResponseWriter, r *http.Request, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	if data != nil {
		err := json.NewEncoder(w).Encode(data)

		if err != nil {
			h.Error(w, r, http.StatusInternalServerError, err)
		}
	}
}

func (h *handler) Error(w http.ResponseWriter, r *http.Request, statusCode int, err error) {
	h.l.Errorf("httpapi: insider_trade handler error: %v, status code: %v, request: %v", err, statusCode, r)
	h.Respond(w, r, statusCode, nil)
}
