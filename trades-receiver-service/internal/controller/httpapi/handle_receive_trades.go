package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/alexbobkovv/insider-trades/trades-receiver-service/internal/entity"
)

type InsiderTrades struct {
	SecEntities *[]struct {
		Cik           *int    `json:"Cik"`
		Name          *string `json:"Name"`
		TradingSymbol *string `json:"TradingSymbol"`
	} `json:"SecEntities"`
	SecFilings *[]struct {
		ID          *string `json:"Id"`
		FilingURL   *string `json:"FilingUrl"`
		AccessionP1 *int    `json:"AccessionP1"`
		AccessionP2 *int    `json:"AccessionP2"`
		AccessionP3 *int    `json:"AccessionP3"`
		FilingType  *int    `json:"FilingType"`
		ReportedOn  *string `json:"ReportedOn"`
		Issuer      *int    `json:"Issuer"`
		Issuer_     *string `json:"_Issuer,omitempty"`
		Owner       *int    `json:"Owner"`
		Owner_      *string `json:"_Owner,omitempty"`
	} `json:"SecFilings"`
	HeldOfficerPositions *[]struct {
		ID            *string `json:"Id"`
		Officer       *int    `json:"Officer"`
		Company       *int    `json:"Company"`
		PositionTitle *string `json:"PositionTitle"`
		ObservedOn    *string `json:"ObservedOn"`
	} `json:"HeldOfficerPositions,omitempty"`
	SecurityTransactionHoldings *[]struct {
		ID                                *string  `json:"Id"`
		FromFiling                        *string  `json:"FromFiling"`
		FromFiling_                       *string  `json:"_FromFiling,omitempty"`
		EntryType                         *int     `json:"EntryType"`
		QuantityOwnedFollowingTransaction *float64 `json:"QuantityOwnedFollowingTransaction"`
		DirectIndirect                    *int     `json:"DirectIndirect"`
		SecurityTitle                     *string  `json:"SecurityTitle"`
		SecurityType                      *int     `json:"SecurityType"`
		AcquiredDisposed                  *int     `json:"AcquiredDisposed,omitempty"`
		Quantity                          *float64 `json:"Quantity,omitempty"`
		PricePerSecurity                  *float64 `json:"PricePerSecurity,omitempty"`
		TransactionDate                   *string  `json:"TransactionDate,omitempty"`
		TransactionCode                   *int     `json:"TransactionCode,omitempty"`
		ConversionOrExercisePrice         *float64 `json:"ConversionOrExercisePrice,omitempty"`
		ExercisableDate                   *string  `json:"ExercisableDate,omitempty"`
		ExpirationDate                    *string  `json:"ExpirationDate,omitempty"`
		UnderlyingSecurityTitle           *string  `json:"UnderlyingSecurityTitle,omitempty"`
		UnderlyingSecurityQuantity        *float64 `json:"UnderlyingSecurityQuantity,omitempty"`
	} `json:"SecurityTransactionHoldings"`
}

// TODO list transaction codes
const (
	PurchaseTransactionCode = 0
)

// TODO parse entities
func (h *handler) receiveTrades(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var trades InsiderTrades
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&trades); err != nil {
		h.l.Info("receiveTrades handler: failed to decode json to struct: ", err, "\nrequest body: ", r.Body)
		h.Respond(w, r, http.StatusCreated, nil)
		return
	}
	if trades.HeldOfficerPositions != nil {
		for _, transaction := range *trades.HeldOfficerPositions {
			if transaction.Officer != nil {
				h.l.Info("Got struct:", *transaction.Officer)

			}
		}

	}
	// TODO fill fields
	err := h.s.Receive(r.Context(), &entity.Transaction{})

	if err != nil {
		h.l.Info("receiveTrades handler: failed to receive entity: ", err, "\nrequest body: ", r.Body)
		h.Respond(w, r, http.StatusCreated, nil)
		return
	}

	h.Respond(w, r, http.StatusCreated, nil)
}
