package models

type AIAnalysis struct {
	RootCause string `json:"rootCause"`
    Impact string `json:"impact"`
    Fix string `json:"fix"`
    Severity string `json:"severity"`
    Confidence int `json:"confidence"`
}