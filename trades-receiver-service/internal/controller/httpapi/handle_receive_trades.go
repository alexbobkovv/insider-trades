package httpapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/alexbobkovv/insider-trades/trades-receiver-service/internal/entity"
	"github.com/go-playground/validator"
)

type (
	SecEntity struct {
		Cik           *int    `json:"Cik" validate:"required"`
		Name          *string `json:"Name" validate:"required"`
		TradingSymbol *string `json:"TradingSymbol"`
	}

	SecFilings struct {
		ID          *string `json:"Id"`
		FilingURL   *string `json:"FilingUrl" validate:"required,url"`
		AccessionP1 *int    `json:"AccessionP1"`
		AccessionP2 *int    `json:"AccessionP2"`
		AccessionP3 *int    `json:"AccessionP3"`
		FilingType  *int    `json:"FilingType"`
		ReportedOn  *string `json:"ReportedOn" validate:"required"`
		Issuer      *int    `json:"Issuer" validate:"required"`
		Issuer_     *string `json:"_Issuer"`
		Owner       *int    `json:"Owner" validate:"required"`
		Owner_      *string `json:"_Owner"`
	}

	HeldOfficerPosition struct {
		ID            *string `json:"Id"`
		Officer       *int    `json:"Officer" validate:"required"`
		Company       *int    `json:"Company" validate:"required"`
		PositionTitle *string `json:"PositionTitle" validate:"required"`
		ObservedOn    *string `json:"ObservedOn"`
	}

	SecurityTransactionHolding struct {
		ID                                *string  `json:"Id"`
		FromFiling                        *string  `json:"FromFiling"`
		FromFiling_                       *string  `json:"_FromFiling"`
		EntryType                         *int     `json:"EntryType"`
		QuantityOwnedFollowingTransaction *float64 `json:"QuantityOwnedFollowingTransaction"`
		DirectIndirect                    *int     `json:"DirectIndirect"`
		SecurityTitle                     *string  `json:"SecurityTitle"`
		SecurityType                      *int     `json:"SecurityType"`
		AcquiredDisposed                  *int     `json:"AcquiredDisposed"`
		Quantity                          *float64 `json:"Quantity"`
		PricePerSecurity                  *float64 `json:"PricePerSecurity"`
		TransactionDate                   *string  `json:"TransactionDate"`
		TransactionCode                   *int     `json:"TransactionCode"`
		ConversionOrExercisePrice         *float64 `json:"ConversionOrExercisePrice"`
		ExercisableDate                   *string  `json:"ExercisableDate"`
		ExpirationDate                    *string  `json:"ExpirationDate"`
		UnderlyingSecurityTitle           *string  `json:"UnderlyingSecurityTitle"`
		UnderlyingSecurityQuantity        *float64 `json:"UnderlyingSecurityQuantity"`
	}

	InsiderTrades struct {
		SecEntities                 *[]SecEntity                  `json:"SecEntities" validate:"dive,required"`
		SecFilings                  *[]SecFilings                 `json:"SecFiling" validate:"dive,required"`
		HeldOfficerPositions        *[]HeldOfficerPosition        `json:"HeldOfficerPositions"`
		SecurityTransactionHoldings *[]SecurityTransactionHolding `json:"SecurityTransactionHoldings" validate:"dive,required"`
	}
)

func (h *handler) validateTrades(trades *InsiderTrades) error {
	if trades == nil {
		return errors.New("validateTrades: expected non nil trades struct")
	}

	v := validator.New()
	err := v.Struct(trades)

	if err != nil {
		// Validation syntax error
		if err, ok := err.(*validator.InvalidValidationError); ok {
			return fmt.Errorf("validateTrades: validation syntax error: %w", err)
		}

		validationErrors := make([]string, 0)

		for _, validationError := range err.(validator.ValidationErrors) {

			switch validationError.Tag() {
			case "required":
				validationErrors = append(validationErrors, fmt.Sprintf("%s is a required field", validationError.Field()))
			case "max":
				validationErrors = append(validationErrors, fmt.Sprintf("%s must be a maximum of %s in length", validationError.Field(), validationError.Param()))
			case "url":
				validationErrors = append(validationErrors, fmt.Sprintf("%s must be a valid URL", validationError.Field()))
			default:
				validationErrors = append(validationErrors, fmt.Sprintf("something wrong on %s; %s", validationError.Field(), validationError.Tag()))
			}
		}

		return fmt.Errorf("validateTrades: validation errors: %v", strings.Join(validationErrors, ", "))

	}

	if len(*trades.SecFilings) != 1 {
		return errors.New("validateTrades: failed to validate sec filing: empty or more than one filing")
	}

	return nil
}

func (h *handler) logRequestBody(r *http.Request) {
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.l.Errorf("error reading request body: %v", err)
		return
	}
	h.l.Info("logRequestBody: %v", string(buf))
	reader := ioutil.NopCloser(bytes.NewBuffer(buf))
	r.Body = reader
}

// receiveTrades godoc
// @Summary     receiveTrades from external api
// @Description receiveTrades from external api
// @Tags  	    trades
// @ID          receiveTrades
// @Accept      json
// @Produce     json
// @Param  		request body InsiderTrades true "Insider trades request"
// @Success     201 {object} nil
// @Failure     500 {object} nil
// @Router      /insider-trades/receiver [post]
func (h *handler) receiveTrades(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var trades InsiderTrades
	defer func() {
		if err := r.Body.Close(); err != nil {
			h.l.Errorf("receiveTrades: failed to close body: %v", err)
		}
	}()

	h.logRequestBody(r)

	if err := json.NewDecoder(r.Body).Decode(&trades); err != nil {
		h.l.Info("receiveTrades handler: failed to decode json to struct: ", err, "\nrequest body: ", r.Body)
		h.Respond(w, r, http.StatusCreated, nil)
		return
	}

	if err := h.validateTrades(&trades); err != nil {
		h.l.Errorf("receiveTrades handler: failed to validate trades struct: %v", err)
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

	sFiling := &(*trades.SecFilings)[0]

	trade.SecF = h.fillSecFiling(sFiling, trades.HeldOfficerPositions)

	var err error
	trade.Sth, err = h.fillSecurityTransactionHoldings(trades.SecurityTransactionHoldings)
	if err != nil {
		h.l.Errorf("receiveTrades: %v", err)
	}

	if err := h.s.Receive(r.Context(), &trade); err != nil {
		h.l.Error("receiveTrades: ", err, " request body: ", r.Body)
		h.Respond(w, r, http.StatusCreated, nil)
		return
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

func (h *handler) fillSecFiling(sFiling *SecFilings, positions *[]HeldOfficerPosition) *entity.SecFiling {
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

func (h *handler) fillSecurityTransactionHoldings(holdings *[]SecurityTransactionHolding) ([]*entity.SecurityTransactionHoldings, error) {
	var holdingEntities []*entity.SecurityTransactionHoldings

	for _, holding := range *holdings {
		if holding.SecurityTitle == nil || holding.Quantity == nil ||
			holding.PricePerSecurity == nil || holding.TransactionDate == nil ||
			holding.TransactionCode == nil {
			continue
		}
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

	if len(holdingEntities) == 0 {
		return nil, errors.New("fillSecurityTransactionHoldings: expected more than 0 holding entities")
	}

	return holdingEntities, nil
}
