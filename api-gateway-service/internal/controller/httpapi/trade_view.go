package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/alexbobkovv/insider-trades/api-gateway-service/config"
	"github.com/alexbobkovv/insider-trades/api-gateway-service/internal/cache"
	"github.com/alexbobkovv/insider-trades/api-gateway-service/internal/service"
	"github.com/alexbobkovv/insider-trades/pkg/logger"
	"github.com/gorilla/mux"
)

const (
	tradeViewsURL = "/api-gateway/v1/trade-views"
)

type handler struct {
	s     service.Gateway
	l     *logger.Logger
	cfg   *config.Config
	cache *cache.TradeCache
}

func NewHandler(service service.Gateway, logger *logger.Logger, config *config.Config, tradeCache *cache.TradeCache) *handler {
	return &handler{s: service, l: logger, cfg: config, cache: tradeCache}
}

// Register handlers
// Comment for swaggo/swag
// @title       Insider-trades api-gateway-service API
// @version     1.0
// @description Provides trade views for frontend
// @host        localhost:8082
// @BasePath    /
// @accept json
// @produce json
// @schemes http https
func (h *handler) Register(router *mux.Router) http.Handler {

	router.Use(h.setHeadersMiddleware)
	router.HandleFunc(tradeViewsURL, h.listTradeViews).Methods("GET", "OPTIONS")

	return router
}

func (h *handler) setHeadersMiddleware(next http.Handler) http.Handler {
	const methodName = "(h *handler) setHeadersMiddleware"

	if h.cfg.HTTPServer.AllowOrigin == "" {
		h.l.Fatalf("%s: please specify CORS allow origin feild in config.yml file(\"allow_origin: '*'\" to allow all hosts(not recomended))", methodName)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", h.cfg.HTTPServer.AllowOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Expose-Headers", "X-next-cursor")
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
