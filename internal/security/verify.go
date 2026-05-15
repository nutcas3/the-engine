package security

import (
	"context"
	"fmt"
	"strings"
)

// SecurityPolicy represents security validation policies
type SecurityPolicy struct {
	RequireImageSignature bool     `json:"require_image_signature"`
	MaxCVESeverity        string   `json:"max_cve_severity"`
	AllowedRegistries     []string `json:"allowed_registries"`
	BlockUntrustedImages  bool     `json:"block_untrusted_images"`
}

// ImageVerificationResult represents the result of image security verification
type ImageVerificationResult struct {
	ImageRef       string   `json:"image_ref"`
	SignatureValid bool     `json:"signature_valid"`
	CVEsFound      []string `json:"cves_found"`
	Trusted        bool     `json:"trusted"`
	Warnings       []string `json:"warnings"`
}

// VerifyImage performs comprehensive image security verification
func VerifyImage(ctx context.Context, imageRef string) (*ImageVerificationResult, error) {
	result := &ImageVerificationResult{
		ImageRef: imageRef,
		Warnings: []string{},
	}

	// 1. Verify image signature (Cosign integration)
	signatureValid, err := verifyImageSignature(ctx, imageRef)
	if err != nil {
		result.Warnings = append(result.Warnings, fmt.Sprintf("Signature verification failed: %v", err))
		result.SignatureValid = false
	} else {
		result.SignatureValid = signatureValid
	}

	// 2. Scan for CVEs
	cves, err := scanImageCVEs(ctx, imageRef)
	if err != nil {
		result.Warnings = append(result.Warnings, fmt.Sprintf("CVE scanning failed: %v", err))
		result.CVEsFound = []string{}
	} else {
		result.CVEsFound = cves
	}

	// 3. Check if image is from trusted registry
	trusted := isTrustedRegistry(imageRef)
	result.Trusted = trusted

	if !trusted {
		result.Warnings = append(result.Warnings, "Image from untrusted registry")
	}

	return result, nil
}

// verifyImageSignature checks if the image has a valid Cosign signature
func verifyImageSignature(ctx context.Context, imageRef string) (bool, error) {
	// In production, this would use the actual Cosign library
	// For now, we'll simulate signature verification

	// Mock: Images with "trusted" in the name are considered signed
	if strings.Contains(imageRef, "trusted") || strings.Contains(imageRef, "official") {
		return true, nil
	}

	// Mock: Images from certain registries are considered signed
	trustedRegistries := []string{
		"ghcr.io/",
		"docker.io/library/",
		"quay.io/",
		"gcr.io/",
	}

	for _, registry := range trustedRegistries {
		if strings.HasPrefix(imageRef, registry) {
			return true, nil
		}
	}

	// Default to unsigned for unknown images
	return false, nil
}

// scanImageCVEs scans an image for known CVEs
func scanImageCVEs(ctx context.Context, imageRef string) ([]string, error) {
	// In production, this would integrate with Trivy, Grype, or similar scanners
	// For now, we'll simulate CVE detection

	// Mock: Images with "vulnerable" in the name have CVEs
	if strings.Contains(imageRef, "vulnerable") {
		return []string{
			"CVE-2024-1234: Critical vulnerability in base image",
			"CVE-2024-5678: High severity vulnerability in application",
		}, nil
	}

	// Mock: Old versions have CVEs
	if strings.Contains(imageRef, "old") || strings.Contains(imageRef, "legacy") {
		return []string{
			"CVE-2023-9999: Medium severity vulnerability in outdated package",
		}, nil
	}

	// Default to no CVEs for clean images
	return []string{}, nil
}

// isTrustedRegistry checks if the image comes from a trusted registry
func isTrustedRegistry(imageRef string) bool {
	trustedRegistries := []string{
		"ghcr.io/",
		"docker.io/library/",
		"quay.io/",
		"gcr.io/",
		"registry.k8s.io/",
	}

	for _, registry := range trustedRegistries {
		if strings.HasPrefix(imageRef, registry) {
			return true
		}
	}

	return false
}

// ApplySecurityPolicy applies security policies to an image verification result
func ApplySecurityPolicy(result *ImageVerificationResult, policy SecurityPolicy) error {
	// Check signature requirement
	if policy.RequireImageSignature && !result.SignatureValid {
		return fmt.Errorf("image signature required but not valid")
	}

	// Check CVE severity
	if len(result.CVEsFound) > 0 {
		for _, cve := range result.CVEsFound {
			if strings.Contains(cve, "Critical") && policy.MaxCVESeverity != "critical" {
				return fmt.Errorf("critical CVE found: %s", cve)
			}
			if strings.Contains(cve, "High") && policy.MaxCVESeverity == "low" {
				return fmt.Errorf("high severity CVE found: %s", cve)
			}
		}
	}

	// Check trusted registry requirement
	if policy.BlockUntrustedImages && !result.Trusted {
		return fmt.Errorf("image from untrusted registry is blocked")
	}

	return nil
}

// GetDefaultSecurityPolicy returns the default security policy
func GetDefaultSecurityPolicy() SecurityPolicy {
	return SecurityPolicy{
		RequireImageSignature: true,
		MaxCVESeverity:        "medium",
		AllowedRegistries: []string{
			"ghcr.io/",
			"docker.io/library/",
			"quay.io/",
			"gcr.io/",
		},
		BlockUntrustedImages: true,
	}
}

// ValidateDeploymentConfig validates deployment configuration against security policies
func ValidateDeploymentConfig(ctx context.Context, imageRef string, team string) error {
	policy := GetDefaultSecurityPolicy()

	// Adjust policy based on team (dev teams might be more lenient)
	if team == "dev" || team == "engineering" {
		policy.RequireImageSignature = false // Allow unsigned images for dev
		policy.MaxCVESeverity = "high"       // Allow higher CVE severity for dev
	}

	// Verify the image
	result, err := VerifyImage(ctx, imageRef)
	if err != nil {
		return fmt.Errorf("image verification failed: %w", err)
	}

	// Apply security policy
	err = ApplySecurityPolicy(result, policy)
	if err != nil {
		return fmt.Errorf("security policy violation: %w", err)
	}

	return nil
}

// GetSecurityRecommendations provides security recommendations based on verification results
func GetSecurityRecommendations(result *ImageVerificationResult) []string {
	var recommendations []string

	if !result.SignatureValid {
		recommendations = append(recommendations, "Sign the image with Cosign for supply chain security")
	}

	if !result.Trusted {
		recommendations = append(recommendations, "Use images from trusted registries")
	}

	if len(result.CVEsFound) > 0 {
		recommendations = append(recommendations, "Update base image to fix known vulnerabilities")
	}

	if len(result.Warnings) > 0 {
		recommendations = append(recommendations, "Review and address security warnings")
	}

	return recommendations
}
