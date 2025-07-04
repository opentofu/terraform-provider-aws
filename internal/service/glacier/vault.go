// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package glacier

import (
	"context"
	"fmt"
	"log"

	"github.com/YakDriver/regexache"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/glacier"
	"github.com/aws/aws-sdk-go-v2/service/glacier/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	"github.com/hashicorp/terraform-provider-aws/internal/flex"
	"github.com/hashicorp/terraform-provider-aws/internal/sdkv2"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKResource("aws_glacier_vault", name="Vault")
// @Tags(identifierAttribute="id")
func resourceVault() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceVaultCreate,
		ReadWithoutTimeout:   resourceVaultRead,
		UpdateWithoutTimeout: resourceVaultUpdate,
		DeleteWithoutTimeout: resourceVaultDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"access_policy": sdkv2.IAMPolicyDocumentSchemaOptional(),
			names.AttrARN: {
				Type:     schema.TypeString,
				Computed: true,
			},
			names.AttrLocation: {
				Type:     schema.TypeString,
				Computed: true,
			},
			names.AttrName: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 255),
					validation.StringMatch(regexache.MustCompile(`^[0-9A-Za-z_.-]+$`),
						"only alphanumeric characters, hyphens, underscores, and periods are allowed"),
				),
			},
			"notification": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"events": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								ValidateFunc: validation.StringInSlice([]string{
									"ArchiveRetrievalCompleted",
									"InventoryRetrievalCompleted",
								}, false),
							},
						},
						"sns_topic": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: verify.ValidARN,
						},
					},
				},
			},
			names.AttrTags:    tftags.TagsSchema(),
			names.AttrTagsAll: tftags.TagsSchemaComputed(),
		},
	}
}

func resourceVaultCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).GlacierClient(ctx)

	name := d.Get(names.AttrName).(string)
	input := glacier.CreateVaultInput{
		VaultName: aws.String(name),
	}

	_, err := conn.CreateVault(ctx, &input)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "creating Glacier Vault (%s): %s", name, err)
	}

	d.SetId(name)

	if err := createTags(ctx, conn, d.Id(), getTagsIn(ctx)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting Glacier Vault (%s) tags: %s", d.Id(), err)
	}

	if v, ok := d.GetOk("access_policy"); ok {
		policy, err := structure.NormalizeJsonString(v.(string))
		if err != nil {
			return sdkdiag.AppendFromErr(diags, err)
		}

		input := glacier.SetVaultAccessPolicyInput{
			Policy: &types.VaultAccessPolicy{
				Policy: aws.String(policy),
			},
			VaultName: aws.String(d.Id()),
		}

		_, err = conn.SetVaultAccessPolicy(ctx, &input)

		if err != nil {
			return sdkdiag.AppendErrorf(diags, "setting Glacier Vault (%s) access policy: %s", d.Id(), err)
		}
	}

	if v, ok := d.GetOk("notification"); ok && len(v.([]any)) > 0 && v.([]any)[0] != nil {
		input := glacier.SetVaultNotificationsInput{
			VaultName:               aws.String(d.Id()),
			VaultNotificationConfig: expandVaultNotificationConfig(v.([]any)[0].(map[string]any)),
		}

		_, err := conn.SetVaultNotifications(ctx, &input)

		if err != nil {
			return sdkdiag.AppendErrorf(diags, "setting Glacier Vault (%s) notifications: %s", d.Id(), err)
		}
	}

	return append(diags, resourceVaultRead(ctx, d, meta)...)
}

func resourceVaultRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).GlacierClient(ctx)

	output, err := findVaultByName(ctx, conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] Glacier Vault (%s) not found, removing from state", d.Id())
		d.SetId("")
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading Glacier Vault (%s): %s", d.Id(), err)
	}

	d.Set(names.AttrARN, output.VaultARN)
	d.Set(names.AttrLocation, fmt.Sprintf("/%s/vaults/%s", meta.(*conns.AWSClient).AccountID(ctx), d.Id()))
	d.Set(names.AttrName, output.VaultName)

	accessPolicy, err := findVaultAccessPolicyByName(ctx, conn, d.Id())
	switch {
	case tfresource.NotFound(err):
		d.Set("access_policy", nil)
	case err != nil:
		return sdkdiag.AppendErrorf(diags, "reading Glacier Vault (%s) access policy: %s", d.Id(), err)
	default:
		policy, err := verify.PolicyToSet(d.Get("access_policy").(string), aws.ToString(accessPolicy.Policy))
		if err != nil {
			return sdkdiag.AppendFromErr(diags, err)
		}

		d.Set("access_policy", policy)
	}

	notificationConfig, err := findVaultNotificationsByName(ctx, conn, d.Id())
	switch {
	case tfresource.NotFound(err):
		d.Set("notification", nil)
	case err != nil:
		return sdkdiag.AppendErrorf(diags, "reading Glacier Vault (%s) notifications: %s", d.Id(), err)
	default:
		tfMap := map[string]any{}

		if v := notificationConfig.Events; v != nil {
			tfMap["events"] = v
		}

		if v := notificationConfig.SNSTopic; v != nil {
			tfMap["sns_topic"] = aws.ToString(v)
		}

		if err := d.Set("notification", []any{tfMap}); err != nil {
			return sdkdiag.AppendErrorf(diags, "setting notification: %s", err)
		}
	}

	return diags
}

func resourceVaultUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).GlacierClient(ctx)

	if d.HasChange("access_policy") {
		if v, ok := d.GetOk("access_policy"); ok {
			policy, err := structure.NormalizeJsonString(v.(string))
			if err != nil {
				return sdkdiag.AppendFromErr(diags, err)
			}

			input := glacier.SetVaultAccessPolicyInput{
				Policy: &types.VaultAccessPolicy{
					Policy: aws.String(policy),
				},
				VaultName: aws.String(d.Id()),
			}

			_, err = conn.SetVaultAccessPolicy(ctx, &input)

			if err != nil {
				return sdkdiag.AppendErrorf(diags, "setting Glacier Vault (%s) access policy: %s", d.Id(), err)
			}
		} else {
			input := glacier.DeleteVaultAccessPolicyInput{
				VaultName: aws.String(d.Id()),
			}

			_, err := conn.DeleteVaultAccessPolicy(ctx, &input)

			if err != nil {
				return sdkdiag.AppendErrorf(diags, "deleting Glacier Vault (%s) access policy: %s", d.Id(), err)
			}
		}
	}

	if d.HasChange("notification") {
		if v, ok := d.GetOk("notification"); ok && len(v.([]any)) > 0 && v.([]any)[0] != nil {
			input := glacier.SetVaultNotificationsInput{
				VaultName:               aws.String(d.Id()),
				VaultNotificationConfig: expandVaultNotificationConfig(v.([]any)[0].(map[string]any)),
			}

			_, err := conn.SetVaultNotifications(ctx, &input)

			if err != nil {
				return sdkdiag.AppendErrorf(diags, "setting Glacier Vault (%s) notifications: %s", d.Id(), err)
			}
		} else {
			input := glacier.DeleteVaultNotificationsInput{
				VaultName: aws.String(d.Id()),
			}

			_, err := conn.DeleteVaultNotifications(ctx, &input)

			if err != nil {
				return sdkdiag.AppendErrorf(diags, "deleting Glacier Vault (%s) notifications: %s", d.Id(), err)
			}
		}
	}

	return append(diags, resourceVaultRead(ctx, d, meta)...)
}

func resourceVaultDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).GlacierClient(ctx)

	log.Printf("[DEBUG] Deleting Glacier Vault: %s", d.Id())
	input := glacier.DeleteVaultInput{
		VaultName: aws.String(d.Id()),
	}
	_, err := conn.DeleteVault(ctx, &input)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting Glacier Vault (%s): %s", d.Id(), err)
	}

	return diags
}

func findVaultByName(ctx context.Context, conn *glacier.Client, name string) (*glacier.DescribeVaultOutput, error) {
	input := glacier.DescribeVaultInput{
		VaultName: aws.String(name),
	}

	output, err := conn.DescribeVault(ctx, &input)

	if errs.IsA[*types.ResourceNotFoundException](err) {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return output, nil
}

func findVaultAccessPolicyByName(ctx context.Context, conn *glacier.Client, name string) (*types.VaultAccessPolicy, error) {
	input := glacier.GetVaultAccessPolicyInput{
		VaultName: aws.String(name),
	}

	output, err := conn.GetVaultAccessPolicy(ctx, &input)

	if errs.IsA[*types.ResourceNotFoundException](err) {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil || output.Policy == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return output.Policy, nil
}

func findVaultNotificationsByName(ctx context.Context, conn *glacier.Client, name string) (*types.VaultNotificationConfig, error) {
	input := glacier.GetVaultNotificationsInput{
		VaultName: aws.String(name),
	}

	output, err := conn.GetVaultNotifications(ctx, &input)

	if errs.IsA[*types.ResourceNotFoundException](err) {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil || output.VaultNotificationConfig == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return output.VaultNotificationConfig, nil
}

func expandVaultNotificationConfig(tfMap map[string]any) *types.VaultNotificationConfig {
	if tfMap == nil {
		return nil
	}

	apiObject := &types.VaultNotificationConfig{}

	if v, ok := tfMap["events"].(*schema.Set); ok && v.Len() > 0 {
		apiObject.Events = flex.ExpandStringValueSet(v)
	}

	if v, ok := tfMap["sns_topic"].(string); ok && v != "" {
		apiObject.SNSTopic = aws.String(v)
	}

	return apiObject
}
