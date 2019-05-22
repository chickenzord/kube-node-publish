package config

import (
	"os"
	"path/filepath"
	"strconv"
)

// Config represents app configs
type Config struct {
	DomainName  string
	CheckPeriod string
	InCluster   bool
	KubeConfig  string
	LogLevel    string
}

func getStrEnv(key string, defValue string) string {
	if value, found := os.LookupEnv(key); found {
		return value
	}
	return defValue
}

func getIntEnv(key string, defValue int) int {
	str := getStrEnv(key, string(defValue))
	if result, err := strconv.Atoi(str); err == nil {
		return result
	}
	return 0
}

func getBoolEnv(key string) bool {
	str := getStrEnv(key, "false")
	return str == "true" || str == "yes" || str == "on"
}

// NewConfigFromEnv returns new app config from environment variables
func NewConfigFromEnv() Config {
	return Config{
		DomainName:  getStrEnv("DOMAIN_NAME", ""),
		CheckPeriod: getStrEnv("CHECK_PERIOD", "60s"),
		InCluster:   getBoolEnv("IN_CLUSTER"),
		KubeConfig:  getStrEnv("KUBECONFIG", defaultKubeConfig()),
		LogLevel:    getStrEnv("LOG_LEVEL", "debug"),
	}
}

func defaultKubeConfig() string {
	home := getStrEnv("HOME", "")
	return filepath.Join(home, ".kube", "config")
}
