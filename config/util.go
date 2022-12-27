package config

import (
	"time"

	"github.com/saiya/mesh_for_home_server/logger"
)

func parseDurationOrDefault(str string, fallback time.Duration) time.Duration {
	if str == "" {
		return fallback
	}

	d, err := time.ParseDuration(str)
	if err != nil {
		logger.Get().Warnf("Cannot parse duration string \"%s\" in the configuration file: %v", str, err)
		return fallback
	}
	return d
}
