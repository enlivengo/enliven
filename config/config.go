package config

// Stores the config that everything will want to get ahold of
var config Config

// Config represents string kvps of application configuration
type Config map[string]string

// CreateConfig overwrites the current config with whatever is passed in
func CreateConfig(suppliedConfig Config) {
	config = suppliedConfig
}

// MergeConfig takes a default config and merges a supplied one into it.
func MergeConfig(existingConfig Config, suppliedConfig Config) Config {
	for key, value := range suppliedConfig {
		existingConfig[key] = value
	}
	return existingConfig
}

// UpdateConfig merges and adds config to the enliven config
func UpdateConfig(suppliedConfig Config) Config {
	config = MergeConfig(config, suppliedConfig)
	return config
}

// GetConfig returns the config map
func GetConfig() Config {
	if config == nil {
		config = Config{}
	}
	return config
}
