package service

import (
	"context"
	"encoding/json"
	"og_Analyzer/internal/models"

	"github.com/sashabaranov/go-openai"
)

type AnalyZer struct {
	client *openai.Client
}

func NewAnalyZer(apiKey string) *AnalyZer {
	// Groq uses the OpenAI format but a different BaseURL
	config := openai.DefaultConfig(apiKey)
	config.BaseURL = "https://api.groq.com/openai/v1"

	client := openai.NewClientWithConfig(config)
	return &AnalyZer{client: client}
}

func (s *AnalyZer) AnalyzeError(ctx context.Context, logMessage string) (models.AIAnalysis, error) {
	// Use 'llama-3.3-70b-versatile' or 'mixtral-8x7b-32768' - both are fast and free on Groq
	resp, err := s.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: "llama-3.3-70b-versatile",
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: `You are an expert SRE and DevOps incident analyzer. Return ONLY valid JSON. Schema: {"rootCause": "string",  "impact": string,  "fix": string,  "severity": "LOW" | "MEDIUM" | "HIGH" | "CRITICAL","confidence": number}	Rules:- confidence must be integer from 0-100
,- do not include markdown,- do not include explanations,- output raw JSON only		`,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: logMessage,
				}	,
			},
		},
	)

	if err != nil {
		return models.AIAnalysis{}, err
	}

	contentt:=resp.Choices[0].Message.Content 
	var Analysis models.AIAnalysis 
			err=json.Unmarshal([]byte(contentt),&Analysis)
			if err!=nil{
				return models.AIAnalysis{},err
			}

			return Analysis,nil
}
