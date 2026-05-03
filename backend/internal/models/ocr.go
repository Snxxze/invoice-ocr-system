package models

type OCRSummary struct {
	Rental      string `json:"rental"`
	Electricity string `json:"electricity"`
	Water       string `json:"water"`
	CableTV     string `json:"cable_tv"`
	Total       string `json:"total"`
}

type OCRBox struct {
	Text       string      `json:"text"`
	Confidence float64     `json:"confidence"`
	Box        [][]float64 `json:"box"`
}

type OCRResponse struct {
	Summary OCRSummary `json:"summary"`
	Data    []OCRBox   `json:"data"`
}
