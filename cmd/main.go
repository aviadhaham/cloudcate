package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go/aws"

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

func main() {
	profiles, err := getProfiles()
	if err != nil {
		log.Fatalf("Failed to get profiles: %v", err)
	}

	regions, err := getRegions()
	if err != nil {
		log.Fatalf("Failed to get regions: %v", err)
	}

	// searchTerm := "git-pages"

	for _, profile := range profiles {
		cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(profile))
		if err != nil {
			log.Fatalf("failed to load configuration for profile, %v", err)
		}

		for _, region := range regions {

			cfg.Region = region

			elbv2Client := elasticloadbalancingv2.NewFromConfig(cfg)
			input := &elasticloadbalancingv2.DescribeLoadBalancersInput{}
			output, err := elbv2Client.DescribeLoadBalancers(context.TODO(), input)
			if err != nil {
				fmt.Println(err)
			}

			jsonData, err := json.Marshal(output)
			if err != nil {
				fmt.Println("Error: ", err)
			}
			jsonString := string(jsonData)

			// for _, lb :=

			fmt.Printf("Profile: %s, Region: %s\n", profile, region)
			fmt.Println(jsonString, "\n******split*****\n")
		}
	}

	// fetch the aws account from the LoadBalancerArn field returned in each loadbalancer
}
