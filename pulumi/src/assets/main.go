package assets

import (
	"github.com/pulumi/pulumi-ns1/sdk/go/ns1"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
	"github.com/tinkerbell/infrastructure/src/internal"
)

type AssetsDns struct {
	Cname string
}

func CreateAssetsDns(ctx *pulumi.Context, infrastructure internal.Infrastructure) (AssetsDns, error) {
	// Create DNS record for Teleport
	// This verification record is required for Cloudflare to issue a certificate for this domain
	_, err := ns1.NewRecord(ctx, "assets-cname-verification", &ns1.RecordArgs{
		Zone:   pulumi.String(infrastructure.Zone.Zone),
		Domain: pulumi.String("_f213f2a073773bab6da1e978b94e1a92.assets"),
		Type:   pulumi.String("CNAME"),
		Answers: ns1.RecordAnswerArray{
			ns1.RecordAnswerArgs{
				Answer: pulumi.String("_69acff72294eb94eaac26bef9c1fdb06.zjfbrrwmzc.acm-validations.aws"),
			},
		},
	})

	if err != nil {
		return AssetsDns{}, err
	}

	return AssetsDns{
		Cname: "test",
	}, nil
}
