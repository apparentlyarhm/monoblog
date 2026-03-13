package server

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ProxyHost    string `envconfig:"PROXY_HOST" required:"true"`
	HmacSecret   string `envconfig:"HMAC_SECRET" required:"true"`
	GlobalApiKey string `envconfig:"GLOBAL_API_KEY" required:"true"`
	ViewEndpoint string `envconfig:"VIEW_ENDPOINT" required:"true"`
}

// Process the env to read the env into a struct
func Load() (Config, error) {
	var cfg Config
	// The first argument is a prefix, which we'll leave empty.

	err := envconfig.Process("", &cfg)
	if err != nil {
		return cfg, fmt.Errorf("failed to load config: %w", err)
	}

	// just protection against empty string or something, just in case...
	// using log here doesnt make sense since we are not doing anything
	fmt.Printf("[ENV] len HMAC SECRET: %v\n", len(cfg.HmacSecret))
	fmt.Printf("[ENV] Proxy: %v\n", cfg.ProxyHost)

	return cfg, nil
}
