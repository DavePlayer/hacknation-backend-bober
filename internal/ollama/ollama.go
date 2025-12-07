package ollama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type OllamaGenerateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type OllamaGenerateResponse struct {
	Model      string `json:"model"`
	CreatedAt  string `json:"created_at"`
	Response   string `json:"response"`
	Done       bool   `json:"done"`
	DoneReason string `json:"done_reason"`
	// reszta pól nie jest konieczna, możesz dodać jak chcesz:
	// Context        []int  `json:"context"`
	// TotalDuration  int64  `json:"total_duration"`
	// ...
}

func CallOllama(prompt string) (string, error) {
	reqBody := OllamaGenerateRequest{
		Model:  "llama3.1:8b", // albo co tam masz
		Prompt: prompt,
		Stream: false,
	}

	b, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewReader(b))
	if err != nil {
		return "", fmt.Errorf("http post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ollama status %d: %s", resp.StatusCode, string(body))
	}

	var og OllamaGenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&og); err != nil {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("decode response: %w; body=%s", err, string(body))
	}

	return og.Response, nil
}
