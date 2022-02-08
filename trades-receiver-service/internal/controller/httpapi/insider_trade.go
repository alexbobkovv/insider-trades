package httpapi

import "insidertradesreceiver/internal/service"

type handler struct {
	service service.InsiderTrade
}

func NewHandler(s service.InsiderTrade) *handler {
	return &handler{service: s}
}