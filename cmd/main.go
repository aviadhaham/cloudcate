package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
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
		AllRegions: aws.Bool(true),
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

func findLoadBalancer(config aws.Config, region string, searchValue string) ([]string, error) {
	config.Region = region

	elbv2Client := elasticloadbalancingv2.NewFromConfig(config)
	input := &elasticloadbalancingv2.DescribeLoadBalancersInput{}
	output, err := elbv2Client.DescribeLoadBalancers(context.TODO(), input)
	if err != nil && output == nil {
		if strings.Contains(err.Error(), "InvalidClientTokenId") || strings.Contains(err.Error(), "no identity-based policy allows the elasticloadbalancing:DescribeLoadBalancers action") {
			return nil, err
		}
		// fmt.Println("error describing load balancers", err)
	}

	loadBalancers := output.LoadBalancers
	for _, lb := range loadBalancers {
		if strings.Contains(*lb.LoadBalancerArn, searchValue) {
			lbArnSlice := strings.Split(*lb.LoadBalancerArn, ":")
			return lbArnSlice, nil
		}
	}

	return nil, err
}

func findResourceInRegion(profile string, cfg aws.Config, region string, resourceType string, resourceName string) {
	switch resourceType {
	case "loadbalancer":
		lbArnSlice, _ := findLoadBalancer(cfg, region, resourceName)
		// if lbArnSlice == nil {
		// 	fmt.Printf("no load balancer was found: %s", resourceName)
		// }
		// if err != nil {
		// 	fmt.Printf("%s", err)
		// }
		if lbArnSlice != nil {
			fmt.Printf("Region: %s\nAWS Account: %s\nLB Details: %s", lbArnSlice[3], lbArnSlice[4], lbArnSlice[5])
		}
	}
}

func main() {
	profiles, err := getProfiles()
	if err != nil {
		log.Fatalf("Failed to get profiles: %v", err)
	}

	regions, err := getRegions()
	if err != nil {
		log.Fatalf("Failed to get regions: %v", err)
	}

	// resourceSearchFunctions := map[string]interface{}{
	// 	"LB": findLoadBalancer,
	// }

	resourceName := "common-svc-np-nginx-v1"
	resourceType := "loadbalancer"

	var wg sync.WaitGroup

	for _, profile := range profiles {
		cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(profile))
		if err != nil {
			log.Fatalf("failed to load configuration for profile, %v", err)
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
