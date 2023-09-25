package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"

	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

func getProfiles() ([]string, error) {
	// read .aws/creds file
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

func getRegions(profile string) ([]string, error) {

	// use .aws/credentials file to get profiles, but use only the first one in the file
	// hardcoded region to us-east-1, because there's no chance it's not going to be active
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(profile))
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
		log.Fatalf("failed to describe regions, %v", err)
		return nil, err
	}

	regionsList := []string{}
	for _, region := range resp.Regions {
		regionsList = append(regionsList, *region.RegionName)
	}

	return regionsList, err
}

func getAwsAccount(cfg aws.Config, region string) string {
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

func findLoadBalancer(config aws.Config, region string, searchValue string) ([]string, error) {
	config.Region = region

	elbv2Client := elasticloadbalancingv2.NewFromConfig(config)
	input := &elasticloadbalancingv2.DescribeLoadBalancersInput{}
	output, err := elbv2Client.DescribeLoadBalancers(context.TODO(), input)
	if err != nil {
		var accessDeniedErr *awshttp.ResponseError
		if errors.As(err, &accessDeniedErr) && accessDeniedErr.HTTPStatusCode() == 403 {
			return nil, nil
		}
		fmt.Printf("Unable to list load balancers, %v", err)
	}

	filteredLoadBalancers := []string{}
	if output != nil {
		for _, lb := range output.LoadBalancers {
			if strings.Contains(*lb.LoadBalancerArn, searchValue) {
				filteredLoadBalancers = append(filteredLoadBalancers, *lb.LoadBalancerArn)
			}
		}
	}

	return filteredLoadBalancers, nil
}

func findS3Bucket(config aws.Config, region string, searchValue string) []string {
	config.Region = region

	s3Client := s3.NewFromConfig(config)
	output, err := s3Client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		var accessDeniedErr *awshttp.ResponseError
		if errors.As(err, &accessDeniedErr) && accessDeniedErr.HTTPStatusCode() == 403 {
			return nil
		}
		fmt.Printf("Unable to list buckets, %v", err)
		return nil
	}

	filteredS3Buckets := []string{}
	if output != nil {
	for _, bucket := range output.Buckets {
		if strings.Contains(*bucket.Name, searchValue) {
				filteredS3Buckets = append(filteredS3Buckets, *bucket.Name)
			}
		}
	}
	return filteredS3Buckets
}

func findDns(config aws.Config, region string, searchValue string) map[string][]types.ResourceRecordSet {

	config.Region = region

	route53Client := route53.NewFromConfig(config)
	hostedZones, err := route53Client.ListHostedZones(context.TODO(), &route53.ListHostedZonesInput{})
	if err != nil {
		var accessDeniedErr *awshttp.ResponseError
		if errors.As(err, &accessDeniedErr) && accessDeniedErr.HTTPStatusCode() == 403 {
			return nil
		}
		fmt.Printf("Unable to list hosted zones, %v", err)
		return nil
	}

	filteredRecordsMap := make(map[string][]types.ResourceRecordSet)
	for _, zone := range hostedZones.HostedZones {
		recordSets, err := route53Client.ListResourceRecordSets(context.TODO(), &route53.ListResourceRecordSetsInput{
			HostedZoneId: zone.Id,
		})
		if err != nil {
			fmt.Println("Error listing resource record sets:", err)
			continue
		}

		for _, record := range recordSets.ResourceRecordSets {
			if strings.Contains(aws.ToString(record.Name), searchValue) {
				filteredRecordsMap[*zone.Name] = append(filteredRecordsMap[*zone.Name], record)
			}

		}
	}
	return filteredRecordsMap
}

func findResourceInRegion(profile string, cfg aws.Config, region string, resourceType string, resourceName string) {
	associatedAwsAccount := getAwsAccount(cfg, region)

	switch resourceType {
	case "loadbalancer":
		lbSlice, _ := findLoadBalancer(cfg, region, resourceName)
		for _, lb := range lbSlice {
			fmt.Printf("\nFound LB:\n%s\n    in Region: %s\n    in AWS Account: %s (profile '%s')\n\n", lb, region, associatedAwsAccount, profile)
		}
	case "s3":
		bucketsSlice := findS3Bucket(cfg, region, resourceName)
		for _, bucket := range bucketsSlice {
			fmt.Printf("\nFound S3 bucket: %s\n    in Region: %s\n    in AWS Account: %s (profile '%s')\n\n", bucket, region, associatedAwsAccount, profile)
		}
	case "dns":
		dnsRecordsMap := findDns(cfg, region, resourceName)
		if len(dnsRecordsMap) == 0 {
			fmt.Printf("\n%s\n************** No DNS records were found for '%s' in AWS account '%s' (%s) **************\n", strings.Repeat("-", 120), resourceName, associatedAwsAccount, profile)
			break
		}
		for zoneName, dnsRecords := range dnsRecordsMap {
			fmt.Printf("\n%s\n************** Zone name %s (from AWS account '%s' - profile '%s') **************\n", strings.Repeat("-", 120), zoneName, associatedAwsAccount, profile)
			for _, dnsRecord := range dnsRecords {
				fmt.Printf("DNS record '%s' || type '%s'\n", *dnsRecord.Name, dnsRecord.Type)
			}
		}
	}
}

func main() {
	profiles, err := getProfiles()
	if err != nil {
		log.Fatalf("Failed to get profiles: %v", err)
	}

	regions, err := getRegions(profiles[0])
	if err != nil {
		log.Fatalf("Failed to get regions: %v", err)
	}

	resourceGlobality := map[string]bool{
		"loadbalancer": false,
		"s3":           true,
		"dns":          true,
	}

	var resourceName string
	var resourceType string
	fmt.Print("\nEnter resource string/substring you wish to search name: ")
	fmt.Scanln(&resourceName)
	fmt.Print("\nEnter the resource type [e.g., dns, s3, loadbalancer]: ")
	fmt.Scanln(&resourceType)

	var wg sync.WaitGroup

	for _, profile := range profiles {
		cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(profile))
		if err != nil {
			log.Fatalf("failed to load configuration for profile, %v", err)
		}

		if resourceGlobality[resourceType] {
			wg.Add(1)
			go func(profile string, cfg aws.Config, region string) {
				defer wg.Done()
				findResourceInRegion(profile, cfg, region, resourceType, resourceName)
			}(profile, cfg, regions[0])
			continue
		}

		for _, region := range regions {
			wg.Add(1)
			go func(profile string, cfg aws.Config, region string) {
				defer wg.Done()
				findResourceInRegion(profile, cfg, region, resourceType, resourceName)
			}(profile, cfg, region)
		}
	}
	wg.Wait()
}
