package service

import (
	"context"
	"fmt"
	"net/url"

	awsc "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/serdarozerr/request-reply/internal/config"
)

func NewSQSClient(ctx context.Context, cfg *config.AWSConfig) (*sqs.Client, error) {
	var opts []func(*awsc.LoadOptions) error

	opts = append(opts, awsc.WithRegion(cfg.Region))
	opts = append(opts, awsc.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.SecretAccessKey, "")))

	awsCfg, err := awsc.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return nil, err
	}

	var sqsOpts []func(*sqs.Options)
	if cfg.Endpoint != "" {
		// Validate endpoint format
		if _, err := url.Parse(cfg.Endpoint); err != nil {
			return nil, fmt.Errorf("invalid SQS endpoint URL: %w", err)
		}

		sqsOpts = append(sqsOpts, func(o *sqs.Options) {
			o.BaseEndpoint = &cfg.Endpoint
		})
	}

	// Create and return SQS client with options
	return sqs.NewFromConfig(awsCfg, sqsOpts...), nil
}
