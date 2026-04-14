package utils

import "testing"

func TestConfigureJWTSecretFromEnvRotatesValidationKey(t *testing.T) {
	originalSecret := string(GetSecret())
	t.Cleanup(func() {
		SetSecret(originalSecret)
	})

	t.Setenv("JWT_SECRET", "first-secret")
	if err := ConfigureJWTSecretFromEnv(); err != nil {
		t.Fatalf("ConfigureJWTSecretFromEnv() error = %v", err)
	}

	token, err := GenerateToken(12, "+15550001111")
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	t.Setenv("JWT_SECRET", "second-secret")
	if err := ConfigureJWTSecretFromEnv(); err != nil {
		t.Fatalf("ConfigureJWTSecretFromEnv() after rotation error = %v", err)
	}

	if _, err := ValidateToken(token); err == nil {
		t.Fatal("ValidateToken() succeeded with a token signed by the previous secret")
	}
}

func TestConfigureJWTSecretRejectsEmptySecret(t *testing.T) {
	originalSecret := string(GetSecret())
	t.Cleanup(func() {
		SetSecret(originalSecret)
	})

	SetSecret("")
	if _, err := GenerateToken(7, "+15550002222"); err == nil {
		t.Fatal("GenerateToken() succeeded without a configured JWT secret")
	}

	if err := ConfigureJWTSecret(""); err == nil {
		t.Fatal("ConfigureJWTSecret() accepted an empty secret")
	}
}
