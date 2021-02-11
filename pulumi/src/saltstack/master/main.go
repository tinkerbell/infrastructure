package master

import (
	"fmt"

	"github.com/pulumi/pulumi-equinix-metal/sdk/go/equinix"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi/config"
)

// This is made available to the CreateSaltMaster function through
// the stack configuratiom.
type SaltMasterConfig struct {
	Facilities []equinix.Facility
	Plan       equinix.Plan
}

type SaltMaster struct {
	Device equinix.Device
}

func CreateSaltMaster(ctx *pulumi.Context) (SaltMaster, error) {
	metalConfig := config.New(ctx, "equinix-metal")
	projectId := metalConfig.Require("projectId")

	stackConfig := config.New(ctx, "")
	saltMasterConfig := &SaltMasterConfig{}
	stackConfig.RequireObject("saltMaster", saltMasterConfig)

	facilitiesStringInput := pulumi.StringArray{}

	for _, x := range saltMasterConfig.Facilities {
		facilitiesStringInput = append(facilitiesStringInput, x)
	}

	deviceArgs := equinix.DeviceArgs{
		ProjectId:       pulumi.String(projectId),
		Hostname:        pulumi.String(fmt.Sprintf("%s-%s", ctx.Stack(), "salt-master")),
		Plan:            saltMasterConfig.Plan,
		Facilities:      facilitiesStringInput,
		OperatingSystem: equinix.OperatingSystemUbuntu2004,
		Tags: pulumi.StringArray{
			pulumi.String("role:salt-master"),
		},
		BillingCycle: equinix.BillingCycleHourly,
		UserData:     pulumi.String(cloudInitConfig),
	}

	device, err := equinix.NewDevice(ctx, "salt-master", &deviceArgs)
	if err != nil {
		return SaltMaster{}, err
	}

	ctx.Export("saltMasterIp", &device.AccessPublicIpv4)

	return SaltMaster{
		Device: *device,
	}, nil
}
