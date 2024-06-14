package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func FindElasticIp(config aws.Config, region string, searchValue string) ([]types.Address, error) {
	config.Region = region

	ec2Client := ec2.NewFromConfig(config)
	input := &ec2.DescribeAddressesInput{}
	output, err := ec2Client.DescribeAddresses(context.TODO(), input)

	if err != nil {
		fmt.Printf("Unable to list elastic IPs, %v", err)
	}

	filteredElasticIps := []types.Address{}
	if output != nil {
		for _, address := range output.Addresses {
			if strings.Contains(*address.PublicIp, searchValue) {
				filteredElasticIps = append(filteredElasticIps, address)
				// fmt.Println(*address.PublicIp)
			}
		}
	}

	return filteredElasticIps, err
}
