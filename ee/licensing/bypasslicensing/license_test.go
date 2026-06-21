package bypasslicensing

import (
	"testing"

	"github.com/SigNoz/signoz/pkg/types/licensetypes"
	"github.com/SigNoz/signoz/pkg/valuer"
	"github.com/stretchr/testify/require"
)

func TestEnterpriseBypassLicenseMatchesSelfHostedContract(t *testing.T) {
	orgID := valuer.MustNewUUID("0196f794-ff30-7bee-a5f4-ef5ad315715e")

	license := enterpriseBypassLicense(orgID)

	require.Equal(t, BypassLicenseKey, license.Key)
	require.Equal(t, licensetypes.PlanNameEnterprise, license.PlanName)
	require.Equal(t, "valid", license.Status.StringValue())
	require.Equal(t, licenseStateActivated, license.State)
	require.Equal(t, int64(1), license.ValidFrom)
	require.Equal(t, int64(-1), license.ValidUntil)

	require.Equal(t, licenseStatusValid, license.Data["status"])
	require.Equal(t, licenseStateActivated, license.Data["state"])
	require.Equal(t, licensePlatformSelfHosted, license.Data["platform"])
	require.Equal(t, "self-hosted-enterprise", license.Data["plan_id"])
	require.Contains(t, license.Data, "event_queue")
	require.Contains(t, license.Data, "plan")
}

func TestEnterpriseBypassFeaturesDisableCloudOnlyFeatures(t *testing.T) {
	features := enterpriseBypassFeatures()

	featureMap := make(map[valuer.String]bool, len(features))
	for _, feature := range features {
		featureMap[feature.Name] = feature.Active
	}

	require.False(t, featureMap[licensetypes.ChatSupport])
	require.False(t, featureMap[licensetypes.Gateway])
	require.False(t, featureMap[licensetypes.PremiumSupport])
	require.True(t, featureMap[licensetypes.SSO])
	require.True(t, featureMap[licensetypes.Onboarding])
	require.True(t, featureMap[licensetypes.AnomalyDetection])
}
