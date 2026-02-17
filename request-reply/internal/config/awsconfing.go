package config

import (
	"encoding/json"
	"log/slog"
	"os"
)

type AWSConfig struct {
	AccessKeyID     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`
	Region          string `json:"region"`

	// Optional SQS endpoint for local testing (e.g., LocalStack)
	Endpoint string `json:"endpoint"`
}

func NewAwsConfig(filePath string) *AWSConfig {

	var config AWSConfig

	fp, err := os.Open(filePath)
	if err != nil {
		slog.Error("Can't open config file", "error", err)

		panic(err)
	}

	err = json.NewDecoder(fp).Decode(&config)
	if err != nil {
		slog.Error("Can't decode config file", "error", err)
		panic(err)
	}

	if config.SecretAccessKey == "" || config.AccessKeyID == "" || config.Region == "" {
		slog.Error("AWS config is missing required fields", "config", config)
		panic("AWS config is missing required fields")
	}

	return &config

}
