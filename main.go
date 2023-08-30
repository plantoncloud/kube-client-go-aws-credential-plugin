package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type ExecCredential struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Status     Status `json:"status"`
}

type Status struct {
	Token string `json:"token"`
}

type StaticCredentialsProvider struct {
	Value aws.Credentials
}

func (s *StaticCredentialsProvider) Retrieve(ctx context.Context) (aws.Credentials, error) {
	return s.Value, nil
}

func main() {
	if len(os.Args) < 4 {
		log.Fatalf("requires three arguments 1. AWS Access Key ID 2. AWS Secret Access Key 3. AWS Region in the same order.")
	}

	awsAccessKeyID := os.Args[1]
	awsSecretAccessKey := os.Args[2]
	awsRegion := os.Args[3]

	credProvider := &StaticCredentialsProvider{
		Value: aws.Credentials{
			AccessKeyID:     awsAccessKeyID,
			SecretAccessKey: awsSecretAccessKey,
		},
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(awsRegion),
		config.WithCredentialsProvider(credProvider),
	)

	if err != nil {
		log.Fatalf("failed to load aws config with error %s", err.Error())
	}

	client := sts.NewFromConfig(cfg)
	input := &sts.GetSessionTokenInput{}
	output, err := client.GetSessionToken(context.TODO(), input)

	if err != nil {
		log.Fatalf("failed to get session token with error %s", err.Error())
	}

	execCredential := ExecCredential{
		APIVersion: "client.authentication.k8s.io/v1",
		Kind:       "ExecCredential",
		Status:     Status{Token: *output.Credentials.SessionToken},
	}

	execCredentialJSON, err := json.Marshal(execCredential)
	if err != nil {
		log.Fatalf("failed to convert exec-credentials to json with error %s", err.Error())
	}

	fmt.Println(string(execCredentialJSON))
}
