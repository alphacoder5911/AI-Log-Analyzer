package models

type LogEntry struct {
    ID        string `json:"id"`
    Service   string `json:"service"`
    Level     string `json:"level"` // INFO, ERROR, CRITICAL
    Message   string `json:"message"`
    Timestamp string `json:"timestamp"`
}