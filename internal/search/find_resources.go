package search

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/aviadhaham/cloudcate/internal/config"
	"github.com/aviadhaham/cloudcate/internal/services"

	"github.com/aws/aws-sdk-go-v2/aws"
	aws_config "github.com/aws/aws-sdk-go-v2/config"
)

func findResourcesInRegion(profile string, cfg aws.Config, region string, resourceSubType string, resourceType string, resourceName string) ([]interface{}, error) {
	associatedAwsAccount := GetAwsAccount(cfg, region)
	var results []interface{}

	switch resourceType {
	case "loadbalancer":
		lbSlice, _ := services.FindLoadBalancer(cfg, region, resourceName)
		for _, lb := range lbSlice {
			results = append(results, LoadBalancerSearchResult{
				SearchResultNonGlobal: SearchResultNonGlobal{
					SearchResult: SearchResult{
						Account: associatedAwsAccount,
						Profile: profile,
					},
					Region: region,
				},
				LoadBalancerName:    lb[0],
				LoadBalancerDnsName: lb[1],
			})
		}
	case "ec2":
		instances, err := services.FindEc2(cfg, region, resourceName)
		if err != nil {
			return nil, fmt.Errorf("error finding EC2 instances: %v", err)
		}
		if len(instances) == 0 {
			return nil, nil
		}
		for _, instance := range instances {
			ec2SearchResult := Ec2SearchResult{
				SearchResultNonGlobal: SearchResultNonGlobal{
					SearchResult: SearchResult{
						Account: associatedAwsAccount,
						Profile: profile,
					},
					Region: region,
				},
				InstanceId: *instance.InstanceId,
			}

			if instance.Tags != nil {
				for _, tag := range instance.Tags {
					if *tag.Key == "Name" {
						ec2SearchResult.InstanceName = *tag.Value
						break
					}
				}
			}

			if instance.PrivateIpAddress != nil {
				ec2SearchResult.PrivateIpAddress = *instance.PrivateIpAddress
			}
			if instance.PrivateDnsName != nil {
				ec2SearchResult.PrivateDnsName = *instance.PrivateDnsName
			}
			if instance.PublicDnsName != nil {
				ec2SearchResult.PublicDnsName = *instance.PublicDnsName
			}
			if instance.PublicIpAddress != nil {
				ec2SearchResult.PublicIpAddress = *instance.PublicIpAddress
			}

			results = append(results, ec2SearchResult)
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
					},
					HostedZoneName: zoneName,
					DnsRecordName:  aws.ToString(dnsRecord.Name),
				})
			}
		}
	case "iam":
		if resourceSubType == "user" {
			users, err := services.FindIamUser(cfg, region, resourceName)
			if err != nil {
				log.Printf("profile '%s': %v\n", profile, err)
			}
			if users == nil {
				return nil, fmt.Errorf("no IAM users found")
			}
			for _, user := range users {
				if user != "" {
					results = append(results, IamUserSearchResult{
						SearchResult: SearchResult{
							Account: associatedAwsAccount,
							Profile: profile,
						},
						UserName: user,
					})
				}
			}
		}
		if resourceSubType == "key" {
			accessKeys, err := services.FindIamUserKey(cfg, region, resourceName)
			if err != nil {
				log.Printf("profile '%s': %v\n", profile, err)
			}
			if accessKeys == nil {
				return nil, fmt.Errorf("no IAM user access keys users found")
			}
			for user, key := range accessKeys {
				results = append(results, IamUserKeySearchResult{
					IamUserSearchResult: IamUserSearchResult{
						SearchResult: SearchResult{
							Account: associatedAwsAccount,
							Profile: profile,
						},
						UserName: user,
					},
					AccessKey: key,
				})
			}
		}

	case "elastic_ip":
		addresses, err := services.FindElasticIp(cfg, region, resourceName)
		if err != nil {
			return nil, fmt.Errorf("error finding elastic IP addresses: %v", err)
		}

		for _, address := range addresses {
			elasticIpSearchResult := ElasticIpSearchResult{
				SearchResultNonGlobal: SearchResultNonGlobal{
					SearchResult: SearchResult{
						Account: associatedAwsAccount,
						Profile: profile,
					},
					Region: region,
				},
				PublicIp: *address.PublicIp,
			}

			if address.InstanceId != nil && *address.InstanceId != "" {
				elasticIpSearchResult.InstanceId = *address.InstanceId
			}

			results = append(results, elasticIpSearchResult)
		}
	case "cloudfront":
		distributions, err := services.FindCloudfront(cfg, region, resourceName)
		if err != nil && len(distributions) == 0 {
			return nil, fmt.Errorf("error finding cloudfront distributions: %v", err)
		}

		for _, distribution := range distributions {
			cloudFrontSearchResult := CloudfrontSearchResult{
				SearchResult: SearchResult{
					Account: associatedAwsAccount,
					Profile: profile,
				},
				DistributionArn: *distribution.ARN,
				DistributionId:  *distribution.Id,
				DomainName:      *distribution.DomainName,
			}

			results = append(results, cloudFrontSearchResult)
		}
	}
	return results, nil
}

func FindResources(profiles []string, servicesGlobality map[string]bool, resourceType string, resourceSubType string, resourceName string) ([]interface{}, error) {
	var results []interface{}
	var wg sync.WaitGroup
	resultChan := make(chan []interface{})

	for _, profile := range profiles {
		cfg, err := aws_config.LoadDefaultConfig(context.TODO(), aws_config.WithSharedConfigProfile(profile))
		if err != nil {
			return nil, fmt.Errorf("failed to load configuration for profile, %v", err)
		}

		regions, err := GetRegions(profile)
		if err != nil {
			log.Fatalf("Failed to get regions: %v", err)
		}

		if config.ServicesGlobality[resourceType] {
			wg.Add(1)
			go func(profile string, cfg aws.Config, region string) {
				defer wg.Done()
				res, err := findResourcesInRegion(profile, cfg, region, resourceSubType, resourceType, resourceName)
				if err != nil {
					log.Printf("profile '%s', error searching for resources in region %s: %v", profile, region, err)
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
				res, err := findResourcesInRegion(profile, cfg, region, resourceSubType, resourceType, resourceName)
				if err != nil {
					log.Printf("profile '%s', error searching for resources in region %s: %v", profile, region, err)
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
