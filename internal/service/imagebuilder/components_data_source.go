// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package imagebuilder

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/imagebuilder"
	awstypes "github.com/aws/aws-sdk-go-v2/service/imagebuilder/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/enum"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	"github.com/hashicorp/terraform-provider-aws/internal/namevaluesfilters"
	tfslices "github.com/hashicorp/terraform-provider-aws/internal/slices"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKDataSource("aws_imagebuilder_components", name="Components")
func dataSourceComponents() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceComponentsRead,

		Schema: map[string]*schema.Schema{
			names.AttrARNs: {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			names.AttrFilter: namevaluesfilters.Schema(),
			names.AttrNames: {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			names.AttrOwner: {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: enum.Validate[awstypes.Ownership](),
			},
		},
	}
}

func dataSourceComponentsRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).ImageBuilderClient(ctx)

	input := &imagebuilder.ListComponentsInput{}

	if v, ok := d.GetOk(names.AttrFilter); ok {
		input.Filters = namevaluesfilters.New(v.(*schema.Set)).ImageBuilderFilters()
	}

	if v, ok := d.GetOk(names.AttrOwner); ok {
		input.Owner = awstypes.Ownership(v.(string))
	}

	components, err := findComponents(ctx, conn, input)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading Image Builder Components: %s", err)
	}

	d.SetId(meta.(*conns.AWSClient).Region(ctx))
	d.Set(names.AttrARNs, tfslices.ApplyToAll(components, func(v awstypes.ComponentVersion) string {
		return aws.ToString(v.Arn)
	}))
	d.Set(names.AttrNames, tfslices.ApplyToAll(components, func(v awstypes.ComponentVersion) string {
		return aws.ToString(v.Name)
	}))

	return diags
}

func findComponents(ctx context.Context, conn *imagebuilder.Client, input *imagebuilder.ListComponentsInput) ([]awstypes.ComponentVersion, error) {
	var output []awstypes.ComponentVersion

	pages := imagebuilder.NewListComponentsPaginator(conn, input)
	for pages.HasMorePages() {
		page, err := pages.NextPage(ctx)

		if err != nil {
			return nil, err
		}

		output = append(output, page.ComponentVersionList...)
	}

	return output, nil
}
