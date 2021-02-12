package githubactions

import (
	"fmt"

	"github.com/pulumi/pulumi-equinix-metal/sdk/go/equinix"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi/config"
	"github.com/tinkerbell/infrastructure/src/internal"
)

type GitHubConfig struct {
	Runners []GitHubActionRunnerConfig
}

// GitHubActionRunnerConfig is the struct we allow in the stack configuration
// to describe the GitHubActionRunner we provision
type GitHubActionRunnerConfig struct {
	Facility equinix.Facility
	Plan     equinix.Plan
}

// GitHubActionRunner is the return struct for CreateSaltMaster
type GitHubActionRunners struct {
	Devices []equinix.Device
}

// CreateGitHubActionRunner Provisions a GitHub Action Runner
func CreateGitHubActionRunners(ctx *pulumi.Context, infrastructure internal.Infrastructure) (GitHubActionRunners, error) {
	metalConfig := config.New(ctx, "equinix-metal")
	projectID := metalConfig.Require("projectId")

	stackConfig := config.New(ctx, "")
	githubActionsConfig := &GitHubConfig{}
	stackConfig.RequireObject("github", githubActionsConfig)

	runners := &GitHubActionRunners{}

	for i, runner := range githubActionsConfig.Runners {
		deviceArgs := equinix.DeviceArgs{
			ProjectId: pulumi.String(projectID),
			Hostname:  pulumi.String(fmt.Sprintf("github-action-runner-%d", i)),
			Plan:      runner.Plan,
			Facilities: pulumi.StringArray{
				runner.Facility,
			},
			OperatingSystem: equinix.OperatingSystemUbuntu2004,
			Tags: pulumi.StringArray{
				pulumi.String("role:github-action-runner"),
			},
			BillingCycle: equinix.BillingCycleHourly,
			UserData: pulumi.All(infrastructure.SaltMasterIp).ApplyT(func(args []interface{}) string {
				return cloudInitConfig(&MinionConfig{
					MasterIp: args[0].(string),
				})
			}).(pulumi.StringOutput),
		}

		device, err := equinix.NewDevice(ctx, fmt.Sprintf("github-action-runner-%d", i), &deviceArgs)

		if err != nil {
			return GitHubActionRunners{}, err
		}

		runners.Devices = append(runners.Devices, *device)
	}

	return *runners, nil
}
