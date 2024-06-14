package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
)

func FindDns(config aws.Config, region string, searchValue string) map[string][]types.ResourceRecordSet {

	config.Region = region

	route53Client := route53.NewFromConfig(config)
	hostedZones, err := route53Client.ListHostedZones(context.TODO(), &route53.ListHostedZonesInput{})
	if err != nil {
		var accessDeniedErr *http.ResponseError
		if errors.As(err, &accessDeniedErr) && accessDeniedErr.HTTPStatusCode() == 403 {
			return nil
		}
		fmt.Printf("Unable to list hosted zones, %v", err)
		return nil
	}

	searchValue = strings.ToLower(searchValue)

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
