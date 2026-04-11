package certificates

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"time"

	"github.com/adamjames870/seacert/internal/database/sqlc"
	"github.com/adamjames870/seacert/internal/domain"
	"github.com/adamjames870/seacert/internal/dto"
	"github.com/google/generative-ai-go/genai"
)

const (
	basePrompt = `Extract certificate details from this image. 
Return strictly a single JSON object (not an array). 
If a field is not visible, return null.
Use the following JSON structure:
{
  "cert-type-name": "string",
  "cert-number": "string",
  "issuer-name": "string",
  "issued-date": "YYYY-MM-DD",
  "expiry-date": "YYYY-MM-DD",
  "remarks": "string",
  "cert-type-id": "string",
  "issuer-id": "string"
}`
	modelName = "gemini-2.5-flash" // Deprecated: using environment variable GEMINI_MODEL_NAME
)

func GenerateCacheKey(instructions string, certTypes []sqlc.CertificateType, issuers []sqlc.Issuer, modelName string) (string, error) {
	type payload struct {
		Instructions string   `json:"instructions"`
		CertTypes    []string `json:"cert_types"`
		Issuers      []string `json:"issuers"`
		Version      int      `json:"version"`
		Model        string   `json:"model"`
	}

	ctStrings := make([]string, len(certTypes))
	for i, ct := range certTypes {
		ctStrings[i] = fmt.Sprintf("%s:%s", ct.ID.String(), ct.Name)
	}
	slices.Sort(ctStrings)

	iStrings := make([]string, len(issuers))
	for i, iss := range issuers {
		iStrings[i] = fmt.Sprintf("%s:%s", iss.ID.String(), iss.Name)
	}
	slices.Sort(iStrings)

	p := payload{
		Instructions: instructions,
		CertTypes:    ctStrings,
		Issuers:      iStrings,
		Version:      3,
		Model:        modelName,
	}

	b, err := json.Marshal(p)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(b)
	return hex.EncodeToString(hash[:]), nil
}

func ExtractCertificateData(ctx context.Context, logger *slog.Logger, repo domain.Repository, client *genai.Client, modelName string, fileBytes []byte, mimeType string, certTypes []sqlc.CertificateType, issuers []sqlc.Issuer) (*dto.ExtractedCertificate, error) {
	if client == nil {
		return nil, fmt.Errorf("Gemini client is not initialized")
	}

	var promptBuilder strings.Builder
	promptBuilder.WriteString(basePrompt)
	promptBuilder.WriteString("\n\n")
	promptBuilder.WriteString("Please match the 'issuer-name' and 'cert-type-name' against these known lists. If a match is found, return the ID from the list in 'issuer-id' or 'cert-type-id' respectively.\n")

	promptBuilder.WriteString("\nKnown Certificate Types:\n")
	for _, ct := range certTypes {
		promptBuilder.WriteString(fmt.Sprintf("- %s (ID: %s)\n", ct.Name, ct.ID.String()))
	}

	promptBuilder.WriteString("\nKnown Issuers:\n")
	for _, i := range issuers {
		promptBuilder.WriteString(fmt.Sprintf("- %s (ID: %s)\n", i.Name, i.ID.String()))
	}

	fullInstructions := promptBuilder.String()

	// Count tokens for the would-be cached static prompt
	tempModel := client.GenerativeModel(modelName)
	countResp, err := tempModel.CountTokens(ctx, genai.Text(fullInstructions))
	if err != nil {
		logger.Error("Failed to count tokens for instructions", "error", err)
	}
	totalTokens := countResp.TotalTokens
	logger.Info("Instruction token count", "count", totalTokens)

	cacheKey, err := GenerateCacheKey(fullInstructions, certTypes, issuers, modelName)
	if err != nil {
		logger.Error("Failed to generate cache key", "error", err)
	}

	var cachedContentName string
	if cacheKey != "" && totalTokens >= 2048 {
		cache, err := repo.GetPromptCache(ctx, cacheKey)
		if err == nil {
			logger.Info("Using existing Gemini context cache", "cache_key", cacheKey, "gemini_cache_name", cache.GeminiCacheName)
			cachedContentName = cache.GeminiCacheName
		} else {
			logger.Info("Cache miss or expired, creating new Gemini context cache", "cache_key", cacheKey)

			// Create new cached content
			cc, err := client.CreateCachedContent(ctx, &genai.CachedContent{
				Model: modelName,
				Contents: []*genai.Content{
					{
						Parts: []genai.Part{
							genai.Text(fullInstructions),
						},
						Role: "user",
					},
				},
				Expiration: genai.ExpireTimeOrTTL{TTL: time.Hour * 24}, // Cache for 24 hours
			})
			if err != nil {
				logger.Error("Failed to create Gemini cached content", "error", err)
				// Fallback to standard request (continue without cachedContentName)
			} else {
				cachedContentName = cc.Name
				err = repo.UpsertPromptCache(ctx, sqlc.UpsertPromptCacheParams{
					CacheKey:        cacheKey,
					ModelName:       modelName,
					GeminiCacheName: cachedContentName,
					ExpiresAt:       cc.Expiration.ExpireTime,
				})
				if err != nil {
					logger.Error("Failed to save prompt cache to database", "error", err)
				}
				logger.Info("Created and saved new Gemini context cache", "cache_key", cacheKey, "gemini_cache_name", cachedContentName)
			}
		}
	}

	model := client.GenerativeModel(modelName)
	model.ResponseMIMEType = "application/json"

	var prompt []genai.Part
	if cachedContentName != "" {
		model.CachedContentName = cachedContentName
	} else {
		// Fallback or skipped cache: include instructions in the prompt
		if totalTokens < 2048 {
			logger.Info("Skipping Gemini context cache (tokens < 2048)", "tokens", totalTokens)
		}
		prompt = append(prompt, genai.Text(fullInstructions))
	}

	prompt = append(prompt, genai.Blob{
		MIMEType: mimeType,
		Data:     fileBytes,
	})

	resp, err := model.GenerateContent(ctx, prompt...)
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("no content generated by Gemini")
	}

	part := resp.Candidates[0].Content.Parts[0]
	text, ok := part.(genai.Text)
	if !ok {
		return nil, fmt.Errorf("unexpected response part type: %T", part)
	}

	var extracted dto.ExtractedCertificate
	rawText := []byte(text)

	// Attempt to unmarshal as a single object first
	if err := json.Unmarshal(rawText, &extracted); err != nil {
		// If that fails, try unmarshaling as an array of objects
		var list []dto.ExtractedCertificate
		if errArray := json.Unmarshal(rawText, &list); errArray == nil && len(list) > 0 {
			extracted = list[0]
		} else {
			return nil, fmt.Errorf("failed to unmarshal Gemini response: %w. Response: %s", err, string(text))
		}
	}

	return &extracted, nil
}
