package httpapi

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/alexbobkovv/insider-trades/pkg/types/cursor"
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
func (h *handler) listTransactions(w http.ResponseWriter, r *http.Request) {
	const methodName = "(h *handler) listTrades"

	queryParams := r.URL.Query()

	reqCursorStr := queryParams.Get("cursor")
	reqCursor, err := cursor.NewFromEncodedString(reqCursorStr)
	if err != nil {
		h.Error(w, r, http.StatusBadRequest, fmt.Errorf("%s: failed to parse cursor: %w", methodName, err))
		return
	}

	const defaultLimit = 20

	var reqLimit uint32
	reqLimitStr := queryParams.Get("limit")
	if reqLimitStr == "" {
		reqLimit = defaultLimit
	} else {
		limitInt, err := strconv.Atoi(reqLimitStr)
		if err != nil {
			h.Error(w, r, http.StatusBadRequest, fmt.Errorf("%s: failed to typecast limit to int: %w", methodName, err))
			return
		} else {
			reqLimit = uint32(limitInt)
		}
	}

	transactions, nextCursor, err := h.s.ListTransactions(r.Context(), reqCursor, reqLimit)
	if err != nil {
		h.Error(w, r, http.StatusInternalServerError, fmt.Errorf("%s: %w", methodName, err))
		return
	}
	if !nextCursor.IsEmpty() {
		w.Header().Add("next_cursor", nextCursor.GetEncoded())
	}

	h.Respond(w, r, http.StatusOK, transactions)
}
