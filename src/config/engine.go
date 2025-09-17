package config

// EngineConfig contains engine-specific configuration
type EngineConfig struct {
	UseAltScreen     bool
	TargetFPS        int
	EnableDebugMode  bool
	MaxMessageBuffer int
}

// DefaultEngineConfig returns the default engine configuration
func DefaultEngineConfig() EngineConfig {
	return EngineConfig{
		UseAltScreen:     true,
		TargetFPS:        60,
		EnableDebugMode:  false,
		MaxMessageBuffer: 100,
	}
}

// Global engine config instance
var Engine = DefaultEngineConfig()
