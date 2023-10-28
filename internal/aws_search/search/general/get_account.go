package aws_search

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

func GetAwsAccount(cfg aws.Config, region string) string {
	cfg.Region = region

	stsClient := sts.NewFromConfig(cfg)
	identity, err := stsClient.GetCallerIdentity(context.TODO(), &sts.GetCallerIdentityInput{})
	if err != nil {
		var accessDeniedErr *awshttp.ResponseError
		if errors.As(err, &accessDeniedErr) && accessDeniedErr.HTTPStatusCode() == 403 {
			return ""
		}
		fmt.Printf("\nUnable to get caller identity: %v\nRegion: %s\n", err, region)
	}
	return *identity.Account
}
