package master

import (
	"fmt"

	"github.com/pulumi/pulumi-equinix-metal/sdk/go/equinix"
	"github.com/pulumi/pulumi-ns1/sdk/go/ns1"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi/config"
	"github.com/tinkerbell/infrastructure/src/internal"
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

func CreateSaltMaster(ctx *pulumi.Context, infrastructure internal.Infrastructure) (SaltMaster, error) {
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

	// Create DNS record for Teleport
	_, err = ns1.NewRecord(ctx, "teleport", &ns1.RecordArgs{
		Zone:   pulumi.String(infrastructure.Zone.Zone),
		Domain: pulumi.String(fmt.Sprintf("teleport.%s", infrastructure.Zone.Zone)),
		Type:   pulumi.String("A"),
		Answers: ns1.RecordAnswerArray{
			ns1.RecordAnswerArgs{
				Answer: &device.AccessPublicIpv4,
			},
		},
	})

	if err != nil {
		return SaltMaster{}, err
	}

	return SaltMaster{
		Device: *device,
	}, nil
}
