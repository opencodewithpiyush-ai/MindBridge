package config

// PublicToInternal maps OpenAI-style model names to internal gateway names.
var PublicToInternal = map[string]string{
	"gpt-5.5":                "gateway-gpt-5-5",
	"gpt-5.4":                "gateway-gpt-5-4",
	"gpt-5.3":                "gateway-gpt-5-3",
	"gpt-5.1":                "gateway-gpt-5-1",
	"gpt-5":                  "gateway-gpt-5",
	"gpt-4o":                 "gateway-gpt-4o",
	"gpt-4o-mini":            "gateway-gpt-4o-mini",
	"grok-4":                 "gateway-grok-4",
	"claude-sonnet-4-6":      "gateway-claude-sonnet-4-6",
	"claude-opus-4-5":        "gateway-claude-opus-4-5",
	"claude-opus-4-1":        "gateway-claude-opus-4-1",
	"deepseek-v4-pro":        "gateway-deepseek-v4-pro",
	"deepseek-v4-flash":      "gateway-deepseek-v4-flash",
	"deepseek-r1":            "gateway-deepseek-r1",
	"gemini-3.1-pro":         "gateway-gemini-3-1-pro",
	"gemini-3-pro":           "gateway-gemini-3-pro",
	"gemini-2.5-flash":       "gateway-gemini-2.5-flash",
	"qwen-3-max":             "gateway-qwen-3-max",
	"llama-3.3-70b-versatile":"gateway-llama-3-3-70b-versatile",
	"kimi-k2":                "gateway-deepinfra-kimi-k2",
}

// InternalToPublic maps internal gateway names back to OpenAI-style names.
var InternalToPublic = map[string]string{
	"gateway-gpt-5-5":                "gpt-5.5",
	"gateway-gpt-5-4":                "gpt-5.4",
	"gateway-gpt-5-3":                "gpt-5.3",
	"gateway-gpt-5-1":                "gpt-5.1",
	"gateway-gpt-5":                  "gpt-5",
	"gateway-gpt-4o":                 "gpt-4o",
	"gateway-gpt-4o-mini":            "gpt-4o-mini",
	"gateway-grok-4":                 "grok-4",
	"gateway-claude-sonnet-4-6":      "claude-sonnet-4-6",
	"gateway-claude-opus-4-5":        "claude-opus-4-5",
	"gateway-claude-opus-4-1":        "claude-opus-4-1",
	"gateway-deepseek-v4-pro":        "deepseek-v4-pro",
	"gateway-deepseek-v4-flash":      "deepseek-v4-flash",
	"gateway-deepseek-r1":            "deepseek-r1",
	"gateway-gemini-3-1-pro":         "gemini-3.1-pro",
	"gateway-gemini-3-pro":           "gemini-3-pro",
	"gateway-gemini-2.5-flash":       "gemini-2.5-flash",
	"gateway-qwen-3-max":             "qwen-3-max",
	"gateway-llama-3-3-70b-versatile":"llama-3.3-70b-versatile",
	"gateway-deepinfra-kimi-k2":      "kimi-k2",
}

// ResolveModel converts a public model name to internal name if needed.
func ResolveModel(name string) string {
	if name == "" {
		return "gateway-claude-opus-4-1"
	}
	if internal, ok := PublicToInternal[name]; ok {
		return internal
	}
	return name
}
