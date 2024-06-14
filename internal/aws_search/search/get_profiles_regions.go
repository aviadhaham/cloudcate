package aws_search

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	config "github.com/aviadhaham/cloudcate/internal/aws_search/config"
	aws_config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go/aws"
)

func GetProfiles() ([]string, error) {
	// read .aws/credetials file
	f, err := os.Open(fmt.Sprintf("%s/.aws/credentials", os.Getenv("HOME")))
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	profileList := []string{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "[") && strings.Contains(line, "]") {
			profile := strings.Trim(scanner.Text(), "[]")
			profileList = append(profileList, profile)
		}
	}

	return profileList, err
}

func GetRegions(profile string) ([]string, error) {

	// use .aws/credentials file to get profiles, but use only the first one in the file
	// hardcoded region to us-east-1, because there's no chance it's not going to be active
	cfg, err := aws_config.LoadDefaultConfig(context.TODO(), aws_config.WithSharedConfigProfile(profile))
	if err != nil {
		log.Fatalf("failed to load configuration for profile, %v", err)
	}
	cfg.Region = "us-east-1"
	client := ec2.NewFromConfig(cfg)

	input := &ec2.DescribeRegionsInput{
		Filters: []ec2types.Filter{
			{
				Name:   aws.String("opt-in-status"),
				Values: []string{"opt-in-not-required", "opted-in"},
			},
		},
	}

	resp, err := client.DescribeRegions(context.TODO(), input)
	if err != nil {
		log.Printf("profile '%s', failed to describe regions, %v", profile, err)
		// Return a hardcoded list of regions instead of terminating the application
		return config.AwsFullRegionsList, nil
	}

	regionsList := []string{}
	for _, region := range resp.Regions {
		regionsList = append(regionsList, *region.RegionName)
	}

	return regionsList, nil
}
