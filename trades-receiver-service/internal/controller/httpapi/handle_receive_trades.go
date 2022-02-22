package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/alexbobkovv/insider-trades/trades-receiver-service/internal/entity"
)

type (
	SecEntity struct {
		Cik           *int    `json:"Cik"`
		Name          *string `json:"Name"`
		TradingSymbol *string `json:"TradingSymbol"`
	}

	SecFiling struct {
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
	}

	HeldOfficerPosition struct {
		ID            *string `json:"Id"`
		Officer       *int    `json:"Officer"`
		Company       *int    `json:"Company"`
		PositionTitle *string `json:"PositionTitle"`
		ObservedOn    *string `json:"ObservedOn"`
	}

	SecurityTransactionHolding struct {
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
	}

	InsiderTrades struct {
		SecEntities                 *[]SecEntity                  `json:"SecEntities"`
		SecFilings                  *[]SecFiling                  `json:"SecFilings"`
		HeldOfficerPositions        *[]HeldOfficerPosition        `json:"HeldOfficerPositions,omitempty"`
		SecurityTransactionHoldings *[]SecurityTransactionHolding `json:"SecurityTransactionHoldings"`
	}
)

func (h *handler) receiveTrades(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var trades InsiderTrades
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&trades); err != nil {
		h.l.Info("receiveTrades handler: failed to decode json to struct: ", err, "\nrequest body: ", r.Body)
		h.Respond(w, r, http.StatusCreated, nil)
		return
	}

	var trade entity.Trade

	for _, sEntity := range *trades.SecEntities {
		if sEntity.TradingSymbol == nil {
			trade.Ins = h.fillInsider(&sEntity)

		} else {
			trade.Cmp = h.fillCompany(&sEntity)

		}
	}

	if len(*trades.SecFilings) != 1 {
		h.l.Info("receiveTrades handler: failed to receive sec filing: ", "empty or more than one filing, ", "request body: ", r.Body)
		h.Respond(w, r, http.StatusCreated, nil)
		return
	}

	sFiling := &(*trades.SecFilings)[0]

	trade.SecF = h.fillSecFiling(sFiling, trades.HeldOfficerPositions)
	trade.Sth = h.fillSecurityTransactionHoldings(trades.SecurityTransactionHoldings)

	if err := h.s.Receive(r.Context(), &trade); err != nil {
		h.l.Error("service receive, failed to receive trade: ", err, " request body: ", r.Body)
	}

	h.Respond(w, r, http.StatusCreated, nil)
}

func (h *handler) fillInsider(sEntity *SecEntity) *entity.Insider {
	return &entity.Insider{
		Cik:  *sEntity.Cik,
		Name: *sEntity.Name,
	}
}

func (h *handler) fillCompany(sEntity *SecEntity) *entity.Company {
	return &entity.Company{
		Cik:    *sEntity.Cik,
		Name:   *sEntity.Name,
		Ticker: *sEntity.TradingSymbol,
	}
}

func (h *handler) fillSecFiling(sFiling *SecFiling, positions *[]HeldOfficerPosition) *entity.SecFiling {
	var officerPositions []string
	var positionString string

	if len(*positions) == 0 {
		positionString = "10% owner"
	} else {
		for _, pos := range *positions {
			officerPositions = append(officerPositions, *pos.PositionTitle)
		}
		positionString = strings.Join(officerPositions, ", ")
	}

	return &entity.SecFiling{
		FilingType:      sFiling.FilingType,
		URL:             *sFiling.FilingURL,
		OfficerPosition: &positionString,
		ReportedOn:      *sFiling.ReportedOn,
	}
}

func (h *handler) fillSecurityTransactionHoldings(holdings *[]SecurityTransactionHolding) []*entity.SecurityTransactionHoldings {
	var holdingEntities []*entity.SecurityTransactionHoldings

	for _, holding := range *holdings {
		holdingEntities = append(holdingEntities, &entity.SecurityTransactionHoldings{
			QuantityOwnedFollowingTransaction: holding.QuantityOwnedFollowingTransaction,
			SecurityTitle:                     *holding.SecurityTitle,
			SecurityType:                      holding.SecurityType,
			Quantity:                          int(*holding.Quantity),
			PricePerSecurity:                  *holding.PricePerSecurity,
			TransactionDate:                   *holding.TransactionDate,
			TransactionCode:                   *holding.TransactionCode,
		})

	}

	return holdingEntities
}
