package aws_search

import (
	"context"
	"fmt"
	"log"
	"sync"

	general "github.com/aviadhaham/cloudcate-service/internal/aws_search/search/general"
	services "github.com/aviadhaham/cloudcate-service/internal/aws_search/search/services"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

func findResourcesInRegion(profile string, cfg aws.Config, region string, resourceType string, resourceName string) ([]interface{}, error) {
	associatedAwsAccount := general.GetAwsAccount(cfg, region)
	var results []interface{}

	switch resourceType {
	case "loadbalancer":
		lbSlice, _ := services.FindLoadBalancer(cfg, region, resourceName)
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
		bucketsSlice := services.FindS3Bucket(cfg, region, resourceName)
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
		dnsRecordsMap := services.FindDns(cfg, region, resourceName)
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

func FindResources(profiles []string, regions []string, resourceGlobality map[string]bool, resourceType string, resourceName string) ([]interface{}, error) {
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
