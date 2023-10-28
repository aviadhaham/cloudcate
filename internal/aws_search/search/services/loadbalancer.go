package aws_search

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
)

func FindLoadBalancer(config aws.Config, region string, searchValue string) ([]string, error) {
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
			if strings.Contains(*lb.LoadBalancerArn, searchValue) || strings.Contains(*lb.DNSName, searchValue) {
				filteredLoadBalancers = append(filteredLoadBalancers, *lb.LoadBalancerArn)
				fmt.Println(*lb.LoadBalancerArn)
			}
		}
	}

	return filteredLoadBalancers, nil
}
