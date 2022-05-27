package httpapi

import (
	"fmt"
	"net/http"
	"strconv"
)

// TODO refactor
// getAllTransactions godoc
// @Summary     Get all transactions
// @Description Get all transactions objects with cursor pagination
// @Tags  	    trades
// @ID          getAllTransactions
// @Accept      json
// @Produce     json
// @Param  		cursor path string false "pagination"
// @Param  		limit path int false "limit"
// @Success     200 {object} entity.Transaction
// @Failure     404 {object} nil
// @Failure     500 {object} nil
// @Router      /trades/api/v1 [get]
func (h *handler) getAllTransactions(w http.ResponseWriter, r *http.Request) {

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
			h.Error(w, r, http.StatusBadRequest, fmt.Errorf("getAllTransactions handler: failed to typecast limit to int: %w", err))
			return
		} else {
			limit = limitInt
		}
	}

	// TODO getAll
	transactions, nextCursor, err := h.s.GetAll(r.Context(), cursor, limit)
	if err != nil {
		h.Error(w, r, http.StatusInternalServerError, fmt.Errorf("getAllTransactions handler: %w", err))
		return
	}
	if nextCursor != "" {
		w.Header().Add("next_cursor", nextCursor)
	}
	h.Respond(w, r, http.StatusOK, transactions)
}
