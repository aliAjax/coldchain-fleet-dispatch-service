package config

import (
	"os"
	"time"
)

type Config struct {
	HTTPAddr          string
	DispatchInterval  time.Duration
	AlertScanInterval time.Duration
	ExcursionWindow   time.Duration
}

func Load() Config {
	return Config{
		HTTPAddr:          getEnv("HTTP_ADDR", "127.0.0.1:18080"),
		DispatchInterval:  getDuration("DISPATCH_INTERVAL", 2*time.Second),
		AlertScanInterval: getDuration("ALERT_SCAN_INTERVAL", 2*time.Second),
		ExcursionWindow:   getDuration("ALERT_EXCURSION_WINDOW", 30*time.Minute),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getDuration(key string, fallback time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return fallback
	}
	return d
}
