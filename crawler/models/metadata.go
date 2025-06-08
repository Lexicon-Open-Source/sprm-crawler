package models

type UrlFrontierMetadata struct {
	Title         string `json:"title"`
	Scenario      string `json:"scenario"`
	PackageNumber string `json:"package_number"`
	PackageName   string `json:"package_name"`
	StartDate     string `json:"start_date"`
	EndDate       string `json:"end_date"`
	Duration      string `json:"duration"`
	Status        string `json:"status"`
}
