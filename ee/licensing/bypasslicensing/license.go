package bypasslicensing

import (
	"time"

	"github.com/SigNoz/signoz/pkg/types/licensetypes"
	"github.com/SigNoz/signoz/pkg/valuer"
)

// BypassLicenseKey is the synthetic license key used for self-hosted enterprise bypass.
const BypassLicenseKey = "local-bypass-license"

const (
	licenseStatusValid        = "VALID"
	licenseStateActivated     = "ACTIVATED"
	licensePlatformSelfHosted = "SELF_HOSTED"
)

// enterpriseBypassFeatures returns enterprise features enabled for self-hosted deployments.
// Cloud-only capabilities stay disabled to avoid external service integrations.
func enterpriseBypassFeatures() []*licensetypes.Feature {
	features := make([]*licensetypes.Feature, 0, len(licensetypes.EnterprisePlan))
	for _, feature := range licensetypes.EnterprisePlan {
		cloned := *feature
		switch feature.Name {
		case licensetypes.ChatSupport, licensetypes.Gateway, licensetypes.PremiumSupport:
			cloned.Active = false
		default:
			cloned.Active = true
		}
		features = append(features, &cloned)
	}

	return features
}

// enterpriseBypassLicense returns a synthetic enterprise license that never expires.
func enterpriseBypassLicense(organizationID valuer.UUID) *licensetypes.License {
	now := time.Now()
	nowString := now.UTC().Format(time.RFC3339)
	features := enterpriseBypassFeatures()

	return &licensetypes.License{
		ID:  organizationID,
		Key: BypassLicenseKey,
		Data: map[string]interface{}{
			"status":      licenseStatusValid,
			"state":       licenseStateActivated,
			"platform":    licensePlatformSelfHosted,
			"valid_from":  float64(1),
			"valid_until": float64(-1),
			"free_until":  "",
			"created_at":  nowString,
			"updated_at":  nowString,
			"plan_id":     "self-hosted-enterprise",
			"plan": map[string]interface{}{
				"created_at":  nowString,
				"description": "Self-hosted enterprise",
				"is_active":   true,
				"name":        licensetypes.PlanNameEnterprise.StringValue(),
				"updated_at":  nowString,
			},
			"event_queue": map[string]interface{}{
				"event":        "",
				"status":       "",
				"scheduled_at": "",
				"created_at":   nowString,
				"updated_at":   nowString,
			},
			"features": features,
		},
		PlanName:        licensetypes.PlanNameEnterprise,
		Features:        features,
		Status:          valuer.NewString(licenseStatusValid),
		State:           licenseStateActivated,
		ValidFrom:       1,
		ValidUntil:      -1,
		CreatedAt:       now,
		UpdatedAt:       now,
		LastValidatedAt: now,
		OrganizationID:  organizationID,
	}
}
