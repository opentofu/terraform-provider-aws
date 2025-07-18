// Code generated by internal/generate/tagstests/main.go; DO NOT EDIT.

package sesv2_test

import (
	"context"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	tfstatecheck "github.com/hashicorp/terraform-provider-aws/internal/acctest/statecheck"
	tfsesv2 "github.com/hashicorp/terraform-provider-aws/internal/service/sesv2"
)

func expectFullResourceTags(ctx context.Context, resourceAddress string, knownValue knownvalue.Check) statecheck.StateCheck {
	return tfstatecheck.ExpectFullResourceTags(tfsesv2.ServicePackage(ctx), resourceAddress, knownValue)
}

func expectFullDataSourceTags(ctx context.Context, resourceAddress string, knownValue knownvalue.Check) statecheck.StateCheck {
	return tfstatecheck.ExpectFullDataSourceTags(tfsesv2.ServicePackage(ctx), resourceAddress, knownValue)
}
