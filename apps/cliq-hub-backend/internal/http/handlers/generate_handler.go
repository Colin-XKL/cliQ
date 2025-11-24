package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"cliq-hub-backend/internal/errors"
	"cliq-hub-backend/internal/llm"
	"cliq-hub-backend/internal/version"
	"repo/cliqfile"
)

type GenerateHandler struct {
	client    llm.Client
	debugMode bool
}

func NewGenerateHandler(c llm.Client, debugMode bool) *GenerateHandler {
	return &GenerateHandler{client: c, debugMode: debugMode}
}

type GenerateRequest struct {
	CommandExample string `json:"command_example" binding:"required"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	Author         string `json:"author"`
	Encoding       string `json:"encoding"` // "plain" or "base64"
}

type GenerateResponse struct {
	YAML     string            `json:"yaml"`
	Encoding string            `json:"encoding"`
	Meta     map[string]string `json:"meta"`
}

func (h *GenerateHandler) Handle(c *gin.Context) {
	var req GenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.New("invalid_input", "invalid JSON or missing fields"))
		return
	}
	enc := strings.ToLower(strings.TrimSpace(req.Encoding))
	if enc == "" {
		enc = "plain"
	}
	if enc != "plain" && enc != "base64" {
		c.JSON(http.StatusBadRequest, errors.New("invalid_input", "encoding must be 'plain' or 'base64'"))
		return
	}

	content, err := h.client.GenerateCliqfileFromPrompt(c.Request.Context(), llm.GenerateRequest{
		CommandExample: req.CommandExample,
		Name:           req.Name,
		Description:    req.Description,
		Author:         req.Author,
	})
	if err != nil {
		errResp := errors.New("llm_error", err.Error())
		if h.debugMode {
			errResp = errResp.WithMeta("llm_request", req)
		}
		c.JSON(http.StatusBadGateway, errResp)
		return
	}

	raw := cliqfile.StripFences(content)
	t, err := cliqfile.Parse([]byte(raw))
	if err != nil {
		errResp := errors.New("llm_output_invalid", "failed to parse YAML from LLM")
		if h.debugMode {
			errResp = errResp.WithMeta("llm_request", req).WithMeta("llm_output", raw)
		}
		c.JSON(http.StatusBadGateway, errResp)
		return
	}

	// ensure defaults present if LLM omitted
	if t.Version == "" {
		t.Version = version.TemplateVersion
	}
	if t.CliqTemplateVersion == "" {
		t.CliqTemplateVersion = version.CliqTemplateSpecVersion
	}
	if t.Author == "" && req.Author != "" {
		t.Author = req.Author
	}
	if t.Name == "" && req.Name != "" {
		t.Name = req.Name
	}
	if t.Description == "" && req.Description != "" {
		t.Description = req.Description
	}

	// Re-marshal to validate logic on the final object
	out, err := cliqfile.GenerateYAML(t)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.New("marshal_error", err.Error()))
		return
	}

	if validationErrors, err := cliqfile.Validate([]byte(out)); err != nil || len(validationErrors) > 0 {
		var errMsg string
		if err != nil {
			errMsg = err.Error()
		} else {
			var msgBuilder strings.Builder
			for _, ve := range validationErrors {
				msgBuilder.WriteString(fmt.Sprintf("%s: %s; ", ve.Field, ve.Message))
			}
			errMsg = msgBuilder.String()
		}

		if h.debugMode {
			log.Printf("Validation Error: %s", errMsg)
		}
		errResp := errors.New("validation_error", errMsg)
		if h.debugMode {
			errResp = errResp.WithMeta("llm_request", req).WithMeta("llm_output", content)
		}
		c.JSON(http.StatusUnprocessableEntity, errResp)
		return
	}

	resp := GenerateResponse{
		YAML:     out,
		Encoding: enc,
		Meta: map[string]string{
			"name":                  t.Name,
			"version":               t.Version,
			"cliq_template_version": t.CliqTemplateVersion,
		},
	}
	if enc == "base64" {
		resp.YAML = cliqfile.Base64Encode(out)
	}
	c.JSON(http.StatusOK, resp)
}
