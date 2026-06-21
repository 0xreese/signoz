package bypasslicensing

import (
	"context"

	"github.com/SigNoz/signoz/pkg/factory"
	"github.com/SigNoz/signoz/pkg/licensing"
	"github.com/SigNoz/signoz/pkg/types/licensetypes"
	"github.com/SigNoz/signoz/pkg/valuer"
)

type provider struct {
	settings factory.ScopedProviderSettings
	stopChan chan struct{}
}

func NewProviderFactory() factory.ProviderFactory[licensing.Licensing, licensing.Config] {
	return factory.NewProviderFactory(factory.MustNewName("bypass"), func(ctx context.Context, providerSettings factory.ProviderSettings, config licensing.Config) (licensing.Licensing, error) {
		return New(ctx, providerSettings, config)
	})
}

func New(_ context.Context, ps factory.ProviderSettings, _ licensing.Config) (licensing.Licensing, error) {
	settings := factory.NewScopedProviderSettings(ps, "github.com/SigNoz/signoz/ee/licensing/bypasslicensing")
	return &provider{
		settings: settings,
		stopChan: make(chan struct{}),
	}, nil
}

func (provider *provider) Start(context.Context) error {
	<-provider.stopChan
	return nil
}

func (provider *provider) Stop(context.Context) error {
	close(provider.stopChan)
	return nil
}

func (provider *provider) Validate(context.Context) error {
	return nil
}

func (provider *provider) Activate(ctx context.Context, organizationID valuer.UUID, _ string) error {
	provider.settings.Logger().DebugContext(ctx, "license activation bypassed", "org_id", organizationID.StringValue())
	return nil
}

func (provider *provider) GetActive(_ context.Context, organizationID valuer.UUID) (*licensetypes.License, error) {
	return enterpriseBypassLicense(organizationID), nil
}

func (provider *provider) Refresh(ctx context.Context, organizationID valuer.UUID) error {
	provider.settings.Logger().DebugContext(ctx, "license refresh bypassed", "org_id", organizationID.StringValue())
	return nil
}

func (provider *provider) Checkout(_ context.Context, _ valuer.UUID, _ *licensetypes.PostableSubscription) (*licensetypes.GettableSubscription, error) {
	return &licensetypes.GettableSubscription{RedirectURL: ""}, nil
}

func (provider *provider) Portal(_ context.Context, _ valuer.UUID, _ *licensetypes.PostableSubscription) (*licensetypes.GettableSubscription, error) {
	return &licensetypes.GettableSubscription{RedirectURL: ""}, nil
}

func (provider *provider) GetFeatureFlags(_ context.Context, organizationID valuer.UUID) ([]*licensetypes.Feature, error) {
	return enterpriseBypassFeatures(), nil
}

func (provider *provider) Collect(_ context.Context, orgID valuer.UUID) (map[string]any, error) {
	return licensetypes.NewStatsFromLicense(enterpriseBypassLicense(orgID)), nil
}
