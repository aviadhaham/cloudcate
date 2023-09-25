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

func getRegions() ([]string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

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

	if output != nil {
		loadBalancers := output.LoadBalancers
		for _, lb := range loadBalancers {
			if strings.Contains(*lb.LoadBalancerArn, searchValue) {
				lbArnSlice := strings.Split(*lb.LoadBalancerArn, ":")
				return lbArnSlice, nil
			}
		}
	}

	return nil, err
}

func findS3Bucket(config aws.Config, region string, searchValue string) string {
	config.Region = region

	s3Client := s3.NewFromConfig(config)
	output, err := s3Client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		var accessDeniedErr *awshttp.ResponseError
		if errors.As(err, &accessDeniedErr) && accessDeniedErr.HTTPStatusCode() == 403 {
			return ""
		}
		fmt.Printf("Unable to list buckets, %v", err)
		return ""
	}

	for _, bucket := range output.Buckets {
		if strings.Contains(*bucket.Name, searchValue) {
			return *bucket.Name
		}
	}
	return ""
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
		lbArnSlice, _ := findLoadBalancer(cfg, region, resourceName)
		if lbArnSlice != nil {
			fmt.Printf("Found LB:\n%s\n    in Region: %s\n    in AWS Account: %s (profile '%s')\n\n", lbArnSlice[5], region, associatedAwsAccount, profile)
		}
	case "s3":
		bucketName := findS3Bucket(cfg, region, resourceName)
		if bucketName != "" {
			fmt.Printf("S3 bucket: %s -> AWS account: %s (profile '%s')\n", bucketName, associatedAwsAccount, profile)
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

	// regions := []string{"us-east-2, us-east-1, us-west-1, us-west-2, af-south-1, ap-east-1, ap-south-2, ap-southeast-3, ap-southeast-4, ap-south-1, ap-northeast-3, ap-northeast-2, ap-southeast-1, ap-southeast-2, ap-northeast-1, ca-central-1, eu-central-1, eu-west-1, eu-west-2, eu-south-1, eu-west-3, eu-south-2, eu-north-1, eu-central-2, il-central-1, me-south-1, me-central-1, sa-east-1"}
	regions := []string{"us-east-1", "us-west-1", "us-west-2", "af-south-1", "ap-east-1", "ap-south-2", "ap-southeast-3", "ap-southeast-4", "ap-south-1", "ap-northeast-3", "ap-northeast-2", "ap-southeast-1", "ap-southeast-2", "ap-northeast-1", "ca-central-1", "eu-central-1", "eu-west-1", "eu-west-2", "eu-south-1", "eu-west-3", "eu-south-2", "eu-north-1", "eu-central-2", "il-central-1", "me-south-1", "me-central-1", "sa-east-1"}

	resourceGlobality := map[string]bool{
		"loadbalancer": false,
		"s3":           true,
		"dns":          true,
	}

	// resourceName := "1954630134"
	// resourceType := "loadbalancer"

	//resourceName := "marketplace"
	//resourceType := "s3"

	resourceName := "rancher"
	resourceType := "dns"

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
