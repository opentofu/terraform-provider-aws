// Code generated by internal/generate/servicepackage/main.go; DO NOT EDIT.

package cognitoidp

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/types"
	"github.com/hashicorp/terraform-provider-aws/names"
)

type servicePackage struct{}

func (p *servicePackage) FrameworkDataSources(ctx context.Context) []*types.ServicePackageFrameworkDataSource {
	return []*types.ServicePackageFrameworkDataSource{
		{
			Factory:  newUserGroupDataSource,
			TypeName: "aws_cognito_user_group",
			Name:     "User Group",
		},
		{
			Factory:  newUserGroupsDataSource,
			TypeName: "aws_cognito_user_groups",
			Name:     "User Groups",
		},
		{
			Factory:  newUserPoolDataSource,
			TypeName: "aws_cognito_user_pool",
			Name:     "User Pool",
		},
	}
}

func (p *servicePackage) FrameworkResources(ctx context.Context) []*types.ServicePackageFrameworkResource {
	return []*types.ServicePackageFrameworkResource{
		{
			Factory:  newManagedUserPoolClientResource,
			TypeName: "aws_cognito_managed_user_pool_client",
			Name:     "Managed User Pool Client",
		},
		{
			Factory:  newUserPoolClientResource,
			TypeName: "aws_cognito_user_pool_client",
			Name:     "User Pool Client",
		},
	}
}

func (p *servicePackage) SDKDataSources(ctx context.Context) []*types.ServicePackageSDKDataSource {
	return []*types.ServicePackageSDKDataSource{
		{
			Factory:  dataSourceUserPoolClient,
			TypeName: "aws_cognito_user_pool_client",
			Name:     "User Pool Client",
		},
		{
			Factory:  dataSourceUserPoolClients,
			TypeName: "aws_cognito_user_pool_clients",
			Name:     "User Pool Clients",
		},
		{
			Factory:  dataSourceUserPoolSigningCertificate,
			TypeName: "aws_cognito_user_pool_signing_certificate",
			Name:     "User Pool Signing Certificate",
		},
		{
			Factory:  dataSourceUserPools,
			TypeName: "aws_cognito_user_pools",
			Name:     "User Pools",
		},
	}
}

func (p *servicePackage) SDKResources(ctx context.Context) []*types.ServicePackageSDKResource {
	return []*types.ServicePackageSDKResource{
		{
			Factory:  resourceIdentityProvider,
			TypeName: "aws_cognito_identity_provider",
			Name:     "Identity Provider",
		},
		{
			Factory:  resourceResourceServer,
			TypeName: "aws_cognito_resource_server",
			Name:     "Resource Server",
		},
		{
			Factory:  resourceRiskConfiguration,
			TypeName: "aws_cognito_risk_configuration",
			Name:     "Risk Configuration",
		},
		{
			Factory:  resourceUser,
			TypeName: "aws_cognito_user",
			Name:     "User",
		},
		{
			Factory:  resourceUserGroup,
			TypeName: "aws_cognito_user_group",
			Name:     "User Group",
		},
		{
			Factory:  resourceUserInGroup,
			TypeName: "aws_cognito_user_in_group",
			Name:     "Group User",
		},
		{
			Factory:  resourceUserPool,
			TypeName: "aws_cognito_user_pool",
			Name:     "User Pool",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: names.AttrARN,
			},
		},
		{
			Factory:  resourceUserPoolDomain,
			TypeName: "aws_cognito_user_pool_domain",
			Name:     "User Pool Domain",
		},
		{
			Factory:  resourceUserPoolUICustomization,
			TypeName: "aws_cognito_user_pool_ui_customization",
			Name:     "User Pool UI Customization",
		},
	}
}

func (p *servicePackage) ServicePackageName() string {
	return names.CognitoIDP
}

// NewClient returns a new AWS SDK for Go v2 client for this service package's AWS API.
func (p *servicePackage) NewClient(ctx context.Context, config map[string]any) (*cognitoidentityprovider.Client, error) {
	cfg := *(config["aws_sdkv2_config"].(*aws.Config))
	optFns := []func(*cognitoidentityprovider.Options){
		cognitoidentityprovider.WithEndpointResolverV2(newEndpointResolverV2()),
		withBaseEndpoint(config[names.AttrEndpoint].(string)),
		withExtraOptions(ctx, p, config),
	}

	return cognitoidentityprovider.NewFromConfig(cfg, optFns...), nil
}

// withExtraOptions returns a functional option that allows this service package to specify extra API client options.
// This option is always called after any generated options.
func withExtraOptions(ctx context.Context, sp conns.ServicePackage, config map[string]any) func(*cognitoidentityprovider.Options) {
	if v, ok := sp.(interface {
		withExtraOptions(context.Context, map[string]any) []func(*cognitoidentityprovider.Options)
	}); ok {
		optFns := v.withExtraOptions(ctx, config)

		return func(o *cognitoidentityprovider.Options) {
			for _, optFn := range optFns {
				optFn(o)
			}
		}
	}

	return func(*cognitoidentityprovider.Options) {}
}

func ServicePackage(ctx context.Context) conns.ServicePackage {
	return &servicePackage{}
}
