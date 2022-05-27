package httpapi

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/alexbobkovv/insider-trades/pkg/types/cursor"
)

func (h *handler) listTradeViews(w http.ResponseWriter, r *http.Request) {
	const methodName = "(h *handler) listTradeViews"

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

	tradeViews, nextCursor, err := h.s.ListTrades(r.Context(), reqCursor, reqLimit)
	if err != nil {
		h.Error(w, r, http.StatusInternalServerError, fmt.Errorf("%s: %w", methodName, err))
		return
	}
	if !nextCursor.IsEmpty() {
		w.Header().Add("next_cursor", nextCursor.GetEncoded())
	}

	h.Respond(w, r, http.StatusOK, tradeViews)
}
