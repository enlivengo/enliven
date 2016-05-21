package enliven

// Config represents string kvps of application configuration
type Config map[string]string

// MergeConfig takes a default config and merges a supplied one into it.
func MergeConfig(defaultConfig Config, suppliedConfig Config) Config {
	for key, value := range suppliedConfig {
		defaultConfig[key] = value
	}

	return defaultConfig
}
