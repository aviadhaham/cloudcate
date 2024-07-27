package services

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func FindVpc(config aws.Config, region string, searchValue string) ([]types.Vpc, error) {
	config.Region = region

	ec2Client := ec2.NewFromConfig(config)
	input := &ec2.DescribeVpcsInput{}
	filteredVpcs := []types.Vpc{}

	paginator := ec2.NewDescribeVpcsPaginator(ec2Client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			return nil, err
		}

		for _, vpc := range page.Vpcs {
			if vpc.CidrBlock != nil && strings.Contains(*vpc.CidrBlock, searchValue) || strings.Contains(*vpc.VpcId, searchValue) {
				filteredVpcs = append(filteredVpcs, vpc)
			} else {
				for _, tag := range vpc.Tags {
					if tag.Key != nil && strings.Contains(strings.ToLower(*tag.Key), searchValue) || tag.Value != nil && strings.Contains(strings.ToLower(*tag.Value), searchValue) {
						filteredVpcs = append(filteredVpcs, vpc)
						break
					}
				}
			}
		}
	}

	return filteredVpcs, nil
}
