package aws_search

import (
	"context"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func FindEc2(config aws.Config, region string, searchValue string) ([]types.Instance, error) {
	config.Region = region

	ec2Client := ec2.NewFromConfig(config)
	input := &ec2.DescribeInstancesInput{}

	output, err := ec2Client.DescribeInstances(context.TODO(), input)

	if err != nil {
		log.Printf("Unable to list instances, %v", err)
		return nil, err
	}

	filteredInstances := []types.Instance{}
	if output != nil {
		for _, reservation := range output.Reservations {
			for _, instance := range reservation.Instances {
				if instance.PrivateDnsName != nil && strings.Contains(*instance.PrivateDnsName, searchValue) ||
					instance.PrivateIpAddress != nil && strings.Contains(*instance.PrivateIpAddress, searchValue) ||
					instance.PublicDnsName != nil && strings.Contains(*instance.PublicDnsName, searchValue) ||
					instance.PublicIpAddress != nil && strings.Contains(*instance.PublicIpAddress, searchValue) {
					filteredInstances = append(filteredInstances, instance)
				}
				for _, tag := range instance.Tags {
					if tag.Key != nil && strings.Contains(*tag.Key, searchValue) || tag.Value != nil && strings.Contains(*tag.Value, searchValue) {
						filteredInstances = append(filteredInstances, instance)
						break
					}
				}
			}
		}
	}

	return filteredInstances, nil
}
