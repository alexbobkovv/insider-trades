package httpapi

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/alexbobkovv/insider-trades/trades-receiver-service/internal/service"
	"github.com/alexbobkovv/insider-trades/trades-receiver-service/pkg/logger"

	"github.com/gorilla/mux"
)

const (
	receiverURL = "/insider-trades/receiver"
	tradesURL   = "/trades"
	rootURL     = "/"
)

type handler struct {
	s service.InsiderTrade
	l *logger.Logger
}

func NewHandler(service service.InsiderTrade, logger *logger.Logger) *handler {
	return &handler{s: service, l: logger}
}

func (h *handler) Register(router *mux.Router) http.Handler {
	router.HandleFunc(receiverURL, h.receiveTrades).Methods("POST")
	router.HandleFunc(tradesURL, h.HandleGetTrades).Methods("GET")
	router.HandleFunc(rootURL, h.HandleHomePage).Methods("GET")

	return router
}

func (h *handler) HandleHomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Under Construction")
}

func (h *handler) HandleGetTrades(w http.ResponseWriter, r *http.Request) {

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
