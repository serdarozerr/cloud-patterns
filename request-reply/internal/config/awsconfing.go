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
	Name			string `json:"queue_name"`
	// Optional SQS endpoint for local testing (e.g., LocalStack)
	QueueURL string `json:"queue_url"`
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

	if config.Name =="" && config.QueueURL==""{
		slog.Error("Both queue name and queue url fields are emtpy, provide one at least")
		panic("Both queue name and queue url fields are emtpy, provide one at least")
	}

	return &config

}

