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

	viewsCache, nextCur, err := h.cache.ListTrades(r.Context(), reqCursor, reqLimit)
	if err == nil && len(viewsCache) != 0 {

		if !nextCur.IsEmpty() {
			w.Header().Add("X-next-cursor", nextCur.GetEncoded())
		}

		h.Respond(w, r, http.StatusOK, viewsCache)
		return
	}

	tradeViews, nextCursor, err := h.s.ListTrades(r.Context(), reqCursor, reqLimit)
	if err != nil {
		h.Error(w, r, http.StatusInternalServerError, fmt.Errorf("%s: %w", methodName, err))
		return
	}
	if !nextCursor.IsEmpty() {
		w.Header().Add("X-next-cursor", nextCursor.GetEncoded())
	}

	h.Respond(w, r, http.StatusOK, tradeViews)
}
