package httpapi

import (
	"encoding/json"
	"insidertradesreceiver/internal/entity"
	"net/http"
)

type InsiderTrades struct {
	SecEntities []struct {
		Cik           int    `json:"Cik"`
		Name          string `json:"Name"`
		TradingSymbol string `json:"TradingSymbol"`
	} `json:"SecEntities"`
	SecFilings []struct {
		ID          string `json:"Id"`
		FilingURL   string `json:"FilingUrl"`
		AccessionP1 int    `json:"AccessionP1"`
		AccessionP2 int    `json:"AccessionP2"`
		AccessionP3 int    `json:"AccessionP3"`
		FilingType  int    `json:"FilingType"`
		ReportedOn  string `json:"ReportedOn"`
		Issuer      int    `json:"Issuer"`
		_Issuer     string `json:"_Issuer,omitempty"`
		Owner       int    `json:"Owner"`
		_Owner      string `json:"_Owner,omitempty"`
	} `json:"SecFilings"`
	HeldOfficerPositions []struct {
		ID            string `json:"Id"`
		Officer       string `json:"Officer"`
		Company       string `json:"Company"`
		PositionTitle string `json:"PositionTitle"`
		ObservedOn    string `json:"ObservedOn"`
	} `json:"HeldOfficerPositions,omitempty"`
	SecurityTransactionHoldings []struct {
		ID                                string  `json:"Id"`
		FromFiling                        string  `json:"FromFiling"`
		_FromFiling                       string  `json:"_FromFiling,omitempty"`
		EntryType                         int     `json:"EntryType"`
		QuantityOwnedFollowingTransaction float64 `json:"QuantityOwnedFollowingTransaction"`
		DirectIndirect                    int     `json:"DirectIndirect"`
		SecurityTitle                     string  `json:"SecurityTitle"`
		SecurityType                      int     `json:"SecurityType"`
		AcquiredDisposed                  int     `json:"AcquiredDisposed,omitempty"`
		Quantity                          float64 `json:"Quantity,omitempty"`
		PricePerSecurity                  float64 `json:"PricePerSecurity,omitempty"`
		TransactionDate                   string  `json:"TransactionDate,omitempty"`
		TransactionCode                   int     `json:"TransactionCode,omitempty"`
		ConversionOrExercisePrice         float64 `json:"ConversionOrExercisePrice,omitempty"`
		ExercisableDate                   string  `json:"ExercisableDate,omitempty"`
		ExpirationDate                    string  `json:"ExpirationDate,omitempty"`
		UnderlyingSecurityTitle           string  `json:"UnderlyingSecurityTitle,omitempty"`
		UnderlyingSecurityQuantity        int     `json:"UnderlyingSecurityQuantity,omitempty"`
	} `json:"SecurityTransactionHoldings"`
}

func (h *handler) receiveTrades(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var trades InsiderTrades
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&trades); err != nil {
		h.Error(w, r, http.StatusBadRequest, err)
		return
	}

	h.s.Receive(r.Context(), &entity.InsiderTrade{})
}
