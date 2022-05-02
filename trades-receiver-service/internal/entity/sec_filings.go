package entity

type SecFiling struct {
	ID              string
	FilingType      *int64
	URL             string
	InsiderID       string
	CompanyID       string
	OfficerPosition *string
	ReportedOn      string
}
