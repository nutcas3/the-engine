package security

import (
	"context"
	"testing"
)

func TestVerifyImage(t *testing.T) {
	result, err := VerifyImage(context.Background(), "trusted/test-image:latest")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result.ImageRef != "trusted/test-image:latest" {
		t.Errorf("Expected trusted/test-image:latest, got %s", result.ImageRef)
	}
}

func TestVerifyImageSignature(t *testing.T) {
	valid, err := verifyImageSignature(context.Background(), "trusted/test-image:latest")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !valid {
		t.Error("Expected valid signature for trusted image")
	}
}

func TestScanImageCVEs(t *testing.T) {
	cves, err := scanImageCVEs(context.Background(), "trusted/test-image:latest")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(cves) != 0 {
		t.Error("Expected no CVEs for trusted image")
	}
}

func TestIsTrustedRegistry(t *testing.T) {
	trusted := isTrustedRegistry("ghcr.io/test/image:latest")
	if !trusted {
		t.Error("Expected true for ghcr.io registry")
	}

	trusted = isTrustedRegistry("docker.io/library/test-image:latest")
	if !trusted {
		t.Error("Expected true for docker.io registry")
	}

	trusted = isTrustedRegistry("unknown-registry.com/test/image:latest")
	if trusted {
		t.Error("Expected false for unknown registry")
	}
}

func TestApplySecurityPolicy(t *testing.T) {
	policy := GetDefaultSecurityPolicy()
	result := &ImageVerificationResult{
		ImageRef:       "trusted/test-image:latest",
		SignatureValid: true,
		Trusted:        true,
		CVEsFound:      []string{},
	}

	err := ApplySecurityPolicy(result, policy)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestGetDefaultSecurityPolicy(t *testing.T) {
	policy := GetDefaultSecurityPolicy()
	if !policy.RequireImageSignature {
		t.Error("Expected RequireImageSignature to be true")
	}
	if policy.MaxCVESeverity == "" {
		t.Error("Expected non-empty MaxCVESeverity")
	}
}

func TestValidateDeploymentConfig(t *testing.T) {
	err := ValidateDeploymentConfig(context.Background(), "ghcr.io/test/image:latest", "dev")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestGetSecurityRecommendations(t *testing.T) {
	result := &ImageVerificationResult{
		ImageRef:       "untrusted/test-image:latest",
		SignatureValid: false,
		Trusted:        false,
		CVEsFound:      []string{},
	}

	recommendations := GetSecurityRecommendations(result)
	if len(recommendations) == 0 {
		t.Error("Expected some recommendations")
	}
}
