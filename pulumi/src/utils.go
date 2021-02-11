package main

import (
	"fmt"

	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

func StackName(ctx *pulumi.Context, name string) string {

	return fmt.Sprintf("%s-%s", ctx.Stack(), name)
}
