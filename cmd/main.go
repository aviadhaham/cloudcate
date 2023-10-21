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

	"github.com/gin-gonic/gin"

	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type SearchResult struct {
	Account string `json:"account"`
	Profile string `json:"profile"`
	Region  string `json:"region"`
}

type LoadBalancerSearchResult struct {
	SearchResult
	LoadBalancerArn string `json:"load_balancer_arn"`
}

type S3SearchResult struct {
	SearchResult
	BucketName string `json:"bucket_name"`
}

type DNSSearchResult struct {
	SearchResult
	HostedZoneName string `json:"hosted_zone_name"`
	DnsRecordName  string `json:"dns_record_name"`
	DnsRecordType  string `json:"dns_record_type"`
}

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

func findResourcesInRegion(profile string, cfg aws.Config, region string, resourceType string, resourceName string) ([]interface{}, error) {
	associatedAwsAccount := getAwsAccount(cfg, region)
	var results []interface{}

	switch resourceType {
	case "loadbalancer":
		lbSlice, _ := findLoadBalancer(cfg, region, resourceName)
		for _, lb := range lbSlice {
			results = append(results, LoadBalancerSearchResult{
				SearchResult: SearchResult{
					Account: associatedAwsAccount,
					Profile: profile,
					Region:  region,
				},
				LoadBalancerArn: lb,
			})
		}
	case "s3":
		bucketsSlice := findS3Bucket(cfg, region, resourceName)
		if bucketsSlice == nil {
			return nil, fmt.Errorf("no S3 buckets found")
		}

		for _, bucket := range bucketsSlice {
			results = append(results, S3SearchResult{
				SearchResult: SearchResult{
					Account: associatedAwsAccount,
					Profile: profile,
					Region:  region,
				},
				BucketName: bucket,
			})
		}
	case "dns":
		dnsRecordsMap := findDns(cfg, region, resourceName)
		if len(dnsRecordsMap) == 0 {
			return nil, fmt.Errorf("no DNS records found")
		}
		for zoneName, dnsRecords := range dnsRecordsMap {
			for _, dnsRecord := range dnsRecords {
				results = append(results, DNSSearchResult{
					SearchResult: SearchResult{
						Account: associatedAwsAccount,
						Profile: profile,
						Region:  region,
					},
					HostedZoneName: zoneName,
					DnsRecordName:  aws.ToString(dnsRecord.Name),
				})
			}
		}
	}
	return results, nil
}

func findResources(profiles []string, regions []string, resourceGlobality map[string]bool, resourceType string, resourceName string) ([]interface{}, error) {
	var results []interface{}
	var wg sync.WaitGroup
	resultChan := make(chan []interface{})

	for _, profile := range profiles {
		cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(profile))
		if err != nil {
			return nil, fmt.Errorf("failed to load configuration for profile, %v", err)
		}

		if resourceGlobality[resourceType] {
			wg.Add(1)
			go func(profile string, cfg aws.Config, region string) {
				defer wg.Done()
				res, err := findResourcesInRegion(profile, cfg, region, resourceType, resourceName)
				if err != nil {
					log.Printf("error searching for resources in region %s: %v", region, err)
					return
				}
				resultChan <- res
			}(profile, cfg, regions[0])
			continue
		}

		for _, region := range regions {
			wg.Add(1)
			go func(profile string, cfg aws.Config, region string) {
				defer wg.Done()
				res, err := findResourcesInRegion(profile, cfg, region, resourceType, resourceName)
				if err != nil {
					log.Printf("error searching for resources in region %s: %v", region, err)
					return
				}
				resultChan <- res
			}(profile, cfg, region)
		}
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for res := range resultChan {
		results = append(results, res...)
	}

	return results, nil
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

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Next()
	})

	r.GET("/search", func(c *gin.Context) {
		resourceName := c.Query("resource_name")
		resourceType := c.Query("resource_type")

		results, err := findResources(profiles, regions, resourceGlobality, resourceType, resourceName)
		if err != nil {
			log.Fatalf("Failed to search resources: %v", err)
		}

		// Convert the results to a slice of maps
		// var data []map[string]interface{}
		// for _, result := range results {
		// 	data = append(data, map[string]interface{}{
		// 		"account": result.(S3SearchResult).Account,
		// 		"profile": result.(S3SearchResult).Profile,
		// 		"region":  result.(S3SearchResult).Region,
		// 		"bucket":  result.(S3SearchResult).BucketName,
		// 	})
		// }

		c.JSON(http.StatusOK, gin.H{
			"results": results,
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
