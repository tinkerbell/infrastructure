package main

import (
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
	tinkdns "github.com/tinkerbell/infrastructure/src/dns/tinkerbell.org"
	runners "github.com/tinkerbell/infrastructure/src/github-actions"
	"github.com/tinkerbell/infrastructure/src/internal"
	"github.com/tinkerbell/infrastructure/src/saltstack/master"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		zone, err := tinkdns.ManageDns(ctx)
		if err != nil {
			return err
		}

		infrastructure := internal.Infrastructure{
			Zone: zone,
		}

		saltMaster, err := master.CreateSaltMaster(ctx, infrastructure)

		if err != nil {
			return err
		}

		infrastructure.SaltMasterIp = saltMaster.Device.AccessPrivateIpv4

		_, err = runners.CreateGitHubActionRunners(ctx, infrastructure)

		if err != nil {
			return err
		}

		return nil
	})
}
