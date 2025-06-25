package models

type Metadata struct {
	ID                 string              `json:"id"`
	Accused            string              `json:"accused"`
	IDNumber           string              `json:"id_number"`
	Gender             string              `json:"gender"`
	Nationality        string              `json:"nationality"`
	State              string              `json:"state"`
	Category           string              `json:"category"`
	Employer           string              `json:"employer"`
	Position           string              `json:"position"`
	Court              string              `json:"court"`
	Judge              string              `json:"judge"`
	Officer            string              `json:"officer"`
	DefenseAttorney    string              `json:"defense_attorney"`
	PastConvictions    string              `json:"past_convictions"`
	SentencingDate     string              `json:"sentencing_date"`
	Appeal             string              `json:"appeal"`
	ProcurementDetails []ProcurementDetail `json:"procurement_details"`
}

type ProcurementDetail struct {
	Number      string `json:"number"`
	Summary     string `json:"summary"`
	Offenses    string `json:"offenses"`
	Punishments string `json:"punishments"`
}
