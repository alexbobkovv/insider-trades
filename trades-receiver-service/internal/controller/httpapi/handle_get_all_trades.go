package httpapi

import (
	"fmt"
	"net/http"
	"strconv"
)

func (h *handler) getAllTrades(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	queryParams := r.URL.Query()

	var cursor string
	cursors, present := queryParams["cursor"]
	if !present || len(cursors) != 1 {
		// TODO cursor
		cursor = ""
	} else {
		cursor = cursors[0]
	}

	var limit int
	limits, present := queryParams["limit"]
	if !present || len(limits) != 1 {
		limit = 20
	} else {
		limitInt, err := strconv.Atoi(limits[0])
		if err != nil {
			h.Error(w, r, http.StatusBadRequest, fmt.Errorf("getAllTrades handler: failed to typecast limit to int: %w", err))
		} else {
			limit = limitInt
		}
	}

	// TODO getAll
	transactions, nextCursor, err := h.s.GetAll(r.Context(), cursor, limit)
	if err != nil {
		h.Error(w, r, http.StatusInternalServerError, fmt.Errorf("getAllTrades handler: %w", err))
	}
	if nextCursor != "" {
		w.Header().Add("next_cursor", nextCursor)
	}
	h.Respond(w, r, http.StatusOK, transactions)
}
