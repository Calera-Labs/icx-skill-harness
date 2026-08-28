package byok

import "strings"

// NewProvider constructs the appropriate LLMClient based on provider name
func NewProvider(providerName, apiKey, model, baseURL string) LLMClient {
	p := strings.ToLower(strings.TrimSpace(providerName))
	switch p {
	case "anthropic", "claude":
		if model == "" {
			model = "claude-3-5-sonnet-20241022"
		}
		return NewAnthropicClient(apiKey, model, baseURL)

	case "openai", "deepseek", "groq", "ollama", "vllm":
		if baseURL == "" {
			switch p {
			case "deepseek":
				baseURL = "https://api.deepseek.com/v1"
				if model == "" {
					model = "deepseek-chat"
				}
			case "groq":
				baseURL = "https://api.groq.com/openai/v1"
				if model == "" {
					model = "llama-3.3-70b-versatile"
				}
			case "ollama":
				baseURL = "http://localhost:11434/v1"
				if model == "" {
					model = "llama3"
				}
			default:
				baseURL = "https://api.openai.com/v1"
				if model == "" {
					model = "gpt-4o-mini"
				}
			}
		}
		return NewOpenAIClient(apiKey, model, baseURL)

	case "gemini", "google":
		fallthrough
	default:
		if model == "" {
			model = "gemini-3.5-flash-lite"
		}
		return NewGeminiClient(apiKey, model)
	}
}
