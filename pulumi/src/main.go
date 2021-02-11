package main

import (
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
	"github.com/tinkerbell/infrastructure/src/saltstack/master"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := master.CreateSaltMaster(ctx)

		if err != nil {
			return err
		}

		return nil
	})
}
