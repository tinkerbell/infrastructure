package assets

import (
	"github.com/pulumi/pulumi-ns1/sdk/go/ns1"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
	"github.com/tinkerbell/infrastructure/src/internal"
)

type DNS struct {
	Cname *ns1.Record
}

func CreateAssetsDNS(ctx *pulumi.Context, infrastructure internal.Infrastructure) (DNS, error) {
	// This verification record is required for Cloudflare to issue a certificate for this domain
	_, err := ns1.NewRecord(ctx, "assets-cname-verification", &ns1.RecordArgs{
		Zone:   pulumi.String(infrastructure.Zone.Zone),
		Domain: pulumi.String("_f213f2a073773bab6da1e978b94e1a92.assets.tinkerbell.org"),
		Type:   pulumi.String("CNAME"),
		Answers: ns1.RecordAnswerArray{
			ns1.RecordAnswerArgs{
				Answer: pulumi.String("_69acff72294eb94eaac26bef9c1fdb06.zjfbrrwmzc.acm-validations.aws"),
			},
		},
	})
	if err != nil {
		return DNS{}, err
	}

	// assets.tinkerbell.org
	// Because this Cloudfront/S3 distribution isn't configured on a CNCF
	// account, it's been created outwidth this automation.
	// So we're just linking them together. It's not ideal.
	cname, err := ns1.NewRecord(ctx, "assets-cname", &ns1.RecordArgs{
		Zone:   pulumi.String(infrastructure.Zone.Zone),
		Domain: pulumi.String("assets.tinkerbell.org"),
		Type:   pulumi.String("CNAME"),
		Answers: ns1.RecordAnswerArray{
			ns1.RecordAnswerArgs{
				Answer: pulumi.String("d3b2o8hr2qjhyx.cloudfront.net"),
			},
		},
	})
	if err != nil {
		return DNS{}, err
	}

	return DNS{
		Cname: cname,
	}, nil
}
