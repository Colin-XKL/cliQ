package config

import (
    "errors"
    "os"
)

type Config struct {
    LLMBaseURL string
    LLMAPIKey  string
    LLMModel   string
}

func Load() (*Config, error) {
    cfg := &Config{
        LLMBaseURL: os.Getenv("LLM_BASE_URL"),
        LLMAPIKey:  os.Getenv("LLM_API_KEY"),
        LLMModel:   os.Getenv("LLM_MODEL"),
    }
    if cfg.LLMAPIKey == "" {
        return nil, errors.New("LLM_API_KEY is required")
    }
    if cfg.LLMModel == "" {
        return nil, errors.New("LLM_MODEL is required")
    }
    return cfg, nil
}