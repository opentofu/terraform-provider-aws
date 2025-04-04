// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package servicecatalog

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	awstypes "github.com/aws/aws-sdk-go-v2/service/servicecatalog/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKDataSource("aws_servicecatalog_portfolio_constraints", name="Portfolio Constraints")
func dataSourcePortfolioConstraints() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourcePortfolioConstraintsRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(PortfolioConstraintsReadyTimeout),
		},

		Schema: map[string]*schema.Schema{
			"accept_language": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      acceptLanguageEnglish,
				ValidateFunc: validation.StringInSlice(acceptLanguage_Values(), false),
			},
			"details": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"constraint_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						names.AttrDescription: {
							Type:     schema.TypeString,
							Computed: true,
						},
						names.AttrOwner: {
							Type:     schema.TypeString,
							Computed: true,
						},
						"portfolio_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"product_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						names.AttrType: {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"portfolio_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"product_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourcePortfolioConstraintsRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).ServiceCatalogClient(ctx)

	output, err := waitPortfolioConstraintsReady(ctx, conn, d.Get("accept_language").(string), d.Get("portfolio_id").(string), d.Get("product_id").(string), d.Timeout(schema.TimeoutRead))

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "describing Service Catalog Portfolio Constraints: %s", err)
	}

	if len(output) == 0 {
		return sdkdiag.AppendErrorf(diags, "getting Service Catalog Portfolio Constraints: no results, change your input")
	}

	acceptLanguage := d.Get("accept_language").(string)

	if acceptLanguage == "" {
		acceptLanguage = acceptLanguageEnglish
	}

	d.Set("accept_language", acceptLanguage)
	d.Set("portfolio_id", d.Get("portfolio_id").(string))
	d.Set("product_id", d.Get("product_id").(string))

	if err := d.Set("details", flattenConstraintDetails(output)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting details: %s", err)
	}

	d.SetId(portfolioConstraintsID(d.Get("accept_language").(string), d.Get("portfolio_id").(string), d.Get("product_id").(string)))

	return diags
}

func flattenConstraintDetail(apiObject awstypes.ConstraintDetail) map[string]any {
	tfMap := map[string]any{}

	if v := apiObject.ConstraintId; v != nil {
		tfMap["constraint_id"] = aws.ToString(v)
	}

	if v := apiObject.Description; v != nil {
		tfMap[names.AttrDescription] = aws.ToString(v)
	}

	if v := apiObject.Owner; v != nil {
		tfMap[names.AttrOwner] = aws.ToString(v)
	}

	if v := apiObject.PortfolioId; v != nil {
		tfMap["portfolio_id"] = aws.ToString(v)
	}

	if v := apiObject.ProductId; v != nil {
		tfMap["product_id"] = aws.ToString(v)
	}

	if v := apiObject.Type; v != nil {
		tfMap[names.AttrType] = aws.ToString(v)
	}

	return tfMap
}

func flattenConstraintDetails(apiObjects []awstypes.ConstraintDetail) []any {
	if len(apiObjects) == 0 {
		return nil
	}

	var tfList []any

	for _, apiObject := range apiObjects {
		tfList = append(tfList, flattenConstraintDetail(apiObject))
	}

	return tfList
}
