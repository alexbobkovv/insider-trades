package entity

type SecFiling struct {
	ID              string
	FilingType      *int
	URL             string
	InsiderID       string
	CompanyID       string
	OfficerPosition *string
	ReportedOn      string
}
