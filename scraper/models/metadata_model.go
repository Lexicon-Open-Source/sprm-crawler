package models

type Metadata struct {
	ID                 string              `json:"id"`
	Title              string              `json:"title"`
	Address            string              `json:"address"`
	City               string              `json:"city"`
	Status             string              `json:"status"`
	NPWP               string              `json:"npwp"`
	Province           string              `json:"province"`
	Injunction         Injunction          `json:"injunction"`
	ProcurementDetails []ProcurementDetail `json:"procurement_details"`
}

type Injunction struct {
	Number        string `json:"number"`
	Rule          string `json:"rule"`
	Description   string `json:"description"`
	StartDate     string `json:"start_date"`
	EndDate       string `json:"end_date"`
	PublishedDate string `json:"published_date"`
	Duration      string `json:"duration"`
}

type ProcurementDetail struct {
	TenderID        string `json:"tender_id"`
	PackageName     string `json:"package_name"`
	ProcurementType string `json:"procurement_type"`
	Ceiling         string `json:"ceiling"`
	Unit            string `json:"unit"`
	FiscalYear      string `json:"fiscal_year"`
	InstitutionArea string `json:"institution_area"`
	EstimatedPrice  string `json:"estimated_price"`
}
