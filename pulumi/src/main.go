package main

import (
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
	assets "github.com/tinkerbell/infrastructure/src/assets"
	tinkdns "github.com/tinkerbell/infrastructure/src/dns/tinkerbell.org"
	runners "github.com/tinkerbell/infrastructure/src/github-actions"
	"github.com/tinkerbell/infrastructure/src/internal"
	"github.com/tinkerbell/infrastructure/src/saltstack/master"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		zone, err := tinkdns.ManageDNS(ctx)
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

		infrastructure.SaltMasterIP = saltMaster.Device.AccessPrivateIpv4

		_, err = runners.CreateGitHubActionRunners(ctx, infrastructure)

		if err != nil {
			return err
		}

		_, err = assets.CreateAssetsDNS(ctx, infrastructure)
		if err != nil {
			return err
		}

		return nil
	})
}
