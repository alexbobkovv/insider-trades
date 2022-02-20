package httpapi

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

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

func NewHandler(service service.InsiderTrade, logger *logger.Logger) (*handler, error) {
	return &handler{s: service, l: logger}, nil
}

func (h *handler) Register(router *mux.Router) http.Handler {
	// router.HandleFunc(receiverURL, h.HandlePostTrades).Methods("POST")
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

// TODO implement method
func (h *handler) HandlePostTrades(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatalf("can't read request")
	}
	log.Println(string(data))
	logFile, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Println("Can't open a file")
	}
	defer logFile.Close()
	writer := bufio.NewWriter(logFile)
	n, err := writer.WriteString(string(data))
	if err != nil {
		log.Println("Can't write in file")
	}
	log.Println("Bytes wrote:", n)
	defer writer.Flush()
	h.Respond(w, r, http.StatusCreated, nil)
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
	h.l.Info("Httpapi handler error:", err, "Status code:", statusCode, "request:", r)
	h.Respond(w, r, statusCode, nil)
}
