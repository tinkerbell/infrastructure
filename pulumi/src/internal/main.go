package internal

import (
	ns1 "github.com/pulumi/pulumi-ns1/sdk/go/ns1"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

// We'll use this type to pass around components of
// our infrastructure that may be needed by other providers:
// with DNS zone being the first use-case.
type Infrastructure struct {
	Zone         ns1.LookupZoneResult
	SaltMasterIP pulumi.StringOutput
}
