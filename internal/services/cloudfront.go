package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
)

func FindCloudfront(config aws.Config, region string, searchValue string) ([]types.DistributionSummary, error) {
	config.Region = region

	cfClient := cloudfront.NewFromConfig(config)
	output, err := cfClient.ListDistributions(context.TODO(), &cloudfront.ListDistributionsInput{})
	if err != nil {
		var accessDeniedErr *http.ResponseError
		if errors.As(err, &accessDeniedErr) && accessDeniedErr.HTTPStatusCode() == 403 {
			return nil, err
		}
		fmt.Printf("Unable to list cloudfront distributions %v", err)
		return nil, err
	}

	filteredCfDistributions := []types.DistributionSummary{}
	if output != nil {
		for _, distribution := range output.DistributionList.Items {
			if strings.Contains(*distribution.DomainName, searchValue) || strings.Contains(*distribution.Id, searchValue) {
				filteredCfDistributions = append(filteredCfDistributions, distribution)
			}
		}
	}

	return filteredCfDistributions, err
}
