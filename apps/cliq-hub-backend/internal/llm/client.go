package llm

import (
	"context"
	_ "embed"
	"fmt"
	"strings"
	"text/template"
	"time"

	openai "github.com/sashabaranov/go-openai"

	"cliq-hub-backend/internal/config"
)

const userPromptTemplate = `
Given a CLI command example and optional metadata, generate a complete cliqfile YAML.

Requirements:
- Fields: name, description, version ("1.0"), author, cliq_template_version ("1.0"), cmds (with name, description, command, variables).
- Return RAW YAML ONLY (no code fences, no extra text).

Input:
command_example: {{.CommandExample}}
name: {{.Name}}
description: {{.Description}}
author: {{.Author}}
`

//go:embed assets/cliqfile_syntax.md
var cliqfileSyntaxDoc string

type Client interface {
	GenerateCliqfileWithRetry(ctx context.Context, req GenerateRequest, maxRounds int, validator func(string) error) (string, error)
}

type client struct {
	oa    *openai.Client
	model string
}

type GenerateRequest struct {
	CommandExample string
	Name           string
	Description    string
	Author         string
}

func NewClient(cfg *config.Config) (Client, error) {
	conf := openai.DefaultConfig(cfg.LLMAPIKey)
	if cfg.LLMBaseURL != "" {
		conf.BaseURL = cfg.LLMBaseURL
	}
	c := openai.NewClientWithConfig(conf)
	return &client{oa: c, model: cfg.LLMModel}, nil
}

func (c *client) GenerateCliqfileWithRetry(ctx context.Context, req GenerateRequest, maxRounds int, validator func(string) error) (string, error) {
	// Define the system prompt that includes the CLIQ file syntax documentation
	systemPrompt := fmt.Sprintf("You generate ONLY valid cliqfile YAML per schema. No prose. No markdown fences.\n\nCLIQFILE SYNTAX DOCUMENTATION:\n%s", cliqfileSyntaxDoc)

	// Parse the template and execute it with the request data
	tmpl, err := template.New("userPrompt").Parse(userPromptTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse user prompt template: %w", err)
	}

	var userPromptBuilder strings.Builder
	if err := tmpl.Execute(&userPromptBuilder, req); err != nil {
		return "", fmt.Errorf("failed to execute user prompt template: %w", err)
	}
	userPrompt := userPromptBuilder.String()

	messages := []openai.ChatCompletionMessage{
		{Role: openai.ChatMessageRoleSystem, Content: systemPrompt},
		{Role: openai.ChatMessageRoleUser, Content: userPrompt},
	}

	var lastContent string

	if maxRounds < 1 {
		maxRounds = 1
	}

	for i := 0; i < maxRounds; i++ {
		// Use a per-request timeout
		reqCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
		resp, err := c.oa.CreateChatCompletion(reqCtx, openai.ChatCompletionRequest{
			Model:       c.model,
			Messages:    messages,
			Temperature: 0,
		})
		cancel()

		if err != nil {
			return "", err
		}
		if len(resp.Choices) == 0 {
			return "", fmt.Errorf("empty LLM response")
		}
		content := resp.Choices[0].Message.Content
		lastContent = content

		// Validate
		if validator != nil {
			if err := validator(content); err != nil {
				// Prepare for next round
				messages = append(messages, openai.ChatCompletionMessage{
					Role:    openai.ChatMessageRoleAssistant,
					Content: content,
				})
				messages = append(messages, openai.ChatCompletionMessage{
					Role:    openai.ChatMessageRoleUser,
					Content: fmt.Sprintf("The previous output was invalid. Error: %s. Please fix it and return the complete valid YAML.", err.Error()),
				})
				continue
			}
		}

		return content, nil
	}

	// If we reached here, we ran out of rounds
	// We return the last content anyway, so the caller can validate it and report specific errors.
	return lastContent, nil
}
