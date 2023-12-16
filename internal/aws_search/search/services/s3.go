package aws_search

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func FindS3Bucket(config aws.Config, region string, searchValue string) []string {
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

	searchValue = strings.ToLower(searchValue)

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
