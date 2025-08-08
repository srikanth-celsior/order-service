package utils

import (
	"context"
	"fmt"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

// GetSecret fetches a secret value from Google Cloud Secret Manager
func GetSecret(secretName string, projectID string) (string, error) {
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create secretmanager client: %w", err)
	}
	defer client.Close()

	accessRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: fmt.Sprintf("projects/%s/secrets/%s/versions/latest", projectID, secretName),
	}

	result, err := client.AccessSecretVersion(ctx, accessRequest)
	if err != nil {
		return "", fmt.Errorf("failed to access secret version: %w", err)
	}

	return string(result.Payload.Data), nil
}

// GetSecrets fetches multiple secrets from Google Cloud Secret Manager
func GetSecrets(secretNames []string, projectID string) (map[string]string, error) {
	secrets := make(map[string]string)
	for _, name := range secretNames {
		val, err := GetSecret(name, projectID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch secret '%s': %w", name, err)
		}
		secrets[name] = val
	}
	return secrets, nil
}
