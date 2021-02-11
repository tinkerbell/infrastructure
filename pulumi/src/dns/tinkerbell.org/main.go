package tinkerbelldotorg

import (
	ns1 "github.com/pulumi/pulumi-ns1/sdk/go/ns1"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

func ManageDns(ctx *pulumi.Context) (ns1.LookupZoneResult, error) {
	zone, err := ns1.LookupZone(ctx, &ns1.LookupZoneArgs{
		Zone: "tinkerbell.org",
	})

	if err != nil {
		return ns1.LookupZoneResult{}, err
	}

	// We can add Records here that don't have corresponding code across the
	// other parts of this repository

	return *zone, nil
}
